package proxy

import (
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/proxy/direct"
	"github.com/zly-app/crawler/proxy/static"
)

var proxyCreator = map[string]func(app zapp_core.IApp) core.IProxy{
	"direct": direct.NewDirectProxy,
	"static": static.NewStaticProxy,
}

func NewProxy(app zapp_core.IApp, proxyType string) core.IProxy {
	creator, ok := proxyCreator[proxyType]
	if !ok {
		logger.Log.Fatal("proxy.type 未定义", zap.String("type", proxyType))
	}
	return creator(app)
}

// 注册代理创造者, 重复注册会报错并结束程序
func RegistryProxyCreator(proxyType string, creator func(app zapp_core.IApp) core.IProxy) {
	if _, ok := proxyCreator[proxyType]; ok {
		logger.Log.Fatal("重复注册proxy建造者", zap.String("type", proxyType))
	}
	proxyCreator[proxyType] = creator
}
