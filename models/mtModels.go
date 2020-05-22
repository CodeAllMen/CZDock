package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type MtReqData struct {
	ID                int64  `orm:"pk;auto;column(id)"` //自增ID
	MtType            string `orm:"column(mt_type)"`    // 0 欢迎短信   1  扣费短信
	SubscriptionID    string `orm:"column(subscription_id)"`
	Msisdn            string
	Keyword           string
	ShortCode         string
	SendTime          string `orm:"column(send_time)"`
	ResponseStatus    string `orm:"column(response_status)"`
	ResponseMessageID string `orm:"column(response_message_id)"`
	ResponseErrorCode string `orm:"column(response_error_code)"`

	MtOperator    string
	MtTime        string
	TransactionID string `orm:"column(transaction_id)"`
	MtStatus      string
}

func (mt *MtReqData) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(mt)
	if err != nil {
		logs.Error("MtReqData Insert ERROR:", err.Error())
	}
	return err
}

func (mt *MtReqData) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(mt)
	if err != nil {
		logs.Error("MtReqData Update( ERROR:", err.Error())
	}
	return err
}

func (mt *MtReqData) GetMtByMtID(mtID int) error {
	o := orm.NewOrm()
	err := o.QueryTable("mt_req_data").Filter("id", mtID).One(mt)
	if err != nil {
		logs.Error("MtReqData  GetMtByMtID  ERROR: ", err.Error(), "~~~~", mt)
	}
	return err
}
