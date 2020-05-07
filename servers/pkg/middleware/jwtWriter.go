package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/context"
	xjwt "github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtWriter 将jwt信息写入到请求中
func JwtWriter() Handler {
	return func(ctx IMiddleContext) {
		ctx.Next()
		conf := ctx.ServerConf().GetJWTConf()

		if conf.Disable {
			return
		}
		setJwtResponse(ctx, conf, ctx.User().GetJwt())
	}
}

func setJwtResponse(ctx context.IContext, jwtAuth *xjwt.JWTAuth, data interface{}) {
	if data == nil {
		return
	}
	jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
	if err != nil {
		ctx.Response().AbortWithError(500, fmt.Errorf("jwt配置出错：%v", err))
		return
	}
	setToken(ctx, jwtAuth, jwtToken)
}

//setToken 设置jwt到响应头或cookie中
func setToken(ctx context.IContext, jwt *xjwt.JWTAuth, token string) {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		ctx.Response().SetHeader(jwt.Name, token)
	default:
		expireTime := time.Now().Add(time.Duration(time.Duration(jwt.ExpireAt)*time.Second - 8*60*60*time.Second))
		expireVal := expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")

		if jwt.Domain != "" {
			ctx.Response().SetHeader("Set-Cookie", fmt.Sprintf("%s=%s;domain=%s;path=/;expires=%s;", jwt.Name, token, jwt.Domain, expireVal))
			return
		}
		ctx.Response().SetHeader("Set-Cookie", fmt.Sprintf("%s=%s;path=/;expires=%s;", jwt.Name, token, expireVal))
	}
}
