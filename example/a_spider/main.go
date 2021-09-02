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
