package main

import (
	logger "github.com/sirupsen/logrus"
	"webscoket-test/model"
	_ "webscoket-test/routers"
	"webscoket-test/service"
)

func main() {
	go func() {
		service.InitSchedule()
	}()
	//beego.BConfig.WebConfig.DirectoryIndex = true
	//beego.BConfig.WebConfig.StaticDir["/"] = "swagger"
	////beego.SeySetLogFuncCall(false)
	////自定义错误页面
	//beego.Run()

	deviceInfo, err := model.FindDeviceInfoByUserName("faker-username")
	logger.Infof("deviceInfo %s ,%v", deviceInfo.Name, err)
	select {}
}
