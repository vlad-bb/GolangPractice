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

func NewLruCache(capacity int) *CustomLRUCache {
	return &CustomLRUCache{
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
		c.removeOldest()
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

func (c *CustomLRUCache) removeOldest() {
	if c.queue == nil || c.queue.Len() == 0 {
		return
	}
	item := c.queue.Back()
	if item == nil {
		return
	}
	c.queue.Remove(item)
	if item.Value != nil {
		delete(c.items, item.Value.(*Item).Key)
	}
}
