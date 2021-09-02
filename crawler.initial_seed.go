package crawler

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"
)

const onceTriggerTimeLayout = "2006-01-02 15:04:05"

// 检查提交初始化种子
func (c *Crawler) CheckSubmitInitialSeed() {
	if c.conf.Spider.UseScheduler { // 交给调度器管理
		return
	}

	expression := c.conf.Spider.SubmitInitialSeedOpportunity
	switch expression {
	case "", "none":
		return
	case "start":
		c.SubmitInitialSeed()
		return
	}

	// 一次性触发器
	targetTime, err := time.ParseInLocation(onceTriggerTimeLayout, expression, time.Local)
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
				c.SubmitInitialSeed()
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
				c.SubmitInitialSeed()
			}
		}
	}()
}

// 提交种子
func (c *Crawler) SubmitInitialSeed() {
	if c.conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		empty, err := c.CheckQueueIsEmpty(c.conf.Spider.Name)
		if err != nil {
			c.app.Error("检查队列是否为空失败", zap.Error(err))
			return
		}

		if !empty {
			c.app.Debug("队列非空忽略初始化种子提交")
			return
		}
	}

	c.app.Info("开始提交初始化种子")
	if err := utils.Recover.WrapCall(c.spider.SubmitInitialSeed); err != nil {
		c.app.Error("提交初始化种子失败", zap.String("error", utils.Recover.GetRecoverErrorDetail(err)))
		return
	}
	c.app.Info("初始化种子提交完成")
}
