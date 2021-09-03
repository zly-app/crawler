package static

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// socks5代理
type Socks5Proxy struct {
	dialerCtx proxy.ContextDialer
	dialer    proxy.Dialer
}

func (d *Socks5Proxy) EnableProxy(transport *http.Transport) { transport.DialContext = d.dialContext }
func (d *Socks5Proxy) dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if d.dialerCtx != nil {
		return d.dialerCtx.DialContext(ctx, network, address)
	}
	return d.dialer.Dial(network, address)
}
func (d *Socks5Proxy) DisableProxy(transport *http.Transport) { transport.DialContext = nil }
func (d *Socks5Proxy) ChangeProxy(transport *http.Transport)  {}
func (d *Socks5Proxy) Close() error                           { return nil }

// http代理
type HttpProxy struct {
	u *url.URL
}

func (h *HttpProxy) EnableProxy(transport *http.Transport)         { transport.Proxy = h.proxy }
func (h *HttpProxy) proxy(request *http.Request) (*url.URL, error) { return h.u, nil }
func (h *HttpProxy) DisableProxy(transport *http.Transport)        { transport.Proxy = nil }
func (h *HttpProxy) ChangeProxy(transport *http.Transport)         {}
func (h *HttpProxy) Close() error                                  { return nil }

func NewStaticProxy(app zapp_core.IApp) core.IProxy {
	conf := newProxyConfig()
	confKey := fmt.Sprintf("services.%s.proxy", config.NowServiceType)
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		app.Fatal("proxy.static配置错误", zap.Error(err))
	}

	// 解析地址
	u, err := url.Parse(conf.Address)
	if err != nil {
		app.Fatal("proxy.static配置的address无法解析", zap.Error(err))
	}

	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		if conf.User != "" || conf.Password != "" {
			u.User = url.UserPassword(conf.User, conf.Password)
		}
		return &HttpProxy{u: u}
	case "socks5", "socks5h":
		var auth *proxy.Auth
		if conf.User != "" || conf.Password != "" {
			auth = &proxy.Auth{User: conf.User, Password: conf.Password}
		}

		dialer, err := proxy.SOCKS5("tcp", u.Host, auth, nil)
		if err != nil {
			app.Fatal("proxy.sock5生成失败", zap.Error(err))
		}
		if d, ok := dialer.(proxy.ContextDialer); ok {
			return &Socks5Proxy{dialerCtx: d}
		}
		return &Socks5Proxy{dialer: dialer}
	}
	app.Fatal("proxy.socks5配置的scheme不支持", zap.String("scheme", u.Scheme))
	return nil
}
