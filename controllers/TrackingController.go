package controllers

import (
	"github.com/MobileCPX/PreBaseLib/databaseutil/redis"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/logs"
	mrand "math/rand"
	"strconv"
	"strings"
	"time"
)

type TrackingController struct {
	BaseController
}

func (c *TrackingController) InsertAffClick() {
	// serviceID 类型   876551:KW,PT;8765351:KW1,P4T;
	track := new(models.AffTrack)
	returnStr := ""
	defer func() {
		if returnStr == "false" {
			track.Update()
		} else if returnStr == "" {
			track.Insert()
		}
	}()
	reqTrack := new(tracking.Track)
	reqTrack, err := reqTrack.BodyToTrack(c.Ctx.Request.Body)
	if err != nil {
		logs.Error("InsertAffClick 解析json数据失败")
		c.StringResult("false")
	}

	track.Track = *reqTrack

	keywordList := splitServiceIDToKeyword(track.ServiceID)

	if !strings.Contains(track.ServiceID, "-") {
		track.ServiceID = c.getServiceID(keywordList)
	}

	if track.ServiceID == "" {
		logs.Error("InsertAffClick ServiceID 为空")
		c.StringResult("false")
		// c.RedirectURL("https://google.com")
	}

	serviceConf, isExist := c.serviceConfig(track.ServiceID)
	if !isExist {
		logs.Error("InsertAffClick  ServiceID 不存在,", track.ServiceID)
		c.StringResult("false")
	}

	track.ServiceName = serviceConf.ServiceName
	track.Keyword = serviceConf.Keyword
	track.ShortCode = serviceConf.ShortCode

	trackID, err := track.Insert()
	if err != nil {
		logs.Error("InsertAffClick  插入click数据失败")
		c.StringResult("false")
	}
	returnStr = strconv.Itoa(int(trackID)) // 返回自增ID

	c.StringResult(returnStr)
}

func splitServiceIDToKeyword(serviceID string) []string {
	// eg.  "876551:KW,PT;8765351:KW1,P4T;" ---> ["876551-KW","876551-PT","8765351-KW1","8765351-P4T"]
	var shordCodeKeywordList []string
	strSplList := strings.Split(serviceID, ";")
	for _, v := range strSplList {
		splitKeywordAndShordList := strings.Split(v, ":")
		shordCode := ""
		if len(splitKeywordAndShordList) == 2 {
			shordCode = splitKeywordAndShordList[0]
			keywordList := strings.Split(splitKeywordAndShordList[1], ",")
			for _, keyword := range keywordList {
				shordCodeKeywordList = append(shordCodeKeywordList, shordCode+"-"+keyword)
			}
		}
	}
	return shordCodeKeywordList
}

func (c *TrackingController) getServiceID(keywordList []string) string {
	keywordList = randomSortList(keywordList)

	for _, v := range keywordList {
		serviceConf, isExist := c.serviceConfig(v)
		if isExist {
			// 检查今日订阅数是否超过订阅限制
			nowTime, _ := util.GetNowTime()
			limitSubNum := serviceConf.LimitSubNum
			todaySubNumRedisKey := nowTime + "-" + v
			redisUtil := redis.RedUtil{}
			todaySubNum, _ := redisUtil.GetValueByKeyInt(todaySubNumRedisKey)
			if todaySubNum <= limitSubNum {
				return v
			}
		}
	}

	return ""
}

func randomSortList(list []string) []string {
	UnixNano := time.Now().UnixNano()

	var newList []string
	seed := mrand.New(mrand.NewSource(UnixNano))
	sendList := seed.Perm(len(list))
	for i := range sendList {
		newList = append(newList, list[sendList[i]])
	}
	return newList
}
