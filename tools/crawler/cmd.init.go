package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/tools/utils"
)

//go:embed embed_assets/*
var embedFiles embed.FS

func embedFilesRelease(projectName, basePath string) {
	templateArgs := utils.MakeTemplateArgs(projectName, "dev")

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
		} else if strings.HasSuffix(path, ".tfile") {
			path = strings.TrimSuffix(path, ".tfile") + ".template"
		} else if strings.HasSuffix(path, ".template") {
			path = strings.TrimSuffix(path, ".template")
			bs = []byte(utils.RenderTemplate(string(bs), templateArgs))
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
func CmdInit(cl *cli.Context) error {
	if cl.Args().Len() != 1 {
		logger.Log.Fatal("必须也只能写入一个工程名")
	}
	projectName := cl.Args().Get(0)

	if projectName == "." {
		workdir := utils.MustGetWorkdir()
		utils.DirMustEmpty(workdir)
		projectName = utils.MustGetDirName(workdir)
	} else {
		utils.MustCreateDirOrDirIsEmpty(projectName, 666)
		// 进入工程目录
		if err := os.Chdir(projectName); err != nil {
			logger.Log.Fatal("进入工程目录失败", zap.String("projectName", projectName), zap.Error(err))
		}
	}

	embedFilesRelease(projectName, "embed_assets")
	fmt.Println("初始化成功")
	return nil
}
