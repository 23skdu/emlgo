//go:build !js || !wasm
// +build !js !wasm

package jit

import (
	"container/list"
	"sync"
)

const maxCacheEntries = 1024

type jitCacheEntry struct {
	key string
	fn  Func
}

type jitCache struct {
	mu      sync.RWMutex
	items   map[string]*list.Element
	order   *list.List
	maxSize int
}

var defaultCache = newJITCache(maxCacheEntries)

func newJITCache(maxSize int) *jitCache {
	return &jitCache{
		items:   make(map[string]*list.Element),
		order:   list.New(),
		maxSize: maxSize,
	}
}

func (c *jitCache) get(key string) (Func, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*jitCacheEntry).fn, true
	}
	return nil, false
}

func (c *jitCache) put(key string, fn Func) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return
	}
	if c.order.Len() >= c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*jitCacheEntry).key)
		}
	}
	entry := &jitCacheEntry{key: key, fn: fn}
	elem := c.order.PushFront(entry)
	c.items[key] = elem
}

func (c *jitCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*list.Element)
	c.order.Init()
}

// CompileCached compiles an expression and caches the result.
// Subsequent calls with the same expression return the cached function.
func CompileCached(expr string) (Func, error) {
	if fn, ok := defaultCache.get(expr); ok {
		return fn, nil
	}
	c := NewCompiler()
	fn, err := c.Compile(expr)
	if err != nil {
		return nil, err
	}
	defaultCache.put(expr, fn)
	return fn, nil
}

// ClearJITCache clears the JIT expression cache.
func ClearJITCache() {
	defaultCache.clear()
}
