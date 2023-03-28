package queue

import (
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/queue/memory"
	"github.com/zly-app/crawler/queue/redis"
)

var queueCreator = map[string]func(app zapp_core.IApp) core.IQueue{
	"memory": memory.NewMemoryQueue,
	"redis":  redis.NewRedisQueue,
}

func NewQueue(app zapp_core.IApp, queryType string) core.IQueue {
	creator, ok := queueCreator[queryType]
	if !ok {
		logger.Log.Fatal("queue.type 未定义", zap.String("type", queryType))
	}
	return creator(app)
}

// 注册队列创造者, 重复注册会报错并结束程序
func RegistryQueueCreator(queueType string, creator func(app zapp_core.IApp) core.IQueue) {
	if _, ok := queueCreator[queueType]; ok {
		logger.Log.Fatal("重复注册queue建造者", zap.String("type", queueType))
	}
	queueCreator[queueType] = creator
}
