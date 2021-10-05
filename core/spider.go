package core

import (
	"github.com/zly-app/crawler/core/dom"
)

type (
	// 解析方法
	ParserMethod = func(seed *Seed) error
	// 检查方法
	CheckMethod = func(seed *Seed) error
)

const (
	// 解析方法名前缀
	ParserMethodNamePrefix = "Parser"
	// 检查方法名前缀
	CheckMethodNamePrefix = "Check"
)

type ISpider interface {
	// 初始化
	Init() error
	// 提交初始化种子
	SubmitInitialSeed() error
	// 关闭
	Close() error
}

// 给spider用的工具
type ISpiderTool interface {
	Crawler() ICrawler
	/**创建种子
	  url 抓取连接
	  parserMethod 解析方法, 可以是方法名或方法实体
	*/
	NewSeed(url string, parserMethod interface{}) *Seed
	// 提交种子
	SubmitSeed(seed *Seed)
	// 保存结果
	SaveResult(data interface{})

	/*
	   **放入种子
	    seed 种子
	    front 是否放在队列前面
	*/
	PutSeed(seed *Seed, front bool)
	/*
	   **放入种子原始数据
	    raw 种子原始数据
	    parserFuncName 解析函数名
	    front 是否放在队列前面
	*/
	PutRawSeed(raw string, parserFuncName string, front bool)
	/*
	   **放入错误种子
	    seed 种子
	    isParserError 是否为解析错误
	*/
	PutErrorSeed(seed *Seed, isParserError bool)
	/*
	   **放入错误种子原始数据
	    raw 种子原始数据
	    isParserError 是否为解析错误
	*/
	PutErrorRawSeed(raw string, isParserError bool)

	// 添加一些元素到集合中, 返回添加的数量, 已存在的元素不会计数
	SetAdd(items ...string) int
	// 判断集合是否包含某个元素
	SetHasItem(item string) bool
	// 从集合中移除一些元素, 返回成功移除的数量, 元素不存在不会计数也不会报错
	SetRemove(items ...string) int
	// 获取集合大小
	GetSetSize() int
	// 生成相对于在当前种子页面上的某个连接的实际连接
	UrlJoin(seed *Seed, link string) string
	// 获取dom
	GetDom(seed *Seed) *dom.Dom
	// 获取xmlDom
	GetXmlDom(seed *Seed) *dom.XmlDom
	// 获取jsonDom
	GetJsonDom(seed *Seed) *dom.JsonDom
}
