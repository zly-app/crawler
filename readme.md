# 分布式爬虫框架服务

> 提供用于 https://github.com/zly-app/zapp 的服务


<!-- TOC -->

- [分布式爬虫框架服务](#%E5%88%86%E5%B8%83%E5%BC%8F%E7%88%AC%E8%99%AB%E6%A1%86%E6%9E%B6%E6%9C%8D%E5%8A%A1)
- [说明](#%E8%AF%B4%E6%98%8E)
    - [示例](#%E7%A4%BA%E4%BE%8B)
- [配置](#%E9%85%8D%E7%BD%AE)
    - [配置参考](#%E9%85%8D%E7%BD%AE%E5%8F%82%E8%80%83)
    - [使用redis作为队列](#%E4%BD%BF%E7%94%A8redis%E4%BD%9C%E4%B8%BA%E9%98%9F%E5%88%97)
    - [使用redis作为集合](#%E4%BD%BF%E7%94%A8redis%E4%BD%9C%E4%B8%BA%E9%9B%86%E5%90%88)
    - [使用代理](#%E4%BD%BF%E7%94%A8%E4%BB%A3%E7%90%86)

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
	"fmt"
	"github.com/zly-app/zapp"
	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

// 一个spider
type Spider struct {
	core.ISpiderTool
}

// 初始化
func (s *Spider) Init(tool core.ISpiderTool) error {
	s.ISpiderTool = tool
	return nil
}

// 提交初始化种子
func (s *Spider) SubmitInitialSeed() error {
	seed := s.NewSeed("https://www.baidu.com/", s.Parser) // 创建种子并指定解析方法
	s.SubmitSeed(seed)                                    // 提交种子
	return nil
}

// 解析方法
func (s *Spider) Parser(seed *core.Seed) error {
	fmt.Println(string(seed.HttpResponseBody)) // 打印响应body
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
# 使用调度器管理提交初始化种子的时机, 多进程时必须启用
UseScheduler = false
# 是否自动管理cookie
AutoCookie = false

# 框架配置
[services.crawler.frame]
# 请求超时, 毫秒
RequestTimeout = 20000
# 最大尝试请求次数
RequestMaxAttemptCount = 5
```

## 使用redis作为队列

```toml
[services.crawler.queue]
type = 'redis' 		# 使用redis作为队列, 默认是memory
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

## 使用redis作为集合

> 没错, 和队列的配置内容相似

```toml
[services.crawler.set]
type = 'redis' 		# 使用redis作为队列, 默认是memory
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

## 使用代理

```toml
[services.crawler.proxy]
type = 'static' 	# 静态代理, 支持 http, https, socks5, socks5h
address = 'socks5://127.0.0.1:1080'  # 代理地址
User = ''			# 用户名, 可选
Password = ''		# 密码, 可选
```
