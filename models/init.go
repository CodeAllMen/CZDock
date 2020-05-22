package models

import (
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/postback"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(postback.Postback), new(AffTrack), new(MtReqData), new(mo.Mo),
		new(notification.Notification), new(click.HourClick), new(ReqData))
}
