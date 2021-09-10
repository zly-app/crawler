package main

import (
	_ "embed"

	"github.com/urfave/cli/v2"
)

var (
	//go:embed configs/scheduler.toml
	schedulerConfigContent string // 调度器配置内容
	//go:embed configs/spider_base.toml
	spiderBaseConfigContent string // 爬虫继承的配置内容
)

// 在当前位置初始化工程
func CmdInit(context *cli.Context) error {
	DirMustEmpty(".")
	MustMkdir(0666, "configs", "spiders")
	MustWriteFile("configs/spider_base.toml", []byte(spiderBaseConfigContent), 0666)
	MustWriteFile("configs/scheduler.toml", []byte(schedulerConfigContent), 0666)
	return nil
}
