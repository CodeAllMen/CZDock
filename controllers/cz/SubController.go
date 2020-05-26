/**
  create by yy on 2020/5/22
*/

package cz

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/models"
	"github.com/angui001/CZDock/service"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"strconv"
)

type SubController struct {
	CZBaseController
}

func (sub *SubController) OperatorLookup() {
	var (
		err           error
		serviceConfig models.ServiceInfo
		ok            bool
		errCode       int
		redirectUrl   string
	)

	trackIdStr := sub.GetString("tid")
	msisdn := sub.GetString("msisdn")

	if trackIdStr == "" {
		trackIdStr = getTrackId()
	}

	fmt.Println(trackIdStr)

	// 首先应该 获取对应的 服务名称之类的
	// 点击表是存了的，所以直接从点击表获取

	track := new(models.AffTrack)
	trackId, _ := strconv.Atoi(trackIdStr)
	track.TrackID = int64(trackId)

	if err = track.GetOne(tracking.ByTrackID); err != nil {
		err = libs.NewReportError(err)
	}

	if serviceConfig, ok = sub.serviceConfig(track.ServiceID); !ok {
		fmt.Println("获取 service config 出错")
	}

	// 开始流程
	if err, errCode, redirectUrl = service.OperatorLookupService(&serviceConfig, track, msisdn); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		sub.Data["json"] = libs.Success("failed")
	} else {
		sub.Data["json"] = libs.Success("ok")
	}

	// 介绍一下errCode
	// 因为有很多种情况，所以用errCode 来判断具体的错误
	// 0 是正常，依次执行,
	// 1 是用户手机号出现了已订阅情况，也就是在订阅期限内，自动跳转到内容页
	// 2 执行正确，然后根据 对方的要求 如果满足条件就是2，进行跳转
	// 其他情况，未完待续
	switch errCode {
	case 1:
		sub.RedirectURL(serviceConfig.ContentUrl)
	case 2:
		sub.RedirectURL(redirectUrl)
	default:
		// 默认是跳谷歌，但是为了确认错误，跳到错误页
		sub.RedirectURL("")
	}

	sub.ServeJSON()
}

// operator-lookup的回调 控制器
func (sub *SubController) OperatorLookupCallBack() {
	var (
		bodyData []byte
		err      error
	)

	if bodyData, err = ioutil.ReadAll(sub.Ctx.Request.Body); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	// 存日志，便于后续进行数据 提取 和 操作
	fmt.Println("operator-lookup callback data ========> ", string(bodyData))

	sub.ServeJSON()
}

// 开始订阅
func (sub *SubController) StartSub() {
	//

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
