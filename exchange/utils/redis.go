package utils

import (
	"common/database"
	"exchange/config"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func initRedis() {
	conf := config.GetConf().CacheRedis
	rdb = database.ConnectRedis(conf.Host, conf.Pass, conf.Node)
}

func GetRedis() *redis.Client {
	return rdb
}
