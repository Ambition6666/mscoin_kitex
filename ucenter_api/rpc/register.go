package rpc

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/ucenter/register"
	"ucenter/config"
)

var registerCli register.Client

func initRegister() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	registerCli, err = register.NewClient(config.ServerName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetRegisterClient() *register.Client {
	return &registerCli
}
