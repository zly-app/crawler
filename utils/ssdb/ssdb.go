package ssdb

import (
	"fmt"
	"net"

	"github.com/seefan/gossdb/v2"
	rconf "github.com/seefan/gossdb/v2/conf"
	"github.com/seefan/gossdb/v2/pool"
	zapp_core "github.com/zly-app/zapp/core"
)

func NewSsdb(app zapp_core.IApp, confKey string) (*pool.Connectors, error) {
	conf := newSsdbConfig()
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		return nil, fmt.Errorf("配置错误: %v", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", conf.Address)
	if err != nil {
		return nil, fmt.Errorf("无法解析addres: %v", err)
	}

	p, err := gossdb.NewPool(&rconf.Config{
		Host:           addr.IP.String(),
		Port:           addr.Port,
		Password:       conf.Password,
		ReadTimeout:    conf.ReadTimeout / 1e3,
		WriteTimeout:   conf.WriteTimeout / 1e3,
		ConnectTimeout: conf.DialTimeout / 1e3,
		MinPoolSize:    conf.MinIdleConns,
		MaxPoolSize:    conf.PoolSize,
		AutoClose:      true,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}
