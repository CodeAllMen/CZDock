/**
  create by yy on 2020/5/25
*/

package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/angui001/CZDock/global"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// 检查电话号码订阅状态
// true  为 已订阅
// false 为 未订阅
func checkMsisdnSubStatus(msisdn string, serviceConfig *models.ServiceInfo) (ok bool) {
	// 首先根据 电话号码查询数据库
	// 这里拿到
	var (
		err error
	)

	ok = false

	moT := &mo.Mo{}
	if moT, err = moT.GetMoByMsisdnShortCodeAndKeywordID(msisdn, serviceConfig.ServerOrder, serviceConfig.ServerOrder); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	if moT.ID != 0 {
		// 检查时间范围
		ctime := moT.SubTime
		// 获取本地时区，
		loc, _ := time.LoadLocation("Local")
		// 指定时间模板
		l := "2006-01-02 15:04:05"
		t, _ := time.ParseInLocation(l, ctime, loc)
		// 订阅过期时间
		unsub := t.AddDate(0, 0, 7)
		// 如果在过期时间内，则跳转到内容站，否则跳转到支付页面
		if unsub.After(time.Now()) {
			logs.Info("用户", msisdn, "未超出期限")
			ok = true
		} else {
			logs.Info("用户", msisdn, "超出期限，跳转支付页面")
			ok = false
		}
	} else {
		logs.Info("用户未订阅 :", msisdn, "跳转支付页面")
	}

	return
}

// digest
func generateDigest(postData map[string]string, keyOrigin string) (digest string) {
	var (
		post string
		keys []string
	)
	// 首先取出所有键值
	for key := range postData {
		keys = append(keys, key)
	}

	// 排序
	sort.Strings(keys)
	fmt.Println(keys)

	for _, k := range keys {
		post = post + postData[k]
	}

	fmt.Println(post)

	key := []byte(keyOrigin)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(post))
	digest = hex.EncodeToString(mac.Sum(nil))

	return
}

func OperatorLookupService(serviceConfig *models.ServiceInfo, track *models.AffTrack, msisdn string) (err error, errCode int, other string) {
	// 先检测用户的手机号 是否已经订阅
	// 如果已经订阅则返回错误信息和代码
	var (
		ok             bool
		result         []byte
		operatorLookup models.OperatorLookupResult
	)

	if ok = checkMsisdnSubStatus(msisdn, serviceConfig); ok {
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
	postData["redirect"] = "0"
	postData["url_callback"] = serviceConfig.DockUrl + "/sub/operator_lookup"
	fmt.Println(track.TrackID)
	postData["request_id"] = fmt.Sprintf("%v", track.TrackID)

	// 生成 digest
	postData["digest"] = generateDigest(postData, serviceConfig.MerchantPassword)
	fmt.Println(postData["digest"])

	// 请求之前就创建阻塞map 和 需要的数据
	other = fmt.Sprintf("%v", track.TrackID)

	global.SubLock.Mux.Lock()
	global.SubLock.ChanMap[other] = make(chan int, 1)
	global.SubLock.TrackMap[other] = track
	global.SubLock.ServiceConfMap[other] = serviceConfig
	global.SubLock.MarkMap[other] = make(chan int, 1)
	global.SubLock.Mux.Unlock()

	if result, err = sendRequest(postData, serviceConfig.ServerUrl); err != nil {
		err = libs.NewReportError(err)
		return
	}
	// request := httplib.Post(serviceConfig.ServerUrl)

	// if response, err = request.JSONBody(postData); err != nil {
	// 	err = libs.NewReportError(err)
	// 	return
	// }
	//
	// if result, err = response.Bytes(); err != nil {
	// 	err = libs.NewReportError(err)
	// 	return
	// }

	fmt.Println("operator-lookup data ===========> ", string(result))

	// 数据解析 xml 数据，然后
	// 这里应该进行重定向 到 xml数据里的 redirect url
	// 这里的重定向应该判断一下 返回的参数才行
	// errCode = 2
	if err = xml.Unmarshal(result, &operatorLookup); err != nil {
		err = libs.NewReportError(err)
		return
	}
	fmt.Println("SubService operator lookup: ", operatorLookup)
	// 不管如何，数据应该入库一次
	// 数据入库操作

	switch operatorLookup.ActionResult.Status {
	case 3:
		errCode = 2
		other = operatorLookup.ActionResult.Url
	case 5:
		// 当为5的时候 就进行其他处理
		errCode = 3
		other = fmt.Sprintf("%v", track.TrackID)
	default:
		errCode = 0
	}

	return
}

func StartSubService(serviceConfig *models.ServiceInfo, track *models.AffTrack, msisdn, operator string) (redirectUrl string, err error) {
	var (
		result            []byte
		promptContentArgs []byte
		operatorLookup    models.OperatorLookupResult
	)

	fmt.Println("msisdn =======================> ", msisdn, "  operator =======> ", operator)

	// make params
	postData := make(map[string]string)
	postData["action"] = "start-subscription"
	postData["merchant"] = serviceConfig.MerchantId
	postData["order"] = serviceConfig.ServerOrder
	postData["request_id"] = fmt.Sprintf("subscription_%v", track.TrackID)
	postData["service_name"] = serviceConfig.ServiceName
	postData["url_callback"] = serviceConfig.DockUrl + "/notification"

	// 构造短信内容
	smsMap := make(map[string]map[string]string)
	smsMap["text"] = make(map[string]string)
	smsMap["text"]["en"] = fmt.Sprintf(serviceConfig.PromptContentArgsEn, serviceConfig.ContentUrl)
	smsMap["text"]["cs"] = fmt.Sprintf(serviceConfig.PromptContentArgsCs, serviceConfig.ContentUrl)
	if promptContentArgs, err = json.Marshal(&smsMap); err != nil {
		err = libs.NewReportError(err)
		return
	}

	postData["prompt_content_args"] = string(promptContentArgs)
	postData["operator"] = operator
	postData["msisdn"] = msisdn
	postData["channel"] = "web"
	postData["amount"] = "99"
	// redirect url when finish all
	postData["url_return"] = serviceConfig.StartSubReturnUrl + "?msisdn=" + msisdn

	postData["digest"] = generateDigest(postData, serviceConfig.MerchantPassword)

	fmt.Println("send subscription request: ")
	if result, err = sendRequest(postData, serviceConfig.ServerUrl); err != nil {
		err = libs.NewReportError(err)
		return
	}

	fmt.Println("unmarshal subscription xml: ")
	// 解析为结构体
	if err = xml.Unmarshal(result, &operatorLookup); err != nil {
		err = libs.NewReportError(err)
		return
	}

	fmt.Println("subscription result: ", operatorLookup)

	switch operatorLookup.ActionResult.Status {
	case 1:
		err = libs.NewReportError(errors.New("subscription failure request"))
	case 3:
		redirectUrl = operatorLookup.ActionResult.Url
	case 5:
		// 5 则是交给后续处理
		fmt.Println("后续处理")

	}

	fmt.Println("start-subscription data ===========> ", string(result))

	return
}

func sendRequest(values map[string]string, URL string) ([]byte, error) {
	// 这里添加post的body内容
	data := make(url.Values)
	for k, v := range values { // 遍历需要发送的数据
		data[k] = []string{v}
	}

	// 把post表单发送给目标服务器
	logs.Info("把post表单发送给目标服务器")
	res, err := http.PostForm(URL, data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return responseData, err
}
