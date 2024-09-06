package domain

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"market/dao"
	"market/model"
)

type ExchangeCoinDomain struct {
	ExchangeCoinDao *dao.ExchangeCoinDao
}

func NewExchangeCoinDomain() *ExchangeCoinDomain {
	return &ExchangeCoinDomain{
		ExchangeCoinDao: dao.NewExchangeCoinDao(),
	}
}

func (d *ExchangeCoinDomain) FindVisible(ctx context.Context) []*model.ExchangeCoin {
	list, err := d.ExchangeCoinDao.FindVisible(ctx)
	if err != nil {
		klog.Error(err)
		return []*model.ExchangeCoin{}
	}
	return list
}

func (d *ExchangeCoinDomain) FindSymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	coin, err := d.ExchangeCoinDao.FindSymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return coin, nil
}
