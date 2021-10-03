package config

import (
	zapp_core "github.com/zly-app/zapp/core"
)

// 默认服务类型
const DefaultServiceType zapp_core.ServiceType = "crawler"

// 当前服务类型
var NowServiceType = DefaultServiceType

var Conf *ServiceConfig

type ServiceConfig struct {
	Spider   SpiderConfig
	Frame    FrameConfig    // 框架配置
	Queue    QueueConfig    // 队列配置
	Proxy    ProxyConfig    // 代理配置
	Set      SetConfig      // 集合配置
	Pipeline PipelineConfig // 管道配置
}

func NewConfig(app zapp_core.IApp) *ServiceConfig {
	Conf = &ServiceConfig{
		Spider:   newSpiderConfig(app),
		Frame:    newFrameConfig(),
		Queue:    newQueueConfig(),
		Proxy:    newProxyConfig(),
		Set:      newSetConfig(),
		Pipeline: newPipelineConfig(),
	}
	return Conf
}
func (conf *ServiceConfig) Check() (err error) {
	if err = conf.Spider.Check(); err != nil {
		return err
	}
	if err = conf.Frame.Check(); err != nil {
		return err
	}
	if err = conf.Queue.Check(); err != nil {
		return err
	}
	if err = conf.Proxy.Check(); err != nil {
		return err
	}
	if err = conf.Set.Check(); err != nil {
		return err
	}
	if err = conf.Pipeline.Check(); err != nil {
		return err
	}
	return nil
}
