package spider_tool

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/core/dom"
	"github.com/zly-app/crawler/seeds"
)

type SpiderTool struct {
	crawler core.ICrawler
	key     string
}

func (s *SpiderTool) Crawler() core.ICrawler { return s.crawler }

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

func (s *SpiderTool) SetAdd(items ...string) int {
	count, err := s.crawler.Set().Add(s.key, items...)
	if err != nil {
		panic(err)
	}
	return count
}

func (s *SpiderTool) SetHasItem(item string) bool {
	b, err := s.crawler.Set().HasItem(s.key, item)
	if err != nil {
		panic(err)
	}
	return b
}

func (s *SpiderTool) SetRemove(items ...string) int {
	count, err := s.crawler.Set().Remove(s.key, items...)
	if err != nil {
		panic(err)
	}
	return count
}

func (s *SpiderTool) GetSetSize() int {
	size, err := s.crawler.Set().GetSetSize(s.key)
	if err != nil {
		panic(err)
	}
	return size
}

// 生成相对于在当前种子页面上的某个连接的实际连接
func (s *SpiderTool) UrlJoin(seed *core.Seed, link string) string {
	if seed.HttpResponse == nil || seed.HttpResponse.Request == nil || seed.HttpResponse.Request.URL == nil {
		panic("seed不存在页面")
	}

	u, err := seed.HttpResponse.Request.URL.Parse(link)
	if err != nil {
		panic(fmt.Errorf("UrlJoin失败: %v", err))
	}
	return u.String()
}

func (s *SpiderTool) GetDom(seed *core.Seed) *dom.Dom {
	d, err := seed.GetDom()
	if err != nil {
		panic(err)
	}
	return d
}

func (s *SpiderTool) GetXmlDom(seed *core.Seed) *dom.XmlDom {
	d, err := seed.GetXmlDom()
	if err != nil {
		panic(err)
	}
	return d
}

func NewSpiderTool(crawler core.ICrawler) core.ISpiderTool {
	return &SpiderTool{
		crawler: crawler,
		key:     config.Conf.Spider.Name + config.Conf.Frame.SetSuffix,
	}
}
