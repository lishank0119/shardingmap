package shardingmap

import (
	"sync"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	sm := New[string, int]()

	// Test Set & Get
	sm.Set("key1", 100)
	value, ok := sm.Get("key1")
	if !ok || value != 100 {
		t.Errorf("expected 100, got %v", value)
	}

	// Test Update
	sm.Set("key1", 200)
	value, _ = sm.Get("key1")
	if value != 200 {
		t.Errorf("expected updated value 200, got %v", value)
	}

	// Test Delete
	sm.Delete("key1")
	_, ok = sm.Get("key1")
	if ok {
		t.Errorf("expected key1 to be deleted")
	}
}

func TestLen(t *testing.T) {
	sm := New[string, int]()
	sm.Set("key1", 1)
	sm.Set("key2", 2)
	sm.Set("key3", 3)

	if sm.Len() != 3 {
		t.Errorf("expected length 3, got %d", sm.Len())
	}

	sm.Delete("key2")
	if sm.Len() != 2 {
		t.Errorf("expected length 2 after deletion, got %d", sm.Len())
	}
}

func TestForEach(t *testing.T) {
	sm := New[string, int]()
	sm.Set("a", 1)
	sm.Set("b", 2)
	sm.Set("c", 3)

	total := 0
	sm.ForEach(func(k string, v int) {
		total += v
	})

	if total != 6 {
		t.Errorf("expected total 6, got %d", total)
	}
}

func TestCustomShardingFunc(t *testing.T) {
	sm := New[string, int](
		WithShardCount[string, int](4),
		WithShardingFunc[string, int](func(key string) int {
			return len(key) // shard based on key length
		}),
	)

	sm.Set("short", 1)
	sm.Set("longerkey", 2)

	// Check same key always maps to the same shard
	if sm.shardingFunc("short")%4 != sm.shardingFunc("short")%4 {
		t.Errorf("same key should map to the same shard")
	}

	// Check distribution of different keys
	keys := []string{"a", "bb", "ccc", "dddd", "eeeee"}
	distributedShards := make(map[int]bool)
	for _, key := range keys {
		distributedShards[sm.shardingFunc(key)%4] = true
	}

	if len(distributedShards) < 2 {
		t.Errorf("custom sharding function should distribute keys across multiple shards")
	}
}

func TestConcurrency(t *testing.T) {
	sm := New[int, int]()
	wg := sync.WaitGroup{}

	n := 1000
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			sm.Set(i, i*10)
			if val, ok := sm.Get(i); !ok || val != i*10 {
				t.Errorf("concurrent set/get failed for key %d", i)
			}
		}(i)
	}

	wg.Wait()

	if sm.Len() != n {
		t.Errorf("expected length %d after concurrent writes, got %d", n, sm.Len())
	}
}
