package registry

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	//_ "github.com/micro-plat/hydra/registry/registry/etcd"
	_ "github.com/micro-plat/hydra/registry/registry/filesystem"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	_ "github.com/micro-plat/hydra/registry/registry/redis"
	_ "github.com/micro-plat/hydra/registry/registry/zookeeper"
	"github.com/micro-plat/lib4go/logger"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/test/mocks"
)

// 私有变量的测试
// func Test_getAddrByUserPass(t *testing.T) {
// 	type args struct {
// 		addr string
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantU       string
// 		wantP       string
// 		wantAddress string
// 		wantErr     bool
// 	}{
// 		{name: "正确格式的地址", args: args{addr: "root:123456@192.168.5.115:9091"}, wantU: "root", wantP: "123456", wantAddress: "192.168.5.115:9091", wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotU, gotP, gotAddress, err := getAddrByUserPass(tt.args.addr)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getAddrByUserPass() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotU != tt.wantU {
// 				t.Errorf("getAddrByUserPass() gotU = %v, want %v", gotU, tt.wantU)
// 			}
// 			if gotP != tt.wantP {
// 				t.Errorf("getAddrByUserPass() gotP = %v, want %v", gotP, tt.wantP)
// 			}
// 			if gotAddress != tt.wantAddress {
// 				t.Errorf("getAddrByUserPass() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
// 			}
// 		})
// 	}
// }

func TestParse(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name      string
		args      args
		wantProto string
		wantRaddr []string
		wantU     string
		wantP     string
		wantErr   bool
	}{
		{name: "解析zk地址", args: args{address: "zk://192.168.0.101"}, wantProto: "zk", wantRaddr: []string{"192.168.0.101"}, wantU: "", wantP: "", wantErr: false},
		{name: "解析多个zk地址", args: args{address: "zk://192.168.0.101,192.168.0.102"}, wantProto: "zk",
			wantRaddr: []string{"192.168.0.101", "192.168.0.102"}, wantU: "", wantP: "", wantErr: false},
		{name: "解析lm地址", args: args{address: "lm://."}, wantProto: "lm", wantRaddr: []string{"."}, wantU: "", wantP: "", wantErr: false},
		{name: "解析fs地址", args: args{address: "fs://../a/b/c"}, wantProto: "fs", wantRaddr: []string{"../a/b/c"}, wantU: "", wantP: "", wantErr: false},
		{name: "解析etcd地址", args: args{address: "etcd://192.168.0.101:9099"}, wantProto: "etcd", wantRaddr: []string{"192.168.0.101:9099"}, wantU: "", wantP: "", wantErr: false},
		{name: "解析redis地址", args: args{address: "redis://192.168.0.101:6379"}, wantProto: "redis", wantRaddr: []string{"192.168.0.101:6379"}, wantU: "", wantP: "", wantErr: false},
		{name: "解析带有用户名和密码的地址", args: args{address: "redis://root:123456@192.168.0.101:6379"}, wantProto: "redis",
			wantRaddr: []string{"192.168.0.101:6379"}, wantU: "root", wantP: "123456", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProto, gotRaddr, gotU, gotP, err := registry.Parse(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProto != tt.wantProto {
				t.Errorf("Parse() gotProto = %v, want %v", gotProto, tt.wantProto)
			}
			if !reflect.DeepEqual(gotRaddr, tt.wantRaddr) {
				t.Errorf("Parse() gotRaddr = %v, want %v", gotRaddr, tt.wantRaddr)
			}
			if gotU != tt.wantU {
				t.Errorf("Parse() gotU = %v, want %v", gotU, tt.wantU)
			}
			if gotP != tt.wantP {
				t.Errorf("Parse() gotP = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		elem []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "地址拼接", args: args{elem: []string{"", ""}}, want: ""},
		{name: "地址拼接", args: args{elem: []string{"", "a/"}}, want: "/a"},
		{name: "地址拼接", args: args{elem: []string{"a", "b", "c"}}, want: "/a/b/c"},
		{name: "地址拼接", args: args{elem: []string{"..", "a/b", "c/"}}, want: "/../a/b/c"},
		{name: "地址拼接", args: args{elem: []string{"..", "", "\\", "c/"}}, want: "/../c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registry.Join(tt.args.elem...); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRegistry(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()
	log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, meta).GetRequestID())
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantR   registry.IRegistry
		wantErr bool
		err     string
	}{
		{name: "获取zk的注册中心", args: args{address: "zk://192.168.0.101"}, wantErr: false},
		{name: "获取lm的注册中心", args: args{address: "lm://."}, wantErr: false},
		{name: "获取fs的注册中心", args: args{address: "fs://../"}, wantErr: true, err: "配置文件不存在:../registry.test.conf.toml stat ../registry.test.conf.toml: no such file or directory"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotR, err := registry.NewRegistry(tt.args.address, log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && err.Error() != tt.err {
				t.Errorf("NewRegistry() error = %v, wantErr %v", err.Error(), tt.err)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("NewRegistry() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}