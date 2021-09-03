package dom

import (
	"sync"

	"github.com/andybalholm/cascadia"
	"github.com/golang/groupcache/lru"
)

var (
	cssCacheOnce  sync.Once
	cssCache      *lru.Cache
	cssCacheMutex sync.Mutex
)

// 获取css查询器, 表达式错误会panic
func getCssQuery(expr string) cascadia.Selector {
	cssCacheOnce.Do(func() {
		cssCache = lru.New(50)
	})
	cssCacheMutex.Lock()
	defer cssCacheMutex.Unlock()
	if v, ok := cssCache.Get(expr); ok {
		return v.(cascadia.Selector)
	}
	v, err := cascadia.Compile(expr)
	if err != nil {
		panic(err)
	}
	cssCache.Add(expr, v)
	return v
}
