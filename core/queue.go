package core

import (
	"context"
	"errors"
)

// 队列为空错误
var EmptyQueueError = errors.New("empty queue")

// 队列
type IQueue interface {
	/*
		**将种子原始数据放入队列
		 queueName 队列名
		 seed 种子
		 front 是否放在队列前面
		*return 队列长度
	*/
	Put(ctx context.Context, queueName string, raw string, front bool) (int, error)
	/*
		** 弹出一个种子原始数据
		 queueName 队列名
		 front 是否从队列前面弹出
		*return 种子原始数据
	*/
	Pop(ctx context.Context, queueName string, front bool) (string, error)
	// 获取队列长度
	QueueSize(ctx context.Context, queueName string) (int, error)
	// 删除队列
	Delete(ctx context.Context, queueName string) error
	// 关闭
	Close(ctx context.Context) error
}
