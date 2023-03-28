package redis

import (
	"context"
	"fmt"

	"github.com/zly-app/component/redis"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type RedisSet struct {
	client redis.UniversalClient
}

func (r *RedisSet) Add(ctx context.Context, key string, items ...string) (int, error) {
	a := make([]interface{}, len(items))
	for i, item := range items {
		a[i] = item
	}
	count, err := r.client.SAdd(ctx, key, a...).Result()
	return int(count), err
}

func (r *RedisSet) HasItem(ctx context.Context, key, item string) (bool, error) {
	return r.client.SIsMember(ctx, key, item).Result()
}

func (r *RedisSet) Remove(ctx context.Context, key string, items ...string) (int, error) {
	a := make([]interface{}, len(items))
	for i, item := range items {
		a[i] = item
	}
	count, err := r.client.SRem(ctx, key, a...).Result()
	return int(count), err
}

func (r *RedisSet) DeleteSet(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisSet) GetSetSize(ctx context.Context, key string) (int, error) {
	size, err := r.client.SCard(ctx, key).Result()
	return int(size), err
}

func (r *RedisSet) Close(ctx context.Context) error {
	return r.client.Close()
}

func NewRedisSet(app zapp_core.IApp) core.ISet {
	confKey := fmt.Sprintf("services.%s.set", config.NowServiceType)
	conf := redis.NewRedisConfig()
	err := app.GetConfig().Parse(confKey, conf)
	if err != nil {
		app.Fatal("创建set.redis失败: 解析配置失败", zap.Error(err))
	}
	client, err := redis.NewClient(conf)
	if err != nil {
		app.Fatal("创建set.redis失败", zap.Error(err))
	}
	return &RedisSet{client: client}
}
