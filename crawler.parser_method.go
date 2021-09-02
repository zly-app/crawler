package crawler

import (
	"fmt"
	"reflect"

	"github.com/zly-app/crawler/core"
)

// 检查处理程序, 反射获取spider中符合条件的方法
func (c *Crawler) CheckSpiderParserMethod() {
	aType := reflect.TypeOf(c.spider)
	aValue := reflect.ValueOf(c.spider)

	c.parserMethods = make(map[string]core.ParserMethod)
	for i := 0; i < aType.NumMethod(); i++ {
		method := aType.Method(i)
		if method.PkgPath != "" {
			continue
		}
		methodValue := aValue.Method(i)
		if parserMethod, ok := methodValue.Interface().(core.ParserMethod); ok {
			c.parserMethods[method.Name] = parserMethod
		}
	}

	if len(c.parserMethods) == 0 {
		c.app.Fatal("spider不存在任何解析方法,它将无法处理seed")
	}
}

// 检查是否为期望的响应
func (c *Crawler) CheckIsExpectResponse(seed *core.Seed) (*core.Seed, error) {
	if seed.CheckExpectMethod == "" {
		return seed, nil
	}

	checkMethod := c.parserMethods[seed.CheckExpectMethod]
	if err := checkMethod(seed); err != nil {
		return nil, fmt.Errorf("非预期的响应: %v", err)
	}
	return seed, nil
}

// 解析
func (c *Crawler) Parser(seed *core.Seed) error {
	parserMethod := c.parserMethods[seed.ParserMethod] // 解析方法一定存在, 请求中间件已经检查了
	return parserMethod(seed)
}

// 获取解析方法
func (c *Crawler) GetSpiderParserMethod(methodName string) (core.ParserMethod, bool) {
	method, ok := c.parserMethods[methodName]
	return method, ok
}
