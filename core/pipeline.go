package core

type IPipeline interface {
	// 处理
	Process(spiderName string, data interface{}) error
	// 关闭
	Close() error
}
