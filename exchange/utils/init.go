package utils

func Init() {
	initRedis()
	initMysql()
	initRocketMQProducer()
}

func Close() {
	closeRocketMQProducer()
	closeRocketMQConsumer()
}
