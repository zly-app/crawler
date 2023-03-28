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
	switch expression {
	case "", "none":
		return
	case "start":
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
		c.app.Fatal("spider.SubmitInitialSeedOpportunity 配置格式错误", zap.Error(err))
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
		c.app.Error("检查队列是否为空失败", zap.Error(err))
		return false
	}

	if !empty {
		c.app.Debug("队列非空忽略初始化种子提交")
		return false
	}
	return true
}

// 发送提交初始化种子信号
func (c *Crawler) SendSubmitInitialSeedSignal(ctx context.Context) {
	if !c.checkAllowSubmitInitialSeed(ctx) {
		return
	}

	err := c.PutRawSeed(ctx, SubmitInitialSeedSignal, "", true)
	if err != nil {
		c.app.Error("发送提交初始化种子信号失败", zap.Error(err))
		return
	}

	c.app.Info("发送提交初始化种子信号成功")
}

// 提交种子
func (c *Crawler) SubmitInitialSeed(ctx context.Context) error {
	if !c.checkAllowSubmitInitialSeed(ctx) {
		return nil
	}

	c.app.Info("开始提交初始化种子")
	if err := utils.Recover.WrapCall(func() error {
		return c.spider.SubmitInitialSeed(ctx)
	}); err != nil {
		return err
	}
	c.app.Info("初始化种子提交完成")
	return nil
}
