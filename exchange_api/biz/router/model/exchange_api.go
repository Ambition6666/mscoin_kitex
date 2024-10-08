// Code generated by hertz generator. DO NOT EDIT.

package model

import (
	model "exchange_api/biz/handler"
	"github.com/cloudwego/hertz/pkg/app/server"
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
		_exchange := root.Group("/exchange", _exchangeMw()...)
		{
			_asset := _exchange.Group("/asset", _assetMw()...)
			_asset.POST("/add", append(_addMw(), model.Add)...)
			_asset.POST("/current", append(_currentMw(), model.Current)...)
			_asset.POST("/history", append(_historyMw(), model.History)...)
		}
	}
}
