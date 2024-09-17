package consumer

import (
	"encoding/json"
	"exchange/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
	"ucenter/domain"
)

type BitCoinTransactionResult struct {
	Value   float64 `json:"value"`
	Time    int64   `json:"time"`
	Address string  `json:"address"`
	Type    string  `json:"type"`
	Symbol  string  `json:"symbol"`
}

func BitCoinTransaction(topic string) {
	for {
		kafkaData, err := utils.GetRocketMQConsumer().Read(topic)
		if err != nil {
			klog.Error(err)
			continue
		}
		var bt BitCoinTransactionResult
		json.Unmarshal(kafkaData.Data, &bt)
		//解析出来数据 调用domain存储到数据库即可
		transactionDomain := domain.NewMemberTransactionDomain()
		err = transactionDomain.SaveRecharge(bt.Address, bt.Value, bt.Time, bt.Type, bt.Symbol)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			utils.GetRocketMQConsumer().Rput(topic, kafkaData)
		}
	}
}
