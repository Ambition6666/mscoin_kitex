package model

import (
	"common/enums"
	"github.com/jinzhu/copier"
)

type ExchangeOrder struct {
	Id            int64   `gorm:"column:id"`             // 订单主键ID
	OrderId       string  `gorm:"column:order_id"`       // 订单号
	Amount        float64 `gorm:"column:amount"`         // 订单数量
	BaseSymbol    string  `gorm:"column:base_symbol"`    // 基础货币符号
	CanceledTime  int64   `gorm:"column:canceled_time"`  // 取消时间
	CoinSymbol    string  `gorm:"column:coin_symbol"`    // 交易货币符号
	CompletedTime int64   `gorm:"column:completed_time"` // 完成时间
	Direction     int     `gorm:"column:direction"`      // 交易方向（1：买入，2：卖出）
	MemberId      int64   `gorm:"column:member_id"`      // 用户ID
	Price         float64 `gorm:"column:price"`          // 订单价格
	Status        int     `gorm:"column:status"`         // 订单状态
	Symbol        string  `gorm:"column:symbol"`         // 交易对符号
	Time          int64   `gorm:"column:time"`           // 下单时间
	TradedAmount  float64 `gorm:"column:traded_amount"`  // 已成交数量
	Turnover      float64 `gorm:"column:turnover"`       // 成交金额
	Type          int     `gorm:"column:type"`           // 订单类型（0：限价，1：市价）
	UseDiscount   string  `gorm:"column:use_discount"`   // 是否使用折扣（例如VIP折扣）
}

func NewOrder() *ExchangeOrder {
	return &ExchangeOrder{}
}

func (*ExchangeOrder) TableName() string {
	return "exchange_order"
}

// status
const (
	Trading = iota
	Completed
	Canceled
	OverTimed
	Init
)

var StatusMap = enums.Enum{
	Trading:   "TRADING",
	Completed: "COMPLETED",
	Canceled:  "CANCELED",
	OverTimed: "OVERTIMED",
}

// direction
const (
	BUY = iota
	SELL
)

var DirectionMap = enums.Enum{
	BUY:  "BUY",
	SELL: "SELL",
}

// type
const (
	MarketPrice = iota
	LimitPrice
)

var TypeMap = enums.Enum{
	MarketPrice: "MARKET_PRICE",
	LimitPrice:  "LIMIT_PRICE",
}

type ExchangeOrderVo struct {
	OrderId       string  `gorm:"column:order_id"`
	Amount        float64 `gorm:"column:amount"`
	BaseSymbol    string  `gorm:"column:base_symbol"`
	CanceledTime  int64   `gorm:"column:canceled_time"`
	CoinSymbol    string  `gorm:"column:coin_symbol"`
	CompletedTime int64   `gorm:"column:completed_time"`
	Direction     string  `gorm:"column:direction"`
	MemberId      int64   `gorm:"column:member_id"`
	Price         float64 `gorm:"column:price"`
	Status        string  `gorm:"column:status"`
	Symbol        string  `gorm:"column:symbol"`
	Time          int64   `gorm:"column:time"`
	TradedAmount  float64 `gorm:"column:traded_amount"`
	Turnover      float64 `gorm:"column:turnover"`
	Type          string  `gorm:"column:type"`
	UseDiscount   string  `gorm:"column:use_discount"`
}

func (old *ExchangeOrder) ToVo() *ExchangeOrderVo {
	eo := &ExchangeOrderVo{}
	copier.Copy(eo, old)
	eo.Status = StatusMap.Value(old.Status)
	eo.Direction = DirectionMap.Value(old.Direction)
	eo.Type = TypeMap.Value(old.Type)
	return eo
}
