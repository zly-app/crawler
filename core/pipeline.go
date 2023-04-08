package core

import (
	"context"
)

type IPipeline interface {
	Name() string
	// 处理
	Process(ctx context.Context, spiderName string, data interface{}) error
	// 关闭
	Close(ctx context.Context) error
}
