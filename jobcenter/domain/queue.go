package domain

import (
	"common/database"
	"encoding/json"
	"jobcenter/model"
)

const KLINE1M = "kline_1m"

type QueueDomain struct {
	cli *database.RocketMQProducer
}

func (d *QueueDomain) Sync1mKline(data []string, symbol string, period string) {
	kline := model.NewKline(data, period)
	bytes, _ := json.Marshal(kline)
	sendData := database.RocketMQData{
		Topic: KLINE1M,
		Key:   []string{symbol},
		Data:  bytes,
	}
	d.cli.Send(sendData)
}

func (d *QueueDomain) SendRecharge(value float64, address string, time int64) {
	data := make(map[string]any)
	data["value"] = value
	data["address"] = address
	data["time"] = time
	data["type"] = model.RECHARGE
	data["symbol"] = "BTC"
	marshal, _ := json.Marshal(data)
	msg := database.RocketMQData{
		Topic: "BtcTransactionTopic",
		Data:  marshal,
		Key:   []string{address},
	}
	d.cli.Send(msg)
}

func NewQueueDomain(cli *database.RocketMQProducer) *QueueDomain {
	return &QueueDomain{
		cli: cli,
	}
}
