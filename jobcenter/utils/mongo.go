package utils

import (
	"common/database"
	"jobcenter/config"
)

var mdb *database.MongoClient

func initMongo() {
	conf := config.GetConf().Mongo
	mdb = database.ConnectMongo(conf.Username, conf.Password, conf.Database, conf.Url)
}

func GetMongoClient() *database.MongoClient {
	return mdb
}

func closeMongoClient() {
	mdb.Disconnect()
}
