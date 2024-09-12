package consumer

import (
	"common/database"
	"context"
	"encoding/json"
	"exchange/config"
	"exchange/domain"
	"exchange/model"
	"exchange/processor"
	"exchange/utils"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

type RocketmqConsumer struct {
	coinTradeFactory *processor.CoinTradeFactory
	orderDomain      *domain.ExchangeOrderDomain
	consumerManger   *database.RocketMQConsumer
	producer         *database.RocketMQProducer
}

func NewRocketmqConsumer(
	factory *processor.CoinTradeFactory,

) *RocketmqConsumer {
	return &RocketmqConsumer{
		coinTradeFactory: factory,
		orderDomain:      domain.NewExchangeOrderDomain(),
		consumerManger:   utils.GetRocketMQConsumer(),
		producer:         utils.GetRocketMQProducer(),
	}
}

func (k *RocketmqConsumer) Run() {
	k.orderTrading()
}

func (k *RocketmqConsumer) orderTrading() {
	//topic exchange_order_trading
	conf := config.GetConf().Rocketmq
	rmqconf := &rmq.Config{
		Endpoint:      conf.Addr,
		ConsumerGroup: "exchange_api",
		Credentials:   &credentials.SessionCredentials{},
	}
	err := k.consumerManger.AddConsumer(rmqconf, conf.ReadCap, "exchange_order_trading")
	if err != nil {
		klog.Error(err)
		return
	}
	err = k.consumerManger.StartRead("exchange_order_trading")
	if err != nil {
		klog.Error(err)
		return
	}
	go k.readOrderTrading("exchange_order_trading")

	err = k.consumerManger.AddConsumer(rmqconf, conf.ReadCap, "exchange_order_completed")
	if err != nil {
		klog.Error(err)
		return
	}
	err = k.consumerManger.StartRead("exchange_order_completed")
	if err != nil {
		klog.Error(err)
		return
	}
	go k.readOrderComplete("exchange_order_completed")
}

func (k *RocketmqConsumer) readOrderTrading(topic string) {
	for {
		rocketmqData, _ := k.consumerManger.Read(topic)
		order := new(model.ExchangeOrder)
		json.Unmarshal(rocketmqData.Data, &order)
		coinTrade := k.coinTradeFactory.GetCoinTrade(order.Symbol)
		coinTrade.Trade(order)
	}
}

func (k *RocketmqConsumer) readOrderComplete(topic string) {
	for {
		rocketmqData, _ := k.consumerManger.Read(topic)
		var order *model.ExchangeOrder
		json.Unmarshal(rocketmqData.Data, &order)
		klog.Info("读取已完成订单数据成功:", order.OrderId)
		//更新订单
		err := k.orderDomain.UpdateOrderComplete(context.Background(), order)
		if err != nil {
			klog.Error(err)
			//如果失败重新放入  再次进行更新
			k.consumerManger.Rput(topic, rocketmqData)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		//继续放入MQ中 通知钱包服务更新
		for {
			rocketmqData.Topic = "exchange_order_complete_update_success"
			k.producer.Send(rocketmqData)
		}
	}
}
