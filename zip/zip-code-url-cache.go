package zip

import "sync"

var zipCodesToUrl = make(map[string]string)
var lock sync.RWMutex

func GetUrlFromCache(zip string) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()
	url, ok := zipCodesToUrl[zip]
	return url, ok
}

func AcquireLockForCaching() {
	lock.Lock()
}

func ReleaseLockForCaching() {
	lock.Unlock()
}

func WriteToCache(zip, url string) {
	zipCodesToUrl[zip] = url
	lock.Unlock()
}
