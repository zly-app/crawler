package crawler

import (
	"fmt"
	"time"

	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/seeds"
)

func (c *Crawler) Run() {
	// 运行前提交初始化种子
	c.CheckSubmitInitialSeed()

	for {
		select {
		case <-c.app.BaseContext().Done():
			return
		default:
		}

		err := utils.Recover.WrapCall(c.runOnce)
		if err != nil {
			c.app.Error("运行出错, 稍后继续", zap.Int64("waitTime", c.conf.Frame.SpiderErrWaitTime/1000),
				zap.String("error", utils.Recover.GetRecoverErrorDetail(err)))
			time.Sleep(time.Duration(c.conf.Frame.SpiderErrWaitTime) * time.Millisecond)
		}
	}
}

// 开始一次任务
func (c *Crawler) runOnce() error {
	// 获取种子原始数据
	raw, err := c.PopARawSeed()
	if err == core.EmptyQueueError {
		c.app.Info("空队列, 休眠后重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime/1000))
		time.Sleep(time.Duration(c.conf.Frame.EmptyQueueWaitTime) * time.Millisecond)
		return nil
	}
	if err != nil {
		c.app.Error("从队列获取种子失败, 休眠后重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime), zap.Error(err))
		time.Sleep(time.Duration(c.conf.Frame.EmptyQueueWaitTime) * time.Millisecond)
		return nil
	}

	// 开始处理
	err = utils.Recover.WrapCall(func() error {
		return c.seedProcess(raw)
	})
	if err == nil {
		return nil
	}

	switch err {
	case core.InterceptError: // 拦截, 应该立即结束本次任务
		return nil
	case core.ParserError: // 解析错误
		c.PutErrorRawSeed(raw, true)
	default:
		c.PutErrorRawSeed(raw, false)
	}
	return err
}

// 种子处理
func (c *Crawler) seedProcess(raw string) error {
	seed, err := seeds.MakeSeedOfRaw(raw)
	if err != nil {
		c.app.Error("构建种子失败")
		return core.ParserError
	}

	// 请求处理
	seed, err = c.middleware.RequestProcess(c, seed)
	if err != nil {
		return err
	}

	var seedResult *core.Seed
	// 循环尝试下载
	var attempt int
	for {
		// 每次重新生成seed, 因为每次处理可能会修改seed
		seedCopy := *seed
		seedResult, err = c.download(raw, &seedCopy)
		if err == nil {
			break
		}

		if err == core.InterceptError || err == core.ParserError {
			return err
		}
		c.app.Error("尝试下载失败", zap.Int("attempt", attempt), zap.Error(err))
		attempt++
		if attempt > c.conf.Frame.RequestMaxAttemptCount {
			return fmt.Errorf("超过最大尝试次数")
		}
		time.Sleep(time.Duration(c.conf.Frame.RequestRetryWaitTime) * time.Millisecond)
	}

	// 解析
	return c.Parser(seedResult)
}

// 下载完善种子
func (c *Crawler) download(raw string, seed *core.Seed) (*core.Seed, error) {
	// 下载
	seed, err := c.downloader.Download(c, seed)
	if err != nil {
		return nil, err
	}

	// 响应处理
	seed.Raw = raw
	seed, err = c.middleware.ResponseProcess(c, seed)
	if err != nil {
		return nil, err
	}

	// 检查是符合期望的响应
	seed.Raw = raw
	seed, err = c.CheckIsExpectResponse(seed)
	if err != nil {
		return nil, err
	}

	return seed, nil
}
