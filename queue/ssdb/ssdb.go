package ssdb

import (
	"fmt"

	"github.com/seefan/gossdb/v2/client"
	"github.com/seefan/gossdb/v2/pool"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/utils/ssdb"
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
	confKey := fmt.Sprintf("services.%s.queue", config.NowServiceType)
	p, err := ssdb.NewSsdb(app, confKey)
	if err != nil {
		app.Fatal("创建queue.ssdb失败", zap.Error(err))
	}
	return &SsdbQueue{p}
}
