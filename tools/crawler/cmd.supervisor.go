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
	"go.uber.org/zap"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/tools/utils"
)

// 生成supervisor配置
func CmdMakeSupervisorConfig(cl *cli.Context) error {
	projectName := utils.MustEnterProject()

	// 环境
	env := cl.String("env")
	if env == "" {
		logger.Log.Fatal("env为空")
	}
	configFile := fmt.Sprintf("./configs/supervisor_programs.%s.yaml", env)
	templateFile := fmt.Sprintf("./template/supervisor_programs.%s.ini.template", env)

	vi := viper.New()
	vi.SetConfigFile(configFile)
	if err := vi.MergeInConfig(); err != nil {
		logger.Log.Fatal("读取supervisor程序配置文件失败", zap.String("configFile", configFile), zap.Error(err))
	}

	var groups map[string]map[string]string
	if err := vi.Unmarshal(&groups); err != nil {
		logger.Log.Fatal("解析supervisor程序配置文件失败", zap.String("configFile", configFile), zap.Error(err))
	}

	// 读取supervisor程序配置文件模板
	s, err := os.ReadFile(templateFile)
	if err != nil {
		logger.Log.Fatal("读取supervisor程序配置文件模板失败", zap.String("template", templateFile), zap.Error(err))
	}
	spiderConfigTemplate := string(s)

	// supervisor组配置文件模板
	const groupConfigTemplate = `
; 分组
[group:{@group_name}]
; 分组的 spider 列表
programs = {@spider_names}`

	// 删除目录
	path := fmt.Sprintf("supervisor_config/conf.d.%s", env)
	err = os.RemoveAll(path)
	if err != nil {
		logger.Log.Fatal("删除目录失败", zap.String("dir", path), zap.Error(err))
	}

	// 创建配置目录
	utils.MustMkdir(path)

	for groupName, g := range groups {
		var spiderConfigs []string
		var spiderNames []string
		for spiderName, conf := range g {
			if !utils.CheckHasPath(fmt.Sprintf("./spiders/%s", spiderName), true) {
				logger.Log.Warn("spider未找到, 将跳过", zap.String("spiderName", spiderName))
				continue // 可能在别的机器上部署
			}
			// 解析配置
			confValue := strings.Split(conf, ",")
			if len(confValue) == 1 { // 允许不填写调度时机
				confValue = append(confValue, "")
			}
			if len(confValue) != 2 {
				logger.Log.Fatal("spider的配置错误", zap.String("spiderName", spiderName))
			}
			processNum, err := strconv.Atoi(confValue[0])
			if err != nil {
				logger.Log.Fatal("spider的配置错误, 无法获取到进程数", zap.String("spiderName", spiderName), zap.Error(err))
			}
			seedCron := confValue[1]

			// 检查进程数
			if processNum < 1 {
				continue
			}
			if processNum > 99 {
				logger.Log.Warn("spider的process太多, 超过99通常是无意义的", zap.String("spiderName", spiderName), zap.Int("processNum", processNum))
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
			templateArgs := utils.MakeTemplateArgs(projectName, env)
			templateArgs["group_name"] = groupName                                  // 组名
			templateArgs["spider_name"] = spiderName                                // 爬虫名
			templateArgs["spider_dir"] = utils.MustDirJoin("./spiders", spiderName) // 爬虫目录
			templateArgs["process_num"] = processNum                                // 进程数
			templateArgs["seed_cron"] = seedCron                                    // 初始化种子提交时机
			text := utils.RenderTemplate(spiderConfigTemplate, templateArgs)
			spiderConfigs = append(spiderConfigs, text)
		}

		if len(spiderNames) == 0 {
			continue
		}

		templateArgs := utils.MakeTemplateArgs(projectName, env)
		templateArgs["group_name"] = groupName                        // 组名
		templateArgs["spider_names"] = strings.Join(spiderNames, ",") // 爬虫列表
		groupConfigData := utils.RenderTemplate(groupConfigTemplate, templateArgs)
		data := strings.Join(spiderConfigs, "\n\n") + "\n\n" + groupConfigData + "\n\n"
		file := fmt.Sprintf("%s/%s.ini", path, groupName)
		err = os.WriteFile(file, []byte(data), 0666)
		if err != nil {
			logger.Log.Fatal("写入supervisor配置失败", zap.String("file", file), zap.Error(err))
		}
	}

	// 调度器配置
	templateFile = fmt.Sprintf("template/scheduler_config.ini.%s.template", env)
	s, err = os.ReadFile(templateFile)
	if err != nil {
		logger.Log.Fatal("读取调度器程序配置文件模板失败", zap.String("template", templateFile), zap.Error(err))
	}
	templateArgs := utils.MakeTemplateArgs(projectName, env)
	text := utils.RenderTemplate(string(s), templateArgs)
	file := fmt.Sprintf("%s/crawler_scheduler.ini", path)
	err = os.WriteFile(file, []byte(text), 0666)
	if err != nil {
		logger.Log.Fatal("写入调度器程序配置失败", zap.String("file", file), zap.Error(err))
	}

	logger.Log.Info("生成supervisor配置完毕")
	return nil
}
