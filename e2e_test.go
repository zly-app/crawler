package crawler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zly-app/zapp"
	zapp_config "github.com/zly-app/zapp/config"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/core"
)

// e2eSpider 是端到端测试用的 spider：提交一个种子，解析后把响应体送回测试。
type e2eSpider struct {
	core.ISpiderTool
	url  string
	body chan string
}

func (s *e2eSpider) Init(context.Context) error { return nil }

func (s *e2eSpider) SubmitInitialSeed(ctx context.Context) error {
	seed := s.NewSeed(s.url, s.ParserPage)
	s.SubmitSeed(ctx, seed)
	return nil
}

// ParserPage 是解析方法，按 core.ParserMethodNamePrefix 被反射扫描到。
func (s *e2eSpider) ParserPage(_ context.Context, seed *core.Seed) error {
	select {
	case s.body <- string(seed.HttpResponseBody):
	default:
	}
	return nil
}

func (s *e2eSpider) Close(context.Context) error { return nil }

func writeE2EConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "crawler.yaml")
	content := `
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
    queue:
      type: memory
    set:
      type: memory
    proxy:
      type: direct
    pipeline:
      type: stdout
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写配置失败: %v", err)
	}
	return path
}

// TestEndToEndFetch 走完整链路：WithService → Install 注入 spider → Run 启动
// → 提交初始化种子 → 内存队列 → 下载 → 请求/响应中间件 → 反射调用解析方法。
func TestEndToEndFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, "<html><body>hello-e2e</body></html>")
	}))
	defer srv.Close()

	app := zapp.NewApp("crawler-e2e",
		crawler.WithService(),
		zapp.WithConfigOption(zapp_config.WithFiles(writeE2EConfig(t))),
	)

	sp := &e2eSpider{url: srv.URL, body: make(chan string, 1)}
	crawler.Install(app).RegistrySpider(sp)

	// 注意：不能调用 app.Exit() —— zapp 的 Exit 内部会 os.Exit(0)，
	// 在 go test 进程里会触发 "unexpected call to os.Exit"。测试结束由进程回收。
	go app.Run()

	select {
	case body := <-sp.body:
		if !strings.Contains(body, "hello-e2e") {
			t.Fatalf("解析到的响应体不符: %q", body)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("15s 内未完成一次完整抓取")
	}
}
