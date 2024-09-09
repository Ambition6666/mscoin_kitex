// Code generated by hertz generator.

package model

import (
	"common/results"
	"common/tools"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"ucenter_api/config"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _ucMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _register0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _mobileMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _sendcodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _assetMw() []app.HandlerFunc {
	// your code...
	mid := make([]app.HandlerFunc, 0)
	mid = append(mid, Auth(config.AccessSecret))
	return mid
}

func _walletMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _findwalletMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _resetwalletaddressMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _transactionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _findtransactionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _wallet0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _findwalletbysymbolMw() []app.HandlerFunc {
	// your code...
	return nil
}

func Auth(secret string) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		result := results.NewResult()
		result.Fail(4000, "no login")
		token := string(ctx.GetHeader("x-auth-token"))
		if token == "" {
			ctx.JSON(http.StatusOK, result)
			return
		}
		userId, err := tools.ParseToken(token, secret)
		if err != nil {
			ctx.JSON(200, result)
			return
		}
		ctx.Set("userId", userId)
		ctx.Next(c)
	}
}
