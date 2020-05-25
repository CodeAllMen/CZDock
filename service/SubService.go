/**
  create by yy on 2020/5/25
*/

package service

import (
	"fmt"
	"github.com/angui001/CZDock/models"
)

// 检查电话号码订阅状态
// true  为 已订阅
// false 为 未订阅
func checkMsisdnSubStatus(msisdn string) (ok bool) {
	// 首先根据 电话号码查询数据库
	// 这里拿到

	ok = false

	return
}

func StartSubService(serviceConfig *models.ServiceInfo, track *models.AffTrack, msisdn string) (err error, errCode int) {
	// 先检测用户的手机号 是否已经订阅
	// 如果已经订阅则返回错误信息和代码
	var (
		ok bool
	)

	if ok = checkMsisdnSubStatus(msisdn); ok {
		fmt.Println("电话号码已订阅")
		errCode = 1
		return
	}

	// MERCHANT使用提示的msisdn调用API
	// 判断 用户电话号码的 运营商

	return
}
