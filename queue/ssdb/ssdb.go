package ssdb

import (
	"fmt"
	"net"

	"github.com/seefan/gossdb/v2/client"
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

func (s *SsdbQueue) Put(queueName string, raw string, front bool) (int, error) {
	if front {
		size, err := s.pool.GetClient().QPushFront(queueName, raw)
		return int(size), err
	}
	size, err := s.pool.GetClient().QPushBack(queueName, raw)
	return int(size), err
}

func (s *SsdbQueue) Pop(queueName string, front bool) (string, error) {
	var result client.Value
	var err error
	if front {
		result, err = s.pool.GetClient().QPopFront(queueName)
	} else {
		result, err = s.pool.GetClient().QPopBack(queueName)
	}
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", core.EmptyQueueError
	}
	return string(result), nil
}

func (s *SsdbQueue) QueueSize(queueName string) (int, error) {
	size, err := s.pool.GetClient().QSize(queueName)
	return int(size), err
}

func (s *SsdbQueue) Close() error {
	s.pool.Close()
	return nil
}

func (s *SsdbQueue) Delete(queueName string) error {
	return s.pool.GetClient().Del(queueName)
}

func NewSsdbQueue(app zapp_core.IApp) core.IQueue {
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
