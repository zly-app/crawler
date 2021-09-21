package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"

	"github.com/zly-app/crawler/tools/utils"
)

// 在当前位置初始化工程
func CmdInit(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个工程名")
	}
	projectName := context.Args().Get(0)
	utils.MustNoExistPath(projectName)

	utils.MustMkdir(fmt.Sprintf("%s/component", projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/component/component.go", projectName), utils.MustReadEmbedFile(embedFiles, "template/component.go.template"))

	utils.MustMkdir(fmt.Sprintf("%s/configs", projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/spiders.toml", projectName), utils.MustReadEmbedFile(embedFiles, "template/spiders.toml"))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/base_config.toml", projectName), utils.MustReadEmbedFile(embedFiles, "template/base_config.toml"))
	supervisorSchedulerConfig := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/supervisor_scheduler_config.ini")), utils.MakeTemplateArgs(projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/supervisor_scheduler_config.ini", projectName), []byte(supervisorSchedulerConfig))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/supervisor_spider_config.ini", projectName), utils.MustReadEmbedFile(embedFiles, "template/supervisor_spider_config.ini"))

	utils.MustMkdir(fmt.Sprintf("%s/spiders", projectName))

	goModContent := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/go.mod.template")), utils.MakeTemplateArgs(projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/go.mod", projectName), []byte(goModContent))
	return nil
}
