package model

import "common/tools"

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

func NewKline(data []string, period string) *Kline {
	toInt64 := tools.ToInt64(data[0])
	return &Kline{
		Time:         toInt64,
		Period:       period,
		OpenPrice:    tools.ToFloat64(data[1]),
		HighestPrice: tools.ToFloat64(data[2]),
		LowestPrice:  tools.ToFloat64(data[3]),
		ClosePrice:   tools.ToFloat64(data[4]),
		Count:        tools.ToFloat64(data[5]),
		Volume:       tools.ToFloat64(data[6]),
		Turnover:     tools.ToFloat64(data[7]),
		TimeStr:      tools.ToTimeString(toInt64),
	}
}

type OkxKlineRes struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}
