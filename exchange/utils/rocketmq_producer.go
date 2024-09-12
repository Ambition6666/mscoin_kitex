package utils

import (
	"common/database"
	"exchange/config"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
)

var producer *database.RocketMQProducer

func initRocketMQProducer() {
	conf := config.GetConf().Rocketmq
	var err error
	producer, err = database.NewRocketMQProducer(&rmq.Config{
		Endpoint:    conf.Addr,
		Credentials: &credentials.SessionCredentials{},
	}, conf.WriteCap)
	if err != nil {
		panic(err)
	}
	producer.StartWrite()
}

func GetRocketMQProducer() *database.RocketMQProducer {
	return producer
}

func closeRocketMQProducer() {
	producer.Close()
}
