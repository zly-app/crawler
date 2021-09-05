package ssdb

import (
	"fmt"
	"net"

	rconf "github.com/seefan/gossdb/v2/conf"
	"github.com/seefan/gossdb/v2/pool"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/seefan/gossdb/v2"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type SsdbQueue struct {
	pool *pool.Connectors
}

func (s *SsdbQueue) Add(key string, items ...string) (int, error) {
	if len(items) == 1 {
		has, err := s.pool.GetClient().ZExists(key, items[0])
		if has || err != nil {
			return 0, err
		}
		return 1, s.pool.GetClient().ZSet(key, items[0], 1)
	}
	a := make(map[string]int64, len(items))
	for _, item := range items {
		a[item] = 1
	}
	return len(items), s.pool.GetClient().MultiZSet(key, a)
}

func (s *SsdbQueue) HasItem(key, item string) (bool, error) {
	return s.pool.GetClient().ZExists(key, item)
}

func (s *SsdbQueue) Remove(key string, items ...string) (int, error) {
	if len(items) == 1 {
		return 1, s.pool.GetClient().ZDel(key, items[0])
	}
	return len(items), s.pool.GetClient().MultiZDel(key, items...)
}

func (s *SsdbQueue) GetSetSize(key string) (int, error) {
	size, err := s.pool.GetClient().ZSize(key)
	return int(size), err
}

func (s *SsdbQueue) Close() error {
	s.pool.Close()
	return nil
}

func NewSsdbSet(app zapp_core.IApp) core.ISet {
	conf := newSsdbConfig()
	confKey := fmt.Sprintf("services.%s.queue", config.NowServiceType)
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		app.Fatal("queue.ssdb配置错误", zap.Error(err))
	}

	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		app.Fatal("queue.ssdb配置错误, 无法解析addres", zap.Error(err))
	}

	p, err := gossdb.NewPool(&rconf.Config{
		Host:           addr.IP.String(),
		Port:           addr.Port,
		Password:       conf.Password,
		ReadTimeout:    conf.ReadTimeout / 1e3,
		WriteTimeout:   conf.WriteTimeout / 1e3,
		ConnectTimeout: conf.DialTimeout / 1e3,
		MinPoolSize:    conf.MinIdleConns,
		MaxPoolSize:    conf.PoolSize,
		AutoClose:      true,
	})

	return &SsdbQueue{p}
}
