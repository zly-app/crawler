package pipeline

import (
	"strings"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/pipeline/redis_list"
	"github.com/zly-app/crawler/pipeline/stdout"
)

var pipelineCreator = map[string]func(app zapp_core.IApp) core.IPipeline{
	stdout.PipelineName:     stdout.NewStdoutPipeline,
	redis_list.PipelineName: redis_list.NewRedisList,
}

func NewPipeline(app zapp_core.IApp, pipelineType string) core.IPipeline {
	types := strings.Split(pipelineType, ",")
	pipes := make([]core.IPipeline, 0, len(types))
	for _, t := range types {
		creator, ok := pipelineCreator[t]
		if !ok {
			logger.Log.Fatal("pipeline.type 未定义", zap.String("type", pipelineType))
		}
		pipe := creator(app)
		pipes = append(pipes, pipe)
	}
	return NewMultiPipelineWrap(pipes)
}

// 注册管道创造者, 重复注册会报错并结束程序
func RegistryPipelineCreator(pipelineType string, creator func(app zapp_core.IApp) core.IPipeline) {
	if _, ok := pipelineCreator[pipelineType]; ok {
		logger.Log.Fatal("重复注册pipeline建造者", zap.String("type", pipelineType))
	}
	pipelineCreator[pipelineType] = creator
}
