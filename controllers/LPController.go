package controllers

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/astaxie/beego/logs"
)

type LPController struct {
	BaseController
}

func (c *LPController) ThLP() {
	trackID := ""
	var track tracking.Track
	track.AffName = "LP-TEST"
	track.IP = util.GetIPAddress(c.Ctx.Request)
	keyword := c.Ctx.Input.Param(":keyword")
	if keyword == "ge" {
		serviceID := "4556093-GE"
		track.ServiceID = serviceID
		trackID = track.RequestInsertTrack("http://th.eldodigital.com/aff/track")
		fmt.Println(trackID)
		c.Data["trackID"] = trackID
		c.TplName = "th/eldo-game-GE.html"
	} else if keyword == "gd" {
		serviceID := "4556093-GD"
		track.ServiceID = serviceID
		trackID = track.RequestInsertTrack("http://th.eldodigital.com/aff/track")
		fmt.Println(trackID)
		c.Data["trackID"] = trackID
		c.TplName = "th/eldo-game-GD.html"
	} else {
		c.StringResult("404")
	}
}

func (c *LPController) ThReturnPage() {
	logs.Info("ThThanksPage URL: ", c.Ctx.Input.URI())
	if c.GetString("x_api_status") == "success" {
		c.TplName = "th/thank.html"
	} else {
		c.RedirectURL("https://google.com")
	}

}
