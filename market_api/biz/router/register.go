// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	model "market_api/biz/router/model"
	"market_api/ws"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz, ws *ws.WebSocketServer) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	model.Register(r)
	model.RegisterWS(r, ws)
}
