// @APIVersion 1.0.0
// @Title new ipad
// @Description 登录rpc 其他协议 属于半协议产品（仅供内部交流学习使用!）.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"webscoket-test/controllers"
)

func init() {
	ns :=
		web.NewNamespace("/api",
			web.NSNamespace("/ws",
				web.NSInclude(
					&controllers.WebSocketController{},
				),
			),
		)
	web.AddNamespace(ns)
}
