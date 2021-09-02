package core

import (
	"net/http"
)

// 下载器
type IDownloader interface {
	// 下载
	Download(crawler ICrawler, seed *Seed, cookieJar http.CookieJar) (*Seed, error)
	// 关闭
	Close() error
}

// 代理
type IProxy interface {
}
