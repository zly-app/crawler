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

type RedisSet struct {
	client redis.UniversalClient
}

func (r *RedisSet) Add(key string, items ...string) (int, error) {
	a := make([]interface{}, len(items))
	for i, item := range items {
		a[i] = item
	}
	count, err := r.client.SAdd(context.Background(), key, a...).Result()
	return int(count), err
}

func (r *RedisSet) HasItem(key, item string) (bool, error) {
	return r.client.SIsMember(context.Background(), key, item).Result()
}

func (r *RedisSet) Remove(key string, items ...string) (int, error) {
	a := make([]interface{}, len(items))
	for i, item := range items {
		a[i] = item
	}
	count, err := r.client.SRem(context.Background(), key, a...).Result()
	return int(count), err
}

func (r *RedisSet) DeleteSet(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisSet) GetSetSize(key string) (int, error) {
	size, err := r.client.SCard(context.Background(), key).Result()
	return int(size), err
}

func (r *RedisSet) Close() error {
	return r.client.Close()
}

func NewRedisSet(app zapp_core.IApp) core.ISet {
	conf := newRedisConfig()
	confKey := fmt.Sprintf("services.%s.set", config.NowServiceType)
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		app.Fatal("set.redis配置错误", zap.Error(err))
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
	return &RedisSet{client: client}
}
