package config

const (
	defaultSetType = "memory"
)

type SetConfig struct {
	/**集合类型
	  memory 内存, 重启后数据会丢失
	  redis redis实现
	  ssdb ssdb实现
	*/
	Type string
}

func newSetConfig() SetConfig {
	return SetConfig{}
}

func (conf *SetConfig) Check() error {
	if conf.Type == "" {
		conf.Type = defaultSetType
	}
	return nil
}
