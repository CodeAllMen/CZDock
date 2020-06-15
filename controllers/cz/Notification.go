/**
  create by yy on 2020/5/29
*/

package cz

import (
	"encoding/xml"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreBaseLib/splib/common"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/angui001/CZDock/models"
	"github.com/angui001/CZDock/models/dimoco"
	"github.com/angui001/CZDock/util"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"strings"
)

// 接收通知流程
type NotificationController struct {
	CZBaseController
}

type result struct {
	Action            string            `xml:"action"`
	ActionResult      actionResult      `xml:"action_result"`
	Reference         string            `xml:"reference"`
	RequestID         string            `xml:"request_id"`
	Customer          customer          `xml:"customer"`
	PaymentParameters paymentParameters `xml:"payment_parameters"`

	Subscription subscription `xml:"subscription"`

	Transactions      transactions      `xml:"transactions"`
	CustomParameters  customParameters  `xml:"custom_parameters"`
	AdditionalResults additionalResults `xml:"additional_results"`
}

type actionResult struct {
	Status    int    `xml:"status"`
	Code      int    `xml:"code"`
	Detail    string `xml:"detail"`
	DetailPsp string `xml:"detail_psp"`

	RedirectURL redirectURL `xml:"redirect"`
}

type customer struct {
	Msisdn   string `xml:"msisdn"`
	Country  string `xml:"country"`
	Operator string `xml:"operator"`
	IP       string `xml:"ip"`
	Language string `xml:"language"`
}

type paymentParameters struct {
	Channel string `xml:"channel"`
	Method  string `xml:"method"`
	Order   string `xml:"order"`
}
type transactions struct {
	TransactionsID transactionsID `xml:"transaction"`
}
type transactionsID struct {
	ID             string     `xml:"id"`
	Status         string     `xml:"status"`
	Amount         string     `xml:"amount"`
	BilledAmount   string     `xml:"billed_amount"`
	Currency       string     `xml:"currency"`
	SMSMessage     smsMessage `xml:"sms_message"`
	SubscriptionID string     `xml:"subscription_id"`
}
type smsMessage struct {
	ID string `xml:"id"`
}

type subscription struct {
	SubscriptionID string     `xml:"id"`
	Definition     definition `xml:"definition"`
	Status         string     `xml:"status"`
}
type definition struct {
	PeriodType   string `xml:"period_type"`
	PeriodLength int    `xml:"period_length"`
	EventCount   int    `xml:"event_count"`
	Amount       string `xml:"amount"`
	Currency     string `xml:"currency"`
}

type customParameters struct {
	CustomParameters customParameter `xml:"custom_parameter"`
}
type customParameter struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type additionalResults struct {
	AdditionalResult additionalResult `xml:"additional_result"`
}

type additionalResult struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type redirectURL struct {
	URL string `xml:"url"`
}

