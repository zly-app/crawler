package seeds

import (
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// 创建seed
func NewSeed() *core.Seed {
	conf := config.Conf.Spider
	seed := &core.Seed{
		Meta: make(map[string]interface{}),
	}
	seed.Request.Method = conf.RequestMethod
	seed.Request.Params = make(url.Values)
	seed.Request.Headers = make(http.Header)
	seed.Request.UserAgentType = conf.UserAgentType
	seed.Request.Form = make(url.Values)
	seed.Request.PostForm = make(url.Values)
	seed.Request.Trailer = make(http.Header)
	seed.Request.AutoCookie = conf.AutoCookie
	seed.Request.AutoRedirects = conf.AutoRedirects
	seed.Request.Encoding = conf.HtmlEncoding
	seed.Request.Timeout = config.Conf.Frame.RequestTimeout
	return seed
}

// 从原始数据生成seed
func MakeSeedOfRaw(raw string) (*core.Seed, error) {
	seed := &core.Seed{
		Meta: make(map[string]interface{}),
	}
	// 这是为了降低种子保存在队列的占用大小, 不使用用户自定义配置作为默认值是因为用户自定义配置可能会随时变更
	seed.Request.Method = config.DefaultSpiderRequestMethod
	seed.Request.Params = make(url.Values)
	seed.Request.Headers = make(http.Header)
	seed.Request.UserAgentType = config.DefaultSpiderUserAgentType
	seed.Request.Form = make(url.Values)
	seed.Request.PostForm = make(url.Values)
	seed.Request.Trailer = make(http.Header)
	seed.Request.AutoCookie = config.DefaultSpiderAutoCookie
	seed.Request.AutoRedirects = config.DefaultSpiderAutoRedirects
	seed.Request.Encoding = config.DefaultSpiderHtmlEncoding
	seed.Request.Timeout = config.DefaultFrameRequestTimeout
	err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(raw, seed)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

// 将seed编码
//
// 这是为了降低种子保存在队列的占用大小, 不使用用户自定义配置作为默认值是因为用户自定义配置可能会随时变更
func EncodeSeed(seed *core.Seed) (string, error) {
	data := map[string]interface{}{
		"ParserMethod": seed.ParserMethod,
	}
	if seed.CheckExpectMethod != "" {
		data["CheckExpectMethod"] = seed.CheckExpectMethod
	}
	if len(seed.Meta) > 0 {
		data["Meta"] = seed.Meta
	}

	req := make(map[string]interface{})
	if seed.Request.Method != config.DefaultSpiderRequestMethod {
		req["Method"] = seed.Request.Method
	}
	if seed.Request.Url != "" {
		req["Url"] = seed.Request.Url
	}
	if len(seed.Request.Params) > 0 {
		req["Params"] = seed.Request.Params
	}
	if len(seed.Request.Headers) > 0 {
		req["Headers"] = seed.Request.Headers
	}
	if seed.Request.UserAgentType != config.DefaultSpiderUserAgentType {
		req["UserAgentType"] = seed.Request.UserAgentType
	}
	if seed.Request.Body != "" {
		req["Body"] = seed.Request.Body
	}
	if len(seed.Request.Form) > 0 {
		req["Form"] = seed.Request.Form
	}
	if len(seed.Request.PostForm) > 0 {
		req["PostForm"] = seed.Request.PostForm
	}
	if len(seed.Request.Trailer) > 0 {
		req["Trailer"] = seed.Request.Trailer
	}
	if seed.Request.AutoCookie != config.DefaultSpiderAutoCookie {
		req["AutoCookie"] = seed.Request.AutoCookie
	}
	if len(seed.Request.Cookies) > 0 {
		req["Cookies"] = seed.Request.Cookies
	}
	if seed.Request.AutoRedirects != config.DefaultSpiderAutoRedirects {
		req["AutoRedirects"] = seed.Request.AutoRedirects
	}
	if seed.Request.Encoding != config.DefaultSpiderHtmlEncoding {
		req["Encoding"] = seed.Request.Encoding
	}
	if seed.Request.Timeout != config.DefaultFrameRequestTimeout {
		req["Timeout"] = seed.Request.Timeout
	}
	if len(req) > 0 {
		data["Request"] = req
	}

	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(data)
}
