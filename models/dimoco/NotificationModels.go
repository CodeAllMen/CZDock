package dimoco

import (
	"github.com/MobileCPX/PreDimoco/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// GoNotification 订阅，续订、退订通知
type Notification struct {
	ID               int64  `orm:"pk;auto;column(id)" json:"id"`      //自增ID
	SubscriptionID   string `orm:"column(subscription_id)" json:"subscription_id"` // 订阅id
	TransactionID    string `orm:"column(transaction_id)" json:"transaction_id"`
	Action           string
	NotificationType string `orm:"column(notification_type)"` // 通知类型
	Sendtime         string `orm:"column(sendtime);size(30)"` // 点击时间
	ChargeType       string
	RequestID        string `orm:"column(request_id)"`
	ChargeStatus     string
	Order            string

	Msisdn      string `orm:"size(20)"`
	Operator    string
	ServiceID   string
	ServiceType string
	ErrorCode   string
	XMLData     string `orm:"column(xml_data)"`
	Digest      string `orm:"column(digest)"`
	SubStatus   string
}

func (notification *Notification) TableName() string {
	return "charge_notification"
}

func (notification *Notification) Insert() error {
	o := orm.NewOrm()
	nowTime, _ := util.GetNowTimeFormat()
	notification.Sendtime = nowTime
	_, err := o.Insert(notification)
	if err != nil {
		logs.Error("Notification Insert 数据失败，ERROR: ", err.Error())
	}
	return err
}

// GetIdentifyNotificationByTrackID 根据trackID 获取通知信息
func (notification *Notification) GetIdentifyNotificationByTrackID(trackID string) error {
	o := orm.NewOrm()
	err := o.QueryTable("notification").Filter("request_id__istartswith", trackID+"_identify").
		OrderBy("-id").One(notification)
	if err != nil {
		logs.Error("GetIdentiryNotification ERROR", err.Error())
	}
	return err
}

func (notification *Notification) GetUnsubIdentiryNotification(trackID string) error {
	o := orm.NewOrm()
	err := o.QueryTable("notification").Filter("request_id", trackID).
		OrderBy("-id").One(notification)
	if err != nil {
		logs.Error("GetIdentiryNotification ERROR", err.Error())
	}
	return err
}
