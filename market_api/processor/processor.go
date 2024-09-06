package processor

import (
	"common/database"
	"encoding/json"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/cloudwego/kitex/pkg/klog"
	"grpc_common/kitex_gen/market"
	"market_api/biz/model"
	"market_api/config"
)

const KLINE1M = "kline_1m"
const KLINE = "kline"
const TRADE = "trade"
const TradePlateTopic = "exchange_order_trade_plate"
const TradePlate = "tradePlate"

type ProcessData struct {
	Type string //trade 交易 kline k线
	Key  []byte
	Data []byte
}
type MarketHandler interface {
	HandleTrade(symbol string, data []byte)
	HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb)
	HandlerTradePlate(symbol string, plate *model.TradePlateResult)
}
type Processor interface {
	PutThumb(any)
	GetThumb() any
	Process(data *ProcessData)
	AddHandler(h MarketHandler)
}

type DefaultProcessor struct {
	rocketmqCli *database.RocketMQConsumer
	handlers    []MarketHandler
	thumbMap    map[string]*market.CoinThumb
}

var defaultProcessor *DefaultProcessor

func (p *DefaultProcessor) PutThumb(data any) {
	switch data.(type) {
	case []*market.CoinThumb:
		list := data.([]*market.CoinThumb)
		for _, v := range list {
			p.thumbMap[v.Symbol] = v
		}
	}
}

func (p *DefaultProcessor) GetThumb() any {
	return p.thumbMap
}

func NewDefaultProcessor(rcli *database.RocketMQConsumer) *DefaultProcessor {
	return &DefaultProcessor{
		rocketmqCli: rcli,
		handlers:    make([]MarketHandler, 0),
		thumbMap:    make(map[string]*market.CoinThumb),
	}
}

func (p *DefaultProcessor) Init(wh *WebsocketHandler) {
	p.AddHandler(wh)
	//接收kline 1m的同步数据
	go p.startReadFromRocketMQ(KLINE1M, KLINE)
	p.startReadTradePlate(TradePlateTopic)
}

func (p *DefaultProcessor) Process(data *ProcessData) {

	if data.Type == KLINE {
		kline := &model.Kline{}
		json.Unmarshal(data.Data, kline)
		for _, v := range p.handlers {
			v.HandleKLine(string(data.Key), kline, p.thumbMap)
		}
	} else if data.Type == TradePlate {
		tp := &model.TradePlateResult{}
		json.Unmarshal(data.Data, tp)
		for _, v := range p.handlers {
			v.HandlerTradePlate(string(data.Key), tp)
		}
	}
}
func (p *DefaultProcessor) AddHandler(h MarketHandler) {
	p.handlers = append(p.handlers, h)
}

func (p *DefaultProcessor) startReadFromRocketMQ(topic string, tp string) {
	conf := config.GetConf().Rocketmq
	err := p.rocketmqCli.AddConsumer(&rmq.Config{
		Endpoint:    conf.Addr,
		Credentials: &credentials.SessionCredentials{},
	}, conf.ReadCap, topic)
	if err != nil {
		klog.Error(err)
		return
	}

	err = p.rocketmqCli.StartRead(topic)
	if err != nil {
		klog.Error(err)
		return
	}
	go p.dealQueueData(p.rocketmqCli, tp, topic)
}

func (p *DefaultProcessor) dealQueueData(rcli *database.RocketMQConsumer, tp string, topic string) {
	for {
		rocketmqData, _ := rcli.Read(topic)
		pd := &ProcessData{
			Type: tp,
			Key:  []byte(rocketmqData.Key[0]),
			Data: rocketmqData.Data,
		}
		p.Process(pd)
	}
}

func (p *DefaultProcessor) startReadTradePlate(topic string) {
	p.rocketmqCli.StartRead(topic)
	go p.dealQueueData(p.rocketmqCli, TradePlate, topic)
}

func GetDefaultProcessor() *DefaultProcessor {
	if defaultProcessor == nil {
		defaultProcessor = NewDefaultProcessor(database.NewRocketMQConsumer())
	}
	return defaultProcessor
}
