package domain

import (
	"context"
	"market/dao"
	"market/model"
)

type CoinDomain struct {
	CoinDao *dao.CoinDao
}

func (d *CoinDomain) FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error) {
	coin, err := d.CoinDao.FindByUnit(ctx, unit)
	coin.ColdWalletAddress = ""
	return coin, err
}

func (d *CoinDomain) FindCoinById(ctx context.Context, id int64) (*model.Coin, error) {
	coin, err := d.CoinDao.FindById(ctx, id)
	coin.ColdWalletAddress = ""
	return coin, err
}

func (d *CoinDomain) FindAllCoin(ctx context.Context) ([]*model.Coin, error) {
	return d.CoinDao.FindAll(ctx)
}

func NewConnDomain() *CoinDomain {
	return &CoinDomain{
		CoinDao: dao.NewCoinDao(),
	}
}