func (c *NotificationController) Post() {
	fmt.Println("call /notification")
	var resultBody result
	data := c.Ctx.Request.PostFormValue("data")
	digest := c.Ctx.Request.PostFormValue("digest")

	ecoder := xml.Unmarshal([]byte(data), &resultBody)
	if ecoder != nil {
		logs.Error("notification xml 解析错误", ecoder.Error())
	}
	logs.Info("request_id:", resultBody.RequestID)

	fileName := resultBody.RequestID
	if fileName == "" {
		fileName = resultBody.Transactions.TransactionsID.SubscriptionID
	}
	// 将回传的数据存储到文件，便于本地调试
	// libs.WriteDataToFile(fmt.Sprintf("xml_logs/%v", fileName), data, digest)

	fmt.Println("notification case: ", resultBody.Action)
	chargeNotify := new(dimoco.Notification)
	chargeNotify.Action = resultBody.Action
	chargeNotify.SubscriptionID = resultBody.Subscription.SubscriptionID
	chargeNotify.Operator = resultBody.Customer.Operator
	chargeNotify.Msisdn = resultBody.Customer.Msisdn
	chargeNotify.ChargeType = resultBody.Subscription.Status
	chargeNotify.ChargeStatus = resultBody.Transactions.TransactionsID.Status
	chargeNotify.RequestID = strings.ReplaceAll(resultBody.RequestID, "subscription_", "")

	fmt.Println(chargeNotify.RequestID)

	chargeNotify.SubStatus = resultBody.Subscription.Status
	chargeNotify.Order = resultBody.PaymentParameters.Order
	chargeNotify.TransactionID = resultBody.Transactions.TransactionsID.ID
	chargeNotify.XMLData = data
	fmt.Println(chargeNotify, "##############", resultBody.Subscription)

	moT := new(mo.Mo)
	var moBase common.MoBase
	// if chargeNotify.Action == "close-subscription" || chargeNotify.Action == "renew-subscription" {
	// 	_, err := moT.GetMoBySubscriptionID(chargeNotify.SubscriptionID)
	// 	if err != nil {
	// 		moBase.SubscriptionID = chargeNotify.SubscriptionID
	// 		moBase.Msisdn = resultBody.Customer.Msisdn
	// 		moBase.ServiceID = chargeNotify.Order
	// 		moBase.Operator = resultBody.Customer.Operator
	// 	}
	// }

	_, err := moT.GetMoBySubscriptionID(chargeNotify.SubscriptionID)
	if err != nil {
		moBase.SubscriptionID = chargeNotify.SubscriptionID
		moBase.Msisdn = resultBody.Customer.Msisdn
		moBase.ServiceID = chargeNotify.Order
		moBase.Operator = resultBody.Customer.Operator
	}

	switch chargeNotify.Action {
	case "start-subscription":
		// 订阅成功
		fmt.Println("   <======================> 订阅步骤")
		if chargeNotify.SubStatus == "4" || chargeNotify.SubStatus == "3" {

			track := new(models.AffTrack)
			trackID := chargeNotify.RequestID

			if trackID != "" {
				fmt.Println("trackID ===============> 不为空")
				track, _ = models.GetServiceIDByTrackID(trackID)
			} else {
				fmt.Println("trackID ===============> 为空")
			}

			if track.TrackID != 0 {
				moBase.Track = track.Track
				moBase.ServiceID = track.ServiceID
				moBase.TrackID = track.TrackID
			} else {
				moBase.ServiceID = chargeNotify.Order
			}

			moT.ShortCode = chargeNotify.Order
			moT.Keyword = chargeNotify.Order

			moT, chargeNotify.NotificationType = splib.InsertMO(moBase, false, true, "Allterco")
			// 注册电话号码及订阅ID
			serviceConfig := c.getServiceConfig(track.ServiceID)
			registereServer(serviceConfig.ContentUrl, chargeNotify.SubscriptionID)
			registereServer(serviceConfig.ContentUrl, chargeNotify.Msisdn)
			// moT, chargeNotify.NotificationType = splib.InsertMO(chargeNotify, track)
		}
	case "close-subscription":
		moBase = moT.MoBase
		chargeNotify.RequestID = fmt.Sprintf("%v", moT.TrackID)
		chargeNotify.NotificationType, _ = moT.UnsubUpdateMo(chargeNotify.SubscriptionID)

	case "renew-subscription":
		// 交易成功标识
		moBase = moT.MoBase
		if chargeNotify.ChargeStatus == "4" || chargeNotify.ChargeStatus == "5" {
			chargeNotify.NotificationType, _ = moT.AddSuccessMTNum(chargeNotify.SubscriptionID, chargeNotify.TransactionID)
		} else {
			chargeNotify.NotificationType, _ = moT.AddFailedMTNum(chargeNotify.SubscriptionID, chargeNotify.TransactionID)
		}
	case "receive-sms-info":
		_ = chargeNotify.Insert()
		return

	}
	_ = chargeNotify.Insert()

	if chargeNotify.NotificationType != "" && moBase.CampID != 0 {
		nowTime, _ := util.GetNowTimeFormat()
		sendNoti := new(admindata.Notification)
		sendNoti.OfferID = moBase.OfferID
		sendNoti.SubscriptionID = moBase.SubscriptionID
		sendNoti.ServiceID = moBase.ServiceID
		sendNoti.ClickID = moBase.ClickID
		sendNoti.CampID = moBase.CampID
		sendNoti.PubID = moBase.PubID
		sendNoti.PostbackStatus = moT.PostbackStatus
		sendNoti.PostbackMessage = moT.PostbackMessage
		sendNoti.TransactionID = chargeNotify.TransactionID
		sendNoti.AffName = moBase.AffName
		sendNoti.Msisdn = moBase.Msisdn
		sendNoti.Operator = moBase.Operator
		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = chargeNotify.NotificationType
		fmt.Println("===============> 发送数据到后台服务器")
		sendNoti.SendData(admindata.SEC)
	}

	logs.Info("notification", data, digest)
	c.Ctx.WriteString("OK")
}

func registereServer(requestUrl, userName string) {
	// resp := httplib.Get("http://www.c4fungames.com/registere/username?user_name=" + userName)
	fmt.Println(userName)
	str, err := httplib.Get(requestUrl + fmt.Sprintf("/cz/register_msisdn?msisdn=%v", userName)).String()
	if err != nil {
		// error
	}
	fmt.Println(str)

}
