package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"

	redis_utils "github.com/zly-app/crawler/utils/redis"
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
	confKey := fmt.Sprintf("services.%s.set", config.NowServiceType)
	client, err := redis_utils.NewRedis(app, confKey)
	if err != nil {
		app.Fatal("创建set.redis失败", zap.Error(err))
	}
	return &RedisSet{client: client}
}
