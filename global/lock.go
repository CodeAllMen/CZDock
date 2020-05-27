/**
  create by yy on 2020/5/26
*/

package global

import "sync"

// 全局变量，用来进行加锁 防止 资源占用和冲突导致的 报错
var SubLock *SubLockStruct

type SubLockStruct struct {
	// 这里直接上互斥锁
	Mux *sync.Mutex
	// 一个根据 reference 存储对应 channel用于阻塞 的map
	ChanMap map[string]chan int
}

func InitSubLock() {
	SubLock = &SubLockStruct{
		Mux:     &sync.Mutex{},
		ChanMap: make(map[string]chan int),
	}
}
