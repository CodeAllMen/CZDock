package routers

import (
	"github.com/MobileCPX/PreBaseLib/splib/postback"
	"github.com/angui001/CZDock/controllers"
	"github.com/angui001/CZDock/controllers/cz"
	"github.com/astaxie/beego"
)

func init() {
	// 存点击数据
	beego.Router("/aff/track", &controllers.TrackingController{}, "Post:InsertAffClick")
	// set postback
	beego.Router("/set/postback", &postback.SetPostbackController{})

	// 订阅请求
	beego.Router("/sub/req", &cz.SubController{}, "Get:OperatorLookup")
	// operator-lookup 回调
	beego.Router("/sub/operator_lookup", &cz.SubController{}, "Post:OperatorLookupCallBack")

	// notification接口
	beego.Router("/notification", &cz.NotificationController{})

	// 泰国接收MT通知
	beego.Router("/th/", &controllers.NotificationController{}, "*:Mo")
	// 马来接收MO通知
	beego.Router("/my/mo", &controllers.NotificationController{}, "*:Mo")

	// 获取发送短信内容的服务和跳转到内容页面
	beego.Router("/content/:type/?:index", &controllers.GetMessageContentURLController{})

	// 泰国LP页面
	beego.Router("/th/:keyword", &controllers.LPController{}, "Get:ThLP")
	// 泰国欢迎页面
	beego.Router("/th-sub/return/:trackID", &controllers.LPController{}, "Get:ThReturnPage")

}
