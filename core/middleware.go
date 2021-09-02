package core

// 中间件
type IMiddleware interface {
	// 请求处理
	RequestProcess(crawler ICrawler, seed *Seed) (*Seed, error)
	// 响应处理
	ResponseProcess(crawler ICrawler, seed *Seed) (*Seed, error)
	// 关闭
	Close() error
}

// 请求中间件
type IRequestMiddleware interface {
	// 中间件名
	Name() string
	// 处理
	Process(crawler ICrawler, seed *Seed) (*Seed, error)
	// 关闭
	Close() error
}

// 响应中间件
type IResponseMiddleware interface {
	// 中间件名
	Name() string
	// 处理
	Process(crawler ICrawler, seed *Seed) (*Seed, error)
	// 关闭
	Close() error
}

type MiddlewareBase struct{}

func (m *MiddlewareBase) Name() string                                { return "base" }
func (m *MiddlewareBase) Process(ICrawler, seed *Seed) (*Seed, error) { return seed, nil }
func (m *MiddlewareBase) Close() error                                { return nil }
