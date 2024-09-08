package utils

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
