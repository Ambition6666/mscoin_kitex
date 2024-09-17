package rpc

import (
	"exchange/config"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/market/market"
)

var marketCli market.Client

func initMarket() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	marketCli, err = market.NewClient(config.MARKET, client.WithResolver(r), client.WithTransportProtocol(transport.TTHeader), client.WithMetaHandler(transmeta.ClientTTHeaderHandler))
	if err != nil {
		panic(err)
	}
}

func GetMarketClient() market.Client {
	return marketCli
}
