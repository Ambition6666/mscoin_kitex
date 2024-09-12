package handler

import (
	"context"
	"exchange/domain"
	"exchange/model"
	"exchange/rpc"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/jinzhu/copier"
	exchange "grpc_common/kitex_gen/exchange"
	"grpc_common/kitex_gen/market"
	"grpc_common/kitex_gen/ucenter"
)

// OrderImpl implements the last service interface defined in the IDL.
type OrderImpl struct {
	exchangeOrderDomain *domain.ExchangeOrderDomain
}

// FindOrderHistory implements the OrderImpl interface.
func (s *OrderImpl) FindOrderHistory(ctx context.Context, req *exchange.OrderReq) (resp *exchange.OrderRes, err error) {
	exchangeOrders, total, err := s.exchangeOrderDomain.FindHistory(ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		klog.Error("FindOrderHistory: ", err)
		return nil, kerrors.NewBizStatusError(-1, "查询历史订单失败")
	}
	var list []*exchange.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	resp = &exchange.OrderRes{
		List:  list,
		Total: total,
	}
	return
}

// FindOrderCurrent implements the OrderImpl interface.
func (s *OrderImpl) FindOrderCurrent(ctx context.Context, req *exchange.OrderReq) (resp *exchange.OrderRes, err error) {
	// TODO: Your code here...
	exchangeOrders, total, err := s.exchangeOrderDomain.FindTrading(ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		klog.Error("FindOrderCurrent: ", err)
		return nil, kerrors.NewBizStatusError(-1, "查询正在交易的订单失败")
	}
	var list []*exchange.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	resp = &exchange.OrderRes{
		List:  list,
		Total: total,
	}
	return
}

