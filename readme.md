# 分布式爬虫框架服务

> 提供用于 https://github.com/zly-app/zapp 的服务


<!-- TOC -->

- [分布式爬虫框架服务](#%E5%88%86%E5%B8%83%E5%BC%8F%E7%88%AC%E8%99%AB%E6%A1%86%E6%9E%B6%E6%9C%8D%E5%8A%A1)
- [说明](#%E8%AF%B4%E6%98%8E)
    - [示例](#%E7%A4%BA%E4%BE%8B)
- [配置](#%E9%85%8D%E7%BD%AE)
    - [配置参考](#%E9%85%8D%E7%BD%AE%E5%8F%82%E8%80%83)
    - [持久化队列](#%E6%8C%81%E4%B9%85%E5%8C%96%E9%98%9F%E5%88%97)
        - [使用redis作为队列](#%E4%BD%BF%E7%94%A8redis%E4%BD%9C%E4%B8%BA%E9%98%9F%E5%88%97)
        - [使用ssdb作为队列](#%E4%BD%BF%E7%94%A8ssdb%E4%BD%9C%E4%B8%BA%E9%98%9F%E5%88%97)
    - [使用持久化集合](#%E4%BD%BF%E7%94%A8%E6%8C%81%E4%B9%85%E5%8C%96%E9%9B%86%E5%90%88)
        - [使用redis作为集合](#%E4%BD%BF%E7%94%A8redis%E4%BD%9C%E4%B8%BA%E9%9B%86%E5%90%88)
        - [使用ssdb作为集合](#%E4%BD%BF%E7%94%A8ssdb%E4%BD%9C%E4%B8%BA%E9%9B%86%E5%90%88)
    - [使用代理](#%E4%BD%BF%E7%94%A8%E4%BB%A3%E7%90%86)
- [概念](#%E6%A6%82%E5%BF%B5)
    - [种子 seed](#%E7%A7%8D%E5%AD%90-seed)
    - [队列 queue](#%E9%98%9F%E5%88%97-queue)
    - [下载器 downloader](#%E4%B8%8B%E8%BD%BD%E5%99%A8-downloader)
    - [中间件 downloader](#%E4%B8%AD%E9%97%B4%E4%BB%B6-downloader)
- [设计思路](#%E8%AE%BE%E8%AE%A1%E6%80%9D%E8%B7%AF)
    - [进程独立](#%E8%BF%9B%E7%A8%8B%E7%8B%AC%E7%AB%8B)
    - [请求独立](#%E8%AF%B7%E6%B1%82%E7%8B%AC%E7%AB%8B)
    - [抓取过程](#%E6%8A%93%E5%8F%96%E8%BF%87%E7%A8%8B)
    - [配置化](#%E9%85%8D%E7%BD%AE%E5%8C%96)
    - [模块化](#%E6%A8%A1%E5%9D%97%E5%8C%96)
    - [进程管理](#%E8%BF%9B%E7%A8%8B%E7%AE%A1%E7%90%86)
- [一些操作](#%E4%B8%80%E4%BA%9B%E6%93%8D%E4%BD%9C)
    - [这条数据我已经抓过了, 怎么让爬虫不再抓它了](#%E8%BF%99%E6%9D%A1%E6%95%B0%E6%8D%AE%E6%88%91%E5%B7%B2%E7%BB%8F%E6%8A%93%E8%BF%87%E4%BA%86-%E6%80%8E%E4%B9%88%E8%AE%A9%E7%88%AC%E8%99%AB%E4%B8%8D%E5%86%8D%E6%8A%93%E5%AE%83%E4%BA%86)
    - [非请求的seed](#%E9%9D%9E%E8%AF%B7%E6%B1%82%E7%9A%84seed)
- [工程管理工具](#%E5%B7%A5%E7%A8%8B%E7%AE%A1%E7%90%86%E5%B7%A5%E5%85%B7)
    - [命令](#%E5%91%BD%E4%BB%A4)
- [调度器工具](#%E8%B0%83%E5%BA%A6%E5%99%A8%E5%B7%A5%E5%85%B7)

<!-- /TOC -->

# 说明

```text
crawler.WithService()           # 启用服务
crawler.RegistrySpider(...)     # 服务注入spider
```

## 示例

```go
package main

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

// 一个spider
type Spider struct {
	core.ISpiderTool // 必须继承这个接口
}

// 初始化
func (s *Spider) Init() error { return nil }

// 提交初始化种子
func (s *Spider) SubmitInitialSeed() error {
	seed := s.NewSeed("https://www.baidu.com/", s.Parser) // 创建种子并指定解析方法
	s.SubmitSeed(seed)                                    // 提交种子
	return nil
}

// 解析方法, 必须以 Parser 开头
func (s *Spider) Parser(seed *core.Seed) error {
	data := string(seed.HttpResponseBody) // 获取响应body
	s.SaveResult(data)                    // 保存结果
	return nil
}

// 关闭
func (s *Spider) Close() error { return nil }

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService()) // 启用crawler服务
	crawler.RegistrySpider(new(Spider))                   // 注入spider
	app.Run()                                             // 运行
}
```

# 配置

+ 不需要任何配置文件就可以运行
+ 默认服务类型为 `crawler`, 完整配置说明参考 [Config](./config)

## 配置参考

```toml
# 爬虫配置
[services.crawler.spider]
# 爬虫名
Name = 'a_spider'
# 提交初始化种子的时机
SubmitInitialSeedOpportunity = 'start'
# 是否自动管理cookie
AutoCookie = false

# 框架配置
[services.crawler.frame]
# 请求超时, 毫秒
RequestTimeout = 20000
# 最大尝试请求次数
RequestMaxAttemptCount = 5
```

## 持久化队列

### 使用redis作为队列

```toml
[services.crawler.queue]
type = 'redis'      # 使用redis作为队列, 默认是memory
Address = '127.0.0.1:6379' # 地址
UserName = ''       # 用户名, 可选
Password = ''       # 密码, 可选
DB = 0              # db, 只有非集群有效, 可选, 默认0
IsCluster = false   # 是否为集群, 可选, 默认false
MinIdleConns = 1    # 最小空闲连接数, 可选, 默认1
PoolSize = 1        # 客户端池大小, 可选, 默认1
ReadTimeout = 5000  # 读取超时(毫秒, 可选, 默认5000
WriteTimeout = 5000 # 写入超时(毫秒, 可选, 默认5000
DialTimeout = 5000  # 连接超时(毫秒, 可选, 默认5000
```

### 使用ssdb作为队列

```toml
[services.crawler.queue]
type = 'ssdb'       # 使用ssdb作为队列, 默认是memory
Address = '127.0.0.1:8888' # 地址
Password = ''       # 密码, 可选
MinIdleConns = 1    # 最小空闲连接数, 可选, 默认1
PoolSize = 1        # 客户端池大小, 可选, 默认1
ReadTimeout = 5000  # 读取超时(毫秒, 可选, 默认5000
WriteTimeout = 5000 # 写入超时(毫秒, 可选, 默认5000
DialTimeout = 5000  # 连接超时(毫秒, 可选, 默认5000
```

## 使用持久化集合

> 没错, 和队列的配置内容相似

### 使用redis作为集合

```toml
[services.crawler.set]
type = 'redis'      # 使用redis作为集合, 默认是memory
Address = '127.0.0.1:6379' # 地址
UserName = ''       # 用户名, 可选
Password = ''       # 密码, 可选
DB = 0              # db, 只有非集群有效, 可选, 默认0
IsCluster = false   # 是否为集群, 可选, 默认false
MinIdleConns = 1    # 最小空闲连接数, 可选, 默认1
PoolSize = 1        # 客户端池大小, 可选, 默认1
ReadTimeout = 5000  # 读取超时(毫秒, 可选, 默认5000
WriteTimeout = 5000 # 写入超时(毫秒, 可选, 默认5000
DialTimeout = 5000  # 连接超时(毫秒, 可选, 默认5000
```

### 使用ssdb作为集合

```toml
[services.crawler.set]
type = 'ssdb'       # 使用ssdb作为集合, 默认是memory
Address = '127.0.0.1:8888' # 地址
Password = ''       # 密码, 可选
MinIdleConns = 1    # 最小空闲连接数, 可选, 默认1
PoolSize = 1        # 客户端池大小, 可选, 默认1
ReadTimeout = 5000  # 读取超时(毫秒, 可选, 默认5000
WriteTimeout = 5000 # 写入超时(毫秒, 可选, 默认5000
DialTimeout = 5000  # 连接超时(毫秒, 可选, 默认5000
```

## 使用代理

```toml
[services.crawler.proxy]
type = 'static'     # 静态代理, 支持 http, https, socks5, socks5h
address = 'socks5://127.0.0.1:1080'  # 代理地址
User = ''           # 用户名, 可选
Password = ''       # 密码, 可选
```

# 概念

## 种子 `seed`

1. 描述你要下载数据的站点
2. 这个站点如何请求, 请求的表单和body是什么, cookie?, 是否自动跳转, 失败如何重试等
3. 数据拿到后怎么处理
4. 数据在什么时候抓, 数据每隔多久抓一次
5. 这就是`seed`, `seed`描述了一个数据从开始抓取到处理的过程
6. `seed`是无状态的

## 队列 `queue`

1. 我们待抓取的`seed`会存放到队列中, 依次从队列前面拿出一个`seed`开始抓取流程.
2. 如果抓取的数据是一个列表, 如文章列表, 处理程序应该依次遍历并提交包含了文章信息的`seed`
3. 这些`seed`将根据配置放到队列的前面或后面, 然后继续开始下一轮抓取.
4. 队列是框架实现分布式, 并行化, 无状态化的基础.

## 下载器 `downloader`

1. 使用go内置库`http`进行请求
2. 下载器会根据`seed`描述自动构建请求请求体, 请求方法, 请求表单, header, cookie等

## 中间件 `downloader`

1. 中间件包括`请求中间件`和`响应中间件`
2. `请求中间件`的职责是在`downloader`处理`seed`之前检查`seed`的合法性或者判断是否应该请求.
3. `响应中间件`的职责是在`downloader`处理`seed`之后检查数据的合法性或者判断是否应该将`seed`交给处理程序
4. 开发者可以开发自己的中间件

# 设计思路

## 进程独立

1. 在`crawler`的基础设计里, 爬虫运行的最小单元为一个进程, 一个爬虫可能有多进程, 每个进程可以在任何机器上运行.
2. 每个进程同一时间只会处理一个`seed`, 每个进程具有独立的db连接, 独立的`downloader`等, 进程之间互不影响.
3. 你无需关心多进程之间是怎么协调的, 在开发的时候按照单进程开发然后运行时启动多个进程就行了.
4. 多进程需要队列服务支持, 比如`redis`, `ssdb`. 使用`memory`队列开启多进程运行爬虫可能产生意外的结果.

## 请求独立

1. 每个请求都是独立的, `seed`与进程隔离, 进程通过消耗初始种子(初始url)根据处理逻辑生成更多的种子并放入队列, 进程再从队列取出种子你进行请求和解析.
2. 重复的从队列中取出种子处理并保存数据并生成种子放入队列, 直到队列中没有种子为止.

## 抓取过程

1. `seed`在请求前会经过`请求中间件`进行检查.
2. 下载器`downloader`会根据`seed`自动将网站数据下载并写入`seed`中.
3. 下载完成后`seed`会经过`响应中间件`进行检查.
4. 将`seed`交给处理程序, 使用者决定如何对数据进行抽取.
5. 注: 种子抓取过程中如果收到结束信号会将种子放回队列防止种子丢失

## 配置化

1. 在设计上, 尽力将开发中可能存在改变的常量抽离出来形成配置文件, 方便后期调整.

## 模块化

1. 将爬虫的请求, 队列, 代理, 下载器, 配置管理等抽象为单独的模块, 各司其职, 得以解耦合, 方便后期升级维护
2. 使用者也可以根据自己的需求重新设计自己的逻辑替换一些模块.

## 进程管理

开发中...

# 一些操作

## 这条数据我已经抓过了, 怎么让爬虫不再抓它了

1. 将处理完毕(想要拿到的数据已经持久化)的`seed`的唯一标志(一般是url)存入集合.
2. 提交新的`seed`之前检查集合中是否已存在这个唯一标志, 就这么简单
3. 注: 这种方式仍然可能会再次抓取相同的数据, 因为你可能在这个`seed`处理完毕之前又提交了相同唯一标志的`seed`

## 非请求的`seed`

1. 一个`seed`可能不需要请求, 但是`seed`必须要有处理程序

# 工程管理工具

1. 安装
   `go install github.com/zly-app/crawler/tools/crawler@latest`
2. 使用说明
   `crawler help`

## 命令

1. 初始化一个项目
   `crawler init <project_name> && cd <project_name>`
2. 创建一个爬虫
   `crawler create <spider>`
3. 提交初始化种子
   `crawler start [-env env] <spider>`
4. 清空爬虫所有队列
   `crawler clean [-env env] <spider>`
5. 清空爬虫集合数据
   `crawler clean_set [-env env] <spider>`
6. 生成 `supervisor` 配置, 点 [这里](http://supervisord.org/) 进入supervisor官网
    1. `crawler make [-env env]` 会根据文件 `configs/supervisor_programs.toml` 在 `supervisor_config/conf.d` 目录下生成一些 `ini` 文件
    2. 将 `supervisor` 的配置修改包含文件 `<project_dir>/supervisor_config/conf.d/*.ini`

# 调度器工具

1. 安装
   `go build github.com/zly-app/crawler/tools/crawler_scheduler`
2. 将 `supervisor` 的配置修改包含以下文件
   `<project_dir>/supervisor_config/scheduler_config.ini`
