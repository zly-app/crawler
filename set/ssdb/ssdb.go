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
		return zSet(s.pool.GetClient(), key, items[0])
	}
	return multiZSet(s.pool.GetClient(), key, items...)
}

func (s *SsdbQueue) HasItem(key, item string) (bool, error) {
	return s.pool.GetClient().ZExists(key, item)
}

func (s *SsdbQueue) Remove(key string, items ...string) (int, error) {
	if len(items) == 1 {
		return zDel(s.pool.GetClient(), key, items[0])
	}
	return multiZDel(s.pool.GetClient(), key, items...)
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
