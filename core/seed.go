package core

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"github.com/elgs/gojq"

	"github.com/zly-app/crawler/core/dom"
)

var Utf8Boms = [][]byte{
	{0xff, 0xfe, 0, 0}, // LittleEndianBom4
	{0, 0, 0xfe, 0xff}, // BigEndianBom4
	{0xef, 0xbb, 0xbf}, // Utf8Bom3
	{0xfe, 0xff},       // BigEndianBom2
	{0xff, 0xfe},       // LittleEndianBom2
}

// 种子数据
type Seed struct {
	// 原始数据, 构建种子时的数据
	Raw string `json:"-"`
	// 请求时生成的请求体
	HttpRequest *http.Request `json:"-"`
	// 响应
	HttpResponse *http.Response `json:"-"`
	// 响应数据
	HttpResponseBody []byte `json:"-"`
	// 最终的cookies结果
	HttpCookies []*http.Cookie `json:"-"`

	// 请求参数
	Request struct {
		// 请求方法
		Method string
		// 允许使用代理
		AllowProxy bool
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
		// 响应数据编码, 支持 utf8(默认), gbk, gb2312, gb18030, 用于告诉crawler这个页面的编码是什么, crawler会将它转为 utf8.
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

	dom     *dom.Dom
	xmlDom  *dom.XmlDom
	jsonDom *dom.JsonDom
	gojq    *gojq.JQ
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
	case CheckMethod:
		s.CheckExpectMethod = s.getFuncName(t)
	default:
		panic(fmt.Errorf("无法获取方法名: [%T]%v", checkMethod, checkMethod))
	}
}

// 获取dom
func (s *Seed) GetDom() (*dom.Dom, error) {
	if s.dom != nil {
		return s.dom, nil
	}
	if len(s.HttpResponseBody) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewDom(bytes.NewReader(s.HttpResponseBody))
	if err != nil {
		return nil, err
	}
	s.dom = d
	return d, nil
}

// 获取xmlDom
func (s *Seed) GetXmlDom() (*dom.XmlDom, error) {
	if s.xmlDom != nil {
		return s.xmlDom, nil
	}
	if len(s.HttpResponseBody) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewXmlDom(bytes.NewReader(s.HttpResponseBody))
	if err != nil {
		return nil, err
	}
	s.xmlDom = d
	return d, nil
}

// 获取xmlDom
func (s *Seed) GetJsonDom() (*dom.JsonDom, error) {
	if s.jsonDom != nil {
		return s.jsonDom, nil
	}
	if len(s.HttpResponseBody) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewJsonDom(bytes.NewReader(s.HttpResponseBody))
	if err != nil {
		return nil, err
	}
	s.jsonDom = d
	return d, nil
}

// 获取GoJQ github.com/elgs/gojq
func (s *Seed) GetGoJQ() (*gojq.JQ, error) {
	if s.gojq != nil {
		return s.gojq, nil
	}
	if len(s.HttpResponseBody) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := gojq.NewStringQuery(string(s.HttpResponseBody))
	if err != nil {
		return nil, err
	}
	s.gojq = d
	return d, nil
}

// 尝试去掉 utf8 页面的 bom 头
func (s *Seed) TryTrimUtf8Bom() {
	for i := range Utf8Boms {
		if len(s.HttpResponseBody) >= len(Utf8Boms[i]) &&
			bytes.Equal(s.HttpResponseBody[:len(Utf8Boms[i])], Utf8Boms[i]) {
			s.HttpResponseBody = s.HttpResponseBody[len(Utf8Boms[i]):]
			break
		}
	}
}
