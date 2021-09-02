package config

import (
	"errors"

	zapp_core "github.com/zly-app/zapp/core"
)

const (
	// 默认提交初始化种子的时机
	defaultSpiderSubmitInitialSeedOpportunity = "start"
	// 默认使用调度器
	defaultSpiderUseScheduler = false

	// 默认请求方法
	DefaultSpiderRequestMethod = "get"
	// 默认user-agent类型
	DefaultSpiderUserAgentType = "pc"
	// 默认启用cookie
	DefaultSpiderUseCookie = false
	// 默认自动跳转
	DefaultSpiderAutoRedirects = true
	// 默认html编码
	DefaultSpiderHtmlEncoding = "utf8"
	// 默认将http状态码4xx视为无效
	defaultSpiderHttpStatus4xxIsInvalid = true
	// 默认将http状态码5xx视为无效
	defaultSpiderHttpStatus5xxIsInvalid = true
)

type SpiderConfig struct {
	Name string // 爬虫名
	/*
		**提交初始化种子的时机
		 none 无
		 start 启动时
		 YYYY-MM-DD hh:mm:ss 指定时间触发
		 cron表达式
	*/
	SubmitInitialSeedOpportunity string
	// 使用调度器, 提交初始化种子的时机交给调度器管理, 这可以解决多进程运行时每个进程都在提交种子
	UseScheduler bool

	RequestMethod          string // 默认请求方法
	UserAgentType          string // user-agent 类型; pc,android,ios
	UseCookie              bool   // 是否启用cookie
	AutoRedirects          bool   // 是否自动跳转
	HtmlEncoding           string // 默认html编码
	ExpectHttpStatusCode   []int  // 期望的http状态码列表
	InvalidHttpStatusCode  []int  // 无效的http状态码列表, 如果配置了ExpectHttpStatusCode, 则以ExpectHttpStatusCode为准
	HttpStatus4xxIsInvalid bool   // 将http状态码4xx视为无效
	HttpStatus5xxIsInvalid bool   // 将http状态码5xx视为无效
}

func newSpiderConfig(app zapp_core.IApp) SpiderConfig {
	return SpiderConfig{
		Name:                         app.Name(),
		SubmitInitialSeedOpportunity: defaultSpiderSubmitInitialSeedOpportunity,
		UseScheduler:                 defaultSpiderUseScheduler,

		UseCookie:              DefaultSpiderUseCookie,
		AutoRedirects:          DefaultSpiderAutoRedirects,
		HttpStatus4xxIsInvalid: defaultSpiderHttpStatus4xxIsInvalid,
		HttpStatus5xxIsInvalid: defaultSpiderHttpStatus5xxIsInvalid,
	}
}
func (conf *SpiderConfig) Check() error {
	if conf.Name == "" {
		return errors.New("spider.name is empty")
	}
	if conf.RequestMethod == "" {
		conf.RequestMethod = DefaultSpiderRequestMethod
	}
	if conf.UserAgentType == "" {
		conf.UserAgentType = DefaultSpiderUserAgentType
	}
	if conf.HtmlEncoding == "" {
		conf.HtmlEncoding = DefaultSpiderHtmlEncoding
	}
	return nil
}
