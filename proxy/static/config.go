package static

import (
	"errors"
)

type ProxyConfig struct {
	Address  string // 代理地址
	User     string // 用户名
	Password string // 密码
}

func newProxyConfig() *ProxyConfig {
	return &ProxyConfig{}
}

func (conf *ProxyConfig) Check() error {
	if conf.Address == "" {
		return errors.New("代理的address是空的")
	}
	return nil
}
