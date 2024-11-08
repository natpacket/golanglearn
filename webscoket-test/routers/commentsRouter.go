package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["webscoket-test/controllers:WebSocketController"] = append(beego.GlobalControllerRouter["webscoket-test/controllers:WebSocketController"],
        beego.ControllerComments{
            Method: "WebSocketCtrl",
            Router: `/test`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
