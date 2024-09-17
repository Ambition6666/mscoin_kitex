package repo

import "jobcenter/model"

type BtcTransactionRepo interface {
	FindByTxId(txId string) (*model.BitCoinTransaction, error)
	Save(bt *model.BitCoinTransaction) error
}
