package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp"
	zapp_config "github.com/zly-app/zapp/config"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/tools/utils"
)

func makeCrawler(context *cli.Context) (core.IApp, *crawler.Crawler, string, error) {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个爬虫名")
	}
	utils.MustEnterProject()
	spiderName := context.Args().Get(0)

	// 环境
	env := context.String("env")
	if env == "" {
		logger.Log.Fatal("env为空")
	}
	configFile := fmt.Sprintf("./configs/spider_base_config.%s.toml", env)
	spiderFile := fmt.Sprintf("./spiders/%s/configs/config.%s.toml", spiderName, env)

	// 检查spider存在
	if !utils.CheckHasPath(fmt.Sprintf("spiders/%s", spiderName), true) {
		logger.Log.Fatal("spider不存在", zap.String("spiderName", spiderName))
	}

	// 通过zapp创建crawler
	app := zapp.NewApp("crawler",
		zapp.WithConfigOption(zapp_config.WithFiles(configFile, spiderFile)),
	)
	c := crawler.NewCrawler(app)

	return app, c, spiderName, nil
}

// 发送提交初始化种子信号
func CmdInitSeedSignal(context *cli.Context) error {
	app, c, spiderName, err := makeCrawler(context)
	if err != nil {
		return err
	}
	defer app.Exit()

	// 内存队列不能发送提交初始化种子信号
	if strings.ToLower(config.Conf.Queue.Type) == "memory" {
		logger.Log.Fatal("使用memory队列是无意义的")
	}

	// 检查非空队列不提交初始化种子
	if config.Conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		empty, err := c.CheckQueueIsEmpty(spiderName)
		if err != nil {
			logger.Log.Fatal("检查队列是否为空失败", zap.Error(err))
		}
		if !empty {
			logger.Log.Info("队列非空忽略初始化种子提交")
			return nil
		}
	}

	// 放入提交初始化种子信号到队列
	queueName := spiderName + config.Conf.Frame.SeedQueueSuffix
	_, err = c.Queue().Put(queueName, crawler.SubmitInitialSeedSignal, true)
	if err != nil {
		logger.Log.Fatal("放入提交初始化种子信号到队列失败", zap.Error(err))
	}

	logger.Log.Info("发送提交初始化种子信号成功", zap.String("spiderName", spiderName))
	return nil
}

// 清空爬虫所有队列
func CmdCleanSpiderQueue(context *cli.Context) error {
	app, c, spiderName, err := makeCrawler(context)
	if err != nil {
		return err
	}
	defer app.Exit()

	// 内存队列不能清空
	if strings.ToLower(config.Conf.Queue.Type) == "memory" {
		logger.Log.Fatal("使用memory队列是无意义的")
	}

	// 包含完整后缀
	suffixes := append([]string{
		config.Conf.Frame.SeedQueueSuffix,
		config.Conf.Frame.ErrorSeedQueueSuffix,
		config.Conf.Frame.ParserErrorSeedQueueSuffix,
	}, config.Conf.Frame.QueueSuffixes...)

	for _, suffix := range suffixes {
		queueName := spiderName + suffix
		if err = c.Queue().Delete(queueName); err != nil {
			logger.Log.Fatal("删除队列失败", zap.String("queueName", queueName), zap.Error(err))
		}
	}
	logger.Log.Info("清空爬虫所有队列成功", zap.String("spiderName", spiderName))
	return nil
}

// 清空爬虫集合数据
func CmdCleanSpiderSet(context *cli.Context) error {
	app, c, spiderName, err := makeCrawler(context)
	if err != nil {
		return err
	}
	defer app.Exit()

	// 内存集合不能清空
	if strings.ToLower(config.Conf.Set.Type) == "memory" {
		logger.Log.Fatal("使用memory集合是无意义的")
	}

	setName := spiderName + config.Conf.Frame.SetSuffix
	if err = c.Set().DeleteSet(setName); err != nil {
		logger.Log.Fatal("删除集合失败", zap.String("setName", setName), zap.Error(err))
	}
	logger.Log.Info("清空爬虫集合数据成功", zap.String("spiderName", spiderName))
	return nil
}
