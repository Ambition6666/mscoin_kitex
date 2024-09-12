// Code generated by hertz generator.

package handler

import (
	"common/pages"
	"common/results"
	"context"
	"exchange_api/rpc"
	"grpc_common/kitex_gen/exchange"

	model "exchange_api/biz/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// History .
// @router /exchange/asset/history [POST]
func History(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.ExchangeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	history, err := rpc.GetOrderClient().FindOrderHistory(ctx, &exchange.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   c.GetInt64("userId"),
	})
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}
	list := history.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}

	results.NewResult().Deal(pages.New(b, req.PageNo, req.PageSize, history.Total), nil, c)
}

// Current .
// @router /exchange/asset/current [POST]
func Current(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.ExchangeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	current, err := rpc.GetOrderClient().FindOrderCurrent(ctx, &exchange.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   c.GetInt64("userId"),
	})
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}
	list := current.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}

	results.NewResult().Deal(pages.New(b, req.PageNo, req.PageSize, current.Total), nil, c)
}

// Add .
// @router /exchange/asset/add [POST]
func Add(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.ExchangeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	add, err := rpc.GetOrderClient().Add(ctx, &exchange.OrderReq{
		Symbol:    req.Symbol,
		UserId:    c.GetInt64("userId"),
		Direction: req.Direction,
		Type:      req.Type,
		Price:     req.Price,
		Amount:    req.Amount,
	})
	if err != nil {
		results.NewResult().Deal(nil, err, c)
		return
	}
	results.NewResult().Deal(add.OrderId, nil, c)
}
