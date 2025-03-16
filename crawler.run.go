package crawler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zly-app/zapp/filter"
	zapputils "github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/utils"

	"github.com/zly-app/crawler/cookiejar"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/seeds"
)

func (c *Crawler) Run() {
	// 运行前提交初始化种子
	c.CheckSubmitInitialSeed(c.app.BaseContext())

	for {
		select {
		case <-c.app.BaseContext().Done():
			return
		default:
		}

		ctx := utils.Trace.TraceStart(context.Background(), "runOnceSeed")
		err := zapputils.Recover.WrapCall(func() error {
			return c.runOnce(ctx)
		})
		if err != nil {
			c.app.Error("运行出错, 稍后继续", zap.Int64("waitTime", c.conf.Frame.SpiderErrWaitTime/1000),
				zap.String("error", zapputils.Recover.GetRecoverErrorDetail(err)))
			time.Sleep(time.Duration(c.conf.Frame.SpiderErrWaitTime) * time.Millisecond)
			continue
		}

		if c.conf.Frame.NextSeedWaitTime > 0 {
			time.Sleep(time.Duration(c.conf.Frame.NextSeedWaitTime) * time.Millisecond)
		}
	}
}

// 开始一次任务
func (c *Crawler) runOnce(ctx context.Context) error {
	timeCtx, cancel := context.WithTimeout(ctx, time.Duration(c.conf.Frame.SeedProcessTimeout)*time.Millisecond)
	defer cancel()

	// 获取种子原始数据
	raw, err := c.PopARawSeed(timeCtx)
	if err == core.EmptyQueueError {
		utils.Trace.TraceEvent(timeCtx, "emptyQueue")
		c.app.Info(timeCtx, "空队列, 休眠后重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime/1000))
		utils.Trace.TraceEnd(ctx)
		time.Sleep(time.Duration(c.conf.Frame.EmptyQueueWaitTime) * time.Millisecond)
		return nil
	}
	if err != nil {
		utils.Trace.TraceErrEvent(timeCtx, "PopARawSeed", err)
		c.app.Error(timeCtx, "从队列获取种子失败, 休眠后重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime), zap.Error(err))
		utils.Trace.TraceEnd(ctx)
		time.Sleep(time.Duration(c.conf.Frame.EmptyQueueWaitTime) * time.Millisecond)
		return nil
	}

	defer utils.Trace.TraceEnd(ctx)

	// 提交初始化种子信号
	if raw == SubmitInitialSeedSignal {
		if err = c.SubmitInitialSeed(ctx); err != nil {
			return fmt.Errorf("提交初始化种子失败: %v", err)
		}
		return nil
	}

	// 保存原始种子数据
	c.nowRawSeed.Store(raw)
	defer c.nowRawSeed.Store("")

	// 开始处理
	c.app.Info(timeCtx, "开始处理种子")

	seedProcessCtx, chain := filter.GetServiceFilter(timeCtx, string(config.DefaultServiceType), "seedProcess")
	_, err = chain.Handle(seedProcessCtx, raw, func(ctx context.Context, req interface{}) (interface{}, error) {
		raw = req.(string)
		err := c.seedProcess(seedProcessCtx, raw)
		return nil, err
	})
	if err == nil {
		return nil
	}

	switch err {
	case core.InterceptError: // 拦截, 应该立即结束本次任务
		return nil
	default:
		utils.Trace.TraceEvent(ctx, "PutErrorRawSeed", utils.Trace.AttrKey("isParserErr").Bool(err == core.ParserError))
		c.app.Warn(ctx, "将出错seed放入error队列", zap.Bool("isParserErr", err == core.ParserError))
		if err := c.PutErrorRawSeed(ctx, raw, err == core.ParserError); err != nil {
			utils.Trace.TraceErrEvent(ctx, "PutErrorRawSeed", err)
			c.app.Error(ctx, "将出错seed放入error队列失败", zap.Error(err))
		} else {
			utils.Trace.TraceEvent(ctx, "PutErrorRawSeedOk")
			c.app.Info(ctx, "已将出错seed原始数据放入error队列")
		}
	}
	return err
}

// 种子处理
func (c *Crawler) seedProcess(ctx context.Context, raw string) error {
	//ctx = utils.Trace.TraceStart(ctx, "seedProcess")
	//defer utils.Trace.TraceEnd(ctx)

	var seedResult *core.Seed
	var cookieJar http.CookieJar
	// 循环尝试下载
	var attempt int
	for {
		attempt++

		// 每次重新生成seed, 因为每次处理可能会修改seed
		seed, err := seeds.MakeSeedOfRaw(raw)
		if err != nil {
			utils.Trace.TraceErrEvent(ctx, "MakeSeedOfRaw", err)
			c.app.Error("构建种子失败")
			return core.ParserError
		}

		tCtx := utils.Trace.TraceStart(ctx, "downloadLoop", utils.Trace.AttrKey("attempt").Int(attempt))

		// 请求处理
		seed, err = c.middleware.RequestProcess(tCtx, c, seed)
		if err != nil {
			utils.Trace.TraceEnd(tCtx)
			return err
		}

		seedResult, cookieJar, err = c.download(tCtx, seed)
		if err == nil {
			utils.Trace.TraceEnd(tCtx)
			break
		}

		utils.Trace.TraceErrEvent(tCtx, "download", err)
		utils.Trace.TraceEnd(tCtx)

		if err == core.InterceptError || err == core.ParserError {
			utils.Trace.TraceErrEvent(ctx, "download", err)
			return err
		}
		if attempt >= c.conf.Frame.RequestMaxAttemptCount {
			utils.Trace.TraceErrEvent(ctx, "download", err)
			return fmt.Errorf("尝试下载失败, 超过最大尝试次数: %v", err)
		}
		c.app.Error(tCtx, "尝试下载失败, 等待重试", zap.Int64("waitTime", c.conf.Frame.EmptyQueueWaitTime),
			zap.String("attempt", fmt.Sprintf("%d/%d", attempt, c.conf.Frame.RequestMaxAttemptCount)), zap.Error(err))
		time.Sleep(time.Duration(c.conf.Frame.RequestRetryWaitTime) * time.Millisecond)
	}

	// 保存cookieJar
	c.cookieJar = cookieJar
	defer func() {
		c.cookieJar = nil
	}()

	pCtx := utils.Trace.TraceStart(ctx, "Parser")
	defer utils.Trace.TraceEnd(pCtx)

	// 解析
	utils.Trace.TraceEvent(pCtx, "Parser")
	err := zapputils.Recover.WrapCall(func() error {
		return c.Parser(pCtx, seedResult)
	})
	if err == nil {
		return nil
	}

	utils.Trace.TraceErrEvent(pCtx, "Parser", err)
	c.app.Error(pCtx, "解析时出错", zap.String("err", zapputils.Recover.GetRecoverErrorDetail(err)))

	// 尝试将body保存到队列
	utils.Trace.TraceEvent(pCtx, "trySaveParserErrorSeed")
	if err := c.trySaveParserErrorSeed(context.Background(), raw, seedResult.HttpResponseBody); err != nil {
		utils.Trace.TraceErrEvent(pCtx, "trySaveParserErrorSeed", err)
		c.app.Error("尝试保存解析错误的seed失败, 只能放入原始数据", zap.String("err", zapputils.Recover.GetRecoverErrorDetail(err)))
	} else {
		return core.InterceptError // 既然保存成功则拦截处理
	}

	return core.ParserError
}

// 下载完善种子
func (c *Crawler) download(ctx context.Context, seed *core.Seed) (*core.Seed, http.CookieJar, error) {
	ctx = utils.Trace.TraceStart(ctx, "download")
	defer utils.Trace.TraceEnd(ctx)

	cookieJar, _ := cookiejar.New(nil)
	raw := seed.Raw

	// 下载
	seed, err := c.downloader.Download(ctx, c, seed, cookieJar)
	if err != nil {
		rawBase64 := utils.Convert.Base64Encode(raw)
		utils.Trace.TraceErrEvent(ctx, "download", err,
			utils.Trace.AttrKey("raw").String(rawBase64))
		return nil, nil, err
	}

	// 响应处理
	seed, err = c.middleware.ResponseProcess(ctx, c, seed)
	if err != nil {
		return nil, nil, err
	}

	// 检查是符合期望的响应
	seed, err = c.CheckIsExpectResponse(ctx, seed)
	if err != nil {
		utils.Trace.TraceErrEvent(ctx, "CheckIsExpectResponse", err)
		return nil, nil, err
	}

	return seed, cookieJar, nil
}

// 尝试保存解析错误的seed
func (c *Crawler) trySaveParserErrorSeed(ctx context.Context, raw string, body []byte) error {
	seed, err := seeds.MakeSeedOfRaw(raw)
	if err != nil {
		return fmt.Errorf("构建种子失败: %v", err)
	}
	seed.HttpResponseBody = body
	err = c.PutErrorSeed(ctx, seed, true)
	return err
}
