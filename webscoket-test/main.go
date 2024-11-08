package main

import (
	_ "webscoket-test/routers"
	"webscoket-test/test"
)

func main() {
	//beego.BConfig.WebConfig.DirectoryIndex = true
	//beego.BConfig.WebConfig.StaticDir["/"] = "swagger"
	////beego.SeySetLogFuncCall(false)
	////自定义错误页面
	//beego.Run()
	test.TestMysql()
}
