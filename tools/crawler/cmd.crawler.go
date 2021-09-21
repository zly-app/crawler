package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/tools/utils"
)

func makeCrawler(context *cli.Context) (core.IApp, *crawler.Crawler, string, error) {
	if context.Args().Len() != 1 {
		return nil, nil, "", errors.New("必须也只能写入一个爬虫名")
	}
	utils.MustInProjectDir()
	spiderName := context.Args().Get(0)

	// 检查spider存在
	if !utils.CheckHasPath(fmt.Sprintf("spiders/%s", spiderName), true) {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("spider<%s>不存在\n", spiderName))
		os.Exit(1)
	}

	// 进入spider目录
	if err := os.Chdir(fmt.Sprintf("spiders/%s", spiderName)); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("进入spider<%s>目录失败\n", spiderName))
		os.Exit(1)
	}

	// 通过zapp创建crawler
	app := zapp.NewApp("crawler")
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
		_, _ = os.Stderr.WriteString("使用memory队列是无意义的\n")
		os.Exit(1)
		return nil
	}

	// 检查非空队列不提交初始化种子
	if config.Conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		empty, err := c.CheckQueueIsEmpty(spiderName)
		if err != nil {
			panic(err)
		}
		if !empty {
			fmt.Println("队列非空忽略初始化种子提交")
			return nil
		}
	}

	// 放入提交初始化种子信号到队列
	queueName := spiderName + config.Conf.Frame.SeedQueueSuffix
	_, err = c.Queue().Put(queueName, crawler.SubmitInitialSeedSignal, true)
	if err != nil {
		return err
	}

	fmt.Println(spiderName + ": 发送提交初始化种子信号成功")
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
		_, _ = os.Stderr.WriteString("使用memory队列是无意义的\n")
		os.Exit(1)
		return nil
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
			panic(err)
		}
	}
	fmt.Println(spiderName + ": 清空爬虫所有队列成功")
	return nil
}

// 清空爬虫集合数据
func CmdCleanSpiderSet(context *cli.Context) error {
	app, c, spiderName, err := makeCrawler(context)
	if err != nil {
		return err
	}
	defer app.Exit()

	// 内存队列不能清空
	if strings.ToLower(config.Conf.Set.Type) == "memory" {
		_, _ = os.Stderr.WriteString("使用memory集合是无意义的\n")
		os.Exit(1)
		return nil
	}

	queueName := spiderName + config.Conf.Frame.SetSuffix
	if err = c.Queue().Delete(queueName); err != nil {
		panic(err)
	}
	fmt.Println(spiderName + ": 清空爬虫集合数据成功")
	return nil
}
