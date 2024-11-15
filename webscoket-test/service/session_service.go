package service

import (
	"encoding/json"
	"errors"
	"github.com/adrianbrad/queue"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
	"webscoket-test/model"
)

var sessionSvcSessionOnce sync.Once
var sessionSvr *SessionService

type SessionService struct {
	DeviceNameQueue    *queue.Linked[string]
	SessionMap         cmap.ConcurrentMap[string, *model.WSSession]
	WaitSessionMap     cmap.ConcurrentMap[string, *model.WSSession]
	ActiveSessionQueue *queue.Linked[*model.WSSession] //
}

func GetSessionService() *SessionService {
	sessionSvcSessionOnce.Do(func() {
		sessionSvr = &SessionService{
			DeviceNameQueue:    queue.NewLinked([]string{}),
			SessionMap:         cmap.New[*model.WSSession](),
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
	s.saveSession(wsSession)
	s.AddSessionToActiveQueue(wsSession)
}

func (s *SessionService) saveSession(session *model.WSSession) {
	//
	if session != nil {
		return
	}
	key := session.DeviceInfo.Username
	_ = s.DeviceNameQueue.Offer(key)
	s.SessionMap.Set(key, session)
}

func (s *SessionService) removeSession(session *model.WSSession) {
	if session != nil {
		key := session.DeviceInfo.Name
		s.DeviceNameQueue.Remove(key)
		s.SessionMap.Remove(key)
	}
}

func (s *SessionService) RestoreFinderSession(session *model.WSSession) {
	s.saveSession(session)
}

func (s *SessionService) saveActiveSession(session *model.WSSession) error {
	//
	if session != nil {
		return errors.New("SaveActiveSession: session is nil")
	}
	return s.ActiveSessionQueue.Offer(session)
}

func (s *SessionService) GetSession(username string, status int) (*model.WSSession, error) {

	//这里不锁,可能两个线程同时进来拿到同一个session 这里不去考虑多线程抢占问题
	session, err := s.GetSessionWithName(username)
	if err != nil {
		count := s.DeviceNameQueue.Size() * 2
		for i := 0; i < count; i++ {
			username, _ := s.DeviceNameQueue.Get()
			session, err = s.GetSessionWithName(username)
			if err != nil {
				return nil, errors.New("session is nil")
			}
			if session.DeviceInfo.Status == status {
				return session, nil
			} else {
				//还给队列
				_ = s.DeviceNameQueue.Offer(username)
				continue
			}
		}
	} else {
		s.DeviceNameQueue.Remove(username)
		return session, nil
	}

	return nil, errors.New("session is nil")
}

func (s *SessionService) GetSessionWithName(username string) (*model.WSSession, error) {

	session, ok := s.SessionMap.Get(username)
	if ok {
		return session, nil
	}
	return nil, errors.New("GetSessionWithName: session is nil")
}

func (s *SessionService) getSessionWithConn(conn *websocket.Conn) (*model.WSSession, error) {

	for _, key := range s.SessionMap.Keys() {
		session, ok := s.SessionMap.Get(key)
		if ok && session.Conn == conn {
			return session, nil
		}
	}
	return nil, errors.New("GetSessionWithConn: session is nil")
}

func (s *SessionService) RemoveSessionWithConn(conn *websocket.Conn) {

	session, err := s.getSessionWithConn(conn)
	if err != nil {
		//
		s.removeSession(session)
	}
}

func (s *SessionService) AddSessionToActiveQueue(session *model.WSSession) {

	if session == nil {
		return
	}
	if session.DeviceInfo == nil { //心跳
		_, err := s.getSessionWithConn(session.Conn)
		if err != nil {
			//不存在,直接保存
			session.DeviceInfo.LastActiveTime = time.Now() //记录活跃时间
			_ = s.saveActiveSession(session)
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

func (s *SessionService) WaitFreeSession(session *model.WSSession) {
	if session == nil {
		return
	}
	timestamp := time.Now().UnixMilli()
	s.WaitSessionMap.Set(strconv.FormatInt(timestamp, 10), session)
}

func (s *SessionService) FreeSession() {

	for _, key := range s.WaitSessionMap.Keys() {
		session, ok := s.WaitSessionMap.Get(key)
		if ok {
			timestamp, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				continue
			}
			interval := session.DeviceInfo.Interval
			currentTimeMillis := time.Now().UnixMilli()
			delta := currentTimeMillis - timestamp
			if interval < int(delta) {
				s.RestoreFinderSession(session)
			}
		}
	}
}
