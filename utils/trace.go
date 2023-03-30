package utils

import (
	"context"
	"time"

	"github.com/zly-app/zapp/pkg/utils"
)

var Trace = new(traceCli)

type traceCli struct{}

func (c *traceCli) TraceStart(ctx context.Context, name string, attributes ...utils.OtelSpanKV) context.Context {
	// 生成新的 span
	ctx, _ = utils.Otel.StartSpan(ctx, name, attributes...)
	return ctx
}

func (*traceCli) getOtelSpanKVWithDeadline(ctx context.Context) utils.OtelSpanKV {
	deadline, deadlineOK := ctx.Deadline()
	if !deadlineOK {
		return utils.OtelSpanKey("ctx.deadline").Bool(false)
	}
	d := deadline.Sub(time.Now()) // 剩余时间
	return utils.OtelSpanKey("ctx.deadline").String(d.String())
}

func (c *traceCli) TraceEvent(ctx context.Context, name string, attributes ...utils.OtelSpanKV) {
	span := utils.Otel.GetSpan(ctx)
	attr := []utils.OtelSpanKV{
		c.getOtelSpanKVWithDeadline(ctx),
	}
	attr = append(attr, attributes...)
	utils.Otel.AddSpanEvent(span, name, attr...)
}

func (c *traceCli) TraceErrEvent(ctx context.Context, name string, err error, attributes ...utils.OtelSpanKV) {
	span := utils.Otel.GetSpan(ctx)
	attr := []utils.OtelSpanKV{
		c.getOtelSpanKVWithDeadline(ctx),
		utils.OtelSpanKey("err.detail").String(err.Error()),
	}
	attr = append(attr, attributes...)
	utils.Otel.AddSpanEvent(span, name+" err", attr...)
	utils.Otel.MarkSpanAnError(span, true)
}

func (*traceCli) TraceEnd(ctx context.Context) {
	span := utils.Otel.GetSpan(ctx)
	utils.Otel.EndSpan(span)
}

func (c *traceCli) AttrKey(key string) utils.OtelSpanKey {
	return utils.OtelSpanKey(key)
}
