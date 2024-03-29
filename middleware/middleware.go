package middleware

import (
	"context"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/middleware/request_middleware"
	"github.com/zly-app/crawler/middleware/response_middleware"
	"github.com/zly-app/crawler/utils"
)

// 请求中间件
var requestMiddlewares = []core.IRequestMiddleware{
	request_middleware.NewCheckSeedIsValidMiddleware(),
}

// 响应中间件
var responseMiddlewares = []core.IResponseMiddleware{
	response_middleware.NewCheckSeedIsValidMiddleware(),
	response_middleware.NewAutoTrimUtf8Bom(),
}

// 注册请求中间件
func RegistryRequestMiddleware(m core.IRequestMiddleware) {
	requestMiddlewares = append(requestMiddlewares, m)
}

// 注册响应中间件
func RegistryResponseMiddleware(m core.IResponseMiddleware) {
	responseMiddlewares = append(responseMiddlewares, m)
}

type Middleware struct{}

func (m *Middleware) RequestProcess(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	ctx = utils.Trace.TraceStart(ctx, "Middleware.RequestProcess")
	defer utils.Trace.TraceEnd(ctx)

	var err error
	for _, middleware := range requestMiddlewares {
		seed, err = middleware.Process(ctx, crawler, seed)
		if err != nil {
			utils.Trace.TraceErrEvent(ctx, middleware.Name(), err)
			logger.Log.Error("请求中间件检查不通过", zap.String("middleware", middleware.Name()), zap.Error(err))
			return nil, err
		}
		utils.Trace.TraceEvent(ctx, middleware.Name())
	}
	return seed, nil
}

func (m *Middleware) ResponseProcess(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	ctx = utils.Trace.TraceStart(ctx, "Middleware.ResponseProcess")
	defer utils.Trace.TraceEnd(ctx)

	var err error
	for _, middleware := range responseMiddlewares {
		seed, err = middleware.Process(ctx, crawler, seed)
		if err != nil {
			utils.Trace.TraceErrEvent(ctx, middleware.Name(), err)
			logger.Log.Error("响应中间件检查不通过", zap.String("middleware", middleware.Name()), zap.Error(err))
			return nil, err
		}
		utils.Trace.TraceEvent(ctx, middleware.Name())
	}
	return seed, nil
}

func (m *Middleware) Close() {
	var err error
	for _, middleware := range requestMiddlewares {
		if err = middleware.Close(); err != nil {
			logger.Log.Error("关闭请求中间件失败", zap.Error(err))
		}
	}
	for _, middleware := range responseMiddlewares {
		if err = middleware.Close(); err != nil {
			logger.Log.Error("关闭响应中间件失败", zap.Error(err))
		}
	}
}

func NewMiddleware(app zapp_core.IApp) core.IMiddleware {
	return &Middleware{}
}
