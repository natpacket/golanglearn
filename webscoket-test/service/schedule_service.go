package service

import (
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
)

func InitSchedule() {
	crontab := cron.New(cron.WithSeconds())

	crontab.AddFunc("*/10 * * * * ?", func() {
		logger.Debug("开始 [UpdateRegisterInfo] 定时任务...")
		GetSessionService().UpdateRegisterInfo()
	})

	crontab.AddFunc("*/1 * * * * ?", func() { //每秒都去看下设备是否符合释放条件
		logger.Debug("开始 [FreeSeesion] 定时任务...")
		//GetSessionService().FreeSeesion()
	})

	crontab.Start()
	defer crontab.Stop()
	select {}

}
