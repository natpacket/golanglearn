package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"webscoket-test/model"
)

var sessionSvcSessionOnce sync.Once
var sessionSvr *SessionService

type SessionService struct {
}

func GetSessionService() *SessionService {
	sessionSvcSessionOnce.Do(func() {
		sessionSvr = &SessionService{}
	})
	return sessionSvr
}

func (s *SessionService) registerSession(conn *websocket.Conn, data interface{}) {

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("注册数据序列化失败: %v", err)
		return
	}
	var registerData model.RegisterData
	err = json.Unmarshal(bytes, &registerData)
	if err != nil {
		fmt.Printf("registerData: 反序列化失败 %v", err)
		return
	}
	wsSession := &model.WSSession{
		DeviceInfo: &model.DeviceInfo{
			Name:           registerData.Nickname,
			Username:       registerData.WxUsernameForH5,
			LastActiveTime: time.Now(),
		},
		Conn: conn,
	}
	s.SaveSession(wsSession)
}

func (s *SessionService) SaveSession(session *model.WSSession) {
	//
}
