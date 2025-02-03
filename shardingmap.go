package shardingmap

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Shard defines the structure of each shard
type Shard[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

// ShardingMap defines the structure of the entire sharding map
type ShardingMap[K comparable, V any] struct {
	shards       []*Shard[K, V]
	shardCount   int
	shardingFunc func(K) int
}

// Option defines the options for configuring the ShardingMap
type Option[K comparable, V any] func(*ShardingMap[K, V])

// WithShardCount sets the number of shards
func WithShardCount[K comparable, V any](count int) Option[K, V] {
	return func(sm *ShardingMap[K, V]) {
		if count > 0 {
			sm.shardCount = count
		}
	}
}

// WithShardingFunc sets a custom sharding function
func WithShardingFunc[K comparable, V any](fn func(K) int) Option[K, V] {
	return func(sm *ShardingMap[K, V]) {
		if fn != nil {
			sm.shardingFunc = fn
		}
	}
}

// defaultShardingFunc is the default sharding function using FNV hash
func defaultShardingFunc[K comparable](key K) int {
	h := fnv.New32a()
	if _, err := h.Write([]byte(fmt.Sprintf("%v", key))); err != nil {
		panic(fmt.Sprintf("unexpected error hashing key: %v", err))
	}
	return int(h.Sum32())
}

// New creates a new ShardingMap with support for options
func New[K comparable, V any](opts ...Option[K, V]) *ShardingMap[K, V] {
	sm := &ShardingMap[K, V]{
		shardCount:   16, // Default to 16 shards
		shardingFunc: defaultShardingFunc[K],
	}

	for _, opt := range opts {
		opt(sm)
	}

	sm.shards = make([]*Shard[K, V], sm.shardCount)
	for i := 0; i < sm.shardCount; i++ {
		sm.shards[i] = &Shard[K, V]{items: make(map[K]V)}
	}

	return sm
}

// getShard returns the corresponding shard based on the key
func (sm *ShardingMap[K, V]) getShard(key K) *Shard[K, V] {
	index := sm.shardingFunc(key) % sm.shardCount
	return sm.shards[index]
}

// Set sets the key-value pair
func (sm *ShardingMap[K, V]) Set(key K, value V) {
	shard := sm.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// Get retrieves the value corresponding to the key
func (sm *ShardingMap[K, V]) Get(key K) (V, bool) {
	shard := sm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// Delete removes the key
func (sm *ShardingMap[K, V]) Delete(key K) {
	shard := sm.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// Len returns the total number of elements in the ShardingMap
func (sm *ShardingMap[K, V]) Len() int {
	total := 0
	for _, shard := range sm.shards {
		shard.RLock()
		total += len(shard.items)
		shard.RUnlock()
	}
	return total
}

// ForEach iterates over all key-value pairs
func (sm *ShardingMap[K, V]) ForEach(f func(K, V)) {
	for _, shard := range sm.shards {
		shard.RLock()
		for k, v := range shard.items {
			f(k, v)
		}
		shard.RUnlock()
	}
}
