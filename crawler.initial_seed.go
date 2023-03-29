package crawler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"
)

// 提交初始化种子信号
const SubmitInitialSeedSignal = "SubmitInitialSeed"

// 一次性触发时间样式
const OnceTriggerTimeLayout = "2006-01-02 15:04:05"

// 检查提交初始化种子
func (c *Crawler) CheckSubmitInitialSeed(ctx context.Context) {
	expression := c.conf.Spider.SubmitInitialSeedOpportunity

	ctx, span := utils.Otel.StartSpan(ctx, "CheckSubmitInitialSeed",
		utils.OtelSpanKey("expression").String(expression))
	c.app.Info(ctx, "检查提交初始化种子", zap.String("expression", expression))

	switch expression {
	case "", "none":
		utils.Otel.EndSpan(span)
		return
	case "start":
		utils.Otel.EndSpan(span)
		c.SendSubmitInitialSeedSignal(ctx)
		return
	}

	// 一次性触发器
	targetTime, err := time.ParseInLocation(OnceTriggerTimeLayout, expression, time.Local)
	if err == nil {
		interval := targetTime.Sub(time.Now())
		if interval <= 0 {
			return
		}

		go func() {
			timer := time.NewTimer(interval)
			defer timer.Stop()

			select {
			case <-c.app.BaseContext().Done():
			case <-timer.C:
				c.SendSubmitInitialSeedSignal(ctx)
			}
		}()
		return
	}

	// cron表达式
	schedule, err := cron.ParseStandard(expression)
	if err != nil {
		c.app.Fatal(ctx, "spider.SubmitInitialSeedOpportunity 配置格式错误", zap.Error(err))
	}

	go func() {
		targetTime = time.Now()
		for {
			targetTime = schedule.Next(targetTime)

			interval := targetTime.Sub(time.Now())
			if interval <= 0 {
				continue
			}

			timer := time.NewTimer(interval)
			select {
			case <-c.app.BaseContext().Done():
				timer.Stop()
				return
			case <-timer.C:
				c.SendSubmitInitialSeedSignal(ctx)
			}
		}
	}()
}

// 检查是否允许提交初始化种子
func (c *Crawler) checkAllowSubmitInitialSeed(ctx context.Context) bool {
	if !c.conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		return true
	}

	empty, err := c.CheckQueueIsEmpty(ctx, c.conf.Spider.Name)
	if err != nil {
		utils.Otel.AddSpanEvent(utils.Otel.GetSpan(ctx), "CheckQueueIsEmptyErr",
			utils.OtelSpanKey("err.detail").String(err.Error()))
		c.app.Error(ctx, "检查队列是否为空失败", zap.Error(err))
		return false
	}

	utils.Otel.AddSpanEvent(utils.Otel.GetSpan(ctx), "CheckQueueIsEmptyDone",
		utils.OtelSpanKey("empty").Bool(empty))
	if !empty {
		c.app.Debug(ctx, "队列非空忽略初始化种子提交")
		return false
	}
	return true
}

// 发送提交初始化种子信号
func (c *Crawler) SendSubmitInitialSeedSignal(ctx context.Context) {
	ctx, span := utils.Otel.StartSpan(ctx, "SendSubmitInitialSeedSignal",
		utils.OtelSpanKey("StopSubmitInitialSeedIfNotEmptyQueue").Bool(c.conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue))
	defer utils.Otel.EndSpan(span)

	if !c.checkAllowSubmitInitialSeed(ctx) {
		utils.Otel.AddSpanEvent(span, "notAllowSubmitSignal")
		return
	}

	utils.Otel.AddSpanEvent(span, "StartSubmitSignal")
	err := c.PutRawSeed(ctx, SubmitInitialSeedSignal, "", true)
	if err != nil {
		utils.Otel.AddSpanEvent(span, "SubmitSignalErr", utils.OtelSpanKey("err.detail").String(err.Error()))
		utils.Otel.MarkSpanAnError(span, true)
		c.app.Error(ctx, "发送提交初始化种子信号失败", zap.Error(err))
		return
	}

	utils.Otel.AddSpanEvent(span, "SubmitSignalOk")
	c.app.Info(ctx, "发送提交初始化种子信号成功")
}

// 提交种子
func (c *Crawler) SubmitInitialSeed(ctx context.Context) error {
	ctx, span := utils.Otel.StartSpan(ctx, "SubmitInitialSeed")
	defer utils.Otel.EndSpan(span)

	if !c.checkAllowSubmitInitialSeed(ctx) {
		utils.Otel.AddSpanEvent(span, "notAllowSubmit")
		return nil
	}

	utils.Otel.AddSpanEvent(span, "StartSubmit")
	c.app.Info(ctx, "开始提交初始化种子")
	if err := utils.Recover.WrapCall(func() error {
		return c.spider.SubmitInitialSeed(ctx)
	}); err != nil {
		utils.Otel.AddSpanEvent(span, "SubmitErr", utils.OtelSpanKey("err.detail").String(err.Error()))
		utils.Otel.MarkSpanAnError(span, true)
		return err
	}
	utils.Otel.AddSpanEvent(span, "SubmitOk")
	c.app.Info(ctx, "初始化种子提交完成")
	return nil
}
