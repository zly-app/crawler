package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zlyuancn/zstr"

	"github.com/zly-app/crawler/tools/utils"
)

// 在当前位置初始化工程
func CmdInit(context *cli.Context) error {
	if context.Args().Len() != 1 {
		return errors.New("必须也只能写入一个工程名")
	}
	projectName := context.Args().Get(0)
	utils.MustNoExistPath(projectName)

	utils.MustMkdir(fmt.Sprintf("%s/component", projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/component/component.go", projectName), utils.MustReadEmbedFile(embedFiles, "template/component.go.template"))

	utils.MustMkdir(fmt.Sprintf("%s/configs", projectName))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/scheduler.toml", projectName), utils.MustReadEmbedFile(embedFiles, "template/scheduler.toml"))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/spider_config.toml", projectName), utils.MustReadEmbedFile(embedFiles, "template/spider_config.toml"))
	utils.MustWriteFile(fmt.Sprintf("%s/configs/supervisor_spider_config.ini", projectName), utils.MustReadEmbedFile(embedFiles, "template/supervisor_spider_config.ini"))

	utils.MustMkdir(fmt.Sprintf("%s/spiders", projectName))

	goModContent := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/go.mod.template")), zstr.KV{"project_name", projectName})
	utils.MustWriteFile(fmt.Sprintf("%s/go.mod", projectName), []byte(goModContent))
	return nil
}
