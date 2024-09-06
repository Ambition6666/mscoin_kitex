package handler

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/jinzhu/copier"
	"grpc_common/kitex_gen/market"
	"market/domain"
)

// MarketImpl ExchangeRateImpl implements the last service interface defined in the IDL.
type MarketImpl struct {
	exchangeCoinDomain *domain.ExchangeCoinDomain
	marketDomain       *domain.MarketDomain
	coinDomain         *domain.CoinDomain
}

// FindSymbolThumb 查找可见交易币种的缩略图列表
func (m MarketImpl) FindSymbolThumb(ctx context.Context, req *market.MarketReq) (res *market.SymbolThumbRes, err error) {
	//TODO implement me
	exchangeCoins := m.exchangeCoinDomain.FindVisible(ctx)
	coinThumbs := make([]*market.CoinThumb, len(exchangeCoins))
	for i, v := range exchangeCoins {
		ct := &market.CoinThumb{}
		ct.Symbol = v.Symbol
		coinThumbs[i] = ct
	}
	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil
}

// FindSymbolThumbTrend 查找可见交易币种的趋势缩略图列表
func (m MarketImpl) FindSymbolThumbTrend(ctx context.Context, req *market.MarketReq) (res *market.SymbolThumbRes, err error) {
	//TODO implement me
	exchangeCoins := m.exchangeCoinDomain.FindVisible(ctx)
	coinThumbs := m.marketDomain.SymbolThumbTrend(ctx, exchangeCoins)
	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil
}

// FindSymbolInfo 查找指定符号的交易币种信息
func (m MarketImpl) FindSymbolInfo(ctx context.Context, req *market.MarketReq) (res *market.ExchangeCoin, err error) {
	//TODO implement me
	exchangeCoin, err := m.exchangeCoinDomain.FindSymbol(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	mc := &market.ExchangeCoin{}
	if err := copier.Copy(mc, exchangeCoin); err != nil {
		return nil, err
	}
	return mc, nil
}

// FindCoinInfo 查找指定单位的币种信息
func (m MarketImpl) FindCoinInfo(ctx context.Context, req *market.MarketReq) (res *market.Coin, err error) {
	//TODO implement me
	coin, err := m.coinDomain.FindCoinInfo(ctx, req.Unit)
	if err != nil {
		return nil, err
	}
	mc := &market.Coin{}
	if err := copier.Copy(mc, coin); err != nil {
		return nil, err
	}
	return mc, nil
}

// HistoryKline 获取历史 K 线数据，根据请求的分辨率选择相应的时间周期
func (m MarketImpl) HistoryKline(ctx context.Context, req *market.MarketReq) (res *market.HistoryRes, err error) {
	//TODO implement me
	period := "1H"
	if req.Resolution == "60" {
		period = "1H"
	} else if req.Resolution == "30" {
		period = "30m"
	} else if req.Resolution == "15" {
		period = "15m"
	} else if req.Resolution == "5" {
		period = "5m"
	} else if req.Resolution == "1" {
		period = "1m"
	} else if req.Resolution == "1D" {
		period = "1D"
	} else if req.Resolution == "1W" {
		period = "1W"
	} else if req.Resolution == "1M" {
		period = "1M"
	}
	histories, err := m.marketDomain.HistoryKline(ctx, req.Symbol, req.From, req.To, period)
	if err != nil {
		return nil, err
	}
	return &market.HistoryRes{
		List: histories,
	}, nil
}

// FindVisibleExchangeCoins 查找所有可见的交易币种
func (m MarketImpl) FindVisibleExchangeCoins(ctx context.Context, req *market.MarketReq) (res *market.ExchangeCoinRes, err error) {
	//TODO implement me

	list := m.exchangeCoinDomain.FindVisible(ctx)
	err = copier.Copy(&res.List, list)

	if err != nil {
		klog.Error("复制数据失败:", err)
		return nil, kerrors.NewBizStatusError(-1, "exchangecoin复制数据失败")
	}

	return
}

// FindAllCoin 查找所有币种列表
func (m MarketImpl) FindAllCoin(ctx context.Context, req *market.MarketReq) (res *market.CoinList, err error) {
	//TODO implement me
	coinList, err := m.coinDomain.FindAllCoin(ctx)
	if err != nil {
		return nil, err
	}
	var list []*market.Coin
	copier.Copy(&list, coinList)
	return &market.CoinList{
		List: list,
	}, nil
}

// FindCoinById 根据 ID 查找币种信息
func (m MarketImpl) FindCoinById(ctx context.Context, req *market.MarketReq) (res *market.Coin, err error) {
	//TODO implement me
	coin, err := m.coinDomain.FindCoinById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	mc := &market.Coin{}
	if err := copier.Copy(mc, coin); err != nil {
		return nil, err
	}
	return mc, nil
}

func NewMarketImpl() *MarketImpl {
	return &MarketImpl{
		exchangeCoinDomain: domain.NewExchangeCoinDomain(),
		marketDomain:       domain.NewMarketDomain(),
		coinDomain:         domain.NewConnDomain(),
	}
}
