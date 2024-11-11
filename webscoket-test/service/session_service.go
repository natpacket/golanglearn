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

func (s *SessionService) GetSession(username string, status int) (*model.WSSession, error) {

	session, err := s.GetSessionWithName(username)
	if err != nil {

		session, err = s.GetSessionWithStatus(status)
		if err != nil {
			return nil, err
		} else {
			s.RemoveSession(session)
			return session, nil
		}

	} else {
		s.RemoveSession(session)
	}

	return session, nil
}

func (s *SessionService) GetSessionWithName(username string) (*model.WSSession, error) {
	iterator := s.SessionQueue.IteratorWithNoRemoveItem()
	for e := range iterator {
		if e.DeviceInfo.Name == username {
			return e, nil
		}
	}
	return nil, errors.New("GetSessionWithName: session is nil")
}

func (s *SessionService) GetSessionWithStatus(status int) (*model.WSSession, error) {

	iterator := s.SessionQueue.IteratorWithNoRemoveItem()
	for e := range iterator {
		if e.DeviceInfo.Status == status {
			return e, nil
		}
	}

	return nil, errors.New("GetSessionWithStatus: session is nil")
}

func (s *SessionService) RemoveSessionWithConn(conn *websocket.Conn) {
	session, err := s.FindSessionWithConn(conn)
	if err != nil {
		s.RemoveSession(session)

	}
}
func (s *SessionService) RemoveSession(session *model.WSSession) {
	if session != nil {
		s.SessionQueue.Remove(session)
	}
}

func (s *SessionService) ContainsSession(session *model.WSSession) bool {
	return s.SessionQueue.Contains(session)
}

func (s *SessionService) FindSession(session *model.WSSession) (*model.WSSession, error) {
	if !s.ContainsSession(session) {
		return nil, errors.New("FindSession session is nil")
	}
	iterator := s.SessionQueue.IteratorWithNoRemoveItem()
	for e := range iterator {
		if e == session {
			return e, nil
		}
	}
	return nil, errors.New("FindSession session is nil")
}

func (s *SessionService) FindSessionWithConn(conn *websocket.Conn) (*model.WSSession, error) {

	iterator := s.SessionQueue.IteratorWithNoRemoveItem()
	for e := range iterator {
		if e.Conn == conn {
			return e, nil
		}
	}
	return nil, errors.New("FindSessionWithConn session is nil")
}

func (s *SessionService) ContainsActiveSession(session *model.WSSession) bool {
	return s.ActiveSessionQueue.Contains(session)
}

func (s *SessionService) addSessionToActiveQueue(session *model.WSSession) {

	if session == nil {
		return
	}
	if session.DeviceInfo == nil { //心跳
		has := s.ContainsActiveSession(session)
		if !has {
			//不存在,直接保存
			session.DeviceInfo.LastActiveTime = time.Now() //记录活跃时间
			_ = s.SaveActiveSession(session)
		}
	}

}

func (s *SessionService) UpdateRegisterInfo() {

	iterator := s.ActiveSessionQueue.Iterator()
	for e := range iterator {
		deviceInfo, err := model.FindDeviceInfoByUserName(e.DeviceInfo.Username)
		if err != nil {
			e.DeviceInfo.Status = 101
			e.DeviceInfo.Interval = 10 * 1000 //默认十秒
		} else {

			//影响内存中数据?
			e.DeviceInfo.Status = deviceInfo.Status
			e.DeviceInfo.Interval = deviceInfo.Interval

		}
		model.SaveDeviceInfo(e.DeviceInfo)
	}

}

func (s *SessionService) FreeSession() {

}
