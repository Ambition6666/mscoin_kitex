package domain

import (
	"common/tools"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"grpc_common/kitex_gen/market"
	"ucenter/dao"
	"ucenter/model"
	"ucenter/utils"
)

type MemberWalletDomain struct {
	memberWalletDao *dao.MemberWalletDao
	cache           *redis.Client
}

func (d *MemberWalletDomain) FindWalletBySymbol(ctx context.Context, id int64, name string, coin *market.Coin) (*model.MemberWalletCoin, error) {
	mw, err := d.memberWalletDao.FindByIdAndCoinName(ctx, id, name)
	if err != nil {
		return nil, err
	}
	if mw == nil {
		//新建并存储
		mw, walletCoin := model.NewMemberWallet(id, coin)
		err := d.memberWalletDao.Save(ctx, mw)
		if err != nil {
			return nil, err
		}
		return walletCoin, nil
	}
	nwc := &model.MemberWalletCoin{}
	copier.Copy(nwc, mw)
	nwc.Coin = coin
	return nwc, nil
}

func (d *MemberWalletDomain) Freeze(ctx context.Context, userId int64, money float64, symbol string) error {
	mw, err := d.memberWalletDao.FindByIdAndCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	if mw.Balance < money {
		return errors.New("余额不足")
	}
	err = d.memberWalletDao.UpdateFreeze(ctx, utils.GetMysql(), userId, money, symbol)
	return err
}

func (d *MemberWalletDomain) UpdateWalletCoinAndBase(ctx context.Context, baseWallet *model.MemberWallet, coinWallet *model.MemberWallet) error {
	return utils.GetMysql().Transaction(func(tx *gorm.DB) error {
		err := d.memberWalletDao.UpdateWallet(ctx, tx, baseWallet.Id, baseWallet.Balance, baseWallet.FrozenBalance)
		if err != nil {
			return err
		}
		err = d.memberWalletDao.UpdateWallet(ctx, tx, coinWallet.Id, coinWallet.Balance, coinWallet.FrozenBalance)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *MemberWalletDomain) FindWalletByMemIdAndCoinName(ctx context.Context, memberId int64, coinName string) (*model.MemberWallet, error) {
	mw, err := d.memberWalletDao.FindByIdAndCoinName(ctx, memberId, coinName)
	if err != nil {
		return nil, err
	}
	return mw, nil
}

func (d *MemberWalletDomain) FindWalletByMemIdAndCoinId(ctx context.Context, memberId int64, coinId int64) (*model.MemberWallet, error) {
	mw, err := d.memberWalletDao.FindByIdAndCoinId(ctx, memberId, coinId)
	if err != nil {
		return nil, err
	}
	return mw, nil
}

func (d *MemberWalletDomain) FindWalletByMemId(ctx context.Context, userId int64) ([]*model.MemberWallet, error) {
	memberWallets, err := d.memberWalletDao.FindByMemId(ctx, userId)
	return memberWallets, err
}

func (d *MemberWalletDomain) Copy(ctx context.Context, memberWallet *model.MemberWallet, coinInfo *market.Coin) *model.MemberWalletCoin {
	mwc := &model.MemberWalletCoin{}
	copier.Copy(mwc, memberWallet)
	mwc.Coin = &market.Coin{}
	copier.Copy(mwc.Coin, coinInfo)
	cnyRate := d.cache.Get(ctx, "USDT::CNY::RATE").String()
	if memberWallet.CoinName != "USDT" {
		//获取最新的汇率
		usdRate := d.cache.Get(ctx, memberWallet.CoinName+"::USDT::RATE").String()
		if usdRate == "" {
			usdRate = "1"
		}
		mwc.Coin.UsdRate = tools.ToFloat64(usdRate)
		mwc.Coin.CnyRate = tools.MulN(tools.ToFloat64(usdRate), tools.ToFloat64(cnyRate), 10)
	} else {
		mwc.Coin.UsdRate = 1
		mwc.Coin.CnyRate = tools.ToFloat64(cnyRate)
	}
	return mwc
}

func (d *MemberWalletDomain) UpdateAddress(ctx context.Context, mw *model.MemberWallet) error {
	return d.memberWalletDao.UpdateAddress(ctx, mw)
}

func (d *MemberWalletDomain) GetAddress(ctx context.Context, coin_name string) ([]string, error) {
	return d.memberWalletDao.GetAddress(ctx, coin_name)
}
func (d *MemberWalletDomain) FindByAddress(address string) (*model.MemberWallet, error) {
	return d.memberWalletDao.FindByAddress(context.Background(), address)
}

func NewMemberWalletDomain() *MemberWalletDomain {
	return &MemberWalletDomain{
		dao.NewMemberWalletDao(),
		utils.GetRedis(),
	}
}
