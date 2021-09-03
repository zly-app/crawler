package core

import (
	"errors"
	"net/http"
)

var (
	InterceptError = errors.New("intercept error") // 拦截错误, 应该立即结束本次任务
	ParserError    = errors.New("parser error")    // 解析错误, 这种错误应该将seed放入解析错误队列
)

type ICrawler interface {
	// 获取spider
	Spider() ISpider
	// 获取队列
	Queue() IQueue
	// 获取集合
	Set() ISet
	// 获取下载器
	Downloader() IDownloader
	// 代理
	Proxy() IProxy
	// 获取当前的cookieJar, 可能是空的
	CookieJar() http.CookieJar

	/*
	   **放入种子
	    seed 种子
	    front 是否放在队列前面
	*/
	PutSeed(seed *Seed, front bool) error
	/*
	   **放入种子原始数据
	    raw 种子原始数据
	    parserFuncName 解析函数名
	    front 是否放在队列前面
	*/
	PutRawSeed(raw string, parserFuncName string, front bool) error
	/*
	   **放入错误种子
	    seed 种子
	    isParserError 是否为解析错误
	*/
	PutErrorSeed(seed *Seed, isParserError bool) error
	/*
	   **放入错误种子原始数据
	    raw 种子原始数据
	    isParserError 是否为解析错误
	*/
	PutErrorRawSeed(raw string, isParserError bool) error

	// 检查队列是否为空, 如果spiderName为空则取默认值
	CheckQueueIsEmpty(spiderName string) (bool, error)

	// 获取spider解析方法
	GetSpiderParserMethod(methodName string) (ParserMethod, bool)
}
