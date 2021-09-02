# 分布式爬虫框架服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
crawler.WithService()           # 启用服务
crawler.RegistrySpider(...)     # 服务注入spider
```

# 示例

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
	crawler core.ICrawler
}

// 初始化
func (s *Spider) Init(crawler core.ICrawler) error {
	s.crawler = crawler
	return nil
}

// 提交初始化种子
func (s *Spider) SubmitInitialSeed() error {
	seed := s.crawler.NewSeed("https://www.baidu.com/", s.Parser) // 创建种子并指定解析方法
	s.crawler.SubmitSeed(seed)                                    // 提交种子
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
