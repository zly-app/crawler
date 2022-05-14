package main

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/tools/utils"
)

//go:embed embed_assets/*
var embedFiles embed.FS

func embedFilesRelease(projectName, basePath string) {
	templateArgs := utils.MakeTemplateArgs(projectName)

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
		dirs, err := embedFiles.ReadDir(path)
		if err != nil {
			logger.Log.Fatal("读取目录资源失败", zap.String("path", path), zap.Error(err))
		}
		utils.MustMkdir(strings.TrimPrefix(path, basePath)[1:], 666)
		dispatchDirs(path, dirs)
	}
	releaseFile = func(path string, file fs.DirEntry) {
		path = path + "/" + file.Name()
		bs, err := embedFiles.ReadFile(path)
		if err != nil {
			logger.Log.Fatal("读取文件资源失败", zap.String("path", path), zap.Error(err))
		}

		if strings.HasSuffix(path, ".file") {
			path = strings.TrimSuffix(path, ".file")
		} else if strings.HasSuffix(path, ".template") {
			path = strings.TrimSuffix(path, ".template")
			bs = []byte(zstr.Render(string(bs), templateArgs))
		} else if strings.HasSuffix(path, ".t") {
			path += "emplate"
		}

		utils.MustWriteFile(strings.TrimPrefix(path, basePath)[1:], bs, 666)
	}

	dirs, err := embedFiles.ReadDir(basePath)
	if err != nil {
		logger.Log.Fatal("读取资源失败", zap.String("path", basePath), zap.Error(err))
	}
	dispatchDirs(basePath, dirs)
}

// 初始化工程
func CmdInit(context *cli.Context) error {
	if context.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个工程名")
	}
	projectName := context.Args().Get(0)

	if projectName == "." {
		workdir := utils.MustGetWorkdir()
		utils.DirMustEmpty(workdir)
		projectName = utils.MustGetDirName(workdir)
	} else {
		if utils.CheckHasPath(projectName, true) {
			utils.DirMustEmpty(projectName)
		} else {
			utils.MustMkdir(projectName)
		}
		// 进入工程目录
		if err := os.Chdir(projectName); err != nil {
			logger.Log.Fatal("进入工程目录", zap.String("projectName", projectName), zap.Error(err))
		}
	}

	embedFilesRelease(projectName, "embed_assets")
	utils.MustMkdir("spiders")
	return nil
}
