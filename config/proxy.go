package config

import (
	"strings"
)

const (
	// 默认代理类型
	defaultProxyType = "direct"
)

type ProxyConfig struct {
	/**代理类型:
	  direct 直接的, 不使用代理
	  static 静态代理, 支持 http, https, socks5, socks5h
	*/
	Type string
}

func newProxyConfig() ProxyConfig {
	return ProxyConfig{}
}
func (conf *ProxyConfig) Check() error {
	switch strings.ToLower(conf.Type) {
	case "", "direct":
		conf.Type = defaultProxyType
	}
	return nil
}
