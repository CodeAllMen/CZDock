/**
  create by yy on 2020/5/22
*/

package cz

import (
	"encoding/xml"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/angui001/CZDock/global"
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
		other         string
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
		sub.RedirectURL("")
	}

	// 开始流程
	if err, errCode, other = service.OperatorLookupService(&serviceConfig, track, msisdn); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		sub.Data["json"] = libs.Success("failed")
		sub.ServeJSON()
		sub.StopRun()
	}

	// 介绍一下errCode
	// 因为有很多种情况，所以用errCode 来判断具体的错误
	// 0 是正常，依次执行,
	// 1 是用户手机号出现了已订阅情况，也就是在订阅期限内，自动跳转到内容页
	// 2 执行正确，然后根据 对方的要求 如果满足条件就是2，进行跳转
	// 3 这个就是 完全正确的情况，也就是文档里的pending 状态，进行等待同步api的数据
	// 其他情况，未完待续
	switch errCode {
	case 1:
		sub.RedirectURL(serviceConfig.ContentUrl)
	case 2:
		sub.RedirectURL(other)
	case 3:
		// 正常操作
		global.SubLock.Mux.Lock()
		if global.SubLock.ChanMap[other] == nil {
			// 如果是不存在的则自己创建
			global.SubLock.ChanMap[other] = make(chan int)
		}

		global.SubLock.TrackMap[other] = track
		global.SubLock.ServiceConfMap[other] = &serviceConfig

		// 阻塞，等待 同步 回调完成
		<-global.SubLock.ChanMap[other]
		redirectUrl = global.SubLock.RedirectUrlMap[other]
		// 删除map的 元素，防止内存 爆炸
		delete(global.SubLock.ChanMap, other)
		delete(global.SubLock.RedirectUrlMap, other)
		global.SubLock.Mux.Unlock()
		sub.RedirectURL(redirectUrl)
	}

	sub.Data["json"] = "error"

	// 默认返回数据
	sub.ServeJSON()
}

// operator-lookup的回调 控制器
func (sub *SubController) OperatorLookupCallBack() {
	var (
		bodyData               []byte
		err                    error
		operatorLookupCallback models.OperatorLookupCallback
		redirectUrl            string
	)

	if bodyData, err = ioutil.ReadAll(sub.Ctx.Request.Body); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	fmt.Println(string(bodyData))

	// 解析为结构体
	if err = xml.Unmarshal(bodyData, &operatorLookupCallback); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	// 回调之后应该解除阻塞
	if operatorLookupCallback.Result.ActionResult.Status == 0 {
		// 状态为0是 成功 success
		// 解除阻塞
		global.SubLock.Mux.Lock()
		reference := operatorLookupCallback.Result.Reference
		if global.SubLock.ChanMap[reference] == nil {
			// 如果是不存在的则自己创建
			global.SubLock.ChanMap[reference] = make(chan int)
		}

		// 这里是 开始 start_subscription 流程
		if redirectUrl, err = service.StartSubService(
			global.SubLock.ServiceConfMap[reference],
			global.SubLock.TrackMap[reference],
			operatorLookupCallback.Result.CustomParameters.Msisdn,
			operatorLookupCallback.Result.CustomParameters.Operator); err != nil {
			err = libs.NewReportError(err)
			fmt.Println(err)
		}

		global.SubLock.RedirectUrlMap[reference] = redirectUrl

		global.SubLock.ChanMap[reference] <- 1

		delete(global.SubLock.TrackMap, reference)
		delete(global.SubLock.ServiceConfMap, reference)

		global.SubLock.Mux.Unlock()
	}

	// 存日志，便于后续进行数据 提取 和 操作
	fmt.Println("operator-lookup callback data ========> ", string(bodyData))

	sub.Data["json"] = fmt.Sprintf("%v", err)

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
