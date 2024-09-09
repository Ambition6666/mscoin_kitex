package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/ucenter/login"
	"ucenter_api/config"
)

var loginCli login.Client

func initLogin() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	loginCli, err = login.NewClient(config.ServerName, client.WithResolver(r), client.WithTransportProtocol(transport.GRPC), client.WithMetaHandler(transmeta.ClientHTTP2Handler))

	if err != nil {
		panic(err)
	}
}

func GetLoginClient() login.Client {
	return loginCli
}
