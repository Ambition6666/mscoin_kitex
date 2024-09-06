package rpc

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"grpc_common/kitex_gen/market/exchangerate"
	"market/config"
)

var exchangeRateCli exchangerate.Client

func initExchangeRate() {
	r, err := etcd.NewEtcdResolver([]string{config.EtcdAddr})
	if err != nil {
		panic(err)
	}
	exchangeRateCli, err = exchangerate.NewClient(config.ServerName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetExchangeRateClient() exchangerate.Client {
	return exchangeRateCli
}
