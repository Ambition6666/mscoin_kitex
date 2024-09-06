package init

import (
	"common/database"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"jobcenter/config"
)

var producer *database.RocketMQProducer

func initRocketMQ() {
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
