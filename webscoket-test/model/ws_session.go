package model

import (
	"github.com/gorilla/websocket"
)

type WSSession struct {
	DeviceInfo *DeviceInfo
	Conn       *websocket.Conn
}
