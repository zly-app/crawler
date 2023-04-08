package config

const (
	// 默认管道类型
	defaultPipelineType = "stdout"
)

type PipelineConfig struct {
	Type string // 管道类型, 多管道用英文逗号分隔
}

func newPipelineConfig() PipelineConfig {
	return PipelineConfig{}
}

func (conf *PipelineConfig) Check() error {
	if conf.Type == "" {
		conf.Type = defaultPipelineType
	}
	return nil
}
