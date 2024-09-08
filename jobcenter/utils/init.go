package utils

func Init() {
	initMongo()
	initRocketMQ()
	initRedis()
}

func Close() {
	closeMongoClient()
	closeRocketMQProducer()
}
