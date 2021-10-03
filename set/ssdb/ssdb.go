package ssdb

import (
	"fmt"

	"github.com/seefan/gossdb/v2/pool"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/utils/ssdb"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type SsdbSet struct {
	pool *pool.Connectors
}

func (s *SsdbSet) Add(key string, items ...string) (int, error) {
	if len(items) == 1 {
		return zSet(s.pool.GetClient(), key, items[0])
	}
	return multiZSet(s.pool.GetClient(), key, items...)
}

func (s *SsdbSet) HasItem(key, item string) (bool, error) {
	return s.pool.GetClient().ZExists(key, item)
}

func (s *SsdbSet) Remove(key string, items ...string) (int, error) {
	if len(items) == 1 {
		return zDel(s.pool.GetClient(), key, items[0])
	}
	return multiZDel(s.pool.GetClient(), key, items...)
}

func (s *SsdbSet) DeleteSet(key string) error {
	return s.pool.GetClient().Del(key)
}

func (s *SsdbSet) GetSetSize(key string) (int, error) {
	size, err := s.pool.GetClient().ZSize(key)
	return int(size), err
}

func (s *SsdbSet) Close() error {
	s.pool.Close()
	return nil
}

func NewSsdbSet(app zapp_core.IApp) core.ISet {
	confKey := fmt.Sprintf("services.%s.set", config.NowServiceType)
	p, err := ssdb.NewSsdb(app, confKey)
	if err != nil {
		app.Fatal("创建set.ssdb失败", zap.Error(err))
	}
	return &SsdbSet{p}
}
