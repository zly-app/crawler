package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type RedisQueue struct {
	client redis.UniversalClient
}

func (r *RedisQueue) Put(queueName string, raw string, front bool) (int, error) {
	if front {
		size, err := r.client.LPush(context.Background(), queueName, raw).Result()
		return int(size), err
	}
	size, err := r.client.RPush(context.Background(), queueName, raw).Result()
	return int(size), err
}

func (r *RedisQueue) Pop(queueName string, front bool) (result string, err error) {
	if front {
		result, err = r.client.LPop(context.Background(), queueName).Result()
	} else {
		result, err = r.client.RPop(context.Background(), queueName).Result()
	}
	if err == redis.Nil {
		return "", core.EmptyQueueError
	}
	return result, err
}

func (r *RedisQueue) QueueSize(queueName string) (int, error) {
	size, err := r.client.LLen(context.Background(), queueName).Result()
	return int(size), err
}

func (r *RedisQueue) Close() error {
	return r.client.Close()
}

func NewRedisQueue(app zapp_core.IApp) core.IQueue {
	conf := newRedisConfig()
	confKey := fmt.Sprintf("services.%s.queue.redis", config.NowServiceType)
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		app.Fatal("redis队列配置错误", zap.Error(err))
	}

	var client redis.UniversalClient
	if conf.IsCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(conf.Address, ","),
			Username:     conf.UserName,
			Password:     conf.Password,
			MinIdleConns: conf.MinIdleConns,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         conf.Address,
			Username:     conf.UserName,
			Password:     conf.Password,
			DB:           conf.DB,
			MinIdleConns: conf.MinIdleConns,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		})
	}
	return &RedisQueue{client: client}
}
