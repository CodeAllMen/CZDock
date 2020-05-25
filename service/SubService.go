/**
  create by yy on 2020/5/25
*/

package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/httplib"
)

// 检查电话号码订阅状态
// true  为 已订阅
// false 为 未订阅
func checkMsisdnSubStatus(msisdn string) (ok bool) {
	// 首先根据 电话号码查询数据库
	// 这里拿到

	ok = false

	return
}

// digest
func generateDigest(postData map[string]string, keyOrigin string) (digest string) {
	var (
		post string
	)
	for _, data := range postData {
		post = post + data
	}

	key := []byte(keyOrigin)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(post))
	digest = string(mac.Sum(nil))

	return
}

func StartSubService(serviceConfig *models.ServiceInfo, track *models.AffTrack, msisdn string) (err error, errCode int) {
	// 先检测用户的手机号 是否已经订阅
	// 如果已经订阅则返回错误信息和代码
	var (
		ok       bool
		result   []byte
		response *httplib.BeegoHTTPRequest
	)

	if ok = checkMsisdnSubStatus(msisdn); ok {
		fmt.Println("电话号码已订阅")
		errCode = 1
		return
	}

	// MERCHANT使用提示的msisdn调用API
	// 判断 用户电话号码的 运营商
	// 构造参数
	postData := make(map[string]string)
	postData["action"] = "operator-lookup"
	postData["merchant"] = serviceConfig.MerchantId
	postData["msisdn"] = msisdn
	postData["order"] = serviceConfig.ServerOrder
	postData["redirect"] = string(track.TrackID)
	postData["url_callback"] = serviceConfig.DockUrl + "/"

	// 生成 digest
	postData["digest"] = generateDigest(postData, serviceConfig.MerchantPassword)

	request := httplib.Post(serviceConfig.ServerUrl)

	if response, err = request.JSONBody(postData); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	if result, err = response.Bytes(); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	fmt.Println(string(result))

	// 数据解析 xml 数据，然后
	// 这里应该进行重定向 到 xml数据里的 redirect url

	return
}
