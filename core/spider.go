package core

type ParserMethod = func(seed *Seed) error

type ISpider interface {
	// 初始化
	Init(tool ISpiderTool) error
	// 提交初始化种子
	SubmitInitialSeed() error
	// 关闭
	Close() error
}

// 给spider用的工具
type ISpiderTool interface {
	/**创建种子
	  url 抓取连接
	  parserMethod 解析方法, 可以是方法名或方法实体
	*/
	NewSeed(url string, parserMethod interface{}) *Seed
	// 提交种子
	SubmitSeed(seed *Seed)

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
	SetAdd(key string, items ...string) int
	// 判断集合是否包含某个元素
	SetHasItem(key, item string) bool
	// 从集合中移除一些元素, 返回成功移除的数量, 元素不存在不会计数也不会报错
	SetRemove(key string, items ...string) int
	// 获取集合大小
	GetSetSize(key string) int
}
