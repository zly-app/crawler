package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"
	"go.uber.org/zap"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/tools/utils"
)

// 生成supervisor配置
func CmdMakeSupervisorConfig(context *cli.Context) error {
	projectName := utils.MustGetProjectName()
	vi := viper.New()
	vi.SetConfigFile("configs/spiders.toml")
	if err := vi.MergeInConfig(); err != nil {
		logger.Log.Fatal("读取爬虫配置文件失败", zap.String("configFile", "configs/spiders.toml"), zap.Error(err))
	}

	var groups map[string]map[string]string
	if err := vi.Unmarshal(&groups); err != nil {
		logger.Log.Fatal("解析调度器配置文件失败", zap.String("configFile", "configs/scheduler.toml"), zap.Error(err))
	}

	// 获取supervisor爬虫配置文件模板
	s, err := os.ReadFile("configs/supervisor_spider_config.ini")
	if err != nil {
		logger.Log.Fatal("读取supervisor爬虫配置文件模板失败", zap.String("template", "configs/supervisor_spider_config.ini"), zap.Error(err))
	}
	spiderConfigTemplate := string(s)

	// supervisor组配置文件模板
	const groupConfigTemplate = `
[group:@group_name]
programs = @spider_names`

	// 删除目录
	err = os.RemoveAll("configs/supervisor")
	if err != nil {
		logger.Log.Fatal("删除目录失败", zap.String("dir", "configs/supervisor"), zap.Error(err))
	}

	// 创建配置目录
	utils.MustMkdir("configs/supervisor")

	for groupName, g := range groups {
		var spiderConfigs []string
		var spiderNames []string
		for spiderName, conf := range g {
			if !utils.CheckHasPath(fmt.Sprintf("./spiders/%s", spiderName), true) {
				continue // 可能在别的机器上部署
			}
			// 解析配置
			confValue := strings.Split(conf, ",")
			if len(confValue) != 2 {
				logger.Log.Fatal("spider的配置错误", zap.String("spiderName", spiderName))
			}
			processNum, err := strconv.Atoi(confValue[0])
			if err != nil {
				logger.Log.Fatal("spider的配置错误, 无法获取到进程数", zap.String("spiderName", spiderName), zap.Error(err))
			}

			// 检查进程数
			if processNum < 1 {
				continue
			}
			if processNum > 99 {
				logger.Log.Fatal("spider的process太多, 超过99通常是无意义的", zap.String("spiderName", spiderName), zap.Int("processNum", processNum))
			}
			// 检查提交初始化种子的时机
			expression := confValue[1]
			switch expression {
			case "", "none", "start":
			default:
				_, err = time.ParseInLocation(crawler.OnceTriggerTimeLayout, expression, time.Local)
				if err != nil {
					_, err = cron.ParseStandard(expression)
				}
				if err != nil {
					logger.Log.Fatal("spider的配置错误, 提交初始化种子的时机无法解析", zap.String("spiderName", spiderName))
				}
			}

			spiderNames = append(spiderNames, spiderName)
			templateArgs := utils.MakeTemplateArgs(projectName)
			templateArgs["group_name"] = groupName                                  // 组名
			templateArgs["spider_name"] = spiderName                                // 爬虫名
			templateArgs["spider_dir"] = utils.MustDirJoin("./spiders", spiderName) // 爬虫目录
			templateArgs["process_num"] = processNum                                // 进程数
			templateArgs["seed_cron"] = confValue[1]                                // 初始化种子提交时机
			text := zstr.Render(spiderConfigTemplate, templateArgs)
			spiderConfigs = append(spiderConfigs, text)
		}

		if len(spiderNames) == 0 {
			continue
		}

		templateArgs := utils.MakeTemplateArgs(projectName)
		templateArgs["group_name"] = groupName                        // 组名
		templateArgs["spider_names"] = strings.Join(spiderNames, ",") // 爬虫列表
		groupConfigData := zstr.Render(groupConfigTemplate, templateArgs)
		data := strings.Join(spiderConfigs, "\n\n") + "\n\n\n" + groupConfigData
		err = os.WriteFile(fmt.Sprintf("configs/supervisor/%s.ini", groupName), []byte(data), 0666)
		if err != nil {
			logger.Log.Fatal("写入supervisor配置失败", zap.String("file", fmt.Sprintf("configs/supervisor/%s.ini", groupName)), zap.Error(err))
		}
	}

	logger.Log.Info("生成supervisor配置完毕")
	return nil
}