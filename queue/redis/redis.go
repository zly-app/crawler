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

func (r *RedisQueue) Delete(queueName string) error {
	return r.client.Del(context.Background(), queueName).Err()
}

func NewRedisQueue(app zapp_core.IApp) core.IQueue {
	confKey := fmt.Sprintf("services.%s.queue", config.NowServiceType)
	client, err := redis_utils.NewRedis(app, confKey)
	if err != nil {
		app.Fatal("创建query.redis失败", zap.Error(err))
	}
	return &RedisQueue{client: client}
}
