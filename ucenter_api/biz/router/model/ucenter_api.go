// Code generated by hertz generator. DO NOT EDIT.

package model

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"ucenter_api/biz/handler"
	model "ucenter_api/biz/handler"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_uc := root.Group("/uc", _ucMw()...)
		_uc.POST("/login", append(_loginMw(), model.Login)...)
		_uc.POST("/check/login", append(_registerMw(), model.CheckLogin)...)
		{
			_asset := _uc.Group("/asset", _assetMw()...)
			_asset.POST("/wallet", append(_findwalletMw(), handler.FindWallet)...)
			_wallet := _asset.Group("/wallet", _walletMw()...)
			_wallet.POST("/reset-address", append(_resetwalletaddressMw(), handler.ResetWalletAddress)...)
			{
				_transaction := _asset.Group("/transaction", _transactionMw()...)
				_transaction.POST("/all", append(_findtransactionMw(), handler.FindTransaction)...)
			}
			{
				_wallet0 := _asset.Group("/wallet", _wallet0Mw()...)
				_wallet0.POST("/:coinName", append(_findwalletbysymbolMw(), handler.FindWalletBySymbol)...)
			}
		}
		{
			_mobile := _uc.Group("/mobile", _mobileMw()...)
			_mobile.POST("/code", append(_sendcodeMw(), model.SendCode)...)
		}
		{
			_register := _uc.Group("/register", _registerMw()...)
			_register.POST("/phone", append(_register0Mw(), model.Register)...)
		}
	}
}
