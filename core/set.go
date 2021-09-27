package core

type ISet interface {
	// 添加一些元素到集合中, 返回添加的数量, 已存在的元素不会计数
	Add(key string, items ...string) (int, error)
	// 判断集合是否包含某个元素
	HasItem(key, item string) (bool, error)
	// 从集合中移除一些元素, 返回成功移除的数量, 元素不存在不会计数也不会报错
	Remove(key string, items ...string) (int, error)
	// 删除set
	DeleteSet(key string) error
	// 获取集合大小
	GetSetSize(key string) (int, error)
	// 关闭
	Close() error
}
