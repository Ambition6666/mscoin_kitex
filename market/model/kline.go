package model

import (
	"common/tools"
	"grpc_common/kitex_gen/market"
)

type Kline struct {
	Period       string  `bson:"period,omitempty" json:"period"`             // 时间周期
	OpenPrice    float64 `bson:"openPrice,omitempty" json:"openPrice"`       // 开盘价
	HighestPrice float64 `bson:"highestPrice,omitempty" json:"highestPrice"` // 最高价
	LowestPrice  float64 `bson:"lowestPrice,omitempty" json:"lowestPrice"`   // 最低价
	ClosePrice   float64 `bson:"closePrice,omitempty" json:"closePrice"`     // 收盘价
	Time         int64   `bson:"time,omitempty" json:"time"`                 // 时间戳
	Count        float64 `bson:"count,omitempty" json:"count"`               // 成交笔数
	Volume       float64 `bson:"volume,omitempty" json:"volume"`             // 成交量
	Turnover     float64 `bson:"turnover,omitempty" json:"turnover"`         // 成交额
	TimeStr      string  `bson:"timeStr,omitempty" json:"timeStr"`
}

func (*Kline) Table(symbol, period string) string {
	return "exchange_kline_" + symbol + "_" + period
}

func (k *Kline) ToCoinThumb(symbol string, end *Kline) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.Zone = 0
	ct.Change = k.ClosePrice - end.ClosePrice
	ct.Chg = tools.MulN(tools.DivN(ct.Change, ct.Close, 5), 100, 5)
	ct.UsdRate = k.ClosePrice
	ct.BaseUsdRate = 1
	return ct
}
