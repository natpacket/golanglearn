package service

import (
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
)

func InitSchedule() {
	crontab := cron.New(cron.WithSeconds())

	crontab.AddFunc("*/10 * * * * ?", func() { //每10秒都去看下设备是否符合释放条件
		logger.Debug("开始 [UpdateRegisterInfo] 定时任务...")
		GetSessionService().UpdateRegisterInfo()
	})

	crontab.AddFunc("* * * * * ?", func() { //每秒都去看下设备是否符合释放条件
		logger.Debug("开始 [FreeSession] 定时任务...")
		GetSessionService().FreeSession()
	})

	crontab.Start()
	defer crontab.Stop()
	select {}

}
