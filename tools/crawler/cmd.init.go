package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zlyuancn/zstr"
)

// 在当前位置初始化工程
func CmdInit(context *cli.Context) error {
	if context.Args().Len() != 1 {
		return errors.New("必须也只能写入一个工程名")
	}
	projectName := context.Args().Get(0)
	MustNoExistPath(projectName)

	MustMkdir(fmt.Sprintf("%s/component", projectName))
	MustWriteFile(fmt.Sprintf("%s/component/component.go", projectName), MustReadEmbedFile("template/component.go.template"))

	MustMkdir(fmt.Sprintf("%s/configs", projectName))
	MustWriteFile(fmt.Sprintf("%s/configs/scheduler.toml", projectName), MustReadEmbedFile("template/scheduler.toml"))
	MustWriteFile(fmt.Sprintf("%s/configs/spider_base.toml", projectName), MustReadEmbedFile("template/spider_base.toml"))

	MustMkdir(fmt.Sprintf("%s/spiders", projectName))

	goModContent := zstr.Render(string(MustReadEmbedFile("template/go.mod.template")), zstr.KV{"project_name", projectName})
	MustWriteFile(fmt.Sprintf("%s/go.mod", projectName), []byte(goModContent))
	MustWriteFile(fmt.Sprintf("%s/go.sum", projectName), MustReadEmbedFile("template/go.sum.template"))
	return nil
}
