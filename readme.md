# 分布式爬虫框架服务

> 提供用于 https://github.com/zly-app/zapp 的服务

<!-- TOC -->

- [分布式爬虫框架服务](#%E5%88%86%E5%B8%83%E5%BC%8F%E7%88%AC%E8%99%AB%E6%A1%86%E6%9E%B6%E6%9C%8D%E5%8A%A1)
- [示例](#%E7%A4%BA%E4%BE%8B)
- [配置](#%E9%85%8D%E7%BD%AE)
    - [配置参考](#%E9%85%8D%E7%BD%AE%E5%8F%82%E8%80%83)
- [工程管理工具](#%E5%B7%A5%E7%A8%8B%E7%AE%A1%E7%90%86%E5%B7%A5%E5%85%B7)
    - [快速开始](#%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B)
    - [命令说明](#%E5%91%BD%E4%BB%A4%E8%AF%B4%E6%98%8E)
- [调度器工具](#%E8%B0%83%E5%BA%A6%E5%99%A8%E5%B7%A5%E5%85%B7)
- [概念](#%E6%A6%82%E5%BF%B5)
    - [种子 seed](#%E7%A7%8D%E5%AD%90-seed)
    - [元数据meta](#%E5%85%83%E6%95%B0%E6%8D%AEmeta)
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
    - [这条数据我已经抓过了, 怎么让spider不再抓它了](#%E8%BF%99%E6%9D%A1%E6%95%B0%E6%8D%AE%E6%88%91%E5%B7%B2%E7%BB%8F%E6%8A%93%E8%BF%87%E4%BA%86-%E6%80%8E%E4%B9%88%E8%AE%A9spider%E4%B8%8D%E5%86%8D%E6%8A%93%E5%AE%83%E4%BA%86)
    - [保存抓取结果](#%E4%BF%9D%E5%AD%98%E6%8A%93%E5%8F%96%E7%BB%93%E6%9E%9C)
    - [非请求的seed](#%E9%9D%9E%E8%AF%B7%E6%B1%82%E7%9A%84seed)
    - [使用代理](#%E4%BD%BF%E7%94%A8%E4%BB%A3%E7%90%86)
- [分布式 spider](#%E5%88%86%E5%B8%83%E5%BC%8F-spider)

<!-- /TOC -->

---

# 示例

```go
package main

import (
	"context"

	"github.com/zly-app/zapp"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

// 一个spider
type Spider struct {
	core.ISpiderTool // 必须继承这个接口
}

// 初始化
func (s *Spider) Init(ctx context.Context) error {
	return nil
}

// 提交初始化种子
func (s *Spider) SubmitInitialSeed(ctx context.Context) error {
	seed := s.NewSeed("https://www.sogou.com/", s.Parser) // 创建种子并指定解析方法
	seed.SetCheckExpectMethod(s.Check) // 设置检查方法, 可选
	s.SubmitSeed(ctx, seed)                               // 提交种子
	return nil
}

// 解析方法, 必须以 Parser 开头
func (s *Spider) Parser(ctx context.Context, seed *core.Seed) error {
	data := string(seed.HttpResponseBody) // 获取响应body
	s.SaveResult(ctx, data)               // 保存结果
	return nil
}

// 检查方法. 必须以 Check 开头
func (s *Spider) Check(ctx context.Context, seed *core.Seed) error {
	return nil
}

// 关闭
func (s *Spider) Close(ctx context.Context) error { return nil }

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService()) // 启用crawler服务
	crawler.RegistrySpider(new(Spider))                   // 注入spider
	app.Run()                                             // 运行
}
```

# 配置

+ 不需要任何配置文件就可以运行
+ 默认服务类型为 `crawler`, 完整配置说明请阅读 [Config](./config/config.go)

## 配置参考

```yaml
# spider配置
services:
  crawler:
    spider:
      # spider名
      Name: 'a_spider'
      # 提交初始化种子的时机
      SubmitInitialSeedOpportunity: 'start'
      # 是否自动管理cookie
      AutoCookie: false

    # 框架配置
    frame:
      # 请求超时, 毫秒
      RequestTimeout: 20000
      # 最大尝试请求次数
      RequestMaxAttemptCount: 5
```

---

# 工程管理工具

工程管理工具能帮助你快速搭建和管理你的项目

1. 安装
   `go install github.com/zly-app/crawler/tools/crawler@latest`
2. 使用说明
   `crawler help`

## 快速开始

1. 初始化一个项目

   `crawler init mycrawler && cd mycrawler`

2. 创建一个spider

   `crawler create myspider && cd spiders/myspider`

3. 运行

   `go run .`

## 命令说明

+ `init` 初始化一个项目
   
   `crawler init <project_name>` 命令创建一个项目, 文件夹存在时该文件夹必须是空目录.

   `crawler init .` 命令在当前目录创建项目, 当前目录必须是空目录, 项目名为当前文件夹名.

+ `create` 创建一个 spider, alias `cs`

   `crawler create <spider_name>` 命令在 spiders 目录下创建一个 spider, 如果 spider 文件夹存在且不是空目录时会报错.

   spider 的数据是根据工程下 template/spider_template 目录下的模板生成的.

   如果模板文件后缀名为 `.file`, 则会去除这个后缀
   如果模板文件后缀名为 `.template`, 则会去除这个后缀, 并将模板变量渲染为模板变量的值.
   模板变量以 `{@变量名}` 定义.
   模板变量包含如下数据

   | 模板变量     | 说明                                | 值                  |
   | ------------ | ----------------------------------- | ------------------- |
   | project_name | 项目名                              | <项目名>            |
   | project_dir  | 项目绝对路径                        | <项目绝对路径>      |
   | spider_name  | spider名                            | <spider名>          |
   | spider_dir   | spider绝对路径                      | <spider绝对路径>    |
   | env          | 环境名                              | <环境名>            |
   | date         | 日期. 示例: 2006-01-02              | <当前日期>          |
   | time         | 时间. 示例: 15:04:05                | <当前时间>          |
   | date_time    | 日期时间. 示例: 2006-01-02 15:04:05 | <当前日期时间>      |
   | num_cpu      | cpu逻辑处理器数量                   | <cpu逻辑处理器数量> |

+ `start` 立即提交初始化种子. alias `ss`

   `crawler start [-env <env>] <spider_name>` 命令为spider提交初始化种子信号到指定环境. 执行顺序如下:

   1. 加载 `configs/crawler.<env>.yaml` 和 `spiders/<spider_name>/configs/config.<env>.yaml` 配置文件.
   2. 根据配置文件得到分布式队列配置.
   3. 在分布式队列中提交spider的初始化种子信号.

+ `clean` 清空spider所有队列. **慎用**. 执行顺序如下:

   `crawler clean [-env <env>] <spider_name>` 命令清空spider指定环境的所有队列, 也会清空错误队列.

   1. 加载 `configs/crawler.<env>.yaml` 和 `spiders/<spider_name>/configs/config.<env>.yaml` 配置文件.
   2. 根据配置文件得到分布式队列配置.
   3. 在分布式队列中清空spider所有队列.

+ `clean_set` 清空spider集合数据. **慎用**. 执行顺序如下:

   `crawler clean_set [-env <env>] <spider_name>` 命令清空spider指定环境的所有集合数据.

   1. 加载 `configs/crawler.<env>.yaml` 和 `spiders/<spider_name>/configs/config.<env>.yaml` 配置文件.
   2. 根据配置文件得到分布式集合配置.
   3. 在分布式集合中清空spider集合数据.

+ `make_supervisor` 生成`supervisor`配置. alias `make`. 点 [这里](http://supervisord.org/) 了解`supervisor`

   `crawler make_supervisor [-env <env>]` 命令生成 `supervisor` 配置到 `supervisor_config/conf.d.<env>` 目录下. 生成配置之前会删除这个目录. 执行顺序如下:

   1. 删除 `supervisor_config/conf.d.<env>` 目录并重新创建该目录.
   2. 加载 `configs/supervisor/spider_programs.<env>.yaml` 文件, 记录配置的spider组和spider配置. 
   3. 加载 `template/supervisor/spider_programs.<env>.ini.template` 文件, 这个文件作为spider配置模板.
   4. 遍历要配置的spider组和spider, 根据模板变量渲染spider配置模板后写入spider配置和spider组配置到 `supervisor_config/conf.d.<env>/<group_name>.ini` 文件
   5. 加载 `template/supervisor/scheduler_config.ini.<env>.template` 文件作为调度器配置模板, 根据模板变量渲染后写入到 `supervisor_config/conf.d.<env>/crawler_scheduler.ini` 文件

# 调度器工具

调度器管理工具为分布式spider提供统一的初始化种子信号管理工具.

1. cd 到项目目录下
2. 安装
   `go install github.com/zly-app/crawler/tools/crawler_scheduler@latest && mv ${GOPATH}/bin/crawler_scheduler .`

调度器工具默认加载crawler项目下 `configs/crawler.dev.yaml` 配置文件.

可以通过 `-c` 命令指定配置文件, 也可以在 `supervisor` 配置中的`command`指定.

---

# 概念

## 种子 `seed`

1. 描述你要下载数据的站点
2. 这个站点如何请求, 请求的表单和body是什么, cookie?, 是否自动跳转, 失败如何重试等
3. 数据拿到后怎么处理
4. 数据在什么时候抓, 数据每隔多久抓一次
5. 这就是`seed`, `seed`描述了一个数据从开始抓取到处理的过程
6. `seed`是无状态的

## 元数据`meta`

由于`seed`是无状态的, 但是有些抓取情况是需要先抓取a网站, 再抓取b网站, 而抓取b网站的时候需要将a网站获取到的数据作为分析依据. 此时`meta`就派上用场了.

首先抓取a网站并获取到数据, 接着生成一个新的`seed`用于抓取b网站, 并将需要的数据放入到这个`seed`的`meta`中, 再提交这个`seed`. 当b完整抓取完成时, 再解析函数中的`seed`就会带上这个`meta`.

## 队列 `queue`

1. 我们待抓取的`seed`会存放到队列中, 依次从队列前面拿出一个`seed`开始抓取流程.
2. 如果抓取的数据是一个列表, 如文章列表, 处理程序应该依次遍历并提交包含了文章信息的`seed`
3. 这些`seed`将根据配置放到队列的前面或后面, 然后继续开始下一轮抓取.
4. 队列是框架实现分布式, 并行化, 无状态化的基础.

+ 使用redis作为队列

```yaml
services:
  crawler:
    queue:
      type: 'redis' # 类型
      Address: localhost:6379 # 地址: host1:port1,host2:port2
      UserName: "" # 用户名                     
      Password: "" # 密码
      DB: 0 # db, 只有非集群有效
```

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

1. 在`crawler`的基础设计里, `spider`运行的最小单元为一个进程, 一个`spider`可能有多进程, 每个进程可以在任何机器上运行.
2. 每个进程同一时间只会处理一个`seed`, 每个进程具有独立的db连接, 独立的`downloader`等, 进程之间互不影响.
3. 你无需关心多进程之间是怎么协调的, 在开发的时候按照单进程开发然后运行时启动多个进程就行了.
4. 多进程需要分布式队列服务支持, 比如`redis`. 使用`memory`队列开启多进程运行`spider`可能产生意外的结果.

## 请求独立

1. 每个请求都是独立的, `seed`与进程隔离, 进程通过消耗初始种子(初始url)根据处理逻辑生成更多的种子并放入队列, 进程再从队列取出种子你进行请求和解析.
2. 重复的从队列中取出种子处理并保存数据并生成种子放入队列, 直到队列中没有种子为止.

## 抓取过程

1. `seed`在请求前会经过`请求中间件`进行检查.
2. 下载器`downloader`会根据`seed`自动将网站数据下载并写入`seed`中.
3. 下载完成后`seed`会经过`响应中间件`进行检查.
4. 将`seed`交给处理程序, 使用者决定如何对数据进行抽取.

种子抓取过程中如果进程收到结束信号会将种子放回队列防止种子丢失

## 配置化

1. 在设计上, 尽力将开发中可能存在改变的常量抽离出来形成配置文件, 方便后期调整.

## 模块化

1. 将`spider`的请求, 队列, 代理, 下载器, 配置管理等抽象为单独的模块, 各司其职, 得以解耦合, 方便后期升级维护
2. 使用者也可以根据自己的需求重新设计自己的逻辑替换一些模块.

## 进程管理

通过 `supervisor` 进行进程管理. 也可以自行管理.

---

# 一些操作

## 这条数据我已经抓过了, 怎么让spider不再抓它了

1. 将处理完毕(想要拿到的数据已经持久化)的`seed`的唯一标志(一般是url)存入集合`set`.
2. 提交新的`seed`之前检查集合中是否已存在这个唯一标志

注: 这种方式仍然可能会再次抓取相同的数据, 因为你可能在这个`seed`处理完毕之前又提交了相同唯一标志的`seed`. 但是可以避免绝大部分情况下的重复抓取.
问: 为什么不在处理前存入`seed`的唯一标志. 答: 如果先存放抓取完成标记, 则在抓取过程中一旦程序出现问题, 这条数据将永远不会再抓取了, 所以应该在真的抓取完成后(存入数据后)再写入抓取完成标志.

默认是`memory`提供`set`功能, 在程序重启后之前存入的标志都会失效, 应该使用持久化的`set`. 我们提供了`redis`作为`set`, 可以通过修改配置来使用, 后续可能会提供更多支持. 当然使用者可以自行实现`set`并通过`set.RegistrySetCreator`注册.

+ 使用`redis`作为`set`

```yaml
services:
  crawler:
    set:
      type: 'redis' # 类型
      Address: localhost:6379 # 地址: host1:port1,host2:port2
      UserName: "" # 用户名                     
      Password: "" # 密码
      DB: 0 # db, 只有非集群有效
```

## 保存抓取结果

网站抓取完成后可能需要将需要的数据持久化, 可以调用方法`s.SaveResult`使用`使用管道`Pipeline`工具将结果传输到任何持久储存.

注意: 默认提供的`Pipeline`为`std_out`会将结果打印到控制台, 而非持久化存储. 我们提供了`redis`作为`Pipeline`, 可以通过修改配置来使用, 后续可能会提供更多支持. 当然使用者可以自行实现`Pipeline`并通过`pipeline.RegistryPipelineCreator`注册.

+ 使用`redis`作为`Pipeline`

```yaml
services:
  crawler:
    queue:
      type: 'redis' # 类型
      Address: localhost:6379 # 地址: host1:port1,host2:port2
      UserName: "" # 用户名
      Password: "" # 密码
      DB: 0 # db, 只有非集群有效
```

## 非请求的`seed`

一个`seed`可能不需要发起请求, 但是`seed`必须要有其对应的处理程序. 这类`seed`的用途可以是发起提交种子的信号, 也可以将数据存入meta后进行数据分析等.

## 使用代理

一些数据可能需要代理才能正常抓取, 无需更改spider代码, 在配置中加入代理配置即可使用代理.

```yaml
services:
  crawler:
    proxy:
      type: 'static'     # 静态代理, 支持 http, https, socks5, socks5h
      address: 'socks5://127.0.0.1:1080'  # 代理地址
      User: ''           # 用户名, 可选
      Password: ''       # 密码, 可选
```

---

# 分布式 spider

分布式 spider 必须使用第三方队列`queue`和第三方集合`set`. 万幸的是你不需要更改你的spider代码, 只需要在配置中将队列配置和集合配置改为第三方类型既可以享用. 比如目前提供的`redis`队列以及`redis`集合.

你无需在开发时关心它是分布式的还是单机的, 由于接入了队列, 你直接以单进程spider开发的思路去编写你的代码, 通过更改配置使用分布式队列和分布式集合即可以让单进程的spider在多机器上运行.

分布式 spider 的原理是多个相同spider在运行中将其提交的`seed`都交给第三方队列储存. spider 的各个进程又从第三方队列中各自取出一个`seed`来抓取, 并产生新的`seed`提交给第三方队列. 这样各个进程抓取的`seed`各不相同互不干扰.
