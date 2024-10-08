// Code generated by hertz generator.

package handler

import (
	"common/pages"
	"common/results"
	"context"
	"github.com/jinzhu/copier"
	"grpc_common/kitex_gen/ucenter"
	"ucenter_api/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	model "ucenter_api/biz/model"
)

// FindWalletBySymbol .
// @router /uc/asset/wallet/:coinName [POST]
func FindWalletBySymbol(ctx context.Context, c *app.RequestContext) {
	var err error
	var coinName string
	err = c.BindPath(&coinName)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	wallet, err := rpc.GetAssetClient().FindWalletBySymbol(ctx, &ucenter.AssetReq{
		CoinName: coinName,
		UserId:   c.GetInt64("userId"),
	})
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}

	resp := new(model.MemberWallet)

	copier.Copy(resp, wallet)

	results.NewResult().Deal(resp, nil, c)

}

// FindWallet .
// @router /uc/asset/wallet [POST]
func FindWallet(ctx context.Context, c *app.RequestContext) {
	var err error
	assetReq := new(ucenter.AssetReq)
	assetReq.UserId = c.GetInt64("userId")
	walletList, err := rpc.GetAssetClient().FindWallet(ctx, assetReq)
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}

	resp := new(model.MemberWalletList)
	copier.Copy(resp.List, walletList)
	results.NewResult().Deal(resp, nil, c)
}

// ResetWalletAddress .
// @router /uc/asset/wallet/reset-address [POST]
func ResetWalletAddress(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.AssetReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	assetReq := new(ucenter.AssetReq)
	assetReq.UserId = c.GetInt64("userId")
	assetReq.CoinName = req.Unit
	_, err = rpc.GetAssetClient().ResetWalletAddress(ctx, assetReq)
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}
	resp := new(model.Response)

	results.NewResult().Deal(resp, nil, c)
}

// FindTransaction .
// @router /uc/asset/transaction/all [POST]
func FindTransaction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.AssetReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	assetReq := &ucenter.AssetReq{
		UserId:    c.GetInt64("userId"),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Type:      req.Type,
		Symbol:    req.Symbol,
		PageNo:    int64(req.PageNo),
		PageSize:  int64(req.PageSize),
	}

	mts, err := rpc.GetAssetClient().FindTransaction(ctx, assetReq)
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}

	var resp []*model.MemberTransaction
	if err := copier.Copy(&resp, mts.List); err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}
	if resp == nil {
		resp = []*model.MemberTransaction{}
	}
	b := make([]any, len(resp))
	for i := range resp {
		b[i] = resp[i]
	}

	results.NewResult().Deal(pages.New(b, int64(req.PageNo), int64(req.PageSize), mts.Total), nil, c)
}
