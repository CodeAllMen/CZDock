package redisUtil

import (
	"fmt"
	"time"

	log "github.com/cihub/seelog"
	rlib "github.com/garyburd/redigo/redis"
)

var (
	redisPool *rlib.Pool
)

func Open(host string, port int, password string) {
	redisPool = newPool(host, port, password)
}

func GetConn() rlib.Conn {
	return redisPool.Get()
}

func Close() error {
	return redisPool.Close()
}

func newPool(host string, port int, password string) *rlib.Pool {
	return &rlib.Pool{
		MaxIdle:     19999,
		IdleTimeout: 240 * time.Second,
		Dial: func() (rlib.Conn, error) {
			c, err := rlib.Dial("tcp", fmt.Sprintf("%s:%d", host, port), rlib.DialPassword(password))
			if err != nil {
				log.Error("failed to dial redis server:", err)
				return nil, err
			}
			return c, err
		},
	}
}
