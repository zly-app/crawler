package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	zapp_core "github.com/zly-app/zapp/core"
)

func NewRedis(app zapp_core.IApp, confKey string) (redis.UniversalClient, error) {
	conf := newRedisConfig()
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		return nil, fmt.Errorf("配置错误: %v", err)
	}

	if conf.IsCluster {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(conf.Address, ","),
			Username:     conf.UserName,
			Password:     conf.Password,
			MinIdleConns: conf.MinIdleConns,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		})
		return client, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         conf.Address,
		Username:     conf.UserName,
		Password:     conf.Password,
		DB:           conf.DB,
		MinIdleConns: conf.MinIdleConns,
		PoolSize:     conf.PoolSize,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
	})
	return client, nil
}
