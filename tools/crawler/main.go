package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:  "crawler工程管理",
		Usage: "用于管理你的爬虫组",
		Commands: []*cli.Command{
			{
				Name:      "init",
				Usage:     "初始化工程",
				ArgsUsage: "<your_project_name>",
				Action:    CmdInit,
			},
			{
				Name:      "create",
				Aliases:   []string{"cs"},
				Usage:     "创建一个爬虫",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdCreateSpider,
			},
			{
				Name:      "start",
				Aliases:   []string{"ss"},
				Usage:     "发送提交初始化种子信号",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdInitSeedSignal,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "dev",
						Usage: "环境, 将会读取 configs/crawler_config.{@env}.yaml, spiders/{@spider_name}/configs/config.{@env}.yaml 文件",
					},
				},
			},
			{
				Name:      "clean",
				Usage:     "* 清空爬虫所有队列 - 慎用",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdCleanSpiderQueue,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "dev",
						Usage: "环境, 将会读取 configs/crawler_config.{@env}.yaml, spiders/{@spider_name}/configs/config.{@env}.yaml 文件",
					},
				},
			},
			{
				Name:      "clean_set",
				Usage:     "* 清空爬虫集合数据 - 慎用",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdCleanSpiderSet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "dev",
						Usage: "环境, 将会读取 configs/crawler_config.{@env}.yaml, spiders/{@spider_name}/configs/config.{@env}.yaml 文件",
					},
				},
			},
			{
				Name:      "make_supervisor",
				Aliases:   []string{"make"},
				Usage:     "删除supervisor配置后根据模板文件 template/supervisor/spider_programs.{@env}.ini 重新生成supervisor配置, 生成的文件路径为 supervisor_config/conf.d.{@env}/{@spider_name}.ini",
				ArgsUsage: " ",
				Action:    CmdMakeSupervisorConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "dev",
						Usage: "环境, 将会读取 configs/supervisor/spider_programs.{@env}.yaml, template/supervisor/spider_programs.{@env}.ini 文件",
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Log.Fatal("启动失败", zap.Error(err))
	}
}
