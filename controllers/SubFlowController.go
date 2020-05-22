package controllers

import (
	"fmt"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type SubFlowController struct {
	BaseController

	TrackID int

	track *models.AffTrack

	serviceConf models.ServiceInfo
}

func (c *SubFlowController) Prepare() {
	trackID := c.GetString("tid")
	trackIDInt, _ := strconv.Atoi(trackID)
	if trackIDInt == 0 {
		logs.Error("SubFlowController Prepare 将trackID转为INT类型失败", trackID)
		c.RedirectURL("https://google.com")
	}

	c.TrackID = trackIDInt

	c.track = c.getTrackData(trackIDInt)

	// 配置信息
	c.serviceConf = c.getServiceConfig(c.track.Track.ServiceID)
	fmt.Println(c.serviceConf)
}

// /sub/req?tid={track_id}
func (c *SubFlowController) SubReq() {
	// trackID := c.GetString("tid")
	// 自己生成TrackID
	AocUrl := strings.Replace(c.serviceConf.AocURL, "{track_id}", strconv.Itoa(c.TrackID), -1)
	// http://aoc-ic.smsgateway.cc/AIS/4208128/S2/?clickId=1{TRACK_ID}&returnURL=
	// TO AOC
	c.RedirectURL(AocUrl)
}