// Add implements the OrderImpl interface.
func (s *OrderImpl) Add(ctx context.Context, req *exchange.OrderReq) (resp *exchange.AddOrderRes, err error) {
	// TODO: Your code here...
	memberRes, err := rpc.GetMemberClient().FindMemberById(ctx, &ucenter.MemberReq{
		MemberId: req.UserId,
	})
	if err != nil {
		klog.Error("Add: ", err)
		return nil, kerrors.NewBizStatusError(-1, "添加失败")
	}
	if memberRes.TransactionStatus == 0 {
		return nil, kerrors.NewBizStatusError(-1, "此用户已经被禁止交易")
	}
	if req.Type == model.TypeMap[model.LimitPrice] && req.Price <= 0 {
		return nil, kerrors.NewBizStatusError(-1, "限价模式下价格不能小于等于0")
	}
	if req.Amount <= 0 {
		return nil, kerrors.NewBizStatusError(-1, "数量不能小于等于0")
	}
	exchangeCoin, err := rpc.GetMarketClient().FindSymbolInfo(ctx, &market.MarketReq{
		Symbol: req.Symbol,
	})
	if err != nil {
		klog.Error("Add: ", err)
		return nil, kerrors.NewBizStatusError(-1, "添加失败")
	}
	if exchangeCoin.Exchangeable != 1 && exchangeCoin.Enable != 1 {
		return nil, kerrors.NewBizStatusError(-1, "禁止交易的货币")
	}
	//基准币
	baseSymbol := exchangeCoin.GetBaseSymbol()
	//交易币
	coinSymbol := exchangeCoin.GetCoinSymbol()
	cc := baseSymbol
	if req.Direction == model.DirectionMap[model.SELL] {
		//根据交易币查询
		cc = coinSymbol
	}
	coin, err := rpc.GetMarketClient().FindCoinInfo(ctx, &market.MarketReq{
		Unit: cc,
	})
	if err != nil || coin == nil {
		return nil, kerrors.NewBizStatusError(-1, "不支持的货币")
	}
	if req.Type == model.TypeMap[model.MarketPrice] && req.Direction == model.DirectionMap[model.BUY] {
		if exchangeCoin.GetMinTurnover() > 0 && req.Amount < float64(exchangeCoin.GetMinTurnover()) {
			return nil, kerrors.NewBizStatusError(-1, "成交额至少是"+fmt.Sprintf("%d", exchangeCoin.GetMinTurnover()))
		}
	} else {
		if exchangeCoin.GetMaxVolume() > 0 && exchangeCoin.GetMaxVolume() < req.Amount {
			return nil, kerrors.NewBizStatusError(-1, "数量超出"+fmt.Sprintf("%f", exchangeCoin.GetMaxVolume()))
		}
		if exchangeCoin.GetMinVolume() > 0 && exchangeCoin.GetMinVolume() > req.Amount {
			return nil, kerrors.NewBizStatusError(-1, "数量不能低于"+fmt.Sprintf("%f", exchangeCoin.GetMinVolume()))
		}
	}
	//查询用户钱包
	baseWallet, err := rpc.GetAssetClient().FindWalletBySymbol(ctx, &ucenter.AssetReq{
		UserId:   req.UserId,
		CoinName: baseSymbol,
	})
	if err != nil {
		klog.Error("Add: ", err)
		return nil, kerrors.NewBizStatusError(-1, "没有钱包")
	}
	exCoinWallet, err := rpc.GetAssetClient().FindWalletBySymbol(ctx, &ucenter.AssetReq{
		UserId:   req.UserId,
		CoinName: coinSymbol,
	})
	if err != nil {
		klog.Error("Add: ", err)
		return nil, kerrors.NewBizStatusError(-1, "没有钱包")
	}
	if baseWallet.IsLock == 1 || exCoinWallet.IsLock == 1 {
		return nil, kerrors.NewBizStatusError(-1, "钱包已冻结")
	}
	if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.GetMinSellPrice() > 0 {
		if req.Price < exchangeCoin.GetMinSellPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, kerrors.NewBizStatusError(-1, "不能低于最低限价:"+fmt.Sprintf("%f", exchangeCoin.GetMinSellPrice()))
		}
	}
	if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.GetMaxBuyPrice() > 0 {
		if req.Price > exchangeCoin.GetMaxBuyPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, kerrors.NewBizStatusError(-1, "不能低于最高限价:"+fmt.Sprintf("%f", exchangeCoin.GetMaxBuyPrice()))
		}
	}
	//是否启用了市价买卖
	if req.Type == model.TypeMap[model.MarketPrice] {
		if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.EnableMarketBuy == 0 {
			return nil, kerrors.NewBizStatusError(-1, "不支持市价购买")
		} else if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.EnableMarketSell == 0 {
			return nil, kerrors.NewBizStatusError(-1, "不支持市价出售")
		}
	}
	resp.OrderId = req.OrderId
	return
}

// FindByOrderId implements the OrderImpl interface.
func (s *OrderImpl) FindByOrderId(ctx context.Context, req *exchange.OrderReq) (resp *exchange.ExchangeOrderOrigin, err error) {
	// TODO: Your code here...
	exchangeOrder, err := s.exchangeOrderDomain.FindByOrderId(ctx, req.OrderId)
	if err != nil {
		klog.Error("FindByOrderId: ", err)
		return nil, kerrors.NewBizStatusError(-1, "查询失败")
	}
	err = copier.Copy(resp, exchangeOrder)
	if err != nil {
		return nil, kerrors.NewBizStatusError(-1, "查询失败")
	}
	return
}

// CancelOrder implements the OrderImpl interface.
func (s *OrderImpl) CancelOrder(ctx context.Context, req *exchange.OrderReq) (resp *exchange.CancelOrderRes, err error) {
	// TODO: Your code here...
	err = s.exchangeOrderDomain.UpdateOrderStatusCancel(ctx, req.OrderId, int(req.UpdateStatus))
	if err != nil {
		klog.Error("CancelOrder: ", err)
		return nil, kerrors.NewBizStatusError(-1, "取消失败")
	}
	resp.OrderId = req.OrderId
	return
}

func NewOrderImpl() *OrderImpl {
	return &OrderImpl{
		exchangeOrderDomain: domain.NewExchangeOrderDomain(),
	}
}
