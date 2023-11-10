// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"hello/controllers"
)

func init() {
	ns :=
		web.NewNamespace("/v1",
			web.NSNamespace("/cms",
				web.NSInclude(
					&controllers.CMSController{},
				),
			),
		)
	web.AddNamespace(ns)
}
