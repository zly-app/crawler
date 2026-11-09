package crawler

import (
	"github.com/zly-app/zapp"
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// WithService 启用 crawler 服务，必须在 zapp 初始化前调用。
func WithService() zapp.Option {
	service.RegisterCreatorFunc(config.DefaultServiceType, func(app zapp_core.IApp) zapp_core.IService {
		return NewCrawler(app)
	})
	return zapp.WithService(config.DefaultServiceType)
}

// Install 返回 spider 注册器，必须在 zapp 初始化后调用（例如 app.Run 之前）。
func Install(app zapp_core.IApp) *Registrar {
	return &Registrar{app: app}
}

// Registrar 是 crawler 注册器：拥有「向 crawler 服务注入 spider」这条规则。
type Registrar struct {
	app zapp_core.IApp
}

// RegistrySpider 注册spider
func (r *Registrar) RegistrySpider(spider core.ISpider) {
	r.app.InjectService(config.DefaultServiceType, spider)
}
