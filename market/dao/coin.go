package dao

import (
	"context"
	"gorm.io/gorm"
	"market/init"
	"market/model"
)

type CoinDao struct {
	conn *gorm.DB
}

func (d *CoinDao) FindByUnit(ctx context.Context, unit string) (*model.Coin, error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true, Context: ctx})
	coin := &model.Coin{}
	err := session.Model(&model.Coin{}).Where("unit=?", unit).Take(coin).Error
	return coin, err
}
func (d *CoinDao) FindAll(ctx context.Context) (list []*model.Coin, err error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true, Context: ctx})
	err = session.Model(&model.Coin{}).Find(&list).Error
	return
}
func (d *CoinDao) FindById(ctx context.Context, id int64) (*model.Coin, error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true, Context: ctx})
	coin := &model.Coin{}
	err := session.Model(&model.Coin{}).Where("id=?", id).Take(coin).Error
	return coin, err
}

func NewCoinDao() *CoinDao {
	return &CoinDao{
		conn: init.GetMysql(),
	}
}
