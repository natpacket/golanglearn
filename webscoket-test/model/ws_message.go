package model

type WSMessage struct {
	MsgType string
	Data    interface{}
}

type RegisterData struct {
	Nickname        string
	WxUsernameForH5 string
}
