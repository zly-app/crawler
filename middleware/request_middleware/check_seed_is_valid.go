package request_middleware

import (
	"fmt"

	"github.com/zly-app/crawler/core"
)

// 检查种子是否合法
type CheckSeedIsValid struct {
	core.MiddlewareBase
}

func NewCheckSeedIsValidMiddleware() core.IRequestMiddleware {
	return new(CheckSeedIsValid)
}

func (m *CheckSeedIsValid) Name() string { return "CheckSeedIsValid" }
func (m *CheckSeedIsValid) Process(crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	// 检查预期响应是可选的
	if seed.CheckExpectMethod != "" {
		_, ok := crawler.GetSpiderCheckMethod(seed.CheckExpectMethod)
		if !ok {
			return nil, fmt.Errorf("未找到检查预期响应方法: %s", seed.CheckExpectMethod)
		}
	}

	if seed.ParserMethod == "" {
		return nil, fmt.Errorf("种子的解析方法是空的: raw: %s", seed.Raw)
	}
	_, ok := crawler.GetSpiderParserMethod(seed.ParserMethod)
	if !ok {
		return nil, fmt.Errorf("未找到解析方法: %s", seed.ParserMethod)
	}
	return seed, nil
}
