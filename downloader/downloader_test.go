package downloader

import "testing"

// NewDownloader 必须自持 http 客户端（不再依赖包级 Client 单例）
func TestNewDownloaderOwnsClient(t *testing.T) {
	d, ok := NewDownloader(nil).(*Downloader)
	if !ok {
		t.Fatal("NewDownloader 返回类型不是 *Downloader")
	}
	if d.client == nil {
		t.Fatal("Downloader.client 为 nil")
	}
	if d.client.Client == nil {
		t.Fatal("HttpClient.Client 为 nil")
	}
	if d.client.Client.Transport == nil {
		t.Fatal("HttpClient.Transport 为 nil")
	}
}

// 不同实例持有独立的客户端
func TestDownloadersHaveDistinctClients(t *testing.T) {
	a := NewDownloader(nil).(*Downloader)
	b := NewDownloader(nil).(*Downloader)

	if a.client == b.client {
		t.Fatal("两个 Downloader 共享了同一个 http 客户端")
	}
}
