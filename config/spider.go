package config

import (
	"errors"

	zapp_core "github.com/zly-app/zapp/core"
)

const (
	// 默认提交初始化种子的时机
	defaultSpiderSubmitInitialSeedOpportunity = "start"

	// 默认请求方法
	DefaultSpiderRequestMethod = "get"
	// 允许使用代理
	DefaultSpiderAllowProxy = true
	// 默认user-agent类型
	DefaultSpiderUserAgentType = "pc"
	// 默认自动管理cookie
	DefaultSpiderAutoCookie = false
	// 默认自动跳转
	DefaultSpiderAutoRedirects = true
	// 默认html编码
	DefaultSpiderHtmlEncoding = "utf8"
	// 默认将http状态码4xx视为无效
	defaultSpiderHttpStatus4xxIsInvalid = true
	// 默认将http状态码5xx视为无效
	defaultSpiderHttpStatus5xxIsInvalid = true
	// 对某些 ContentType 自动去掉utf8bom头, 多种类型用英文逗号分隔
	defAutoTrimUtf8BomWithContentType = "text/html,text/plain,text/xml"
)

type SpiderConfig struct {
	Name string // 爬虫名
	/*
		**提交初始化种子的时机
		 none或空字符串 交给调度器管理
		 start 启动时
		 YYYY-MM-DD hh:mm:ss 指定时间触发
		 cron表达式
	*/
	SubmitInitialSeedOpportunity string

	RequestMethod                  string // 默认请求方法
	AllowProxy                     bool   // 允许使用代理
	UserAgentType                  string // user-agent 类型; pc,android,ios
	AutoCookie                     bool   // 是否自动管理cookie, 当前任务提交的种子会继承之前的cookies
	AutoRedirects                  bool   // 是否自动跳转
	HtmlEncoding                   string // 默认html编码
	ExpectHttpStatusCode           []int  // 期望的http状态码列表, 示例: [200, 204]
	InvalidHttpStatusCode          []int  // 无效的http状态码列表, 如果配置了ExpectHttpStatusCode, 则以ExpectHttpStatusCode为准. 示例: [404, 500]
	HttpStatus4xxIsInvalid         bool   // 将http状态码4xx视为无效, spider会自动重试
	HttpStatus5xxIsInvalid         bool   // 将http状态码5xx视为无效, spider会自动重试
	AutoTrimUtf8BomWithContentType string // 对某些 ContentType 自动去掉utf8bom头, 多种类型用英文逗号分隔
}

func newSpiderConfig(app zapp_core.IApp) SpiderConfig {
	return SpiderConfig{
		Name:                         app.Name(),
		SubmitInitialSeedOpportunity: defaultSpiderSubmitInitialSeedOpportunity,

		AllowProxy:                     DefaultSpiderAllowProxy,
		AutoCookie:                     DefaultSpiderAutoCookie,
		AutoRedirects:                  DefaultSpiderAutoRedirects,
		HttpStatus4xxIsInvalid:         defaultSpiderHttpStatus4xxIsInvalid,
		HttpStatus5xxIsInvalid:         defaultSpiderHttpStatus5xxIsInvalid,
		AutoTrimUtf8BomWithContentType: defAutoTrimUtf8BomWithContentType,
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
