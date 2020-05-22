package models

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreLancio/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AffTrack 网盟点击追踪
type AffTrack struct {
	TrackID  int64  `orm:"pk;auto;column(track_id)"`  // 自增ID
	Sendtime string `orm:"column(sendtime);size(30)"` // 点击时间

	SessionStatus    string
	SessionId        string
	SessionErrorCode string

	tracking.Track
}

func (track *AffTrack) TableName() string {
	return "aff_track"
}

func (track *AffTrack) TrackQuery() orm.QuerySeter {
	o := orm.NewOrm()
	return o.QueryTable(track.TableName())
}

func (track *AffTrack) Insert() (int64, error) {
	o := orm.NewOrm()
	track.Sendtime, _ = util.GetNowTimeFormat()
	track.ClickTime = track.Sendtime
	trackID, err := o.Insert(track)
	if err != nil {
		logs.Error("新插入点击错误 ", err.Error())
	}
	return trackID, err
}

func (track *AffTrack) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(track)
	if err != nil {
		logs.Error("AffTrack Update 更新点击数据失败，ERROR ", err.Error())
	}
	return err
}

func (track *AffTrack) GetOne(queryType int) error {
	o := orm.NewOrm()
	var query orm.QuerySeter
	switch queryType {
	case tracking.ByTrackID:
		query = tracking.GetTrackByTrackID(o, track.TrackID)
	}
	return query.One(track)
}

func (track *AffTrack) GetAll() {

}

func InsertHourClick() {
	o := orm.NewOrm()
	hourClick := new(click.HourClick)
	nowTime, _ := util.GetNowTimeFormat()
	nowHour := nowTime[:13]
	fmt.Println(nowHour)
	hourTime := hourClick.GetNewestClickDateTime()
	if hourTime == "" {
		hourTime = "2019-07-01"
	}

	totalHourClick := new([]click.HourClick)
	// SQL := fmt.Sprintf("SELECT left(sendtime,13) as hour_time,postback_price, (case service_id when '889-Vodafone' "+
	//	"THEN 3 WHEN '889-Three' THEN 4 WHEN '892-Vodafone' THEN 11 WHEN '892-Three' THEN 12 ELSE 0 END) as"+
	//	" camp_id, offer_id,aff_name,pub_id,count(1) as click_num ,click_status, promoter_id "+
	//	"from aff_track   where service_id <> ''  and left(sendtime,13)>'%s' and left(sendtime,13)<'%s' group by "+
	//	"left(sendtime,13),offer_id,aff_name,pub_id,"+
	//	"service_id,pro_id ,promoter_id,postback_price,click_status order by left(sendtime,13)", hourTime, nowHour)

	SQL := fmt.Sprintf("SELECT left(sendtime,13) as hour_time,postback_price, "+
		" camp_id, offer_id,aff_name,pub_id,count(1) as click_num ,click_status,short_code, keyword, promoter_id "+
		"from aff_track   where service_id <> ''  and left(sendtime,13)>'%s' and left(sendtime,13)<'%s' group by "+
		"left(sendtime,13),offer_id,aff_name,pub_id,"+
		"service_id,pro_id ,promoter_id,camp_id,postback_price,short_code,keyword,click_status order by left(sendtime,13)", hourTime, nowHour)

	num, _ := o.Raw(SQL).QueryRows(totalHourClick)
	fmt.Println(num)

	for _, v := range *totalHourClick {
		if v.CampID == 0 {
			v.CampID = ServiceData[v.ServiceID].CampID
		}

		if v.ClickNum >= 2 && v.CampID != 0 {
			o.Insert(&v)
		}
	}
}
