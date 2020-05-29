/**
  create by yy on 2020/5/29
*/

package util

import "time"

// GetNowTimeFormat 获取当前时间格式
func GetNowTimeFormat() (nowDatetime, nowDate string) {
	time.LoadLocation("UTC")
	// h, _ := time.ParseDuration("1h")
	nowDatetime = time.Now().UTC().Format("2006-01-02 15:04:05")
	nowDate = time.Now().UTC().Format("2006-01-02")
	return
}
