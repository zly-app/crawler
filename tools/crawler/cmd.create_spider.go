package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"

	"github.com/zly-app/crawler/tools/utils"
)

// 创建一个爬虫
func CmdCreateSpider(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个爬虫名")
	}
	projectName := utils.MustGetProjectName()
	spiderName := context.Args().Get(0)

	utils.MustMkdirAndIsCreate(fmt.Sprintf("spiders/%s", spiderName))
	mainGoContent := strings.ReplaceAll(string(utils.MustReadEmbedFile(embedFiles, "template/template/spider/main.go.template")), "{@spider_name}", spiderName)
	mainGoContent = strings.ReplaceAll(mainGoContent, "{@project_name}", projectName)
	utils.MustWriteFile(fmt.Sprintf("spiders/%s/main.go", spiderName), []byte(mainGoContent))

	utils.MustMkdir(fmt.Sprintf("spiders/%s/configs", spiderName))
	templateArgs := utils.MakeTemplateArgs(projectName)
	templateArgs["spider_name"] = spiderName
	spiderDefaultConfigContent := zstr.Render(string(utils.MustReadEmbedFile(embedFiles, "template/template/spider/config.toml")), templateArgs)
	utils.MustWriteFile(fmt.Sprintf("spiders/%s/configs/default.toml", spiderName), []byte(spiderDefaultConfigContent))
	return nil
}
