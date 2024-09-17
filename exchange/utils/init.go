package utils

func Init() {
	initRocketMQProducer()
	initRocketMQConsumer()
	initRedis()
	initMysql()
}

func Close() {
	closeRocketMQProducer()
	closeRocketMQConsumer()
}
