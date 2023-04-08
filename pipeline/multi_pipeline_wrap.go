package pipeline

import (
	"context"

	"github.com/zly-app/zapp/logger"
	zapputils "github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/utils"
)

const multiPipelineName = "multi_pipeline"

type multiPipelineWrap struct {
	pipes []core.IPipeline
}

func (m *multiPipelineWrap) Name() string { return multiPipelineName }

func (m *multiPipelineWrap) Process(ctx context.Context, spiderName string, data interface{}) error {
	var handlers []func() error
	for i := range m.pipes {
		pipe := m.pipes[i]
		handlers = append(handlers, func() error {
			pipeCtx := utils.Trace.TraceStart(ctx, "pipeline."+pipe.Name())
			defer utils.Trace.TraceEnd(pipeCtx)

			utils.Trace.TraceEvent(pipeCtx, "process")
			err := pipe.Process(pipeCtx, spiderName, data)
			if err != nil {
				utils.Trace.TraceErrEvent(pipeCtx, "process", err)
				logger.Log.Error(ctx, "执行pipeline失败", zap.String("pipeline", pipe.Name()), zap.Error(err))
			}
			return err
		})
	}
	err := zapputils.Go.GoAndWait(handlers...)
	return err
}

func (m *multiPipelineWrap) Close(ctx context.Context) error {
	var handlers []func() error
	for i := range m.pipes {
		pipe := m.pipes[i]
		handlers = append(handlers, func() error {
			err := pipe.Close(ctx)
			if err != nil {
				logger.Log.Error(ctx, "关闭pipeline失败", zap.String("pipeline", pipe.Name()), zap.Error(err))
			}
			return err
		})
	}
	err := zapputils.Go.GoAndWait(handlers...)
	return err
}

func NewMultiPipelineWrap(pipes []core.IPipeline) core.IPipeline {
	return &multiPipelineWrap{
		pipes: pipes,
	}
}
