package utils

import (
	"common/database"
	"github.com/redis/go-redis/v9"
	"jobcenter/config"
)

var rdb *redis.Client

func initRedis() {
	conf := config.GetConf().CacheRedis
	rdb = database.ConnectRedis(conf.Host, conf.Pass, conf.Node)
}

func GetRedis() *redis.Client {
	return rdb
}
