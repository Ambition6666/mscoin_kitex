package dao

import (
	"context"
	"gorm.io/gorm"
	"market/model"
	"market/utils"
)

type ExchangeCoinDao struct {
	conn *gorm.DB
}

func (d *ExchangeCoinDao) FindVisible(ctx context.Context) (list []*model.ExchangeCoin, err error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true, Context: ctx})
	err = session.Model(&model.ExchangeCoin{}).Where("visible=?", 1).Find(&list).Error
	return
}

func (d *ExchangeCoinDao) FindSymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true})
	coin := &model.ExchangeCoin{}
	err := session.Model(&model.ExchangeCoin{}).Where("symbol=?", symbol).Take(coin).Error
	return coin, err
}

func NewExchangeCoinDao() *ExchangeCoinDao {
	return &ExchangeCoinDao{
		conn: utils.GetMysql(),
	}
}
