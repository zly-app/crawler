package utils

import (
	"reflect"
	"runtime"
	"strings"
)

var Reflect = new(reflectUtil)

type reflectUtil struct{}

// GetFuncName 获取函数或方法的名称
func (*reflectUtil) GetFuncName(a interface{}) string {
	aValue := reflect.ValueOf(a)
	if aValue.Kind() != reflect.Func {
		panic("a must a func")
	}

	rawName := runtime.FuncForPC(aValue.Pointer()).Name()
	name := strings.TrimSuffix(rawName, ".func1")
	ss := strings.Split(name, ".")
	name = strings.TrimSuffix(ss[len(ss)-1], "-fm")
	return name
}
