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
				Name:   "init",
				Usage:  "在当前位置初始化工程",
				Action: CmdInit,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
