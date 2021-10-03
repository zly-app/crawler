package config

const (
	// 默认管道类型
	defaultPipelineType = "stdout"
)

type PipelineConfig struct {
	Type string // 管道类型
}

func newPipelineConfig() PipelineConfig {
	return PipelineConfig{}
}

func (conf *PipelineConfig) Check() error {
	if len(conf.Type) == 0 {
		conf.Type = defaultPipelineType
	}
	return nil
}
