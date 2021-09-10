package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/zlyuancn/zstr"
)

// 生成supervisor配置
func CmdMakeSupervisorConfig(context *cli.Context) error {
	projectName := MustGetProjectName()
	vi := viper.New()
	vi.SetConfigFile("configs/scheduler.toml")
	if err := vi.MergeInConfig(); err != nil {
		return fmt.Errorf("读取configs/scheduler.toml文件失败: %v", err)
	}

	var groups map[string]map[string]struct {
		Process int    // 进程数
		Seed    string // 提交初始化种子的时机
		Desc    string // spider描述
	}
	if err := vi.Unmarshal(&groups); err != nil {
		return fmt.Errorf("解析configs/scheduler.toml文件失败: %v", err)
	}

	// 获取supervisor爬虫配置文件模板
	s, err := os.ReadFile("configs/supervisor_spider_config.ini")
	if err != nil {
		panic(err)
	}
	spiderConfigTemplate := string(s)

	// supervisor组配置文件模板
	const groupConfigTemplate = `
[group:@group_name]
programs = @spider_names`

	// 删除目录
	err = os.RemoveAll("configs/supervisor")
	if err != nil {
		panic(err)
	}

	// 创建配置目录
	MustMkdir("configs/supervisor")

	for groupName, g := range groups {
		var spiderConfigs []string
		var spiderNames []string
		for spiderName, conf := range g {
			if !CheckHasPath(fmt.Sprintf("./spiders/%s", spiderName), true) {
				panic(fmt.Errorf("spider<%s>不存在", spiderName))
			}
			if conf.Process == 0 {
				continue
			}
			if conf.Process > 99 {
				panic(fmt.Errorf("spider<%s>的process太多, 超过99通常是无意义的", spiderName))
			}
			spiderNames = append(spiderNames, spiderName)

			args := map[string]interface{}{
				"project_name": projectName,                          // 项目名
				"group_name":   groupName,                            // 组名
				"spider_name":  spiderName,                           // 爬虫名
				"spider_dir":   MustDirJoin("./spiders", spiderName), // 爬虫目录
				"process":      conf.Process,                         // 进程数
				"seed":         conf.Seed,                            // 初始化种子提交时机
				"desc":         conf.Desc,                            // 爬虫描述
			}
			text := zstr.Render(spiderConfigTemplate, args)
			spiderConfigs = append(spiderConfigs, text)
		}

		if len(spiderNames) == 0 {
			continue
		}

		args := map[string]interface{}{
			"group_name":   groupName,                      // 组名
			"spider_names": strings.Join(spiderNames, ","), // 爬虫列表
		}
		groupConfigData := zstr.Render(groupConfigTemplate, args)
		data := strings.Join(spiderConfigs, "\n\n") + "\n\n\n" + groupConfigData
		err := os.WriteFile(fmt.Sprintf("configs/supervisor/%s.ini", groupName), []byte(data), 0666)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("生成supervisor配置完毕")
	return nil
}
