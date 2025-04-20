package lru

import (
	"container/list"
)

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}
type Item struct {
	Key   string
	Value string
}

type CustomLRUCache struct {
	capacity int
	queue    *list.List
	items    map[string]*list.Element
}

func NewLruCache(capacity int) CustomLRUCache {
	return CustomLRUCache{
		capacity: capacity,
		queue:    list.New(),
		items:    make(map[string]*list.Element),
	}
}

func (c *CustomLRUCache) Put(key, value string) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value.(*Item).Value = value
		return
	}
	if c.queue.Len() == c.capacity {
		c.clear()
	}
	item := c.queue.PushFront(&Item{Key: key, Value: value})
	c.items[key] = item
	return
}

func (c *CustomLRUCache) Get(key string) (string, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value.(*Item).Value, true
	}
	return "", false
}

func (c *CustomLRUCache) clear() {
	item := c.queue.Back()
	c.queue.Remove(item)
	delete(c.items, item.Value.(*Item).Key)
}
