package config

const (
	defaultQueueType = "memory"
)

type QueueConfig struct {
	/*
		**队列类型
		 redis
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
