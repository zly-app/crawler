package core

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

// 种子数据
type Seed struct {
	// 原始数据, 构建种子时的数据
	Raw string `json:"-"`

	HttpRequest *http.Request `json:"-"`
	// 响应, 注意: 不能使用它的body, 而是应该使用 ResponseBody
	HttpResponse *http.Response `json:"-"`
	// 响应数据
	HttpResponseBody []byte `json:"-"`
	// 最终的cookies结果
	HttpCookies []*http.Cookie `json:"-"`

	// 请求参数
	Request struct {
		// 请求方法
		Method string
		// url
		Url string
		// url参数
		Params url.Values
		// headers
		Headers http.Header
		// user-agent类型
		UserAgentType string

		// 请求body数据
		Body string
		// 表单
		Form url.Values
		// post表单
		PostForm url.Values
		// 附加头
		Trailer http.Header

		// 是否自动管理cookie
		AutoCookie bool
		// 父cookies, 保留提交这个种子的种子的最终cookies
		//
		// 由于headers中无法保留cookie的过期时间等相关信息, 只能额外提供字段来保存cookie的完整信息
		// spider开发者不应该主动修改这个值
		ParentCookies []*http.Cookie
		// cookies, 这里的cookie会覆盖headers的cookie
		Cookies []*http.Cookie
		// 是否自动跳转
		AutoRedirects bool
		// 响应数据编码
		Encoding string
		// 请求超时时间, 毫秒
		Timeout int64
	}

	// 解析方法名
	ParserMethod string
	// 检查期望方法名
	CheckExpectMethod string
	// 元数据
	Meta map[string]interface{}
}

// GetFuncName 获取函数或方法的名称
func (*Seed) getFuncName(a interface{}) string {
	p := reflect.ValueOf(a).Pointer()
	rawName := runtime.FuncForPC(p).Name()
	name := strings.TrimSuffix(rawName, ".func1")
	ss := strings.Split(name, ".")
	name = strings.TrimSuffix(ss[len(ss)-1], "-fm")
	return name
}

// 设置解析方法
func (s *Seed) SetParserMethod(parserMethod interface{}) {
	switch t := parserMethod.(type) {
	case string:
		s.ParserMethod = t
	case ParserMethod:
		s.ParserMethod = s.getFuncName(t)
	default:
		panic(fmt.Errorf("无法获取方法名: [%T]%v", parserMethod, parserMethod))
	}
}

// 设置检查期望响应方法
func (s *Seed) SetCheckExpectMethod(checkMethod interface{}) {
	switch t := checkMethod.(type) {
	case string:
		s.CheckExpectMethod = t
	case ParserMethod:
		s.CheckExpectMethod = s.getFuncName(t)
	default:
		panic(fmt.Errorf("无法获取方法名: [%T]%v", checkMethod, checkMethod))
	}
}
