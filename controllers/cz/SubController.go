/**
  create by yy on 2020/5/22
*/

package cz

import (
	"fmt"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/service"
	"github.com/astaxie/beego"
)

type SubController struct {
	beego.Controller
}

func (sub *SubController) StartSub() {
	var (
		err error
	)

	trackId := sub.GetString("tid")

	fmt.Println(trackId)

	// 开始流程
	if err = service.StartSubService(trackId); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		sub.Data["json"] = libs.Success("failed")
	} else {
		sub.Data["json"] = libs.Success("ok")
	}

	sub.ServeJSON()
}
