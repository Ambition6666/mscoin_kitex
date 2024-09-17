package ws

import (
	"github.com/cloudwego/kitex/pkg/klog"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"net/http"
)

const ROOM = "market"

type WebSocketHandler func(s socketio.Conn) error
type WebSocketErrorHandler func(s socketio.Conn, err error)
type WebSocketDisconnectHandler func(s socketio.Conn, reason string)

type WebSocketServer struct {
	wsServer *socketio.Server
}

func NewWebSocketServer() *WebSocketServer {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})

	ws := &WebSocketServer{
		wsServer: server,
	}

	// 处理连接事件
	ws.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join(ROOM)
		klog.Infof("Client connected: %s", s.ID())
		return nil
	})

	ws.OnError("/", func(s socketio.Conn, e error) {
		klog.Info("连接错误:", e)
	})

	ws.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.LeaveAll()
		klog.Info("关闭连接：", reason)
	})

	return ws
}

func (ws *WebSocketServer) Start() {
	klog.Info("Starting WebSocket server...")
	if err := ws.wsServer.Serve(); err != nil {
		klog.Fatalf("Error starting WebSocket server: %v", err)
	}
}

func (ws *WebSocketServer) Stop() {
	klog.Info("Stopping WebSocket server...")
	ws.wsServer.Close()
}

func (ws *WebSocketServer) OnConnect(path string, handler WebSocketHandler) {
	ws.wsServer.OnConnect(path, handler)
}

func (w *WebSocketServer) OnError(path string, handler WebSocketErrorHandler) {
	w.wsServer.OnError(path, handler)
}

func (w *WebSocketServer) OnDisconnect(path string, handler WebSocketDisconnectHandler) {
	w.wsServer.OnDisconnect(path, handler)
}

func (ws *WebSocketServer) BroadcastToNamespace(path string, event string, data any) {
	go func() {
		ws.wsServer.BroadcastToRoom(path, ROOM, event, data)
	}()
}

func (ws *WebSocketServer) Serve(r *http.Request, w http.ResponseWriter) {
	klog.Infof("Serving WebSocket request: %v", r.URL)
	klog.Info(r.Header)
	ws.wsServer.ServeHTTP(w, r)
}
