package crawler

import (
	"fmt"
	"net/http"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/downloader"
	"github.com/zly-app/crawler/middleware"
	"github.com/zly-app/crawler/proxy"
	"github.com/zly-app/crawler/queue"
	"github.com/zly-app/crawler/set"
	"github.com/zly-app/crawler/spider_tool"
)

type Crawler struct {
	app           zapp_core.IApp
	conf          *config.ServiceConfig
	parserMethods map[string]core.ParserMethod

	spider     core.ISpider
	spiderTool core.ISpiderTool
	queue      core.IQueue
	set        core.ISet
	downloader core.IDownloader
	proxy      core.IProxy
	middleware core.IMiddleware

	cookieJar http.CookieJar // 当前使用的cookieJar
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

	c.CheckSpiderParserMethod()
}

func (c *Crawler) Start() error {
	err := c.spider.Init(c.spiderTool)
	if err != nil {
		return fmt.Errorf("spider初始化失败: %v", err)
	}

	go c.Run()
	return nil
}

func (c *Crawler) Close() error {
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
	if err = c.middleware.Close(); err != nil {
		c.app.Error("关闭中间件时出错", zap.Error(err))
	}
	return nil
}

func (c *Crawler) Spider() core.ISpider         { return c.spider }
func (c *Crawler) Queue() core.IQueue           { return c.queue }
func (c *Crawler) Downloader() core.IDownloader { return c.downloader }
func (c *Crawler) Proxy() core.IProxy           { return c.proxy }
func (c *Crawler) Set() core.ISet               { return c.set }
func (c *Crawler) CookieJar() http.CookieJar    { return c.cookieJar }

func NewCrawler(app zapp_core.IApp) zapp_core.IService {
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
	}
	crawler.spiderTool = spider_tool.NewSpiderTool(crawler)
	return crawler
}
