package crawler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	zapputils "github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/utils"
)

// 提交初始化种子信号
const SubmitInitialSeedSignal = "SubmitInitialSeed"

// 一次性触发时间样式
const OnceTriggerTimeLayout = "2006-01-02 15:04:05"

// 检查提交初始化种子
func (c *Crawler) CheckSubmitInitialSeed(ctx context.Context) {
	expression := c.conf.Spider.SubmitInitialSeedOpportunity

	startCtx := utils.Trace.TraceStart(ctx, "CheckSubmitInitialSeed", utils.Trace.AttrKey("expression").String(expression))
	c.app.Info(startCtx, "检查提交初始化种子", zap.String("expression", expression))

	switch expression {
	case "", "none":
		utils.Trace.TraceEnd(startCtx)
		return
	case "start":
		c.SendSubmitInitialSeedSignal(startCtx)
		utils.Trace.TraceEnd(startCtx)
		return
	}
	utils.Trace.TraceEnd(startCtx)

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
	ctx = utils.Trace.TraceStart(ctx, "checkAllowSubmitInitialSeed",
		utils.Trace.AttrKey("StopSubmitInitialSeedIfNotEmptyQueue").Bool(c.conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue))
	defer utils.Trace.TraceEnd(ctx)

	if !c.conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		return true
	}

	utils.Trace.TraceEvent(ctx, "CheckQueueIsEmpty", utils.Trace.AttrKey("spiderName").String(c.conf.Spider.Name))
	empty, err := c.CheckQueueIsEmpty(ctx, c.conf.Spider.Name)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "CheckQueueIsEmpty", err)
		c.app.Error(ctx, "检查队列是否为空失败", zap.Error(err))
		return false
	}

	utils.Trace.TraceEvent(ctx, "CheckQueueIsEmptyDone", utils.Trace.AttrKey("empty").Bool(empty))
	if !empty {
		c.app.Debug(ctx, "队列非空忽略初始化种子提交")
		return false
	}
	return true
}

// 发送提交初始化种子信号
func (c *Crawler) SendSubmitInitialSeedSignal(ctx context.Context) {
	ctx = utils.Trace.TraceStart(ctx, "SendSubmitInitialSeedSignal")
	defer utils.Trace.TraceEnd(ctx)

	if !c.checkAllowSubmitInitialSeed(ctx) {
		return
	}

	utils.Trace.TraceEvent(ctx, "SubmitInitialSeedSignal", utils.Trace.AttrKey("front").Bool(true))
	err := c.PutRawSeed(ctx, SubmitInitialSeedSignal, "", true)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "SubmitInitialSeedSignal", err)
		c.app.Error(ctx, "发送提交初始化种子信号失败", zap.Error(err))
		return
	}

	c.app.Info(ctx, "发送提交初始化种子信号成功")
}

// 提交种子
func (c *Crawler) SubmitInitialSeed(ctx context.Context) error {
	ctx = utils.Trace.TraceStart(ctx, "SubmitInitialSeed")
	defer utils.Trace.TraceEnd(ctx)

	if !c.checkAllowSubmitInitialSeed(ctx) {
		return nil
	}

	utils.Trace.TraceEvent(ctx, "StartSubmit")
	c.app.Info(ctx, "开始提交初始化种子")
	if err := zapputils.Recover.WrapCall(func() error {
		return c.spider.SubmitInitialSeed(ctx)
	}); err != nil {
		utils.Trace.TraceErrEvent(ctx, "StartSubmit", err)
		return err
	}
	c.app.Info(ctx, "初始化种子提交完成")
	return nil
}
