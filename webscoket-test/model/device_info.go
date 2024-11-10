package model

import (
	"time"
	"webscoket-test/database"
)

type DeviceInfo struct {
	Name           string
	Username       string
	Status         int
	Interval       int
	LastActiveTime time.Time
}

func (DeviceInfo) TableName() string {
	return "acq_finder_pc_device"
}

func FindDeviceInfoByUserName(username string) *DeviceInfo {
	deviceInfo := &DeviceInfo{}
	database.DB.First(deviceInfo, "username = ?", username)
	return deviceInfo
}
