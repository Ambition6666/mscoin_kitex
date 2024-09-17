package dao

import (
	"context"
	"exchange/utils"
	"gorm.io/gorm"
	"ucenter/model"
)

type WithdrawDao struct {
	conn *gorm.DB
}

func (m *WithdrawDao) UpdateTransactionNumber(ctx context.Context, wr model.WithdrawRecord) error {
	session := m.conn.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
	err := session.Model(&model.WithdrawRecord{}).
		Where("id=?", wr.Id).
		Updates(map[string]any{"transaction_number": wr.TransactionNumber, "status": wr.Status}).Error
	return err
}

func (m *WithdrawDao) Save(ctx context.Context, db *gorm.DB, record *model.WithdrawRecord) error {
	gormConn := db.WithContext(ctx)
	err := gormConn.Create(&record).Error
	return err
}
func (m *WithdrawDao) FindByUserId(ctx context.Context, id int64, page int32, size int32) (list []*model.WithdrawRecord, total int64, err error) {
	session := m.conn.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
	db := session.Model(&model.WithdrawRecord{}).Where("member_id=?", id)

	offset := (page - 1) * size
	db.Count(&total)
	db.Order("create_time desc").Offset(int(offset)).Limit(int(size))
	err = db.Find(&list).Error
	return
}
func NewWithdrawDao() *WithdrawDao {
	return &WithdrawDao{
		conn: utils.GetMysql(),
	}
}
