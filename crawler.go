package crawler

import (
	"fmt"
	"net/http"
	"reflect"
	"sync/atomic"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/downloader"
	"github.com/zly-app/crawler/middleware"
	"github.com/zly-app/crawler/pipeline"
	"github.com/zly-app/crawler/proxy"
	"github.com/zly-app/crawler/queue"
	"github.com/zly-app/crawler/set"
	"github.com/zly-app/crawler/spider_tool"
)

var typeOfISpiderTool = reflect.TypeOf((*core.ISpiderTool)(nil)).Elem()

type Crawler struct {
	app           zapp_core.IApp
	conf          *config.ServiceConfig
	parserMethods map[string]core.ParserMethod // 解析方法
	checkMethods  map[string]core.CheckMethod  // 检查方法

	spider     core.ISpider
	spiderTool core.ISpiderTool
	queue      core.IQueue
	set        core.ISet
	downloader core.IDownloader
	proxy      core.IProxy
	middleware core.IMiddleware
	pipeline   core.IPipeline

	cookieJar  http.CookieJar // 当前使用的cookieJar
	nowRawSeed atomic.Value   // 当前的原始种子数据, 用于在程序退出之前回退到队列中
}

func (c *Crawler) Inject(a ...interface{}) {
	if c.spider != nil {
		c.app.Fatal("crawler服务重复注入")
	}

	if len(a) != 1 {
		c.app.Fatal("crawler服务注入数量必须为1个")
	}

	var ok bool
	c.spider, ok = a[0].(core.ISpider)
	if !ok {
		c.app.Fatal("crawler服务注入类型错误, 它必须能转为 crawler/core.ISpider")
	}

	// 检查继承于core.ISpiderTool
	_, ok = a[0].(core.ISpiderTool)
	if !ok {
		c.app.Fatal("crawler服务注入类型错误, 它必须继承 crawler/core.ISpiderTool")
	}
	// 检查是带指针的结构体
	aType := reflect.TypeOf(a[0])
	if aType.Kind() != reflect.Ptr {
		c.app.Fatal("crawler服务注入类型错误, 它必须是带指针的结构体")
	}
	aType = aType.Elem()
	if aType.Kind() != reflect.Struct {
		c.app.Fatal("crawler服务注入类型错误, 它必须是带指针的结构体")
	}
	// 检查ISpiderTool字段
	fieldT, ok := aType.FieldByName("ISpiderTool")
	if !ok {
		c.app.Fatal("crawler服务注入类型错误, 它必须继承 crawler/core.ISpiderTool")
	}
	if !fieldT.Type.AssignableTo(typeOfISpiderTool) {
		c.app.Fatal("crawler服务注入类型错误, 它必须继承 crawler/core.ISpiderTool")
	}
	// 注入
	field := reflect.ValueOf(a[0]).Elem().FieldByName("ISpiderTool")
	field.Set(reflect.ValueOf(c.spiderTool))

	c.ScanSpiderMethod()
}

func (c *Crawler) Start() error {
	err := c.spider.Init()
	if err != nil {
		return fmt.Errorf("spider初始化失败: %v", err)
	}

	go c.Run()
	return nil
}

func (c *Crawler) Close() error {
	c.rollbackRawSeed() // 回退原始种子数据

	err := c.spider.Close()
	if err != nil {
		c.app.Error("spider关闭时出错", zap.Error(err))
	}

	if err = c.queue.Close(); err != nil {
		c.app.Error("关闭队列时出错误", zap.Error(err))
	}
	if err = c.set.Close(); err != nil {
		c.app.Error("关闭集合时出错", zap.Error(err))
	}
	if err = c.downloader.Close(); err != nil {
		c.app.Error("关闭下载器时出错", zap.Error(err))
	}
	if err = c.proxy.Close(); err != nil {
		c.app.Error("关闭代理时出错", zap.Error(err))
	}
	c.middleware.Close()
	return nil
}

// 回退原始种子数据
func (c *Crawler) rollbackRawSeed() {
	rawData := c.nowRawSeed.Load()
	if rawData == nil {
		return
	}

	raw := rawData.(string)
	if raw == "" {
		return
	}

	err := c.PutRawSeed(raw, "", true)
	if err != nil {
		c.app.Error("回退原始种子数据失败", zap.String("raw", raw), zap.Error(err))
	}
}

func (c *Crawler) Spider() core.ISpider         { return c.spider }
func (c *Crawler) SpiderTool() core.ISpiderTool { return c.spiderTool }
func (c *Crawler) Queue() core.IQueue           { return c.queue }
func (c *Crawler) Pipeline() core.IPipeline     { return c.pipeline }
func (c *Crawler) Downloader() core.IDownloader { return c.downloader }
func (c *Crawler) Proxy() core.IProxy           { return c.proxy }
func (c *Crawler) Set() core.ISet               { return c.set }
func (c *Crawler) CookieJar() http.CookieJar    { return c.cookieJar }

func NewCrawler(app zapp_core.IApp) *Crawler {
	conf := config.NewConfig(app)
	confKey := "services." + string(config.NowServiceType)
	if app.GetConfig().GetViper().IsSet(confKey) {
		err := app.GetConfig().ParseServiceConfig(config.NowServiceType, conf)
		if err != nil {
			logger.Log.Panic("服务配置错误", zap.String("serviceType", string(config.NowServiceType)), zap.Error(err))
		}
	}
	err := conf.Check()
	if err != nil {
		logger.Log.Panic("服务配置错误", zap.String("serviceType", string(config.NowServiceType)), zap.Error(err))
	}

	crawler := &Crawler{
		app:        app,
		conf:       conf,
		queue:      queue.NewQueue(app, conf.Queue.Type),
		set:        set.NewSet(app, conf.Set.Type),
		downloader: downloader.NewDownloader(app),
		proxy:      proxy.NewProxy(app, conf.Proxy.Type),
		middleware: middleware.NewMiddleware(app),
		pipeline:   pipeline.NewPipeline(app, conf.Pipeline.Type),
	}
	crawler.spiderTool = spider_tool.NewSpiderTool(crawler)
	return crawler
}
