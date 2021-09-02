package memory

import (
	"container/list"
	"sync"

	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/crawler/core"
)

type MemoryQueue struct {
	queues map[string]*list.List
	mx     sync.Mutex
}

func (m *MemoryQueue) Put(queueName string, raw string, front bool) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	queue, ok := m.queues[queueName]
	if !ok {
		queue = list.New()
		m.queues[queueName] = queue
	}

	if front {
		queue.PushFront(raw)
	} else {
		queue.PushBack(raw)
	}

	return queue.Len(), nil
}

func (m *MemoryQueue) Pop(queueName string, front bool) (string, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	queue, ok := m.queues[queueName]
	if !ok || queue.Len() == 0 {
		return "", core.EmptyQueueError
	}

	var element *list.Element
	if front {
		element = queue.Front()
	} else {
		element = queue.Back()
	}

	raw := queue.Remove(element).(string)
	return raw, nil
}

func (m *MemoryQueue) QueueSize(queueName string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	queue, ok := m.queues[queueName]
	if !ok {
		return 0, nil
	}

	return queue.Len(), nil
}

func (m *MemoryQueue) Close() error { return nil }

func NewMemoryQueue(app zapp_core.IApp) core.IQueue {
	return &MemoryQueue{
		queues: make(map[string]*list.List),
	}
}
