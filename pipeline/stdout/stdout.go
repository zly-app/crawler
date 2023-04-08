package stdout

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

const PipelineName = "stdout"

type Stdout struct {
	Disable bool // 是否关闭
}

func (s *Stdout) Name() string { return PipelineName }
func (s *Stdout) Process(ctx context.Context, spiderName string, data interface{}) (err error) {
	if s.Disable {
		return nil
	}

	var text string
	switch t := data.(type) {
	case nil:
		text = "nil"
	case string:
		text = t
	case []byte:
		text = string(t)
	default:
		text, err = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(data)
		if err != nil {
			return fmt.Errorf("编码失败: %v", err)
		}
	}
	logger.Log.Info(ctx, "pipeline", zap.String("name", PipelineName), zap.String("data", text))
	return nil
}

func (s *Stdout) Close(ctx context.Context) error { return nil }

func NewStdoutPipeline(app zapp_core.IApp) core.IPipeline {
	confKey := fmt.Sprintf("services.%s.pipeline.%s", config.NowServiceType, PipelineName)
	s := &Stdout{}
	err := app.GetConfig().Parse(confKey, s, true)
	if err != nil {
		app.Fatal("创建pipeline失败: 解析配置失败", zap.String("pipeline", PipelineName), zap.Error(err))
	}
	return s
}
