package rpc

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/market/market"
	"ucenter/config"
)

var marketCli market.Client

func initMarket() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	marketCli, err = market.NewClient("market", client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetMarketClient() market.Client {
	return marketCli
}
