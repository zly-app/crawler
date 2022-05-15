package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"

	"github.com/zly-app/crawler/tools/utils"
)

// 创建一个爬虫
func CmdCreateSpider(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个爬虫名")
	}
	projectName := utils.MustEnterProject()
	spiderName := context.Args().Get(0)

	utils.MustMkdirAndIsCreate(fmt.Sprintf("spiders/%s", spiderName))

	templateArgs := utils.MakeTemplateArgs(projectName, "dev")
	templateArgs["spider_name"] = spiderName

	// main.go
	mainGoContent := string(utils.MustReadFile("template/spider/main.go.template"))
	mainGoContent = utils.RenderTemplate(mainGoContent, templateArgs)
	utils.MustWriteFile(fmt.Sprintf("spiders/%s/main.go", spiderName), []byte(mainGoContent))

	// configs
	utils.MustMkdir(fmt.Sprintf("spiders/%s/configs", spiderName))
	spiderDefaultConfigContent := utils.RenderTemplate(string(utils.MustReadFile("template/spider/configs/config.toml")), templateArgs)
	utils.MustWriteFile(fmt.Sprintf("spiders/%s/configs/config.toml", spiderName), []byte(spiderDefaultConfigContent))
	fmt.Println("创建成功")
	return nil
}
