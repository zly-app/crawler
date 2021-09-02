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
}

var Client = NewHttpClient()

func NewHttpClient() *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        5,                // 最大连接数
			MaxIdleConnsPerHost: 5,                // 最大空闲连接数
			IdleConnTimeout:     time.Second * 30, // 空闲连接在关闭自己之前保持空闲的最大时间
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 跳过tls校验
				RootCAs:            x509.NewCertPool(),
			},
		},
	}
	return &HttpClient{client}
}

// 根据seed修改
func (c *HttpClient) UseSeed(seed *core.Seed) *HttpClient {
	if seed.Request.AutoRedirects {
		c.CheckRedirect = nil
	} else {
		c.CheckRedirect = c.closeRedirect
	}
	return c
}

// 关闭重定向
func (c *HttpClient) closeRedirect(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}
