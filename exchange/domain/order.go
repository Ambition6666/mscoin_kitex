package domain

import (
	"common/tools"
	"context"
	"errors"
	"exchange/dao"
	"exchange/model"
	"exchange/utils"
	"fmt"
	"gorm.io/gorm"
	"grpc_common/kitex_gen/market"
	"grpc_common/kitex_gen/ucenter"
	"time"
)

type ExchangeOrderDomain struct {
	exchangeOrderDao *dao.ExchangeOrderDao
}

// 查询历史订单
func (d *ExchangeOrderDomain) FindHistory(
	ctx context.Context,
	symbol string,
	page int64,
	pageSize int64,
	userId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.exchangeOrderDao.FindBySymbolPage(
		ctx,
		symbol,
		userId,
		-1,
		int(page),
		int(pageSize),
		true)
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func NewExchangeOrderDomain() *ExchangeOrderDomain {
	return &ExchangeOrderDomain{
		exchangeOrderDao: dao.NewExchangeOrderDao(),
	}
}

// 查询正在交易订单
func (d *ExchangeOrderDomain) FindTrading(
	ctx context.Context,
	symbol string,
	page int64,
	pageSize int64,
	userId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.exchangeOrderDao.FindBySymbolPage(
		ctx,
		symbol,
		userId,
		model.Trading,
		int(page),
		int(pageSize),
		true)
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

// 添加交易
func (d *ExchangeOrderDomain) AddOrder(ctx context.Context, exchangeOrder *model.ExchangeOrder, coin *market.ExchangeCoin, wallet *ucenter.MemberWallet, coinWallet *ucenter.MemberWallet) (float64, error) {
	exchangeOrder.Status = model.Init
	exchangeOrder.TradedAmount = 0
	exchangeOrder.Time = time.Now().UnixMilli()
	exchangeOrder.OrderId = tools.Unique("E")
	var money float64
	if exchangeOrder.Direction == model.BUY {
		var turnover float64 = 0
		if exchangeOrder.Type == model.MarketPrice {
			turnover = exchangeOrder.Amount
		} else {
			turnover = tools.MulN(exchangeOrder.Amount, exchangeOrder.Price, 5)
		}
		//费率
		fee := tools.MulN(turnover, coin.Fee, 5)
		if wallet.Balance < turnover {
			return 0, errors.New("余额不足")
		}
		if wallet.Balance-turnover < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		//需要冻结的钱 turnover+fee
		money = tools.AddN(turnover, fee, 5)
	} else {
		fee := tools.MulN(exchangeOrder.Amount, coin.Fee, 5)
		if coinWallet.Balance < exchangeOrder.Amount {
			return 0, errors.New("余额不足")
		}
		if wallet.Balance-exchangeOrder.Amount < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		money = tools.AddN(exchangeOrder.Amount, fee, 5)
	}
	err := utils.GetMysql().Transaction(func(tx *gorm.DB) error {
		err := d.exchangeOrderDao.Save(ctx, tx, exchangeOrder)
		if err != nil {
			return err
		}
		err = d.exchangeOrderDao.Save(ctx, tx, exchangeOrder)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return money, nil
}

// FindByOrderId 通过订单id查找
func (d *ExchangeOrderDomain) FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error) {
	exchangeOrder, err := d.exchangeOrderDao.FindByOrderId(ctx, orderId)
	if err == nil && exchangeOrder == nil {
		return nil, errors.New("订单号不存在")
	}
	return exchangeOrder, err
}

// UpdateOrderStatusCancel 更新取消订单
func (d *ExchangeOrderDomain) UpdateOrderStatusCancel(ctx context.Context, orderId string, updateStatus int) error {
	return d.exchangeOrderDao.UpdateOrderStatusCancel(ctx, orderId, model.Canceled, updateStatus, time.Now().UnixMilli())
}

// FindCurrentTradingCount 查询当前交易订单数量
func (d *ExchangeOrderDomain) FindCurrentTradingCount(ctx context.Context, id int64, symbol string, direction string) (int64, error) {
	return d.exchangeOrderDao.FindCount(ctx, id, symbol, model.DirectionMap.Code(direction), model.Trading)
}

// UpdateOrderStatus 更新订单状态
func (d *ExchangeOrderDomain) UpdateOrderStatus(background context.Context, id string, status int) error {
	return d.exchangeOrderDao.UpdateOrderStatus(background, id, status)
}

// FindTradingOrders 查询所有交易中的订单
func (d *ExchangeOrderDomain) FindTradingOrders(ctx context.Context) ([]*model.ExchangeOrder, error) {
	return d.exchangeOrderDao.FindTradingOrders(ctx)
}

// UpdateOrderComplete 更新订单完成状态
func (d *ExchangeOrderDomain) UpdateOrderComplete(ctx context.Context, order *model.ExchangeOrder) error {
	return d.exchangeOrderDao.UpdateOrderComplete(ctx, order.OrderId, order.TradedAmount, order.Turnover, order.Status)
}
