package dao

import (
	"common/tools"
	"context"
	"errors"
	"gorm.io/gorm"
	"market/utils"
	"ucenter/model"
)

type MemberTransactionDao struct {
	conn *gorm.DB
}

func (d *MemberTransactionDao) Save(ctx context.Context, transaction *model.MemberTransaction) error {
	//TODO implement me
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
	err := session.Save(transaction).Error
	return err
}

func (d *MemberTransactionDao) FindTransaction(
	ctx context.Context,
	pageNo int,
	pageSize int,
	memberId int64,
	startTime string,
	endTime string,
	symbol string,
	transactionType string) (list []*model.MemberTransaction, total int64, err error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
	db := session.Model(&model.MemberTransaction{}).Where("member_id=?", memberId)
	if transactionType != "" {
		db.Where("type=?", tools.ToInt64(transactionType))
	}
	if startTime != "" && endTime != "" {
		sTime := tools.ToMill(startTime)
		eTime := tools.ToMill(endTime)
		db.Where("create_time >= ? and create_time <= ?", sTime, eTime)
	}
	if symbol != "" {
		db.Where("symbol=?", symbol)
	}
	offset := (pageNo - 1) * pageSize
	db.Count(&total)
	db.Order("create_time desc").Offset(offset).Limit(pageSize)
	err = db.Find(&list).Error
	return
}

func (d *MemberTransactionDao) FindByAmountAndTime(ctx context.Context, amount float64, address string, time int64) (mt *model.MemberTransaction, err error) {
	session := d.conn.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
	db := session.Model(&model.MemberTransaction{}).Where("amount=? and address=? and create_time=?", amount, address, time)
	err = db.Take(&mt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func NewMemberTransactionDao() *MemberTransactionDao {
	return &MemberTransactionDao{
		conn: utils.GetMysql(),
	}
}
