package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/tools/utils"
)

// 创建一个爬虫
func CmdCreateSpider(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个爬虫名")
	}
	projectName := utils.MustEnterProject()
	spiderName := context.Args().Get(0)
	spiderBasePath := fmt.Sprintf("spiders/%s", spiderName)
	utils.MustCreateDirOrDirIsEmpty(spiderBasePath, 666)

	templateArgs := utils.MakeTemplateArgs(projectName, "dev")
	templateArgs["spider_name"] = spiderName

	templateBasePath := "template/spider_template"

	var dispatchDirs func(path string, dirs []fs.DirEntry)
	var releaseDir func(path string, dir fs.DirEntry)
	var releaseFile func(path string, file fs.DirEntry)
	dispatchDirs = func(path string, dirs []fs.DirEntry) {
		for _, dir := range dirs {
			if dir.IsDir() {
				releaseDir(path, dir)
			} else {
				releaseFile(path, dir)
			}
		}
	}
	releaseDir = func(path string, dir fs.DirEntry) {
		path = path + "/" + dir.Name()
		dirs, err := os.ReadDir(path)
		if err != nil {
			logger.Log.Fatal("读取目录资源失败", zap.String("path", path), zap.Error(err))
		}
		newPath := utils.MustDirJoin(spiderBasePath, strings.TrimPrefix(path, templateBasePath)[1:])
		fmt.Printf("%s >> %s\n", path, newPath)
		utils.MustCreateDirOrDirIsEmpty(newPath, 666)
		dispatchDirs(path, dirs)
	}
	releaseFile = func(path string, file fs.DirEntry) {
		path = path + "/" + file.Name()
		bs, err := os.ReadFile(path)
		if err != nil {
			logger.Log.Fatal("读取文件资源失败", zap.String("path", path), zap.Error(err))
		}
		path = utils.MustDirJoin(spiderBasePath, strings.TrimPrefix(path, templateBasePath)[1:])

		if strings.HasSuffix(path, ".file") {
			path = strings.TrimSuffix(path, ".file")
		} else if strings.HasSuffix(path, ".template") {
			path = strings.TrimSuffix(path, ".template")
			bs = []byte(utils.RenderTemplate(string(bs), templateArgs))
		}

		utils.MustWriteFile(path, bs, 666)
	}
	dirs, err := os.ReadDir(templateBasePath)
	if err != nil {
		logger.Log.Fatal("读取模板资源失败", zap.String("path", templateBasePath), zap.Error(err))
	}
	dispatchDirs(templateBasePath, dirs)
	fmt.Println("创建成功")
	return nil
}
