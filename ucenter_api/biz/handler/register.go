// Code generated by hertz generator.

package handler

import (
	"common/results"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/jinzhu/copier"
	"grpc_common/kitex_gen/ucenter"
	model "ucenter_api/biz/model"
	"ucenter_api/rpc"
)

// Register .
// @router /uc/register/phone [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.Request
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	regReq := &ucenter.RegReq{}

	copier.Copy(regReq, req)

	_, err = rpc.GetRegisterClient().RegisterByPhone(ctx, regReq)

	resp := new(model.Response)

	results.NewResult().Deal(resp, err, c)
}

// SendCode .
// @router /uc/mobile/code [POST]
func SendCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req model.CodeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	codeReq := &ucenter.CodeReq{}
	copier.Copy(&req, codeReq)

	_, err = rpc.GetRegisterClient().SendCode(ctx, codeReq)

	resp := new(model.Response)

	results.NewResult().Deal(resp, err, c)
}
