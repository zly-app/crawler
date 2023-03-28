package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// 下载器
type IDownloader interface {
	// 下载
	Download(ctx context.Context, crawler ICrawler, seed *Seed, cookieJar http.CookieJar) (*Seed, error)
	// 关闭
	Close() error
}

// 代理
type IProxy interface {
	// 对传输机启用代理
	EnableProxy(transport *http.Transport)
	// 对传输机关闭代理
	DisableProxy(transport *http.Transport)
	// 对传输机切换代理
	ChangeProxy(transport *http.Transport)
	// 关闭
	Close() error
}
