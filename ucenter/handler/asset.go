package handler

import (
	"common/bc"
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/jinzhu/copier"
	"grpc_common/kitex_gen/market"
	"grpc_common/kitex_gen/ucenter"
	"ucenter/domain"
	"ucenter/model"
	"ucenter/utils/rpc"
)

type AssetImpl struct {
	memberWalletDomain      *domain.MemberWalletDomain
	memberTransactionDomain *domain.MemberTransactionDomain
}

// 通过货币查询钱包
func (a *AssetImpl) FindWalletBySymbol(ctx context.Context, req *ucenter.AssetReq) (res *ucenter.MemberWallet, err error) {
	//TODO implement me
	info, err := rpc.GetMarketClient().FindCoinInfo(ctx, &market.MarketReq{
		Unit: req.CoinName,
	})
	if err != nil {
		klog.Error("FindWalletBySymbol", err)
		return nil, kerrors.NewBizStatusError(-1, "查询货币信息失败")
	}

	memberWalletCoin, err := a.memberWalletDomain.FindWalletBySymbol(ctx, req.UserId, req.CoinName, info)
	if err != nil {
		klog.Error("FindWalletBySymbol: ", err)
		return nil, kerrors.NewBizStatusError(-1, "通过货币查询失败")
	}
	copier.Copy(res, memberWalletCoin)
	return
}

// 查询钱包
func (a *AssetImpl) FindWallet(ctx context.Context, req *ucenter.AssetReq) (res *ucenter.MemberWalletList, err error) {
	//TODO implement me
	mws, err := a.memberWalletDomain.FindWalletByMemId(ctx, req.UserId)
	if err != nil {
		klog.Error("FindWallet: ", err)
		return nil, kerrors.NewBizStatusError(-1, "查询钱包信息失败")
	}
	list := make([]*model.MemberWalletCoin, 0)
	for _, v := range mws {
		coinInfo, err := rpc.GetMarketClient().FindCoinInfo(ctx, &market.MarketReq{
			Unit: v.CoinName,
		})
		if err != nil {
			return nil, err
		}
		list = append(list, a.memberWalletDomain.Copy(ctx, v, coinInfo))
	}

	copier.Copy(&res, list)
	return
}

// 重新设置钱包地址
func (a *AssetImpl) ResetWalletAddress(ctx context.Context, req *ucenter.AssetReq) (res *ucenter.AssetResp, err error) {
	//TODO implement me
	mw, err := a.memberWalletDomain.FindWalletByMemIdAndCoinName(ctx, req.UserId, req.CoinName)
	if err != nil {
		klog.Error("ResetWalletAddress: ", err)
		return nil, kerrors.NewBizStatusError(-1, "查询钱包信息失败")
	}
	if mw.Address == "" && req.CoinName == "BTC" {
		wallet, err := bc.NewWallet()
		if err != nil {
			klog.Error("ResetWalletAddress: ", err)
			return nil, kerrors.NewBizStatusError(-1, "获取btc钱包失败")
		}
		address := wallet.GetTestAddress()
		mw.Address = string(address)
		mw.AddressPrivateKey = wallet.GetPriKey()
		err = a.memberWalletDomain.UpdateAddress(ctx, mw)
		if err != nil {
			klog.Error("ResetWalletAddress: ", err)
			return nil, kerrors.NewBizStatusError(-1, "更新btc钱包信息失败")
		}
	}
	return
}

// 查询事务
func (a *AssetImpl) FindTransaction(ctx context.Context, req *ucenter.AssetReq) (res *ucenter.MemberTransactionList, err error) {
	//TODO implement me
	mms, _, err := a.memberTransactionDomain.FindTransaction(
		ctx,
		req.PageNo,
		req.PageSize,
		req.UserId,
		req.Symbol,
		req.StartTime,
		req.EndTime,
		req.Type)
	if err != nil {
		klog.Error("FindTransaction: ", err)
		return nil, kerrors.NewBizStatusError(-1, "获取失败")
	}
	copier.Copy(&res.List, mms)
	return
}

// 获取地址
func (a *AssetImpl) GetAddress(ctx context.Context, req *ucenter.AssetReq) (res *ucenter.AddressList, err error) {
	//TODO implement me
	ars, err := a.memberWalletDomain.GetAddress(ctx, req.CoinName)
	if err != nil {
		klog.Error("GetAddress: ", err)
		return nil, kerrors.NewBizStatusError(-1, "获取btc钱包地址失败")
	}

	res.List = ars

	return
}

func NewAssetImpl() *AssetImpl {
	return &AssetImpl{
		memberWalletDomain:      domain.NewMemberWalletDomain(),
		memberTransactionDomain: domain.NewMemberTransactionDomain(),
	}
}
