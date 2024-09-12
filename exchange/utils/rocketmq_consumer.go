package utils

import (
	"common/database"
)

var consumer *database.RocketMQConsumer

func initRocketMQConsumer() {
	var err error
	consumer = database.NewRocketMQConsumer()
	if err != nil {
		panic(err)
	}
}

func GetRocketMQConsumer() *database.RocketMQConsumer {
	return consumer
}

func closeRocketMQConsumer() {
	consumer.Close()
}
