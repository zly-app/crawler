package crawler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zly-app/crawler/core"
)

// 扫描爬虫方法, 反射获取spider中符合条件的方法
func (c *Crawler) ScanSpiderMethod() {
	aType := reflect.TypeOf(c.spider)
	aValue := reflect.ValueOf(c.spider)

	c.parserMethods = make(map[string]core.ParserMethod)
	c.checkMethods = make(map[string]core.CheckMethod)
	for i := 0; i < aType.NumMethod(); i++ {
		methodType := aType.Method(i)
		if methodType.PkgPath != "" {
			continue
		}

		if strings.HasPrefix(methodType.Name, core.ParserMethodNamePrefix) {
			methodValue := aValue.Method(i)
			if method, ok := methodValue.Interface().(core.ParserMethod); ok {
				c.parserMethods[methodType.Name] = method
			}
		} else if strings.HasPrefix(methodType.Name, core.CheckMethodNamePrefix) {
			methodValue := aValue.Method(i)
			if method, ok := methodValue.Interface().(core.CheckMethod); ok {
				c.checkMethods[methodType.Name] = method
			}
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

	checkMethod := c.checkMethods[seed.CheckExpectMethod]
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

// 获取spider检查方法
func (c *Crawler) GetSpiderCheckMethod(methodName string) (core.CheckMethod, bool) {
	method, ok := c.checkMethods[methodName]
	return method, ok
}
