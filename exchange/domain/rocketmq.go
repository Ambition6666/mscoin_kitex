package domain

import (
	"common/database"
	"context"
	"encoding/json"
	"exchange/config"
	"exchange/model"
	"exchange/utils"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RocketmqDomain struct {
	rocketmqClient   *database.RocketMQProducer
	rocketmqConsumer *database.RocketMQConsumer
	orderDomain      *ExchangeOrderDomain
}

func (d *RocketmqDomain) Send(topic string, userId int64, orderId string, money float64, symbol string, direction int, baseSymbol string, coinSymbol string, status int) bool {
	m := make(map[string]any)
	m["userId"] = userId
	m["orderId"] = orderId
	m["money"] = money
	m["symbol"] = symbol
	m["direction"] = direction
	m["baseSymbol"] = baseSymbol
	m["coinSymbol"] = coinSymbol
	m["status"] = status
	marshal, _ := json.Marshal(m)
	data := database.RocketMQData{
		Topic: topic,
		Key:   []string{orderId},
		Data:  marshal,
	}
	klog.Info(string(marshal))
	d.rocketmqClient.Send(data)

	return true
}

type AddOrderResult struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

func (d *RocketmqDomain) WaitAddOrderResult(topic string) {
	conf := config.GetConf().Rocketmq
	err := d.rocketmqConsumer.AddConsumer(&rmq.Config{
		Endpoint:      conf.Addr,
		ConsumerGroup: "exchange_api",
		Credentials:   &credentials.SessionCredentials{},
	}, conf.ReadCap, topic)
	if err != nil {
		klog.Error(err)
		return
	}
	d.rocketmqConsumer.StartRead(topic)
	for {
		data, _ := d.rocketmqConsumer.Read(topic)
		klog.Info("收到订单增加结果:" + string(data.Data))
		var result AddOrderResult
		json.Unmarshal(data.Data, &result)

		order, err := d.orderDomain.exchangeOrderDao.FindByOrderId(context.Background(), result.OrderId)
		if err != nil {
			klog.Error(err)
			err := d.orderDomain.UpdateOrderStatus(context.Background(), result.OrderId, model.Canceled)
			if err != nil {
				klog.Error("更新状态失败:", err)
				d.rocketmqConsumer.Rput(topic, data)
				return
			}
			continue
		}
		err = d.orderDomain.UpdateOrderStatus(context.Background(), result.OrderId, model.Trading)
		if err != nil {
			klog.Error("更新状态失败:", err)
			d.rocketmqConsumer.Rput(topic, data)
			continue
		}
		if order.Status != model.Init {
			klog.Error("订单已经被处理过", order.Status)
			continue
		}
		//订单初始化完成 发送消息到kafka 等待撮合交易引擎进行交易撮合
		for {
			bytes, _ := json.Marshal(order)
			orderData := database.RocketMQData{
				Topic: "exchange_order_trading",
				Key:   []string{order.OrderId},
				Data:  bytes,
			}
			d.rocketmqClient.Send(orderData)
			klog.Info("订单创建成功，发送创建成功消息:", order.OrderId)
			break
		}

	}
}

func NewRocketmqDomain(orderDomain *ExchangeOrderDomain) *RocketmqDomain {
	k := &RocketmqDomain{
		orderDomain:      orderDomain,
		rocketmqClient:   utils.GetRocketMQProducer(),
		rocketmqConsumer: utils.GetRocketMQConsumer(),
	}
	go k.WaitAddOrderResult("exchange_order_init_complete")
	return k
}
