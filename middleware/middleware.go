package middleware

import (
	"context"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/log"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/middleware/request_middleware"
	"github.com/zly-app/crawler/middleware/response_middleware"
	"github.com/zly-app/crawler/utils"
)

// Middleware 拥有中间件链这份数据；默认中间件由构造时注入，用户中间件通过
// Registry* 方法追加。不再有包级可变切片。
type Middleware struct {
	conf                *config.ServiceConfig
	requestMiddlewares  []core.IRequestMiddleware
	responseMiddlewares []core.IResponseMiddleware
}

// RegistryRequestMiddleware 注册请求中间件
func (m *Middleware) RegistryRequestMiddleware(mw core.IRequestMiddleware) {
	m.requestMiddlewares = append(m.requestMiddlewares, mw)
}

// RegistryResponseMiddleware 注册响应中间件
func (m *Middleware) RegistryResponseMiddleware(mw core.IResponseMiddleware) {
	m.responseMiddlewares = append(m.responseMiddlewares, mw)
}

func (m *Middleware) RequestProcess(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	ctx = utils.Trace.TraceStart(ctx, "Middleware.RequestProcess")
	defer utils.Trace.TraceEnd(ctx)

	var err error
	for _, mw := range m.requestMiddlewares {
		seed, err = mw.Process(ctx, crawler, seed)
		if err != nil {
			utils.Trace.TraceErrEvent(ctx, mw.Name(), err)
			log.Log.Error("请求中间件检查不通过", zap.String("middleware", mw.Name()), zap.Error(err))
			return nil, err
		}
		utils.Trace.TraceEvent(ctx, mw.Name())
	}
	return seed, nil
}

func (m *Middleware) ResponseProcess(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	ctx = utils.Trace.TraceStart(ctx, "Middleware.ResponseProcess")
	defer utils.Trace.TraceEnd(ctx)

	var err error
	for _, mw := range m.responseMiddlewares {
		seed, err = mw.Process(ctx, crawler, seed)
		if err != nil {
			utils.Trace.TraceErrEvent(ctx, mw.Name(), err)
			log.Log.Error("响应中间件检查不通过", zap.String("middleware", mw.Name()), zap.Error(err))
			return nil, err
		}
		utils.Trace.TraceEvent(ctx, mw.Name())
	}
	return seed, nil
}

func (m *Middleware) Close() {
	var err error
	for _, mw := range m.requestMiddlewares {
		if err = mw.Close(); err != nil {
			log.Log.Error("关闭请求中间件失败", zap.Error(err))
		}
	}
	for _, mw := range m.responseMiddlewares {
		if err = mw.Close(); err != nil {
			log.Log.Error("关闭响应中间件失败", zap.Error(err))
		}
	}
}

// NewMiddleware 创建中间件链。配置由组合根传入，默认中间件在构造时装配。
func NewMiddleware(_ zapp_core.IApp, conf *config.ServiceConfig) core.IMiddleware {
	return &Middleware{
		conf: conf,
		requestMiddlewares: []core.IRequestMiddleware{
			request_middleware.NewCheckSeedIsValidMiddleware(),
		},
		responseMiddlewares: []core.IResponseMiddleware{
			response_middleware.NewCheckSeedIsValidMiddleware(conf),
			response_middleware.NewAutoTrimUtf8Bom(conf),
		},
	}
}
