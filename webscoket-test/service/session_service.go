package service

import (
	"encoding/json"
	"errors"
	"github.com/adrianbrad/queue"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"
	"webscoket-test/model"
)

var sessionSvcSessionOnce sync.Once
var sessionSvr *SessionService

type SessionService struct {
	SessionQueue       *queue.Linked[*model.WSSession]
	ActiveSessionQueue *queue.Linked[*model.WSSession] //
}

func GetSessionService() *SessionService {
	sessionSvcSessionOnce.Do(func() {
		sessionSvr = &SessionService{
			SessionQueue:       queue.NewLinked([]*model.WSSession{}),
			ActiveSessionQueue: queue.NewLinked([]*model.WSSession{}),
		}
	})
	return sessionSvr
}

func (s *SessionService) registerSession(conn *websocket.Conn, data interface{}) {

	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Printf("注册数据序列化失败: %v", err)
		return
	}
	var registerData model.RegisterData
	err = json.Unmarshal(bytes, &registerData)
	if err != nil {
		logger.Printf("registerData: 反序列化失败 %v", err)
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
	_ = s.SaveSession(wsSession)
	s.addSessionToActiveQueue(wsSession)
}

func (s *SessionService) SaveSession(session *model.WSSession) error {
	//
	if session != nil {
		return errors.New("SaveSession: session is nil")
	}
	return s.SessionQueue.Offer(session)
}

func (s *SessionService) RestoreSession(session *model.WSSession) error {

	return s.SaveSession(session)
}

func (s *SessionService) SaveActiveSession(session *model.WSSession) error {
	//
	if session != nil {
		return errors.New("SaveActiveSession: session is nil")
	}
	return s.ActiveSessionQueue.Offer(session)
}

func (s *SessionService) RestoreActiveSession(session *model.WSSession) error {

	return s.SaveActiveSession(session)
}

func (s *SessionService) GetSession(username string, status int) (*model.WSSession, error) {

	session, err := s.GetSessionWithName(username)
	if err != nil {
		for i := 0; i < 50; i++ {

			session, err = s.SessionQueue.Get()
			if err != nil {
				continue
			}
			if session.DeviceInfo.Status == status {
				return session, nil
			}
			//不符合条件归还
			_ = s.RestoreSession(session)
		}
	}
	return session, nil
}

func (s *SessionService) GetSessionWithName(username string) (*model.WSSession, error) {
	var session *model.WSSession
	iterator := s.SessionQueue.Iterator()
	for e := range iterator {
		if e.DeviceInfo.Name == username {
			//取出但不回收
			session = e
			continue
		}

		//迭代器取出全部元素，现在回收
		_ = s.RestoreSession(e)
		//if err != nil {
		//	fmt.Printf("添加元素失败，原因: %v",err)
		//}
	}

	if session != nil {
		return session, nil
	}

	return nil, errors.New("GetSessionWithName: session is nil")
}

func (s *SessionService) GetSessionWithStatus(status int) (*model.WSSession, error) {

	var session *model.WSSession
	iterator := s.SessionQueue.Iterator()
	for e := range iterator {
		if e.DeviceInfo.Status == status {
			//取出但不回收
			session = e
			continue
		}

		//迭代器取出全部元素，现在回收
		_ = s.RestoreSession(e)
		//if err != nil {
		//	fmt.Printf("添加元素失败，原因: %v",err)
		//}
	}
	if session != nil {
		return session, nil
	}

	return session, errors.New("GetSessionWithStatus: session is nil")
}

func (s *SessionService) RemoveSession(conn *websocket.Conn) {

	iterator := s.SessionQueue.Iterator()
	for e := range iterator {
		if e.Conn != conn {
			//迭代器取出全部元素，现在回收,回收不符合条件的连接
			_ = s.RestoreSession(e)
		}
	}
}

func (s *SessionService) GetSessionWithConn(conn *websocket.Conn) (*model.WSSession, error) {
	var session *model.WSSession
	iterator := s.SessionQueue.Iterator()
	for e := range iterator {
		if e.Conn == conn {
			session = e
			continue
		}
		//迭代器取出全部元素，现在回收
		_ = s.RestoreSession(e)
	}
	if session != nil {
		return session, nil
	}
	return nil, errors.New("GetSessionWithConn session is nil")
}

func (s *SessionService) GetActiveSessionWithConn(conn *websocket.Conn) (*model.WSSession, error) {
	var session *model.WSSession
	iterator := s.ActiveSessionQueue.Iterator()
	for e := range iterator {
		if e.Conn == conn {

			session = e
			continue
		}
		//迭代器取出全部元素，现在回收
		_ = s.RestoreActiveSession(e)
	}
	if session != nil {
		return session, nil
	}
	return nil, errors.New("GetSessionWithConn session is nil")
}

func (s *SessionService) addSessionToActiveQueue(session *model.WSSession) {

	if session == nil {
		return
	}
	if session.DeviceInfo == nil { //心跳
		session, err := s.GetSessionWithConn(session.Conn)
		if err == nil {
			session.DeviceInfo.LastActiveTime = time.Now() //更新活跃时间
			_ = s.RestoreSession(session)
		}

	}

	activeSession, err := s.GetActiveSessionWithConn(session.Conn)
	if err != nil {
		//不存在,直接保存
		_ = s.SaveActiveSession(session)
	} else {
		//归还
		_ = s.RestoreActiveSession(activeSession)
	}

}

func (s *SessionService) UpdateRegisterInfo() {

}
