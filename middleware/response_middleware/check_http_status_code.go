package response_middleware

import (
	"context"
	"fmt"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type CheckHttpStatusCode struct {
	core.MiddlewareBase
}

func NewCheckSeedIsValidMiddleware() core.IRequestMiddleware {
	return new(CheckHttpStatusCode)
}

func (m *CheckHttpStatusCode) Name() string { return "CheckHttpStatusCode" }
func (m *CheckHttpStatusCode) Process(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	if seed.HttpResponse == nil {
		return seed, nil
	}

	// 检查期望值
	if len(config.Conf.Spider.ExpectHttpStatusCode) > 0 {
		for _, expect := range config.Conf.Spider.ExpectHttpStatusCode {
			if expect == seed.HttpResponse.StatusCode {
				return seed, nil
			}
		}
		return nil, fmt.Errorf("收到非期望的http状态码: %d", seed.HttpResponse.StatusCode)
	}

	// 检查4xx值
	if config.Conf.Spider.HttpStatus4xxIsInvalid && seed.HttpResponse.StatusCode >= 400 && seed.HttpResponse.StatusCode < 500 {
		return nil, fmt.Errorf("收到非期望的http状态码: %d", seed.HttpResponse.StatusCode)
	}

	// 检查5xx值
	if config.Conf.Spider.HttpStatus5xxIsInvalid && seed.HttpResponse.StatusCode >= 500 && seed.HttpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("收到非期望的http状态码: %d", seed.HttpResponse.StatusCode)
	}

	// 检查排除值
	for _, exclude := range config.Conf.Spider.InvalidHttpStatusCode {
		if exclude == seed.HttpResponse.StatusCode {
			return nil, fmt.Errorf("收到无效的http状态码: %d", seed.HttpResponse.StatusCode)
		}
	}

	return seed, nil
}
