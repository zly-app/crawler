package main

import (
	"github.com/zly-app/zapp"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

// 一个spider
type Spider struct {
	core.ISpiderTool
}

// 初始化
func (s *Spider) Init() error {
	return nil
}

// 提交初始化种子
func (s *Spider) SubmitInitialSeed() error {
	seed := s.NewSeed("https://www.sogou.com/", s.Parser) // 创建种子并指定解析方法
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
