package spider_tool

import (
	"net/http"
	"strings"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/seeds"
)

type SpiderTool struct {
	crawler core.ICrawler
}

func (s *SpiderTool) NewSeed(url string, parserMethod interface{}) *core.Seed {
	seed := seeds.NewSeed()
	seed.Request.Url = url
	if seed.Request.AutoCookie && s.crawler.CookieJar() != nil && url != "" {
		req, err := http.NewRequest(strings.ToUpper(seed.Request.Method), url, nil)
		if err == nil {
			cookies := s.crawler.CookieJar().Cookies(req.URL) // 获取这个种子要用到的cookies
			seed.Request.ParentCookies = cookies
		}
	}
	seed.SetParserMethod(parserMethod)
	return seed
}

func (s *SpiderTool) SubmitSeed(seed *core.Seed) {
	s.PutSeed(seed, config.Conf.Frame.SubmitSeedToQueueFront)
}

func (s *SpiderTool) PutSeed(seed *core.Seed, front bool) {
	if err := s.crawler.PutSeed(seed, front); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutRawSeed(raw string, parserFuncName string, front bool) {
	if err := s.crawler.PutRawSeed(raw, parserFuncName, front); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutErrorSeed(seed *core.Seed, isParserError bool) {
	if err := s.crawler.PutErrorSeed(seed, isParserError); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutErrorRawSeed(raw string, isParserError bool) {
	if err := s.crawler.PutErrorRawSeed(raw, isParserError); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) SetAdd(key string, items ...string) int {
	panic("implement me")
}

func (s *SpiderTool) SetHasItem(key, item string) bool {
	panic("implement me")
}

func (s *SpiderTool) SetRemove(key string, items ...string) int {
	panic("implement me")
}

func (s *SpiderTool) GetSetSize(key string) int {
	panic("implement me")
}

func NewSpiderTool(crawler core.ICrawler) core.ISpiderTool {
	return &SpiderTool{crawler: crawler}
}
