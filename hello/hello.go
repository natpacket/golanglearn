package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "hello/routers" //路由需要初始化
)

func main() {
	// bee generate routers 生成路由文件
	// bee run -downdoc=true -gendoc=true  //Generate swagger documentation automatically

	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/"] = "swagger"
	}
	//routers.Init()
	web.Run()
}
