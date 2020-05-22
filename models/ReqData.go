package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type ReqData struct {
	ID         int `orm:"pk;auto;column(id)"`
	DataType   string
	Msisdn     string `json:"msisdn"`
	Keyword    string `json:"keyword"`
	Shortcode  string `json:"shortcode"`
	Message    string `json:"message"`
	MoDateTime string `json:"moDateTime"`
	MtDateTime string
	Telco      string `json:"telco"`
	TxId       string `json:"txId"`
	Type       string `json:"type"`
	ClickId    string `json:"clickId"`
	MtMsg      string
	Price      string
	RefCode    string
	MtStatus   string
}

func (reqData *ReqData) Insert() {
	o := orm.NewOrm()
	_, err := o.Insert(reqData)
	if err != nil {
		logs.Error("ReqData  Insert  ERROR: ", err.Error())
	}
}
