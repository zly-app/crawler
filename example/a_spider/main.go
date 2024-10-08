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
