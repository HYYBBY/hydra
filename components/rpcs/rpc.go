package rpcs

import (
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/lib4go/types"
)

//IRequest Component rpc
type IRequest = rpc.IRequest

//IComponentRPC Component Cache
type IComponentRPC interface {
	GetRegularRPC(names ...string) (c IRequest)
	GetRPC(names ...string) (c IRequest, err error)
}

//StandardRPC rpc服务
type StandardRPC struct {
	c          container.IContainer
	platName   string
	systemName string
}

//NewStandardRPC 创建RPC服务代理
func NewStandardRPC(c container.IContainer, platName string, systemName string) *StandardRPC {
	return &StandardRPC{
		c:          c,
		platName:   platName,
		systemName: systemName,
	}
}

//GetRegularRPC 获取正式的没有异常缓存实例
func (s *StandardRPC) GetRegularRPC(names ...string) (c IRequest) {
	c, err := s.GetRPC(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetRPC 获取缓存操作对象
func (s *StandardRPC) GetRPC(names ...string) (c IRequest, err error) {
	name := types.GetStringByIndex(names, 0, rpcNameNode)

	v, err := s.c.GetOrCreate(rpcTypeNode, name, func(i *conf.JSONConf) (interface{}, error) {
		return NewRequest(s.platName, s.systemName, name, i), nil
	})
	if err != nil {
		return nil, err
	}
	return v.(IRequest), nil
}
