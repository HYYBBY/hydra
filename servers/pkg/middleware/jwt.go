package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/servers/pkg/swap"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtAuth jwt
func JwtAuth(f jwt.IJWTAuth) dispatcher.HandlerFunc {
	return func(r swap.IRequest) {
		jwtAuth, ok := f.GetConf()
		if !ok || jwtAuth == nil || jwtAuth.Disable {
			r.Next()
			return
		}

		//检查jwt.token是否正确
		data, err := checkJWT(ctx, jwtAuth.Name, jwtAuth.Secret)
		if err == nil {
			setJWTRaw(ctx, data)
			r.Next()
			return
		}

		//不需要校验的URL自动跳过
		url := ctx.Request.GetService()
		for _, u := range jwtAuth.Exclude {
			if u == url {
				r.Next()
				return
			}
		}

		//jwt.token错误，返回错误码
		r.GetLogger().Error(err.GetError())
		r.Abort(err.GetCode())
		return

	}
}
func setJwtResponse(ctx *dispatcher.Context, cnf *conf.MetadataConf, data interface{}) {
	if data == nil {
		data = getJWTRaw(ctx)
	}
	if data == nil {
		return
	}
	jwtAuth, ok := cnf.GetMetadata("jwt").(*conf.JWTAuth)
	if !ok || jwtAuth.Disable {
		return
	}
	jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
	if err != nil {
		ctx.AbortWithError(500, fmt.Errorf("jwt配置出错：%v", err))
		return
	}
	expireTime := time.Now().Add(time.Duration(time.Duration(jwtAuth.ExpireAt)*time.Second - 8*60*60*time.Second))
	expireVal := expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	//fmt.Println("jwtToken:", jwtToken)
	ctx.Header("Set-Cookie", fmt.Sprintf("%s=%s;path=/;expires=%s;", jwtAuth.Name, jwtToken, expireVal))
	ctx.Header(jwtAuth.Name, jwtToken)
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx *dispatcher.Context, name string, secret string) (data interface{}, err context.IError) {
	token := getToken(ctx, name)
	if token == "" {
		return nil, context.NewError(403, fmt.Errorf("%s未传入jwt.token", name))
	}
	data, er := jwt.Decrypt(token, secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, context.NewError(401, er)
		}
		return data, context.NewError(403, er)
	}
	return data, nil
}
func getToken(ctx *dispatcher.Context, key string) string {
	if cookie, ok := ctx.Request.GetHeader()[key]; ok {
		return cookie
	}
	return ""
}
func setToken(ctx *dispatcher.Context, name string, token string) {
	ctx.Header(name, token)
}
