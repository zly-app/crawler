package crawler

import (
	"github.com/zly-app/zapp"
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// 启用crawler服务
func WithService() zapp.Option {
	service.RegisterCreatorFunc(config.NowServiceType, func(app zapp_core.IApp) zapp_core.IService {
		return NewCrawler(app)
	})
	return zapp.WithService(config.NowServiceType)
}

// 注册spider
func RegistrySpider(spider core.ISpider) {
	zapp.App().InjectService(config.NowServiceType, spider)
}
