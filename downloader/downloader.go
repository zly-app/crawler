package downloader

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
)

type Downloader struct {
	app zapp_core.IApp
}

func (d *Downloader) Download(crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	if seed.Request.Url == "" {
		return seed, nil
	}

	// 超时
	ctx, cancel := context.WithTimeout(d.app.BaseContext(), time.Duration(seed.Request.Timeout)*time.Millisecond)
	defer cancel()

	// 构建req
	req, err := d.MakeRequestOfSeed(ctx, seed)
	if err != nil {
		return nil, err
	}

	// 开始请求
	Client.UseSeed(seed)
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 处理
	seed.HttpRequest = req
	seed.HttpResponse = resp
	seed.HttpResponseBody, _ = ioutil.ReadAll(resp.Body)

	// todo 检查cookie

	return seed, nil
}

// 根据seed构建请求
func (d *Downloader) MakeRequestOfSeed(ctx context.Context, seed *core.Seed) (*http.Request, error) {
	method := strings.ToUpper(seed.Request.Method)
	body := strings.NewReader(seed.Request.Body)
	req, err := http.NewRequestWithContext(ctx, method, seed.Request.Url, body)
	if err != nil {
		d.app.Error("构建请求失败", zap.Error(err))
		return nil, core.ParserError
	}

	// 额外url参数
	if len(seed.Request.Params) > 0 {
		query := req.URL.Query()
		for k, v := range seed.Request.Params {
			query[k] = append(query[k], v...)
		}
		req.URL.RawQuery = query.Encode()
	}

	// headers
	req.Header = d.MakeRequestHeadersOfSeed(seed)

	req.Form = seed.Request.Form
	req.PostForm = seed.Request.PostForm
	req.Trailer = seed.Request.Trailer

	return req, nil
}

// 从seed生成请求headers
func (d *Downloader) MakeRequestHeadersOfSeed(seed *core.Seed) http.Header {
	headers := RandomHeaders(seed.Request.UserAgentType)
	for k, v := range seed.Request.Headers {
		headers.Del(k)
		for _, s := range v {
			headers.Add(k, s)
		}
	}
	return headers
}

func (d *Downloader) Close() error {
	return nil
}

func NewDownloader(app zapp_core.IApp) core.IDownloader {
	return &Downloader{
		app: app,
	}
}
