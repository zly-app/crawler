package set

import (
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/set/memory"
)

var setCreator = map[string]func(app zapp_core.IApp) core.ISet{
	"memory": memory.NewMemorySet,
}

func NewSet(app zapp_core.IApp, setType string) core.ISet {
	creator, ok := setCreator[setType]
	if !ok {
		logger.Log.Fatal("set.type 未定义", zap.String("type", setType))
	}
	return creator(app)
}

// 注册集合创造者, 重复注册会报错并结束程序
func RegistrySetCreator(setType string, creator func(app zapp_core.IApp) core.ISet) {
	if _, ok := setCreator[setType]; ok {
		logger.Log.Fatal("重复注册set建造者", zap.String("queueType", setType))
	}
	setCreator[setType] = creator
}
