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

type RedisQueue struct {
	client redis.UniversalClient
}

func (r *RedisQueue) Put(ctx context.Context, queueName string, raw string, front bool) (int, error) {
	if front {
		size, err := r.client.LPush(ctx, queueName, raw).Result()
		return int(size), err
	}
	size, err := r.client.RPush(ctx, queueName, raw).Result()
	return int(size), err
}

func (r *RedisQueue) Pop(ctx context.Context, queueName string, front bool) (result string, err error) {
	if front {
		result, err = r.client.LPop(ctx, queueName).Result()
	} else {
		result, err = r.client.RPop(ctx, queueName).Result()
	}
	if err == redis.Nil {
		return "", core.EmptyQueueError
	}
	return result, err
}

func (r *RedisQueue) QueueSize(ctx context.Context, queueName string) (int, error) {
	size, err := r.client.LLen(ctx, queueName).Result()
	return int(size), err
}

func (r *RedisQueue) Close(ctx context.Context) error {
	return r.client.Close()
}

func (r *RedisQueue) Delete(ctx context.Context, queueName string) error {
	return r.client.Del(ctx, queueName).Err()
}

func NewRedisQueue(app zapp_core.IApp) core.IQueue {
	confKey := fmt.Sprintf("services.%s.queue", config.NowServiceType)
	conf := redis.NewRedisConfig()
	err := app.GetConfig().Parse(confKey, conf)
	if err != nil {
		app.Fatal("创建query.redis失败: 解析配置失败", zap.Error(err))
	}
	client, err := redis.NewClient(conf)
	if err != nil {
		app.Fatal("创建query.redis失败", zap.Error(err))
	}
	return &RedisQueue{client: client}
}
