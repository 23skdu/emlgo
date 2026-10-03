//go:build !js || !wasm

package jit

import (
	"fmt"
	"sync"
)

const (
	maxCacheEntries = 1024
	defaultNumShards = 16
)

type jitCacheEntry struct {
	key    string
	fn     Func
	handle *CodeHandle
}

type cacheShard struct {
	mu         sync.RWMutex
	items      map[string]*jitCacheEntry
	ring       []string
	ringPos    int
	maxEntries int
}

func newCacheShard(maxEntries int) *cacheShard {
	if maxEntries < 1 {
		maxEntries = 1
	}
	return &cacheShard{
		items:      make(map[string]*jitCacheEntry, maxEntries),
		ring:       make([]string, maxEntries),
		maxEntries: maxEntries,
	}
}

type singleflightCall struct {
	wg  sync.WaitGroup
	val Func
	err error
}

type singleflightGroup struct {
	mu sync.Mutex
	m  map[string]*singleflightCall
}

func (g *singleflightGroup) do(key string, fn func() (Func, error)) (Func, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*singleflightCall)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(singleflightCall)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

type jitCache struct {
	shards    []*cacheShard
	numShards int
	sf        singleflightGroup
	maxSize   int
}

var defaultCache = newJITCache(maxCacheEntries)
var compilePreflightHook func()

func newJITCache(maxSize int) *jitCache {
	if maxSize <= 0 {
		maxSize = maxCacheEntries
	}
	nShards := defaultNumShards
	if maxSize < defaultNumShards {
		nShards = 1
	}
	perShard := (maxSize + nShards - 1) / nShards
	c := &jitCache{
		shards:    make([]*cacheShard, nShards),
		numShards: nShards,
		maxSize:   maxSize,
	}
	for i := 0; i < nShards; i++ {
		c.shards[i] = newCacheShard(perShard)
	}
	return c
}

func hashKey(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func (c *jitCache) getShard(key string) *cacheShard {
	if c.numShards == 1 {
		return c.shards[0]
	}
	return c.shards[hashKey(key)%uint32(c.numShards)]
}

func (c *jitCache) get(key string) (Func, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	if entry, ok := shard.items[key]; ok && entry != nil {
		return entry.fn, true
	}
	return nil, false
}

func (c *jitCache) put(key string, fn Func, handle ...*CodeHandle) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	var h *CodeHandle
	if len(handle) > 0 {
		h = handle[0]
	}

	if entry, ok := shard.items[key]; ok {
		shard.ring[shard.ringPos] = key
		shard.ringPos = (shard.ringPos + 1) % shard.maxEntries
		entry.fn = fn
		if h != nil {
			entry.handle = h
		}
		return
	}

	if len(shard.items) >= shard.maxEntries {
		oldKey := shard.ring[shard.ringPos]
		if oldKey != "" {
			if oldEntry, ok := shard.items[oldKey]; ok {
				if oldEntry.handle != nil {
					_ = oldEntry.handle.Free()
				}
				delete(shard.items, oldKey)
			}
		}
	}

	shard.ring[shard.ringPos] = key
	shard.ringPos = (shard.ringPos + 1) % shard.maxEntries
	shard.items[key] = &jitCacheEntry{
		key:    key,
		fn:     fn,
		handle: h,
	}
}

func (c *jitCache) clear() {
	for _, shard := range c.shards {
		shard.mu.Lock()
		for _, entry := range shard.items {
			if entry.handle != nil {
				_ = entry.handle.Free()
			}
		}
		shard.items = make(map[string]*jitCacheEntry, shard.maxEntries)
		for i := range shard.ring {
			shard.ring[i] = ""
		}
		shard.ringPos = 0
		shard.mu.Unlock()
	}
}

// CompileCached compiles an expression and caches the result.
// Expressions with differing whitespace or grouping are canonicalized to avoid duplicate compiles.
// Subsequent calls with the same expression return the cached function without allocating or acquiring write locks.
func CompileCached(expr string) (Func, error) {
	// 1. Fast path: check cache with raw expr
	if fn, ok := defaultCache.get(expr); ok {
		return fn, nil
	}

	// 2. Canonicalize expression via Parse and FormatExpr
	ast, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}
	canon := FormatExpr(ast)

	// Check cache with canonical form
	if canon != expr {
		if fn, ok := defaultCache.get(canon); ok {
			defaultCache.put(expr, fn)
			return fn, nil
		}
	}

	if compilePreflightHook != nil {
		compilePreflightHook()
	}

	// 3. Singleflight compile on canonical key
	fn, err := defaultCache.sf.do(canon, func() (Func, error) {
		// Recheck cache in case another goroutine just compiled it
		if fn, ok := defaultCache.get(canon); ok {
			return fn, nil
		}
		c := NewCompiler()
		fn, handle, compileErr := c.CompileAST(ast)
		if compileErr != nil {
			return nil, compileErr
		}
		defaultCache.put(canon, fn, handle)
		return fn, nil
	})
	if err != nil {
		return nil, err
	}
	if canon != expr {
		defaultCache.put(expr, fn)
	}
	return fn, nil
}

// ClearJITCache clears the JIT expression cache and frees all allocated executable memory.
func ClearJITCache() {
	defaultCache.clear()
}
