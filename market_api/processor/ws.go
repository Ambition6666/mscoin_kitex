package processor

import (
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	socketio "github.com/googollee/go-socket.io"
	"grpc_common/kitex_gen/market"
	"market_api/biz/model"
	"market_api/ws"
)

type WebsocketHandler struct {
	server *ws.WebSocketServer
}

func (w *WebsocketHandler) HandleTrade(symbol string, data []byte) {
	//TODO implement me
	panic("implement me")
}

func (w *WebsocketHandler) HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) {
	klog.Info("======接收到数据,symbol=", symbol)
	coinThumb := thumbMap[symbol]
	var thumb *market.CoinThumb
	if coinThumb == nil {
		thumb = kline.InitCoinThumb(symbol)
	} else {
		thumb = kline.ToCoinThumb(symbol, coinThumb)
	}
	marshal, _ := json.Marshal(thumb)
	w.server.BroadcastToNamespace("/", "/topic/market/thumb", string(marshal))

	bytes, _ := json.Marshal(kline)
	w.server.BroadcastToNamespace("/", "/topic/market/kline/"+symbol, string(bytes))
}

func (w *WebsocketHandler) HandlerTradePlate(symbol string, plate *model.TradePlateResult) {
	marshal, _ := json.Marshal(plate)
	klog.Info("====买卖盘通知:", symbol, plate.Direction, ":", fmt.Sprintf("%d", len(plate.Items)))
	w.server.BroadcastToNamespace("/", "/topic/market/trade-plate/"+symbol, string(marshal))
}

func (w *WebsocketHandler) OnConnect(s socketio.Conn) error {
	s.SetContext("")
	fmt.Println("connected:", s.ID())
	return nil
}

func NewWebsocketHandler(server *ws.WebSocketServer) *WebsocketHandler {
	return &WebsocketHandler{
		server: server,
	}
}
