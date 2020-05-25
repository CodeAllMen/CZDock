/**
  create by yy on 2020/5/25
*/

package service

import (
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/angui001/CZDock/libs"
	"github.com/angui001/CZDock/models"
	"strconv"
)

func StartSubService(trackIdStr string) (err error) {
	// 首先应该 获取对应的 服务名称之类的
	// 点击表是存了的，所以直接从点击表获取

	track := new(models.AffTrack)
	trackId, _ := strconv.Atoi(trackIdStr)
	track.TrackID = int64(trackId)

	if err = track.GetOne(tracking.ByTrackID); err != nil {
		err = libs.NewReportError(err)
	}

	return
}
