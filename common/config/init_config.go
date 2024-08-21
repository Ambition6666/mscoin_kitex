package config

import (
	"context"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-etcd/etcd"
	etcdServer "github.com/kitex-contrib/config-etcd/server"
)

func InitConfigClient(serviceName string, key string, uniqueId int64, etcdAddr string, cfg any) server.Suite {
	etcdClient, _ := etcd.NewClient(etcd.Options{Node: []string{etcdAddr}})
	suite := etcdServer.NewSuite(serviceName, etcdClient)

	etcdClient.RegisterConfigCallback(context.Background(), key, uniqueId, func(restoreDefault bool, data string, parser etcd.ConfigParser) {
		if !restoreDefault {
			if err := parser.Decode(data, cfg); err != nil {
				// TODO: 处理 error
				panic(err)
			}
		}

	})

	return suite
}
