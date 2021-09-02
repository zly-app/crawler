package middleware

import (
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/middleware/request_middleware"
	"github.com/zly-app/crawler/middleware/response_middleware"
)

type Middleware struct {
	app                 zapp_core.IApp
	requestMiddlewares  []core.IRequestMiddleware
	responseMiddlewares []core.IResponseMiddleware
}

func (m *Middleware) RequestProcess(crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	var err error
	for _, middleware := range m.requestMiddlewares {
		seed, err = middleware.Process(crawler, seed)
		if err != nil {
			m.app.Error("请求中间件检查不通过", zap.String("middleware", middleware.Name()), zap.Error(err))
			return nil, err
		}
	}
	return seed, nil
}

func (m *Middleware) ResponseProcess(crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	var err error
	for _, middleware := range m.responseMiddlewares {
		seed, err = middleware.Process(crawler, seed)
		if err != nil {
			m.app.Error("响应中间件检查不通过", zap.String("middleware", middleware.Name()), zap.Error(err))
			return nil, err
		}
	}
	return seed, nil
}

func (m *Middleware) Close() error {
	return nil
}

func NewMiddleware(app zapp_core.IApp) core.IMiddleware {
	return &Middleware{
		app: app,
		requestMiddlewares: []core.IRequestMiddleware{
			request_middleware.NewCheckSeedIsValidMiddleware(),
		},
		responseMiddlewares: []core.IResponseMiddleware{
			response_middleware.NewCheckSeedIsValidMiddleware(),
		},
	}
}
