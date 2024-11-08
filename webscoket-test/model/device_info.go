package model

import "time"

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
