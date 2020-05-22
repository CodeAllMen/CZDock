package controllers

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreBaseLib/splib/common"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type NotificationController struct {
	BaseController
}

func (c *NotificationController) Mo() {
	logs.Info("NotificationController: ", c.Ctx.Input.URI())

	notificationType := ""
	tracKID := ""
	serviceConf := models.ServiceInfo{} //服务配置
	track := new(models.AffTrack)       // 点击表
	// 先根据subID 查询mo数据
	moT := new(mo.Mo)
	var moBase = common.MoBase{} // mo 基础数据
	var reqData = models.ReqData{}

	datas := c.Ctx.Request.Form
	reqData.Shortcode = datas.Get("shortcode")
	fmt.Println("malai shortcode:", reqData.Shortcode)

	reqData.Msisdn = datas.Get("msisdn")
	reqData.Keyword = datas.Get("keyword")
	reqData.Telco = datas.Get("telco")
	reqData.TxId = datas.Get("txId")
	reqData.Message = datas.Get("message")
	reqData.Type = datas.Get("type")
	reqData.MoDateTime = datas.Get("moDateTime")
	reqData.ClickId = datas.Get("clickId") // --->  1 N  --- trackID
	reqData.Insert()

	moBase.Keyword = strings.ToUpper(reqData.Keyword)
	moBase.ShortCode = reqData.Shortcode

	if len(reqData.ClickId) > 1 {
		tracKID = reqData.ClickId[1:]
	}

	if reqData.Type == "S" { // 订阅成功
		// 判断数据是否是重复数据
		subID := reqData.Shortcode + "-" + reqData.Keyword + "-" + tracKID
		_, _ = moT.GetMoBySubscriptionID(subID)
		if moT.ID != 0 {
			// 重复数据
			return
		}

		// 获取Track data
		track := new(models.AffTrack)
		trackIDInt, _ := strconv.Atoi(tracKID)
		track.TrackID = int64(trackIDInt)
		_ = track.GetOne(tracking.ByTrackID)
		moBase.ServiceID = ""
		if track.TrackID != 0 {
			moBase.Track = track.Track
			moBase.ServiceID = track.ServiceID
			moBase.TrackID = track.TrackID
		} else {
			moBase.ServiceID = reqData.Shortcode + strings.ToUpper(reqData.Keyword)
		}
		var isExist bool
		serviceConf, isExist = c.serviceConfig(moBase.ServiceID)
		if isExist {
			moBase.Operator = serviceConf.Operator
			moBase.Price = serviceConf.Price
			moBase.Country = serviceConf.Country
		}
		moBase.Msisdn = reqData.Msisdn

		//moBase

		moBase.SubscriptionID = subID

		// 存入MO数据
		moT, notificationType = splib.InsertMO(moBase, false, true, "Allterco")

	} else if reqData.Type == "U" {
		moT = new(mo.Mo)
		// 判断数据是否是重复数据
		subID := reqData.Shortcode + "-" + reqData.Keyword + "-" + tracKID
		_, _ = moT.GetMoBySubscriptionID(subID)
		if moT.ID != 0 {
			notificationType, _ = moT.UnsubUpdateMo(subID)
		}
	}

	// 目的 存 订阅和退订通知
	// 订阅  ——————》type=S
	// 存MO数据，回传网盟  --> postbackURL

	// 退订  ----》type=U

	// 订阅成功后将数据传到Admin服务器
	if notificationType != "" {
		notify := new(notification.Notification)
		notify.NotificationType = notificationType
		notify.SubscriptionId = moBase.SubscriptionID
		notify.TransactionID = moBase.SubscriptionID + "-" + util.RandomString(4)
		notify.ServiceID = track.ServiceID
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.PostbackPrice = moT.PostbackPrice

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID
		sendNoti.ServiceID = moT.ServiceID
		sendNoti.ClickID = moT.ClickID
		sendNoti.Msisdn = moT.Msisdn
		sendNoti.CampID = serviceConf.CampID
		sendNoti.PubID = moT.PubID
		sendNoti.PostbackStatus = moT.PostbackStatus
		sendNoti.PostbackMessage = moT.PostbackMessage
		sendNoti.TransactionID = notify.TransactionID
		sendNoti.Keyword = moT.Keyword
		sendNoti.ShortCode = moT.ShortCode
		sendNoti.AffName = moT.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}
		sendNoti.Operator = moT.Operator

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		go sendNoti.SendData(admindata.PROD)
	}

	c.Ctx.WriteString("OK")
}

func (c *NotificationController) MT() {
	var reqData = models.ReqData{}
	notificationType := ""

	datas := c.Ctx.Request.Form

	reqData.Shortcode = datas.Get("shortcode")
	reqData.Msisdn = datas.Get("msisdn")
	reqData.Keyword = datas.Get("keyword")
	reqData.RefCode = datas.Get("refCode")
	reqData.Telco = datas.Get("telco")
	reqData.TxId = datas.Get("txId")
	reqData.Message = datas.Get("message")
	reqData.Type = datas.Get("type")
	reqData.MtDateTime = datas.Get("mtDateTime")
	reqData.MtStatus = datas.Get("status")
	reqData.ClickId = datas.Get("clickId") // --->  1 N  --- trackID
	reqData.Price = datas.Get("price")     // --->  1 N  --- trackID
	reqData.MtMsg = datas.Get("mtMsg")

	reqData.Insert()

	moT := new(mo.Mo)
	// 先根据subID 查询mo数据
	_, _ = moT.GetMoByMsisdnShortCodeAndKeywordID(reqData.Msisdn, reqData.Shortcode, strings.ToUpper(reqData.Keyword))

	if moT.ID != 0 {
		if strings.ToUpper(reqData.MtStatus) == "SUCCESS" {
			notificationType, _ = moT.AddSuccessMTNum(moT.SubscriptionID, reqData.TxId)
		} else {
			notificationType, _ = moT.AddFailedMTNum(moT.SubscriptionID, reqData.TxId)
		}
	}

	if notificationType != "" {
		notify := new(notification.Notification)
		notify.NotificationType = notificationType
		notify.SubscriptionId = moT.SubscriptionID
		notify.TransactionID = reqData.TxId
		notify.ServiceID = moT.ServiceID
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID

		sendNoti.TransactionID = reqData.TxId
		sendNoti.AffName = moT.AffName

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		go sendNoti.SendData(admindata.PROD)
	}

	c.Ctx.WriteString("OK")
}
