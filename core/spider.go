package core

type ParserMethod = func(seed *Seed) error

type ISpider interface {
	// 初始化
	Init(crawler ICrawler) error
	// 提交初始化种子
	SubmitInitialSeed() error
	// 关闭
	Close() error
}
