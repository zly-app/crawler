package redis_list

import (
	"context"
	"fmt"
	"strings"

	"github.com/zly-app/component/redis"
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/pkg/compactor"
	"github.com/zly-app/zapp/pkg/serializer"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

const PipelineName = "redis_list"

const (
	defQueuePrefix = "result:"
	defSerializer  = "sonic_std"
	defCompactor   = "raw"
)

type RedisList struct {
	SaveToQueueFront bool   // 保存到队列前面
	QueuePrefix      string // 队列名前缀
	Serializer       string // 默认序列化器名, 可选 msgpack, sonic, sonic_std, jsoniter_standard, jsoniter, json, yaml
	Compactor        string // 默认压缩器名, 可选 raw, zstd, gzip

	client     redis.UniversalClient
	redisKey   string
	compactor  compactor.ICompactor
	serializer serializer.ISerializer
}

func (r *RedisList) Name() string { return PipelineName }

func (r *RedisList) Process(ctx context.Context, spiderName string, data interface{}) error {
	rawData, err := r.serializer.MarshalBytes(data)
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	comData, err := r.compactor.CompressBytes(rawData)
	if err != nil {
		return fmt.Errorf("压缩失败: %v", err)
	}

	if r.SaveToQueueFront {
		err = r.client.LPush(ctx, r.redisKey, comData).Err()
	} else {
		err = r.client.RPush(ctx, r.redisKey, comData).Err()
	}
	return err
}

func (r *RedisList) Close(ctx context.Context) error {
	return r.client.Close()
}

func NewRedisList(app zapp_core.IApp) core.IPipeline {
	confKey := fmt.Sprintf("services.%s.pipeline.%s", config.DefaultServiceType, PipelineName)
	rl := &RedisList{
		QueuePrefix: defQueuePrefix,
	}
	err := app.GetConfig().Parse(confKey, rl)
	if err != nil {
		app.Fatal("创建pipeline失败: 解析配置失败", zap.String("pipeline", PipelineName), zap.Error(err))
	}

	rl.redisKey = config.Conf.Frame.Namespace + ":" + rl.QueuePrefix + config.Conf.Spider.Name

	var ok bool
	if rl.Serializer == "" {
		rl.Serializer = defSerializer
	}
	rl.Serializer = strings.ToLower(rl.Serializer)
	rl.serializer, ok = serializer.TryGetSerializer(rl.Serializer)
	if !ok {
		app.Fatal("不支持的Serializer", zap.String("pipeline", PipelineName), zap.String("Serializer", rl.Serializer))
	}

	if rl.Compactor == "" {
		rl.Compactor = defCompactor
	}
	rl.Compactor = strings.ToLower(rl.Compactor)
	rl.compactor, ok = compactor.TryGetCompactor(rl.Compactor)
	if !ok {
		app.Fatal("不支持的Compactor", zap.String("pipeline", PipelineName), zap.String("Compactor", rl.Serializer))
	}

	redisConf := redis.NewRedisConfig()
	err = app.GetConfig().Parse(confKey, redisConf)
	if err != nil {
		app.Fatal("创建redis失败: 解析配置失败", zap.String("pipeline", PipelineName), zap.Error(err))
	}
	rl.client, err = redis.NewClient(redisConf)
	if err != nil {
		app.Fatal("创建redis失败", zap.String("pipeline", PipelineName), zap.Error(err))
	}
	return rl
}
