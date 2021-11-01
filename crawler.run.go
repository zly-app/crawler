package crawler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/cookiejar"
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
			continue
		}

		if c.conf.Frame.NextSeedWaitTime > 0 {
			time.Sleep(time.Duration(c.conf.Frame.NextSeedWaitTime) * time.Millisecond)
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

	// 提交初始化种子信号
	if raw == SubmitInitialSeedSignal {
		if err = c.SubmitInitialSeed(); err != nil {
			return fmt.Errorf("提交初始化种子失败: %v", err)
		}
		return nil
	}

	// 保存原始种子数据
	c.nowRawSeed.Store(raw)
	defer c.nowRawSeed.Store("")

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
		if err := c.PutErrorRawSeed(raw, true); err != nil {
			c.app.Error("将出错seed放入error队列失败", zap.Error(err))
		}
	default:
		if err := c.PutErrorRawSeed(raw, false); err != nil {
			c.app.Error("将出错seed放入error队列失败", zap.Error(err))
		}
	}
	return err
}

// 种子处理
func (c *Crawler) seedProcess(raw string) error {
	var seedResult *core.Seed
	var cookieJar http.CookieJar
	// 循环尝试下载
	var attempt int
	for {
		attempt++

		// 每次重新生成seed, 因为每次处理可能会修改seed
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

		seedResult, cookieJar, err = c.download(raw, seed)
		if err == nil {
			break
		}

		if err == core.InterceptError || err == core.ParserError {
			return err
		}
		if attempt >= c.conf.Frame.RequestMaxAttemptCount {
			return fmt.Errorf("尝试下载失败, 超过最大尝试次数: %v", err)
		}
		c.app.Error("尝试下载失败, 等待重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime),
			zap.String("attempt", fmt.Sprintf("%d/%d", attempt, c.conf.Frame.RequestMaxAttemptCount)), zap.Error(err))
		time.Sleep(time.Duration(c.conf.Frame.RequestRetryWaitTime) * time.Millisecond)
	}

	// 保存cookieJar
	c.cookieJar = cookieJar
	defer func() {
		c.cookieJar = nil
	}()

	// 解析
	err := utils.Recover.WrapCall(func() error {
		return c.Parser(seedResult)
	})
	if err == nil {
		return nil
	}

	_, ok := utils.Recover.GetRecoverError(err)
	if !ok {
		return err
	}
	c.app.Error("解析时panic", zap.String("err", utils.Recover.GetRecoverErrorDetail(err)))
	return core.ParserError
}

// 下载完善种子
func (c *Crawler) download(raw string, seed *core.Seed) (*core.Seed, http.CookieJar, error) {
	cookieJar, _ := cookiejar.New(nil)

	// 下载
	seed, err := c.downloader.Download(c, seed, cookieJar)
	if err != nil {
		return nil, nil, err
	}

	// 响应处理
	seed.Raw = raw
	seed, err = c.middleware.ResponseProcess(c, seed)
	if err != nil {
		return nil, nil, err
	}

	// 检查是符合期望的响应
	seed.Raw = raw
	seed, err = c.CheckIsExpectResponse(seed)
	if err != nil {
		return nil, nil, err
	}

	return seed, cookieJar, nil
}
