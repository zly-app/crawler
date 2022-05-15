package crawler

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/zly-app/crawler/core"
	"github.com/zly-app/crawler/seeds"
)

// 弹出一个种子
func (c *Crawler) PopARawSeed() (string, error) {
	for _, suffix := range c.conf.Frame.QueueSuffixes {
		queueName := c.conf.Frame.Namespace + c.conf.Spider.Name + suffix
		raw, err := c.queue.Pop(queueName, true)
		if err == core.EmptyQueueError { // 这个队列为空
			continue
		}
		if err != nil {
			return "", err
		}

		if raw == SubmitInitialSeedSignal {
			return raw, nil
		}
		c.app.Info("从队列取出一个种子", zap.String("queueName", queueName))
		return raw, nil
	}
	return "", core.EmptyQueueError
}

/*
**放入种子
 seed 种子
 front 是否放在队列前面
*/
func (c *Crawler) PutSeed(seed *core.Seed, front bool) error {
	data, err := seeds.EncodeSeed(seed)
	if err != nil {
		return fmt.Errorf("seed编码失败: %v", err)
	}

	return c.PutRawSeed(data, seed.ParserMethod, front)
}

/*
**放入种子原始数据
 raw 种子原始数据
 parserFuncName 解析函数名
 front 是否放在队列前面
*/
func (c *Crawler) PutRawSeed(raw string, parserFuncName string, front bool) error {
	queueName := c.conf.Frame.Namespace + c.conf.Spider.Name + c.conf.Frame.SeedQueueSuffix
	size, err := c.queue.Put(queueName, raw, front)
	if err != nil {
		return fmt.Errorf("将seed放入队列失败: %v", err)
	}

	if raw == SubmitInitialSeedSignal {
		return nil
	}

	c.app.Info("将seed放入队列", zap.String("parserFuncName", parserFuncName), zap.Int("queueSize", size))
	return nil
}

/*
**放入错误种子
 seed 种子
 isParserError 是否为解析错误
*/
func (c *Crawler) PutErrorSeed(seed *core.Seed, isParserError bool) error {
	data, err := seeds.EncodeSeed(seed)
	if err != nil {
		return fmt.Errorf("seed编码失败: %s", err)
	}

	return c.PutErrorRawSeed(data, isParserError)
}

/*
**放入错误种子原始数据
 raw 种子原始数据
 isParserError 是否为解析错误
*/
func (c *Crawler) PutErrorRawSeed(raw string, isParserError bool) error {
	c.app.Warn("将出错seed放入error队列")
	var queueName string
	if isParserError {
		queueName = c.conf.Frame.Namespace + c.conf.Spider.Name + c.conf.Frame.ParserErrorSeedQueueSuffix
	} else {
		queueName = c.conf.Frame.Namespace + c.conf.Spider.Name + c.conf.Frame.ErrorSeedQueueSuffix
	}

	_, err := c.queue.Put(queueName, raw, false)
	if err != nil {
		return err
	}
	return nil
}

// 检查队列是否为空, 如果spiderName为空则取默认值
func (c *Crawler) CheckQueueIsEmpty(spiderName string) (bool, error) {
	if spiderName == "" {
		spiderName = c.conf.Spider.Name
	}

	for _, suffix := range c.conf.Frame.QueueSuffixes {
		if c.conf.Frame.CheckEmptyQueueIgnoreErrorQueue {
			if suffix == c.conf.Frame.ErrorSeedQueueSuffix || suffix == c.conf.Frame.ParserErrorSeedQueueSuffix {
				continue
			}
		}
		queueName := c.conf.Frame.Namespace + spiderName + suffix
		size, err := c.queue.QueueSize(queueName)
		if err != nil {
			return false, err
		}
		if size > 0 {
			return false, nil
		}
	}
	return true, nil
}
