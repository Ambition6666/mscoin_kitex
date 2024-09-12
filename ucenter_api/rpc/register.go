package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/ucenter/register"
	"ucenter_api/config"
)

var registerCli register.Client

func initRegister() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	registerCli, err = register.NewClient(config.ServerName, client.WithResolver(r), client.WithTransportProtocol(transport.TTHeader), client.WithMetaHandler(transmeta.ClientTTHeaderHandler))
	if err != nil {
		panic(err)
	}
}

func GetRegisterClient() register.Client {
	return registerCli
}
