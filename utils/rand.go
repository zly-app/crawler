package utils

import (
	"math/rand"
	"sync"
	"time"
)

var Rand = &randUtil{
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	mx:   sync.Mutex{},
}

type randUtil struct {
	rand *rand.Rand
	mx   sync.Mutex
}

// 随机 [0, max-1] 之间的随机数
func (u *randUtil) Rand(max int64) int64 {
	return u.RandStart(0, max)
}

// 随机返回 [start, end-1] 之间的随机数
func (u *randUtil) RandStart(start, end int64) int64 {
	if end <= start {
		return 0
	}

	u.mx.Lock()
	v := u.rand.Int63n(end - start)
	u.mx.Unlock()
	return v + start
}
