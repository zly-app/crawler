package spider_tool

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/core/dom"
	"github.com/zly-app/crawler/seeds"
	"github.com/zly-app/crawler/utils"
)

type SpiderTool struct {
	crawler    core.ICrawler
	spiderName string
	setKey     string
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

func (s *SpiderTool) SubmitSeed(ctx context.Context, seed *core.Seed) {
	ctx = utils.Trace.TraceStart(ctx, "SubmitSeed")
	defer utils.Trace.TraceEnd(ctx)

	s.PutSeed(ctx, seed, config.Conf.Frame.SubmitSeedToQueueFront)
}

func (s *SpiderTool) SaveResult(ctx context.Context, data interface{}) {
	ctx = utils.Trace.TraceStart(ctx, "SaveResult")
	defer utils.Trace.TraceEnd(ctx)

	dataText, _ := jsoniter.MarshalToString(data)
	utils.Trace.TraceEvent(ctx, "save", utils.Trace.AttrKey("data").String(dataText))
	err := s.crawler.Pipeline().Process(ctx, s.spiderName, data)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "save", err)
		panic(fmt.Errorf("保存结果失败: %v", err))
	}
}

func (s *SpiderTool) PutSeed(ctx context.Context, seed *core.Seed, front bool) {
	if err := s.crawler.PutSeed(ctx, seed, front); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutRawSeed(ctx context.Context, raw string, parserFuncName string, front bool) {
	if err := s.crawler.PutRawSeed(ctx, raw, parserFuncName, front); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutErrorSeed(ctx context.Context, seed *core.Seed, isParserError bool) {
	if err := s.crawler.PutErrorSeed(ctx, seed, isParserError); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) PutErrorRawSeed(ctx context.Context, raw string, isParserError bool) {
	if err := s.crawler.PutErrorRawSeed(ctx, raw, isParserError); err != nil {
		panic(err)
	}
}

func (s *SpiderTool) SetAdd(ctx context.Context, items ...string) int {
	ctx = utils.Trace.TraceStart(ctx, "SetAdd")
	defer utils.Trace.TraceEnd(ctx)

	utils.Trace.TraceEvent(ctx, "SetAdd", utils.Trace.AttrKey("items").StringSlice(items))
	count, err := s.crawler.Set().Add(ctx, s.setKey, items...)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "SetAdd", err)
		panic(err)
	}
	utils.Trace.TraceEvent(ctx, "SetAddOk", utils.Trace.AttrKey("count").Int(count))
	return count
}

func (s *SpiderTool) SetHasItem(ctx context.Context, item string) bool {
	ctx = utils.Trace.TraceStart(ctx, "SetHasItem")
	defer utils.Trace.TraceEnd(ctx)

	utils.Trace.TraceEvent(ctx, "SetHasItem", utils.Trace.AttrKey("item").String(item))
	b, err := s.crawler.Set().HasItem(ctx, s.setKey, item)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "SetHasItem", err)
		panic(err)
	}
	utils.Trace.TraceEvent(ctx, "SetHasItemOk", utils.Trace.AttrKey("result").Bool(b))
	return b
}

func (s *SpiderTool) SetRemove(ctx context.Context, items ...string) int {
	ctx = utils.Trace.TraceStart(ctx, "SetRemove")
	defer utils.Trace.TraceEnd(ctx)

	utils.Trace.TraceEvent(ctx, "SetRemove", utils.Trace.AttrKey("items").StringSlice(items))
	count, err := s.crawler.Set().Remove(ctx, s.setKey, items...)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "SetRemove", err)
		panic(err)
	}
	utils.Trace.TraceEvent(ctx, "SetRemoveOk", utils.Trace.AttrKey("count").Int(count))
	return count
}

func (s *SpiderTool) GetSetSize(ctx context.Context) int {
	ctx = utils.Trace.TraceStart(ctx, "GetSetSize")
	defer utils.Trace.TraceEnd(ctx)

	utils.Trace.TraceEvent(ctx, "GetSetSize")
	size, err := s.crawler.Set().GetSetSize(ctx, s.setKey)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "GetSetSize", err)
		panic(err)
	}
	utils.Trace.TraceEvent(ctx, "GetSetSize", utils.Trace.AttrKey("size").Int(size))
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

func (s *SpiderTool) GetJsonDom(seed *core.Seed) *dom.JsonDom {
	d, err := seed.GetJsonDom()
	if err != nil {
		panic(err)
	}
	return d
}

func NewSpiderTool(crawler core.ICrawler) core.ISpiderTool {
	return &SpiderTool{
		crawler:    crawler,
		spiderName: config.Conf.Spider.Name,
		setKey:     config.Conf.Frame.Namespace + config.Conf.Spider.Name + config.Conf.Frame.SetSuffix,
	}
}
