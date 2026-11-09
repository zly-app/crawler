package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// stubMW 是测试用中间件，方法集同时满足请求/响应中间件接口
type stubMW struct{ name string }

func (s stubMW) Name() string { return s.name }
func (s stubMW) Process(_ context.Context, _ core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	return seed, nil
}
func (s stubMW) Close() error { return nil }

func newTestConf() *config.ServiceConfig {
	return &config.ServiceConfig{
		Spider: config.SpiderConfig{
			Name:                           "test",
			AutoTrimUtf8BomWithContentType: "text/html",
			HttpStatus4xxIsInvalid:         true,
			HttpStatus5xxIsInvalid:         true,
		},
		Frame: config.FrameConfig{Namespace: "ns"},
	}
}

// 默认装配：1 个请求中间件 + 2 个响应中间件，顺序固定
func TestNewMiddlewareDefaults(t *testing.T) {
	m, ok := NewMiddleware(nil, newTestConf()).(*Middleware)
	if !ok {
		t.Fatal("NewMiddleware 返回类型不是 *Middleware")
	}

	if len(m.requestMiddlewares) != 1 {
		t.Fatalf("默认请求中间件数 = %d, want 1", len(m.requestMiddlewares))
	}
	if len(m.responseMiddlewares) != 2 {
		t.Fatalf("默认响应中间件数 = %d, want 2", len(m.responseMiddlewares))
	}
	if got := m.responseMiddlewares[0].Name(); got != "CheckHttpStatusCode" {
		t.Fatalf("响应中间件[0] = %s, want CheckHttpStatusCode", got)
	}
	if got := m.responseMiddlewares[1].Name(); got != "AutoTrimUtf8Bom" {
		t.Fatalf("响应中间件[1] = %s, want AutoTrimUtf8Bom", got)
	}
}

// 注册追加到链尾，互不影响其他实例
func TestRegistryMiddlewareIsPerInstance(t *testing.T) {
	a := NewMiddleware(nil, newTestConf()).(*Middleware)
	b := NewMiddleware(nil, newTestConf()).(*Middleware)

	a.RegistryRequestMiddleware(stubMW{name: "custom-req"})
	a.RegistryResponseMiddleware(stubMW{name: "custom-rsp"})

	if len(a.requestMiddlewares) != 2 || len(a.responseMiddlewares) != 3 {
		t.Fatalf("注册后链长 = %d/%d, want 2/3", len(a.requestMiddlewares), len(a.responseMiddlewares))
	}
	if len(b.requestMiddlewares) != 1 || len(b.responseMiddlewares) != 2 {
		t.Fatalf("另一个实例被污染: %d/%d", len(b.requestMiddlewares), len(b.responseMiddlewares))
	}
}

// 4xx 按配置拦截
func TestResponseProcessRejects4xx(t *testing.T) {
	m := NewMiddleware(nil, newTestConf())
	seed := &core.Seed{HttpResponse: &http.Response{StatusCode: 404}}

	if _, err := m.ResponseProcess(context.Background(), nil, seed); err == nil {
		t.Fatal("4xx 应被拦截，实际放行")
	}
}

// 5xx 按配置拦截
func TestResponseProcessRejects5xx(t *testing.T) {
	m := NewMiddleware(nil, newTestConf())
	seed := &core.Seed{HttpResponse: &http.Response{StatusCode: 503}}

	if _, err := m.ResponseProcess(context.Background(), nil, seed); err == nil {
		t.Fatal("5xx 应被拦截，实际放行")
	}
}

// ExpectHttpStatusCode 优先于 4xx/5xx 规则
func TestResponseProcessExpectCodeWins(t *testing.T) {
	conf := newTestConf()
	conf.Spider.ExpectHttpStatusCode = []int{404}
	m := NewMiddleware(nil, conf)
	seed := &core.Seed{HttpResponse: &http.Response{StatusCode: 404}}

	if _, err := m.ResponseProcess(context.Background(), nil, seed); err != nil {
		t.Fatalf("期望状态码 404 应放行，实际: %v", err)
	}
}

// 无响应时直接放行
func TestResponseProcessSkipWhenNoResponse(t *testing.T) {
	m := NewMiddleware(nil, newTestConf())

	if _, err := m.ResponseProcess(context.Background(), nil, &core.Seed{}); err != nil {
		t.Fatalf("无响应时应放行，实际: %v", err)
	}
}

// 关闭不 panic
func TestClose(t *testing.T) {
	NewMiddleware(nil, newTestConf()).Close()
}
