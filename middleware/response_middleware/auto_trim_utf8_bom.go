package response_middleware

import (
	"context"
	"strings"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

type AutoTrimUtf8Bom struct {
	core.MiddlewareBase
}

func NewAutoTrimUtf8Bom() core.IRequestMiddleware {
	return AutoTrimUtf8Bom{}
}

func (a AutoTrimUtf8Bom) Name() string { return "AutoTrimUtf8Bom" }

func (a AutoTrimUtf8Bom) Process(ctx context.Context, crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	if seed.HttpResponse == nil || config.Conf.Spider.AutoTrimUtf8BomWithContentType == "" {
		return seed, nil
	}

	t := seed.HttpResponse.Header.Get("content-type")
	if t == "" {
		return seed, nil
	}
	ts := strings.Split(t, ";")

	expectTs := strings.Split(config.Conf.Spider.AutoTrimUtf8BomWithContentType, ",")
	checkIsExpect := func(t string) bool {
		t = strings.TrimSpace(t)
		for i := range expectTs {
			if expectTs[i] == t {
				return true
			}
		}
		return false
	}

	for _, t := range ts {
		if checkIsExpect(t) {
			seed.TryTrimUtf8Bom()
			return seed, nil
		}
	}

	return seed, nil
}
