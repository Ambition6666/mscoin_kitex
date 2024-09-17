package consumer

import (
	"common/database"
	"common/enums"
	"common/tools"
	"context"
	"encoding/json"
	"errors"
	"exchange_api/rpc"
	"github.com/cloudwego/kitex/pkg/klog"
	"grpc_common/kitex_gen/exchange"
	"time"
	"ucenter/domain"
	"ucenter/utils"
)

type OrderAdd struct {
	UserId     int64   `json:"userId"`
	OrderId    string  `json:"orderId"`
	Money      float64 `json:"money"`
	Symbol     string  `json:"symbol"`
	Direction  int     `json:"direction"`
	BaseSymbol string  `json:"baseSymbol"`
	CoinSymbol string  `json:"coinSymbol"`
	Status     int     `json:"status"`
}

var InitStatus = 4

func ExchangeOrderAddConsumer(topic string) {
	for {
		RocketMQData, _ := utils.GetRocketMQConsumer().Read(topic)
		//if RocketMQData == nil {
		//	continue
		//}
		var addData OrderAdd
		err := json.Unmarshal(RocketMQData.Data, &addData)

		if err != nil {
			//不是这个消息 消息类型错误
			klog.Error(err)
			continue
		}
		klog.Info("读取到订单添加消息：", string(RocketMQData.Data))
		var orderId string
		if len(RocketMQData.Key) == 0 {
			orderId = ""
		} else {
			orderId = RocketMQData.Key[0]
		}

		if addData.OrderId != orderId {
			klog.Error(errors.New("不合法的消息，订单号不匹配"))
			continue
		}
		if addData.Status != 4 {
			klog.Error(errors.New("不合法的消息，订单状态不是Init, 请注意此订单是否已经被非法修改:" + orderId))
			continue
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		//查询订单信息 如果是正在交易中 继续 否则return
		exchangeOrder, err := rpc.GetOrderClient().FindByOrderId(ctx, &exchange.OrderReq{
			OrderId: orderId,
		})
		if err != nil {
			cancelOrder(ctx, orderId, topic, InitStatus, RocketMQData)
			continue
		}
		if exchangeOrder.GetStatus() != int32(InitStatus) {
			klog.Error(errors.New("订单状态不是Init, 请注意此订单是否已经被非法修改:" + orderId))
			continue
		}
		//transaction := tran.NewTransaction(db)
		//transaction.Action(func(conn msdb.DbConn) error {
		//
		//})

		//查询订单信息 如果是正在交易中 继续 否则return

		if err != nil {
			klog.Error(err)
			klog.Info("已经有别的进程处理此消息")
			continue
		}

		walletDomain := domain.NewMemberWalletDomain()
		if addData.Direction == 0 {
			//buy baseSymbol
			err = walletDomain.Freeze(ctx, addData.UserId, addData.Money, addData.BaseSymbol)
			klog.Error(err)
		} else if addData.Direction == 1 {
			//sell coinSymbol
			err = walletDomain.Freeze(ctx, addData.UserId, addData.Money, addData.CoinSymbol)
			klog.Error(err)
		}

		if err != nil {
			//重新消费
			err := cancelOrder(ctx, orderId, topic, int(exchangeOrder.GetStatus()), RocketMQData)
			if err != nil {
				klog.Error("重新消费失败:", err)
			}
		}
		//都完成后 通知订单进行状态变更 需要保证一定发送成功
		for {
			m := make(map[string]any)
			m["userId"] = addData.UserId
			m["orderId"] = orderId
			marshal, _ := json.Marshal(m)
			data := database.RocketMQData{
				Topic: "exchange_order_init_complete",
				Key:   []string{orderId},
				Data:  marshal,
			}
			utils.GetRocketMQProducer().Send(data)
			break
		}

	}

}

func cancelOrder(ctx context.Context, orderId string, topic string, originStatus int, RocketMQData database.RocketMQData) error {
	_, err := rpc.GetOrderClient().CancelOrder(ctx, &exchange.OrderReq{
		OrderId:      orderId,
		UpdateStatus: int32(originStatus),
	})
	if err != nil {
		utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
		return err
	}
	return nil
}

type ExchangeOrder struct {
	Id            int64   `gorm:"column:id" json:"id"`
	OrderId       string  `gorm:"column:order_id" json:"orderId"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
	BaseSymbol    string  `gorm:"column:base_symbol" json:"baseSymbol"`
	CanceledTime  int64   `gorm:"column:canceled_time" json:"canceledTime"`
	CoinSymbol    string  `gorm:"column:coin_symbol" json:"coinSymbol"`
	CompletedTime int64   `gorm:"column:completed_time" json:"completedTime"`
	Direction     int     `gorm:"column:direction" json:"direction"`
	MemberId      int64   `gorm:"column:member_id" json:"memberId"`
	Price         float64 `gorm:"column:price" json:"price"`
	Status        int     `gorm:"column:status" json:"status"`
	Symbol        string  `gorm:"column:symbol" json:"symbol"`
	Time          int64   `gorm:"column:time" json:"time"`
	TradedAmount  float64 `gorm:"column:traded_amount" json:"tradedAmount"`
	Turnover      float64 `gorm:"column:turnover" json:"turnover"`
	Type          int     `gorm:"column:type" json:"type"`
	UseDiscount   string  `gorm:"column:use_discount" json:"useDiscount"`
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

func ExchangeOrderComplete(topic string) {
	//先接收消息
	for {
		RocketMQData, _ := utils.GetRocketMQConsumer().Read(topic)
		order := new(ExchangeOrder)
		err := json.Unmarshal(RocketMQData.Data, order)
		if err != nil {
			klog.Error("订单信息解析错误:", err)
			return
		}
		klog.Info("开始更新订单信息:", order)

		if order.Status != Completed {
			continue
		}
		klog.Info("收到exchange_order_complete_update_success 消息成功:" + order.OrderId)
		walletDomain := domain.NewMemberWalletDomain()
		if err != nil {
			klog.Error(err)
			klog.Info("有进程已经拿到锁进行处理了")
			continue
		}

		// BTC/USDT
		ctx := context.Background()
		if order.Direction == BUY {
			baseWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, order.MemberId, order.BaseSymbol)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			coinWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, order.MemberId, order.CoinSymbol)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			if order.Type == MarketPrice {
				//市价买 amount USDT 冻结的钱  asset.turnover扣的钱 还回去的钱 amount-asset.turnover
				baseWallet.FrozenBalance = tools.SubFloor(baseWallet.FrozenBalance, order.Amount, 8)
				baseWallet.Balance = tools.AddFloor(baseWallet.Balance, tools.SubFloor(order.Amount, order.Turnover, 8), 8)
				coinWallet.Balance = tools.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
			} else {
				//限价买 冻结的钱是 asset.price*amount  成交了turnover 还回去的钱 asset.price*amount-asset.turnover
				floor := tools.MulFloor(order.Price, order.Amount, 8)
				baseWallet.FrozenBalance = tools.SubFloor(baseWallet.FrozenBalance, floor, 8)
				baseWallet.Balance = tools.AddFloor(baseWallet.Balance, tools.SubFloor(floor, order.Turnover, 8), 8)
				coinWallet.Balance = tools.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
			}
			err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}
		} else {
			//卖 不管是市价还是限价 都是卖的 BTC  解冻amount 得到的钱是 asset.turnover
			coinWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, order.MemberId, order.CoinSymbol)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			baseWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, order.MemberId, order.BaseSymbol)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}

			coinWallet.FrozenBalance = tools.SubFloor(coinWallet.FrozenBalance, order.Amount, 8)
			baseWallet.Balance = tools.AddFloor(baseWallet.Balance, order.Turnover, 8)
			err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
			if err != nil {
				klog.Error(err)
				utils.GetRocketMQConsumer().Rput(topic, RocketMQData)
				time.Sleep(250 * time.Millisecond)
				continue
			}
		}
		klog.Info("更新钱包成功:" + order.OrderId)
	}

}
