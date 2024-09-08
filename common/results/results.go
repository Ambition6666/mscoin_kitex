package results

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"
)

type Result struct {
	Code    BizCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data"`
}
type BizCode int

const (
	SuccessCode BizCode = 0
	BizFailCode BizCode = -999
)

func NewResult() *Result {
	return &Result{}
}

func (r *Result) Fail(code BizCode, msg string) {
	r.Code = code
	r.Message = msg
}

func (r *Result) Success(data any) {
	r.Code = SuccessCode
	r.Message = "success"
	r.Data = data
}

func (r *Result) Deal(data any, err error, c *app.RequestContext) {
	if err != nil {
		if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
			r.Fail(BizFailCode, bizErr.BizMessage())
			c.JSON(http.StatusOK, r)
		} else {
			c.JSON(http.StatusInternalServerError, "服务错误")
		}
		return
	}

	r.Success(data)
	c.JSON(http.StatusOK, r)
}
