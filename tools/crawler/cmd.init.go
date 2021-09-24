package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/tools/utils"
)

// 初始化工程
func CmdInit(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个工程名")
	}
	projectName := context.Args().Get(0)
	utils.MustNoExistPath(projectName)

	utils.MustMkdir(projectName)
	// 进入工程目录
	if err := os.Chdir(projectName); err != nil {
		logger.Log.Fatal("进入工程目录", zap.String("projectName", projectName), zap.Error(err))
	}

	utils.MustMkdir("component")
	utils.MustWriteFile("component/component.go", utils.MustReadEmbedFile(embedFiles, "template/component/component.go.template"))

	utils.MustMkdir("configs")
	utils.MustWriteFile("configs/base_config.toml", utils.MustReadEmbedFile(embedFiles, "template/configs/base_config.toml"))
	utils.MustWriteFile("configs/supervisor_programs.toml", utils.MustReadEmbedFile(embedFiles, "template/configs/supervisor_programs.toml"))

	utils.MustMkdir("supervisor_config")
	supervisorSchedulerConfig := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/supervisor_config/scheduler_config.ini")), utils.MakeTemplateArgs(projectName))
	utils.MustWriteFile("supervisor_config/scheduler_config.ini", []byte(supervisorSchedulerConfig))

	utils.MustMkdir("template")
	utils.MustWriteFile("template/supervisor_program.ini", utils.MustReadEmbedFile(embedFiles, "template/template/supervisor_program.ini"))

	utils.MustMkdir("spiders")

	goModContent := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/go.mod.template")), utils.MakeTemplateArgs(projectName))
	utils.MustWriteFile("go.mod", []byte(goModContent))
	return nil
}
