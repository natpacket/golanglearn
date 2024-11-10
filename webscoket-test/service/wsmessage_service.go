package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
	logger "github.com/sirupsen/logrus"
	"sync"
	"webscoket-test/model"
)

var messageSvcOnce sync.Once
var messageSvr *WSMessageService

type WSMessage struct {
	MsgType string
	Data    interface{}
}

type WSMessageService struct {
	ResultMap cmap.ConcurrentMap[string, any]
}

func GetWSMessageService() *WSMessageService {
	messageSvcOnce.Do(func() {
		messageSvr = &WSMessageService{
			ResultMap: cmap.New[any](),
		}
	})
	return messageSvr
}

func (w *WSMessageService) ProcessMessage(conn *websocket.Conn, data []byte) {
	var wsMessage model.WSMessage
	err := json.Unmarshal(data, &wsMessage)
	if err != nil {
		logger.Debugf("ProcessMessage json.Unmarshal err: %v", err)
	}
	msgType := wsMessage.MsgType
	switch msgType {
	case "register":
		GetSessionService().registerSession(conn, wsMessage.Data)
		break
	case "heart":
		GetSessionService().addSessionToActiveQueue(
			&model.WSSession{
				Conn: conn,
			})
		break
	default:
		data := wsMessage.Data
		w.storageResult(msgType, data)
		logger.Debugf("recive data %s", data)
	}
}

func (w *WSMessageService) storageResult(msgId string, data interface{}) {
	//
	w.ResultMap.Set(msgId, data)
}

func (w *WSMessageService) ReleaseResult(msgId string) {
	//
	w.ResultMap.Remove(msgId)
}

func SendTextMessage(wssession model.WSSession, text string) error {
	//
	return wssession.Conn.WriteMessage(websocket.TextMessage, []byte(text))
}
