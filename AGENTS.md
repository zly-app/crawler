# AGENTS.md - crawler 框架 AI 开发参考

## 项目概述

crawler 是一个基于 [zapp](https://github.com/zly-app/zapp) v1.4.x 的分布式爬虫框架服务，为 zapp 提供服务 `crawler`。它抽象了 `seed`（种子）、`queue`（队列）、`set`（集合）、`downloader`（下载器）、`middleware`（中间件）、`pipeline`（管道）、`proxy`（代理）等模块，支持单进程/多进程/分布式部署。

- 模块路径: `github.com/zly-app/crawler`
- 依赖框架: `github.com/zly-app/zapp` v1.4.x（`go 1.25+`）
- 服务类型: `crawler`
- 插件: 推荐使用 `github.com/zly-app/plugin/otlp`（统一提供 log/trace/metrics 上报）

## 框架规范

所有组件/工具应使用框架提供的标准实现，禁止自行实现类似功能。

### 基础工具速查表

| 需求 | 工具 | 导入路径 | 说明 |
|------|------|----------|------|
| 日志 | `log` | `github.com/zly-app/zapp/log` | 全局日志, `log.Log.Info/Warn/Error/Fatal(...)`, 不要用标准库 `log` |
| 通用工具(并发/Recover/三元/Trace) | `utils` | `github.com/zly-app/zapp/pkg/utils` | `utils.Recover.WrapCall`, `utils.Go.GoAndWait`, `utils.Trace` 等 |
| 遥测上报(log/trace/metrics) | `otlp` | `github.com/zly-app/plugin/otlp` | `otlp.WithPlugin()` |
| 定时任务 | `cron` | `github.com/zly-app/service/cron` | 调度器工具使用 |
| Redis 客户端 | `redis` | `github.com/zly-app/component/redis` | 队列/集合/管道 redis 实现使用 |

### 核心 API

| 需求 | API | 说明 |
|------|-----|------|
| 启用 crawler 服务 | `crawler.WithService()` | 传入 `zapp.NewApp` 的 option |
| 注册 spider | `crawler.RegistrySpider(spider)` | 在 `app.Run()` 前调用 |
| 自定义服务 | `service.RegisterCreatorFunc(config.DefaultServiceType, creator)` | 在 `init()` 中注册 |
| 自定义队列 | `queue.RegistryQueueCreator(type, creator)` | 内置 `memory`, `redis` |
| 自定义集合 | `set.RegistrySetCreator(type, creator)` | 内置 `memory`, `redis` |
| 自定义管道 | `pipeline.RegistryPipelineCreator(type, creator)` | 内置 `stdout`, `redis_list` |
| 自定义代理 | `proxy.RegistryProxyCreator(type, creator)` | 内置 `direct`, `static` |
| 请求中间件 | `middleware.RegistryRequestMiddleware(m)` | 请求前检查 |
| 响应中间件 | `middleware.RegistryResponseMiddleware(m)` | 响应后检查 |

## Spider 开发指南

### 1. 实现 Spider

一个 spider 必须实现 `core.ISpider` 接口（`Init` / `SubmitInitialSeed` / `Close`），并且**必须继承 `core.ISpiderTool`**。

```go
package main

import (
	"context"

	"github.com/zly-app/plugin/otlp"
	"github.com/zly-app/zapp"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

type Spider struct {
	core.ISpiderTool // 必须继承这个接口
}

// 初始化
func (s *Spider) Init(ctx context.Context) error { return nil }

// 提交初始化种子
func (s *Spider) SubmitInitialSeed(ctx context.Context) error {
	seed := s.NewSeed("https://www.sogou.com/", s.Parser)
	seed.SetCheckExpectMethod(s.Check)
	s.SubmitSeed(ctx, seed)
	return nil
}

// 解析方法, 必须以 Parser 开头
func (s *Spider) Parser(ctx context.Context, seed *core.Seed) error {
	body := string(seed.HttpResponseBody)
	s.SaveResult(ctx, body)
	return nil
}

// 检查方法, 必须以 Check 开头
func (s *Spider) Check(ctx context.Context, seed *core.Seed) error { return nil }

// 关闭
func (s *Spider) Close(ctx context.Context) error { return nil }

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService(), otlp.WithPlugin())
	crawler.RegistrySpider(new(Spider))
	app.Run()
}
```

### 2. 方法命名规则

- **解析方法**: 方法名必须以 `Parser` 前缀开头, 如 `Parser`, `ParserList`, `ParserDetail`, 框架会自动扫描注册
- **检查方法**: 方法名必须以 `Check` 前缀开头, 如 `Check`, 用于检查响应是否符合期望

### 3. SpiderTool 可用方法

| 方法 | 说明 |
|------|------|
| `NewSeed(url, parserMethod)` | 创建种子, parserMethod 可以是方法名或方法实体 |
| `SubmitSeed(ctx, seed)` | 提交种子到队列 |
| `SaveResult(ctx, data)` | 保存抓取结果到管道 |
| `PutSeed(ctx, seed, front)` | 放入种子, front 表示是否放队列前面 |
| `PutRawSeed(ctx, raw, parserFuncName, front)` | 放入种子原始数据 |
| `PutErrorSeed(ctx, seed, isParserError)` | 放入错误种子 |
| `SetAdd / SetHasItem / SetRemove / GetSetSize` | 集合去重操作 |
| `UrlJoin(seed, link)` | 生成相对当前页面的绝对连接 |
| `GetDom / GetXmlDom / GetJsonDom` | 获取 DOM 解析器 |

### 4. Seed 结构

```go
type Seed struct {
	Raw               string           // 原始数据
	HttpRequest       *http.Request    // 请求体
	HttpResponse      *http.Response   // 响应
	HttpResponseBody  []byte           // 响应数据
	HttpCookies       []*http.Cookie   // 最终cookies
	Request           struct {         // 请求参数
		Method, Url, Body string
		Params, Form, PostForm url.Values
		Headers, Trailer http.Header
		UserAgentType string
		AllowProxy, AutoCookie, AutoRedirects bool
		Cookies, ParentCookies []*http.Cookie
		Encoding string // utf8/gbk/gb2312/gb18030, 自动转utf8
		Timeout   int64  // 毫秒
	}
	ParserMethod      string
	CheckExpectMethod string
	Meta              map[string]interface{} // 元数据, 跨请求传递数据
}
```

### 5. DOM 解析

- HTML: `s.GetDom(seed)` 支持 `Xpath`/`XpathOne`/`Css`/`CssOne`/`GetAttr`/`InnerText`/`OutputHTML` 等
- XML: `s.GetXmlDom(seed)` 支持 xpath
- JSON: `s.GetJsonDom(seed)` 支持 xpath, `FindOneChild`, `InnerText`
- 也支持 `seed.GetGoJQ()` 使用 gojq 查询 JSON

### 6. Meta 元数据

`seed` 无状态, 需要跨请求传递数据时使用 `Meta`:

```go
// 在 A 页面解析时创建 B 页面的种子并带上数据
newSeed := s.NewSeed(bUrl, s.ParserB)
newSeed.Meta["itemId"] = itemId
s.SubmitSeed(ctx, newSeed)

// 在 B 页面解析时读取
func (s *Spider) ParserB(ctx context.Context, seed *core.Seed) error {
	itemId := seed.Meta["itemId"]
	// ...
}
```

## 配置说明

配置通过 zapp 加载, key 为 `services.crawler.*`。默认无需任何配置文件即可运行（memory 队列 + stdout 管道）。

```yaml
services:
  crawler:
    spider:     # spider配置
      Name: 'a_spider'
      SubmitInitialSeedOpportunity: 'start' # none/start/指定时间/cron表达式
      AutoCookie: false
      AutoRedirects: true
      HtmlEncoding: 'utf8'
      UserAgentType: 'pc'   # pc/android/ios
      ExpectHttpStatusCode: []
      InvalidHttpStatusCode: []
      HttpStatus4xxIsInvalid: true
      HttpStatus5xxIsInvalid: true

    frame:      # 框架配置
      Namespace: 'my_project'
      RequestTimeout: 20000        # 毫秒
      SeedProcessTimeout: 60000
      RequestMaxAttemptCount: 5
      RequestRetryWaitTime: 1000
      EmptyQueueWaitTime: 60000
      NextSeedWaitTime: 0
      SpiderErrWaitTime: 3000
      StopSubmitInitialSeedIfNotEmptyQueue: true

    queue:      # 队列, 默认 memory
      type: 'redis'                # memory/redis
      Address: localhost:6379
      UserName: ""
      Password: ""
      DB: 0

    set:        # 集合(去重), 默认 memory
      type: 'redis'                # memory/redis
      Address: localhost:6379
      UserName: ""
      Password: ""
      DB: 0

    pipeline:   # 管道(结果输出), 默认 stdout, 多管道逗号分隔
      type: 'stdout'               # stdout/redis_list
      redis_list:
        SaveToQueueFront: false
        QueuePrefix: "result:"
        Serializer: "sonic_std"
        Compactor: "raw"

    proxy:      # 代理, 默认 direct
      type: 'static'               # direct/static
      address: 'socks5://127.0.0.1:1080'
      User: ''
      Password: ''
```

完整配置参考: `tools/crawler/embed_assets/configs/crawler.dev.yaml.template`

## 工程管理工具

`tools/crawler` 是工程管理工具, 用于快速搭建和管理项目:

| 命令 | 说明 |
|------|------|
| `init <name>` | 初始化一个项目 |
| `create <spider>` (alias `cs`) | 在 spiders 目录创建一个 spider |
| `start [-env <env>] <spider>` (alias `ss`) | 提交初始化种子信号到队列 |
| `clean [-env <env>] <spider>` | 清空 spider 所有队列 |
| `clean_set [-env <env>] <spider>` | 清空 spider 集合数据 |
| `make_supervisor [-env <env>]` (alias `make`) | 生成 supervisor 配置 |

`tools/crawler_scheduler` 是调度器工具, 为分布式 spider 提供统一的初始化种子信号管理。

## 项目结构

```
crawler/
├── crawler.go              # crawler服务实现
├── crawler.run.go          # 抓取主循环
├── crawler.initial_seed.go # 初始化种子提交
├── crawler.seed.go         # 种子处理
├── crawler.parser_method.go# 解析方法扫描
├── usage.go                # 服务注册与公开API
├── config/                 # 配置结构
├── core/                   # 核心接口与Seed定义
│   └── dom/                # html/xml/json DOM解析器
├── seeds/                  # seed创建/编码
├── queue/                  # 队列接口与实现(memory/redis)
├── set/                    # 集合接口与实现(memory/redis)
├── pipeline/               # 管道接口与实现(stdout/redis_list)
├── proxy/                  # 代理接口与实现(direct/static)
├── middleware/             # 请求/响应中间件
├── downloader/             # 下载器
├── spider_tool/            # ISpiderTool 实现
├── cookiejar/              # cookie管理
├── utils/                  # 框架内部工具
├── tools/                  # 工程管理工具与调度器工具
│   ├── crawler/
│   │   └── embed_assets/   # 项目模板(init/create使用)
│   └── crawler_scheduler/
└── example/a_spider/       # 示例
```

## 禁止事项

- 禁止自行实现日志功能, 应使用 `github.com/zly-app/zapp/log`
- 禁止使用标准库 `log` 包
- 禁止直接使用 `zap.NewProduction()` 等自行创建 logger
- 禁止自行实现日志轮转、日志采集、链路追踪等功能
- 禁止使用未适配 zapp v1.4.x 的旧插件(honey/zipkinotel/prometheus), 统一使用 `otlp`
- 禁止在 `tools/crawler/embed_assets` 模板中引入与主模块 go.mod 不一致的依赖版本

## 版本兼容说明

- zapp v1.4.x 移除了 `zapp/logger` 包(改为 `zapp/log`), 移除了 `pkg/utils` 中的 `utils.Otel`(改为 `utils.Trace`)
- 所有 zly-app 生态包必须使用适配 zapp v1.4.x 的版本
- 升级 zapp 版本时, 需同步更新 `go.mod.tfile` 模板中的 zapp/otlp/crawler 版本
