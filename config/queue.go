package config

const (
	defaultQueueType = "memory"
)

type QueueConfig struct {
	/**队列类型
	  memory 内存, 重启后数据会丢失
	  redis redis实现
	*/
	Type string
}

func newQueueConfig() QueueConfig {
	return QueueConfig{}
}

func (conf *QueueConfig) Check() error {
	if conf.Type == "" {
		conf.Type = defaultQueueType
	}
	return nil
}
