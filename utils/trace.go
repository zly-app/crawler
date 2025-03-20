package utils

import (
	"context"

	"github.com/zly-app/zapp/pkg/utils"
)

var Trace = new(traceCli)

type traceCli struct{}

func (c *traceCli) TraceStart(ctx context.Context, name string, attributes ...utils.OtelSpanKV) context.Context {
	// 生成新的 span
	ctx = utils.Otel.CtxStart(ctx, name, attributes...)
	return ctx
}

func (c *traceCli) TraceEvent(ctx context.Context, name string, attributes ...utils.OtelSpanKV) {
	utils.Otel.CtxEvent(ctx, name, attributes...)
}

func (c *traceCli) TraceErrEvent(ctx context.Context, name string, err error, attributes ...utils.OtelSpanKV) {
	utils.Otel.CtxErrEvent(ctx, name, err, attributes...)
}

func (*traceCli) TraceEnd(ctx context.Context) {
	utils.Otel.CtxEnd(ctx)
}

func (c *traceCli) AttrKey(key string) utils.OtelSpanKey {
	return utils.OtelSpanKey(key)
}
