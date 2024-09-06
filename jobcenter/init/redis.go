package init

import (
	"common/database"
	"github.com/redis/go-redis/v9"
	"market/config"
)

var rdb *redis.Client

func initRedis() {
	conf := config.GetConf().CacheRedis
	rdb = database.ConnectRedis(conf.Host, conf.Pass, conf.Type)
}

func GetRedis() *redis.Client {
	return rdb
}
