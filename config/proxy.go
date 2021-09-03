package config

import (
	"fmt"
	"strings"
)

const (
	// 默认代理类型
	defaultProxyType = "direct"
)

type ProxyConfig struct {
	/**代理类型:
	  direct 直接的, 不使用代理
	  http http或https
	  socks5 sock5
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
	default:
		return fmt.Errorf("不支持的代理类型: %s", conf.Type)
	}
	return nil
}
