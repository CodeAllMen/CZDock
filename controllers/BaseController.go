package controllers

import (
	"github.com/MobileCPX/PreBaseLib/common"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/angui001/CZDock/models"
	"github.com/astaxie/beego/logs"
)

type BaseController struct {
	common.BaseController
}

func (c *BaseController) getServiceConfig(serviceID string) models.ServiceInfo {
	serviceConfig, isExist := c.serviceConfig(serviceID)
	if !isExist {
		logs.Error("服务名称不存在，请检查服务信息，servideID: ", serviceID)
		c.RedirectURL("http://google.com")
	}
	return serviceConfig
}

func (c *BaseController) serviceConfig(serviceID string) (models.ServiceInfo, bool) {
	serviceConfig, isExist := models.ServiceData[serviceID]
	return serviceConfig, isExist
}

func (c *BaseController) getTrackData(trackID int) *models.AffTrack {
	track := new(models.AffTrack)
	track.TrackID = int64(trackID)
	// 通过TrackID 查询
	if err := track.GetOne(tracking.ByTrackID); err != nil {
		logs.Error("getTrackData 通过trackID查询 track 数据失败,",trackID)
		c.RedirectURL("http://google.com")
	}
	return track
}
