package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "crawler工程管理",
		Usage: "用于管理你的爬虫组",
		Commands: []*cli.Command{
			{
				Name:      "init",
				Usage:     "在当前位置初始化工程",
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
			},
			{
				Name:      "clean",
				Aliases:   []string{"cq"},
				Usage:     "* 清空爬虫所有队列 - 慎用",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdCleanSpiderQueue,
			},
			{
				Name:      "clean_set",
				Aliases:   []string{"cs"},
				Usage:     "* 清空爬虫集合数据 - 慎用",
				ArgsUsage: "<your_spider_name>",
				Action:    CmdCleanSpiderSet,
			},
			{
				Name:      "make_supervisor",
				Aliases:   []string{"make"},
				Usage:     "删除supervisor配置后根据 configs/scheduler.toml 重新生成supervisor配置, 生成的文件路径为 configs/supervisor/*.ini",
				ArgsUsage: " ",
				Action:    CmdMakeSupervisorConfig,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
