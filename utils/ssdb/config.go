package ssdb

import (
	"errors"
)

const (
	// 默认最小空闲连接数
	defaultMinIdleConns = 1
	// 默认最大连接池个数
	defaultPoolSize = 1
	// 默认读取超时
	defaultReadTimeout = 5000
	// 默认写入超时
	defaultWriteTimeout = 5000
	// 默认连接超时
	defaultDialTimeout = 5000
)

type SsdbConfig struct {
	Address      string // 地址: host1:port1
	Password     string // 密码
	MinIdleConns int    // 最小空闲连接数
	PoolSize     int    // 客户端池大小
	ReadTimeout  int    // 读取超时(毫秒
	WriteTimeout int    // 写入超时(毫秒
	DialTimeout  int    // 连接超时(毫秒
}

func newSsdbConfig() *SsdbConfig {
	return &SsdbConfig{
		MinIdleConns: defaultMinIdleConns,
		PoolSize:     defaultPoolSize,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		DialTimeout:  defaultDialTimeout,
	}
}

func (conf *SsdbConfig) Check() error {
	if conf.Address == "" {
		return errors.New("ssdb的address为空")
	}
	return nil
}
