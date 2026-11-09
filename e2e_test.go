package crawler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zly-app/zapp"
	zapp_config "github.com/zly-app/zapp/config"
	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
)

// e2eSpider 是一次启动跑多条路径的 spider：成功、重试、解析失败。
type e2eSpider struct {
	core.ISpiderTool
	okURL    string
	retryURL string
	badURL   string
	body     chan string
}

func (s *e2eSpider) Init(context.Context) error { return nil }

func (s *e2eSpider) SubmitInitialSeed(ctx context.Context) error {
	s.SubmitSeed(ctx, s.NewSeed(s.okURL, s.ParserPage))
	s.SubmitSeed(ctx, s.NewSeed(s.retryURL, s.ParserPage))
	s.SubmitSeed(ctx, s.NewSeed(s.badURL, s.ParserPage))
	return nil
}

// ParserPage 是解析方法，按 core.ParserMethodNamePrefix 被反射扫描到。
func (s *e2eSpider) ParserPage(_ context.Context, seed *core.Seed) error {
	if strings.HasSuffix(seed.Request.Url, "/bad") {
		return errors.New("解析失败（测试注入）")
	}
	select {
	case s.body <- string(seed.HttpResponseBody):
	default:
	}
	return nil
}

func (s *e2eSpider) Close(context.Context) error { return nil }

const e2eConfig = `
services:
  crawler:
    spider:
      name: e2e
      submitInitialSeedOpportunity: start
    frame:
      namespace: e2e
      emptyQueueWaitTime: 20
      spiderErrWaitTime: 20
      seedProcessTimeout: 5000
      requestRetryWaitTime: 10
      requestMaxAttemptCount: 5
    queue:
      type: memory
    set:
      type: memory
    proxy:
      type: direct
    pipeline:
      type: stdout
`

// TestCrawlerEndToEnd 在同一个 zapp 实例里跑三条路径：
// 1) 正常抓取；2) 两次 500 后重试成功；3) 解析失败进入解析错误队列。
// zapp 的服务类型只能注册一次，所以多条路径必须共享同一个 App。
func TestCrawlerEndToEnd(t *testing.T) {
	var retryHits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/retry"):
			if atomic.AddInt32(&retryHits, 1) < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, "ok-after-retry")
		case strings.HasPrefix(r.URL.Path, "/bad"):
			_, _ = fmt.Fprint(w, "<html>bad</html>")
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, "<html><body>hello-e2e</body></html>")
		}
	}))
	defer srv.Close()

	sp := &e2eSpider{
		okURL:    srv.URL + "/ok",
		retryURL: srv.URL + "/retry",
		badURL:   srv.URL + "/bad",
		body:     make(chan string, 2),
	}
	app := startCrawler(t, sp)

	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case body := <-sp.body:
			switch {
			case strings.Contains(body, "hello-e2e"):
				seen["ok"] = true
			case strings.Contains(body, "ok-after-retry"):
				seen["retry"] = true
			default:
				t.Fatalf("解析到的响应体不符: %q", body)
			}
		case <-time.After(15 * time.Second):
			t.Fatalf("15s 内只完成 %d/2 次成功抓取", i)
		}
	}
	if !seen["ok"] || !seen["retry"] {
		t.Fatalf("成功路径未全部覆盖: %v", seen)
	}
	if n := atomic.LoadInt32(&retryHits); n != 3 {
		t.Fatalf("重试路径请求次数 = %d, want 3", n)
	}

	// 解析失败的种子必须进入解析错误队列
	svc, ok := app.GetService(config.DefaultServiceType)
	if !ok {
		t.Fatal("crawler 服务未启用")
	}
	c, ok := svc.(*crawler.Crawler)
	if !ok {
		t.Fatalf("服务类型不符: %T", svc)
	}
	queueName := c.Conf().Frame.Namespace + "." + c.Conf().Spider.Name + c.Conf().Frame.ParserErrorSeedQueueSuffix

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if raw, err := c.Queue().Pop(context.Background(), queueName, false); err == nil && raw != "" {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("15s 内解析错误种子未进入队列 %s", queueName)
}

// startCrawler 启动真实的 zapp + crawler 组合根。
func startCrawler(t *testing.T, sp *e2eSpider) zapp_core.IApp {
	t.Helper()

	path := filepath.Join(t.TempDir(), "crawler.yaml")
	if err := os.WriteFile(path, []byte(e2eConfig), 0o644); err != nil {
		t.Fatalf("写配置失败: %v", err)
	}

	app := zapp.NewApp("crawler-e2e",
		crawler.WithService(),
		zapp.WithConfigOption(zapp_config.WithFiles(path)),
	)
	crawler.Install(app).RegistrySpider(sp)

	// zapp 的 Exit 内部会 os.Exit(0)，测试里不调用；进程结束由 go test 回收。
	go app.Run()
	return app
}
