package hw04lrucache

import "sync"

type Key string

type Element struct {
	Key   Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.Lock()
	defer lru.Unlock()
	if item, exist := lru.items[key]; exist {
		lru.queue.MoveToFront(item)
		item.Value.(*Element).Value = value
		return true
	}

	if lru.queue.Len() == lru.capacity {
		lru.clearLast()
	}

	element := &Element{
		Key:   key,
		Value: value,
	}
	item := lru.queue.PushFront(element)
	lru.items[key] = item

	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.Lock()
	defer lru.Unlock()
	item, exist := lru.items[key]
	if !exist {
		return nil, false
	}
	lru.queue.MoveToFront(item)

	return item.Value.(*Element).Value, true
}

func (lru *lruCache) Clear() {
	lru.Lock()
	defer lru.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}

func (lru *lruCache) clearLast() {
	if item := lru.queue.Back(); item != nil {
		lru.queue.Remove(item)
		delete(lru.items, item.Value.(*Element).Key)
	}
}
