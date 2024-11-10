package service

import (
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
)

func InitSchedule() {
	crontab := cron.New(cron.WithSeconds())
	ss := "*/2 * * * * ?"
	_, err := crontab.AddFunc(ss, Test)
	if err != nil {
		logger.Info("err: %v\n", err)
	}
	crontab.Start()
	defer crontab.Stop()
	select {}

}

func Test() {
	logger.Info("test ")
}
