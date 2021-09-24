package config

import (
	"strings"
)

const (
	/** 默认队列后缀
	  vip 表示一个高优先级的抓取任务
	  debug 调试用
	  seed 种子队列
	  error 下载错误队列, 再次抓取可能会成功
	*/
	defaultFrameQueueSuffixes = ":vip,:debug,:seed,:error"
	// 种子队列后缀名
	defaultFrameSeedQueueSuffix = ":seed"
	// 错误种子队列后缀名
	defaultFrameErrorSeedQueueSuffix = ":error"
	// 解析错误种子队列后缀名
	defaultFrameParserErrorSeedQueueSuffix = ":error_parser"
	// 集合后缀名
	defaultFrameSetSuffix = ":set"

	// 默认非空队列不提交初始化种子
	defaultFrameStopSubmitInitialSeedIfNotEmptyQueue = true
	// 默认检查是否为空队列的程序忽略error队列
	defaultFrameCheckEmptyQueueIgnoreErrorQueue = true
	// 默认提交种子到队列前面
	defaultFrameSubmitSeedToQueueFront = true

	// 默认请求超时
	DefaultFrameRequestTimeout = 20000
	// 默认下一个种子等待时间
	defaultFrameNextSeedWaitTime = 0
	// 默认spider错误后等待时间
	defaultFrameSpiderErrWaitTime = 3000
	// 默认空队列等待时间
	defaultFrameEmptyQueueWaitTime = 60000
	// 默认重试等待时间
	defaultFrameRequestRetryWaitTime = 1000
	// 默认最大尝试请求次数
	defaultFrameRequestMaxAttemptCount = 5
)

type FrameConfig struct {
	QueueSuffixes              []string // 队列后缀, 按顺序查找种子
	SeedQueueSuffix            string   // 种子队列后缀名
	ErrorSeedQueueSuffix       string   // 错误种子队列后缀名, 主要是下载器错误
	ParserErrorSeedQueueSuffix string   // 解析错误种子队列后缀名, 包括处理程序错误, 处理程序不存在, 种子解析失败等
	SetSuffix                  string   // 集合后缀名

	StopSubmitInitialSeedIfNotEmptyQueue bool // 非空队列不提交初始化种子
	CheckEmptyQueueIgnoreErrorQueue      bool // 检查是否为空队列的程序忽略error队列
	SubmitSeedToQueueFront               bool // 提交种子到队列前面, (取出种子是从前面开始)

	RequestTimeout         int64 // 请求超时, 单位毫秒
	NextSeedWaitTime       int64 // 下一个seed等待时间, 单位毫秒
	SpiderErrWaitTime      int64 // spider错误后等待时间, 单位毫秒
	EmptyQueueWaitTime     int64 // 空队列等待时间, 单位毫秒
	RequestRetryWaitTime   int64 // 请求重试等待时间, 单位毫秒
	RequestMaxAttemptCount int   // 最大尝试请求次数
}

func newFrameConfig() FrameConfig {
	return FrameConfig{
		StopSubmitInitialSeedIfNotEmptyQueue: defaultFrameStopSubmitInitialSeedIfNotEmptyQueue,
		CheckEmptyQueueIgnoreErrorQueue:      defaultFrameCheckEmptyQueueIgnoreErrorQueue,
		SubmitSeedToQueueFront:               defaultFrameSubmitSeedToQueueFront,
	}
}
func (conf *FrameConfig) Check() error {
	if len(conf.QueueSuffixes) == 0 {
		conf.QueueSuffixes = strings.Split(defaultFrameQueueSuffixes, ",")
	}
	if conf.SeedQueueSuffix == "" {
		conf.SeedQueueSuffix = defaultFrameSeedQueueSuffix
	}
	if conf.ErrorSeedQueueSuffix == "" {
		conf.ErrorSeedQueueSuffix = defaultFrameErrorSeedQueueSuffix
	}
	if conf.ParserErrorSeedQueueSuffix == "" {
		conf.ParserErrorSeedQueueSuffix = defaultFrameParserErrorSeedQueueSuffix
	}
	if conf.SetSuffix == "" {
		conf.SetSuffix = defaultFrameSetSuffix
	}

	if conf.RequestTimeout <= 0 {
		conf.RequestTimeout = DefaultFrameRequestTimeout
	}
	if conf.SpiderErrWaitTime <= 0 {
		conf.SpiderErrWaitTime = defaultFrameSpiderErrWaitTime
	}
	if conf.NextSeedWaitTime <= 0 {
		conf.NextSeedWaitTime = defaultFrameNextSeedWaitTime
	}
	if conf.EmptyQueueWaitTime <= 0 {
		conf.EmptyQueueWaitTime = defaultFrameEmptyQueueWaitTime
	}
	if conf.RequestRetryWaitTime <= 0 {
		conf.RequestRetryWaitTime = defaultFrameRequestRetryWaitTime
	}
	if conf.RequestMaxAttemptCount <= 0 {
		conf.RequestMaxAttemptCount = defaultFrameRequestMaxAttemptCount
	}
	return nil
}
