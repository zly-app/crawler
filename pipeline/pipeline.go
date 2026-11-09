package pipeline

import (
	"strings"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/log"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/pipeline/redis_list"
	"github.com/zly-app/crawler/pipeline/stdout"
)

// Creator 是管道建造者：接收应用与配置，返回管道实例。
type Creator func(app zapp_core.IApp, conf *config.ServiceConfig) core.IPipeline

// pipelineCreators 是管道类型的注册表。它是框架级必要注册点（用户可注册自定义管道），
// 只允许通过 RegistryPipelineCreator 写入，不再由 init 隐式注册。
var pipelineCreators = map[string]Creator{
	stdout.PipelineName:     stdout.NewStdoutPipeline,
	redis_list.PipelineName: redis_list.NewRedisList,
}

func NewPipeline(app zapp_core.IApp, pipelineType string, conf *config.ServiceConfig) core.IPipeline {
	types := strings.Split(pipelineType, ",")
	pipes := make([]core.IPipeline, 0, len(types))
	for _, t := range types {
		creator, ok := pipelineCreators[t]
		if !ok {
			log.Log.Fatal("pipeline.type 未定义", zap.String("type", pipelineType))
		}
		pipe := creator(app, conf)
		pipes = append(pipes, pipe)
	}
	return NewMultiPipelineWrap(pipes)
}

// RegistryPipelineCreator 注册管道建造者, 重复注册会报错并结束程序
func RegistryPipelineCreator(pipelineType string, creator Creator) {
	if _, ok := pipelineCreators[pipelineType]; ok {
		log.Log.Fatal("重复注册pipeline建造者", zap.String("type", pipelineType))
	}
	pipelineCreators[pipelineType] = creator
}
