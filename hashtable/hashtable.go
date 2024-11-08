package hashtable

import (
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

// HashTable 线程安全的泛型哈希表结构
type HashTable[K comparable, V any] struct {
	buckets    []*bucket[K, V] // 桶数组
	size       atomic.Int64    // 使用原子计数器存储元素数量
	bucketSize int             // 桶的数量
	mu         sync.RWMutex    // 用于扩容的读写锁
	resizing   atomic.Bool     // 标记是否正在进行扩容
}

// bucket 定义了哈希桶结构
type bucket[K comparable, V any] struct {
	entries []entry[K, V]
	mu      sync.RWMutex
}

// entry 定义了键值对结构
type entry[K comparable, V any] struct {
	key   K
	value V
}

// New 创建一个新的哈希表实例
func New[K comparable, V any](initialSize int) *HashTable[K, V] {
	if initialSize < 1 {
		initialSize = 16
	}

	ht := &HashTable[K, V]{
		buckets:    make([]*bucket[K, V], initialSize),
		bucketSize: initialSize,
	}

	for i := 0; i < initialSize; i++ {
		ht.buckets[i] = &bucket[K, V]{
			entries: make([]entry[K, V], 0, 8), // 预分配空间
		}
	}

	return ht
}

// hash 计算给定键的哈希值
func (ht *HashTable[K, V]) hash(key K) int {
	keyStr := fmt.Sprintf("%v", key)
	h := fnv.New32a()
	h.Write([]byte(keyStr))
	ht.mu.RLock()
	bucketSize := ht.bucketSize
	ht.mu.RUnlock()
	return int(h.Sum32()) % bucketSize
}

// Put 向哈希表中插入键值对
func (ht *HashTable[K, V]) Put(key K, value V) {
	retry := true
	for retry {
		index := ht.hash(key)
		if index >= len(ht.buckets) {
			// 如果索引超出范围，等待扩容完成后重试
			continue
		}

		bucket := ht.buckets[index]
		bucket.mu.Lock()

		// 检查key是否已存在
		updated := false
		for i := range bucket.entries {
			if bucket.entries[i].key == key {
				bucket.entries[i].value = value
				updated = true
				retry = false
				break
			}
		}

		if !updated {
			// 添加新的键值对
			bucket.entries = append(bucket.entries, entry[K, V]{
				key:   key,
				value: value,
			})
			bucket.mu.Unlock()

			// 增加计数并检查是否需要扩容
			newSize := ht.size.Add(1)
			if float64(newSize)/float64(ht.bucketSize) > 0.75 {
				ht.tryResize()
			}
			retry = false
		} else {
			bucket.mu.Unlock()
		}
	}
}

// Get 从哈希表中获取值
func (ht *HashTable[K, V]) Get(key K) (V, bool) {
	retry := true
	var result V
	var found bool

	for retry {
		index := ht.hash(key)
		if index >= len(ht.buckets) {
			continue
		}

		bucket := ht.buckets[index]
		bucket.mu.RLock()

		for _, e := range bucket.entries {
			if e.key == key {
				result = e.value
				found = true
				retry = false
				break
			}
		}

		bucket.mu.RUnlock()
		if !found {
			retry = false
		}
	}

	return result, found
}

// Delete 从哈希表中删除键值对
func (ht *HashTable[K, V]) Delete(key K) bool {
	retry := true
	deleted := false

	for retry {
		index := ht.hash(key)
		if index >= len(ht.buckets) {
			continue
		}

		bucket := ht.buckets[index]
		bucket.mu.Lock()

		for i, e := range bucket.entries {
			if e.key == key {
				// 删除找到的条目
				bucket.entries = append(bucket.entries[:i], bucket.entries[i+1:]...)
				deleted = true
				ht.size.Add(-1)
				retry = false
				break
			}
		}

		bucket.mu.Unlock()
		if !deleted {
			retry = false
		}
	}

	return deleted
}

// tryResize 尝试扩容哈希表
func (ht *HashTable[K, V]) tryResize() {
	// 如果已经在扩容，直接返回
	if !ht.resizing.CompareAndSwap(false, true) {
		return
	}

	ht.mu.Lock()
	defer func() {
		ht.mu.Unlock()
		ht.resizing.Store(false)
	}()

	// 再次检查是否需要扩容
	currentSize := ht.size.Load()
	if float64(currentSize)/float64(ht.bucketSize) <= 0.75 {
		return
	}

	newSize := ht.bucketSize * 2
	newBuckets := make([]*bucket[K, V], newSize)

	// 初始化新桶
	for i := 0; i < newSize; i++ {
		newBuckets[i] = &bucket[K, V]{
			entries: make([]entry[K, V], 0, 8),
		}
	}

	// 重新哈希所有现有的键值对
	for _, oldBucket := range ht.buckets {
		oldBucket.mu.Lock()
		entries := make([]entry[K, V], len(oldBucket.entries))
		copy(entries, oldBucket.entries)
		oldBucket.mu.Unlock()

		for _, e := range entries {
			// 计算新的哈希值
			h := fnv.New32a()
			h.Write([]byte(fmt.Sprintf("%v", e.key)))
			newIndex := int(h.Sum32()) % newSize

			// 将条目放入新桶
			newBucket := newBuckets[newIndex]
			newBucket.mu.Lock()
			newBucket.entries = append(newBucket.entries, e)
			newBucket.mu.Unlock()
		}
	}

	// 更新哈希表状态
	ht.buckets = newBuckets
	ht.bucketSize = newSize
}

// Size 返回哈希表中的元素数量
func (ht *HashTable[K, V]) Size() int {
	return int(ht.size.Load())
}
