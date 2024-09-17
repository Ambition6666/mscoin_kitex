package consumer

import (
	"encoding/json"
	"exchange/utils"
	"log"
	"ucenter/domain"
	"ucenter/model"
)

func WithdrawConsumer(topic string) {
	withdrawDomain := domain.NewWithdrawDomain()
	for {
		RocketMQData, err := utils.GetRocketMQConsumer().Read(topic)
		if err != nil {
			log.Println(err)
			continue
		}
		var wr model.WithdrawRecord
		json.Unmarshal(RocketMQData.Data, &wr)
		//调用btc rpc进行转账
		err = withdrawDomain.Withdraw(wr)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("提现成功")
		}
	}
}
