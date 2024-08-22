package rpc

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/ucenter/register"
)

var registerCli register.Client

func initRegister() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.40.134:2379"})
	if err != nil {
		panic(err)
	}
	registerCli, err = register.NewClient("ucenter", client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetRegisterClient() register.Client {
	return registerCli
}
