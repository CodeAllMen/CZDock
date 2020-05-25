/**
  create by yy on 2020/5/22
*/

package cz

import (
	"fmt"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type SubController struct {
	beego.Controller
}

func (sub *SubController) StartSub() {
	var (
		err error
	)

	trackId := sub.GetString("tid")

	if trackId == "" {
		trackId = getTrackId()
	}

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


func getTrackId() string {
	var (
		err      error
		response *httplib.BeegoHTTPRequest
		result   []byte
	)

	postData := make(map[string]interface{})

	postData["service_id"] = "CZ197-FG"
	postData["service_name"] = "FG-CZ197-FG"
	postData["PromoterID"] = 3
	postData["offer_id"] = 193
	postData["camp_id"] = 26
	postData["aff_name"] = "AAA"
	postData["PostbackPrice"] = 0.8

	getTrackIdUrl := "http://cz.foxseekmedia.com/aff/track"

	request := httplib.Post(getTrackIdUrl)

	if response, err = request.JSONBody(postData); err != nil {
		fmt.Println(err)
	}

	if result, err = response.Bytes(); err != nil {
		fmt.Println(err)
	}

	return string(result)
}