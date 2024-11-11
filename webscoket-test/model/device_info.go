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

func FindDeviceInfoByUserName(username string) (*DeviceInfo, error) {
	deviceInfo := &DeviceInfo{}
	tx := database.DB.First(deviceInfo, "username = ?", username)
	return deviceInfo, tx.Error
}

func SaveDeviceInfo(deviceInfo *DeviceInfo) {

	database.DB.Save(deviceInfo)

}
