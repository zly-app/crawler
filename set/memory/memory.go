package memory

import (
	"sync"

	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/crawler/core"
)

type MemorySet struct {
	sets map[string]map[string]struct{}
	mx   sync.Mutex
}

func (m *MemorySet) Add(key string, items ...string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	set, ok := m.sets[key]
	if !ok {
		set = make(map[string]struct{}, 100)
		m.sets[key] = set
	}

	var count int
	l := len(set) // 原来的大小
	for _, item := range items {
		set[item] = struct{}{}
		nl := len(set) // 新的大小
		if nl != l {   // 如果是新加入, 那么新的大小一定和原来的大小不同
			count++
			l = nl
		}
	}
	return count, nil
}

func (m *MemorySet) HasItem(key, item string) (bool, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	set, ok := m.sets[key]
	if ok {
		_, ok = set[item]
	}
	return ok, nil
}

func (m *MemorySet) Remove(key string, items ...string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	set, ok := m.sets[key]
	if !ok {
		return 0, nil
	}

	l := len(set) // 原来的大小
	if l == 0 {
		return 0, nil
	}

	var count int
	for _, item := range items {
		delete(set, item)
		nl := len(set) // 新的大小
		if nl == l {
			continue
		}

		// 如果真的删除了, 那么新的大小一定和原来的大小不同
		count++
		l = nl

		if l == 0 {
			break
		}
	}

	return count, nil
}

func (m *MemorySet) DeleteSet(key string) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	delete(m.sets, key)
	return nil
}

func (m *MemorySet) GetSetSize(key string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	set, ok := m.sets[key]
	if !ok {
		return 0, nil
	}

	return len(set), nil
}

func (m *MemorySet) Close() error { return nil }

func NewMemorySet(app zapp_core.IApp) core.ISet {
	return &MemorySet{
		sets: make(map[string]map[string]struct{}),
	}
}
