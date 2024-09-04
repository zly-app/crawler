package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/zly-app/plugin/honey"
	"github.com/zly-app/plugin/prometheus"
	"github.com/zly-app/plugin/zipkinotel"
	"github.com/zly-app/zapp"
	zapp_config "github.com/zly-app/zapp/config"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/service/cron"

	"github.com/zly-app/crawler"
	"github.com/zly-app/crawler/config"
	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/tools/utils"
)

type Scheduler struct {
	app zapp_core.IApp
	cr  core.ICrawler
}

// 启动
func (s *Scheduler) Start() {
	if strings.ToLower(config.Conf.Queue.Type) == "memory" {
		s.app.Fatal("使用memory队列是无意义的")
	}

	var programsConfigFile string
	err := s.app.GetConfig().Parse("crawler_scheduler.spider_programs_file", &programsConfigFile)
	if err != nil {
		s.app.Fatal("获取程序配置文件路径失败", zap.Error(err))
	}

	vi := viper.New()
	vi.SetConfigFile(programsConfigFile)
	if err := vi.MergeInConfig(); err != nil {
		s.app.Fatal("读取程序配置文件失败", zap.String("programsConfigFile", programsConfigFile), zap.Error(err))
	}

	var groups map[string]map[string]string
	if err := vi.Unmarshal(&groups); err != nil {
		s.app.Fatal("解析程序配置文件失败", zap.String("programsConfigFile", programsConfigFile), zap.Error(err))
	}

	for _, g := range groups {
		for spiderName, conf := range g {
			// 这里不需要检查爬虫是否存在, 只需要发信号就行了
			// 解析配置
			confValue := strings.Split(conf, ",")
			if len(confValue) == 1 { // 允许不填写调度时机
				confValue = append(confValue, "")
			}
			if len(confValue) != 2 {
				panic(fmt.Errorf("spider<%s>的配置错误", spiderName))
			}
			processNum, err := strconv.Atoi(confValue[0])
			if err != nil {
				panic(fmt.Errorf("spider<%s>的配置错误, 无法获取到进程数", spiderName))
			}
			// 检查进程数
			if processNum < 1 {
				continue
			}
			if processNum > 99 {
				panic(fmt.Errorf("spider<%s>的process太多, 超过99通常是无意义的", spiderName))
			}

			// 检查提交初始化种子的时机
			expression := confValue[1]
			switch expression {
			case "", "none", "start": // 不使用调度器
				continue
			default:
				t, err := time.ParseInLocation(crawler.OnceTriggerTimeLayout, expression, time.Local)
				if err == nil {
					cron.RegistryTask(cron.NewTaskOfConfig(spiderName, cron.TaskConfig{
						Trigger:  cron.NewOnceTrigger(t),
						Executor: cron.NewExecutor(2, time.Second, 1),
						Handler:  s.SendSubmitInitialSeedSignal,
						Enable:   true,
					}))
					continue
				}
				cron.RegistryTask(cron.NewTaskOfConfig(spiderName, cron.TaskConfig{
					Trigger:  cron.NewCronTrigger(expression),
					Executor: cron.NewExecutor(2, time.Second, 1),
					Handler:  s.SendSubmitInitialSeedSignal,
					Enable:   true,
				}))
			}
		}
	}
}

// 发送提交初始化种子信号
func (s *Scheduler) SendSubmitInitialSeedSignal(ctx cron.IContext) error {
	spiderName := ctx.Task().Name()

	// 检查非空队列不提交初始化种子
	if config.Conf.Frame.StopSubmitInitialSeedIfNotEmptyQueue {
		empty, err := s.cr.CheckQueueIsEmpty(context.Background(), spiderName)
		if err != nil {
			return fmt.Errorf("检查队列是否为空失败, spiderName: %s, err: %v", spiderName, err)
		}
		if !empty {
			ctx.Info("队列非空忽略初始化种子提交", zap.String("spiderName", spiderName))
			return nil
		}
	}

	// 放入提交初始化种子信号到队列
	queueName := config.Conf.Frame.Namespace + "." + spiderName + config.Conf.Frame.SeedQueueSuffix
	_, err := s.cr.Queue().Put(context.Background(), queueName, crawler.SubmitInitialSeedSignal, true)
	if err != nil {
		return fmt.Errorf("提交初始化种子信号放入到队列失败, spiderName: %s, err: %v", spiderName, err)
	}

	ctx.Info("发送提交初始化种子信号成功", zap.String("spiderName", spiderName))
	return nil
}

func main() {
	utils.MustEnterProject()

	app := zapp.NewApp("crawler-scheduler",
		cron.WithService(),
		zapp.WithConfigOption(zapp_config.WithFiles("./configs/crawler.dev.yaml")),

		zipkinotel.WithPlugin(), // trace
		honey.WithPlugin(),      // log
		prometheus.WithPlugin(), // metrics
	)

	s := &Scheduler{
		app: app,
		cr:  crawler.NewCrawler(app),
	}
	s.Start()

	app.Run()
}
