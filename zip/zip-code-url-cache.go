package zip

import "sync"

var zipCodesToDailyUrl = make(map[string]string)
var zipCodesToHourlyUrl = make(map[string]string)
var dailyLock sync.RWMutex
var hourlyLock sync.RWMutex

func GetDailyUrlFromCache(zip string) (string, bool) {
	dailyLock.RLock()
	defer dailyLock.RUnlock()
	url, ok := zipCodesToDailyUrl[zip]
	return url, ok
}

func AcquireDailyLockForCaching() {
	dailyLock.Lock()
}

func ReleaseDailyLockForCaching() {
	dailyLock.Unlock()
}

func WriteToDailyCache(zip, url string) {
	zipCodesToDailyUrl[zip] = url
	dailyLock.Unlock()
}

func GetHourlyUrlFromCache(zip string) (string, bool) {
	hourlyLock.RLock()
	defer hourlyLock.RUnlock()
	url, ok := zipCodesToHourlyUrl[zip]
	return url, ok
}

func AcquireHourlyLockForCaching() {
	hourlyLock.Lock()
}

func ReleaseHourlyLockForCaching() {
	hourlyLock.Unlock()
}

func WriteToHourlyCache(zip, url string) {
	zipCodesToHourlyUrl[zip] = url
	hourlyLock.Unlock()
}
