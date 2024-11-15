package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"webscoket-test/service"
)

// Conn类型表示WebSocket连接。服务器应用程序从HTTP请求处理程序调用Upgrader.Upgrade方法以获取* Conn：
// var upgrader = websocket.Upgrader{}
var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebSocketController struct {
	beego.Controller
	conn *websocket.Conn
}

func (w *WebSocketController) OnMessage(text []byte) {
	logger.Printf("OnMessage: %s %d", string(text), w.conn)
	service.GetWSMessageService().ProcessMessage(w.conn, text)

}

func (w *WebSocketController) OnClose(code int, text string) error {
	logger.Printf("OnClose: %s %d", text, code)
	service.GetSessionService().RemoveSessionWithConn(w.conn)
	err := w.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// @Summary websocket init
// @router /test
func (w *WebSocketController) WebSocketCtrl() {
	log.Printf("WebSocketCtrl... %x", &w)
	// Upgrade from http request to WebSocket.
	conn, err := upgrader.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer conn.Close()

	w.conn = conn
	conn.SetCloseHandler(w.OnClose)
	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	//p是一个[]字节，messageType是一个值为websocket.BinaryMessage或websocket.TextMessage的int。
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			w.OnMessage(msg)
		}
	}
}
