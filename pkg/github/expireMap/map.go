package expireMap

import (
	"fmt"
	"sync"
	"time"
)

const (
	Expiresuffix = "_expire"
	SleepTime    = time.Hour
)

type ExpireMap struct {
	mp1       sync.Map
	mp2       sync.Map
	sleepTime time.Duration
}

func NewExpireMap() (ExpireMap, func()) {
	e := ExpireMap{
		mp1:       sync.Map{},
		mp2:       sync.Map{},
		sleepTime: SleepTime,
	}
	//顺便返回一个定时清理函数，后面注入到app中
	return e, func() {
		e.Clean()
	}
}
func (e *ExpireMap) Store(key int64, value any, expire time.Duration) {
	e.mp1.Store(getKey(key), value)
	e.mp2.Store(getExpireKey(key), time.Now().Add(expire))
}
func (e *ExpireMap) Load(key int64) (val any, exist bool) {
	return e.mp1.Load(getKey(key))
}
func (e *ExpireMap) Delete(key int64) {
	e.mp1.Delete(getKey(key))
	e.mp2.Delete(getExpireKey(key))
}
func (e *ExpireMap) Clean() {
	for {
		e.mp1.Range(func(key, value interface{}) bool {
			deadline, _ := e.mp2.Load(key.(string) + Expiresuffix)
			if checkIfAfter(deadline.(time.Time)) {
				e.mp1.Delete(key.(string))
				e.mp2.Delete(key.(string) + Expiresuffix)
			}
			return true
		})
		//睡眠,等待下次检测
		time.Sleep(e.sleepTime)
	}
}
func checkIfAfter(deadline time.Time) bool {
	if time.Now().Before(deadline) {
		return false
	}
	return true
}
func getKey(key int64) string {
	return fmt.Sprintf("%d", key)
}
func getExpireKey(key int64) string {
	return fmt.Sprintf("%d", key) + Expiresuffix
}
