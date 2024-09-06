package init

func Init() {
	initMongo()
	initMysql()
	initRocketMQ()
	initRedis()
}

func Close() {
	closeMongoClient()
	closeRocketMQProducer()
}
