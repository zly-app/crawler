package seeds

import (
	"testing"

	"github.com/zly-app/crawler/config"
)

// NewSeed 必须按注入的配置填默认值（不再从包级全局读取）
func TestNewSeedUsesInjectedConfig(t *testing.T) {
	conf := &config.ServiceConfig{
		Spider: config.SpiderConfig{
			RequestMethod:  "post",
			AllowProxy:     false,
			UserAgentType:  "android",
			AutoCookie:     true,
			AutoRedirects:  false,
			HtmlEncoding:   "gbk",
		},
		Frame: config.FrameConfig{RequestTimeout: 1234},
	}

	seed := NewSeed(conf)

	if seed.Request.Method != "post" {
		t.Fatalf("Method = %s, want post", seed.Request.Method)
	}
	if seed.Request.AllowProxy {
		t.Fatal("AllowProxy 应为 false")
	}
	if seed.Request.UserAgentType != "android" {
		t.Fatalf("UserAgentType = %s, want android", seed.Request.UserAgentType)
	}
	if !seed.Request.AutoCookie {
		t.Fatal("AutoCookie 应为 true")
	}
	if seed.Request.AutoRedirects {
		t.Fatal("AutoRedirects 应为 false")
	}
	if seed.Request.Encoding != "gbk" {
		t.Fatalf("Encoding = %s, want gbk", seed.Request.Encoding)
	}
	if seed.Request.Timeout != 1234 {
		t.Fatalf("Timeout = %d, want 1234", seed.Request.Timeout)
	}
	if seed.Meta == nil {
		t.Fatal("Meta 不应为 nil")
	}
}

// 两个不同配置生成的种子互不影响
func TestNewSeedIsolated(t *testing.T) {
	a := &config.ServiceConfig{Spider: config.SpiderConfig{RequestMethod: "get"}}
	b := &config.ServiceConfig{Spider: config.SpiderConfig{RequestMethod: "post"}}

	if NewSeed(a).Request.Method != "get" {
		t.Fatal("a 的 Method 不应被 b 影响")
	}
	if NewSeed(b).Request.Method != "post" {
		t.Fatal("b 的 Method 不正确")
	}
}
