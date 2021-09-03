package downloader

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/zly-app/crawler/core"
)

type HttpClient struct {
	*http.Client
	transport *http.Transport
}

var Client = NewHttpClient()

func NewHttpClient() *HttpClient {
	transport := &http.Transport{
		MaxIdleConns:        5,                // 最大连接数
		MaxIdleConnsPerHost: 5,                // 最大空闲连接数
		IdleConnTimeout:     time.Second * 30, // 空闲连接在关闭自己之前保持空闲的最大时间
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过tls校验
			RootCAs:            x509.NewCertPool(),
		},
	}
	client := &http.Client{
		Transport: transport,
	}
	return &HttpClient{
		Client:    client,
		transport: transport,
	}
}

// 根据seed修改
func (c *HttpClient) UseSeed(crawler core.ICrawler, seed *core.Seed) {
	// 重定向
	if seed.Request.AutoRedirects {
		c.CheckRedirect = nil
	} else {
		c.CheckRedirect = c.closeRedirect
	}

	// 代理
	if seed.Request.AllowProxy {
		crawler.Proxy().EnableProxy(c.transport)
	} else {
		crawler.Proxy().DisableProxy(c.transport)
	}
}

// 切换代理
func (c *HttpClient) ChangeProxy(crawler core.ICrawler, seed *core.Seed) {
	if seed.Request.AllowProxy {
		crawler.Proxy().ChangeProxy(c.transport)
	}
}

// 关闭重定向
func (c *HttpClient) closeRedirect(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}
