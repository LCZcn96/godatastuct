package hashtable

import (
	"fmt"
	"sync"
	"testing"
)

// TestBasicOperations 测试基本的CRUD操作
func TestBasicOperations(t *testing.T) {
	// 创建新的哈希表实例
	ht := New[string, int](16)

	// 测试插入操作
	t.Run("Put操作测试", func(t *testing.T) {
		ht.Put("one", 1)
		ht.Put("two", 2)
		ht.Put("three", 3)

		if size := ht.Size(); size != 3 {
			t.Errorf("期望大小为3, 实际为 %d", size)
		}
	})

	// 测试获取操作
	t.Run("Get操作测试", func(t *testing.T) {
		// 测试存在的键
		if val, exists := ht.Get("one"); !exists || val != 1 {
			t.Errorf("期望值为1, 实际为 %d, exists = %v", val, exists)
		}

		// 测试不存在的键
		if _, exists := ht.Get("nonexistent"); exists {
			t.Error("不存在的键不应该返回存在")
		}
	})

	// 测试更新操作
	t.Run("更新操作测试", func(t *testing.T) {
		ht.Put("one", 100)
		if val, _ := ht.Get("one"); val != 100 {
			t.Errorf("更新后期望值为100, 实际为 %d", val)
		}
	})

	// 测试删除操作
	t.Run("Delete操作测试", func(t *testing.T) {
		// 删除存在的键
		if !ht.Delete("two") {
			t.Error("删除存在的键应该返回true")
		}

		// 确认键已被删除
		if _, exists := ht.Get("two"); exists {
			t.Error("已删除的键不应该存在")
		}

		// 删除不存在的键
		if ht.Delete("nonexistent") {
			t.Error("删除不存在的键应该返回false")
		}
	})
}

// TestDifferentTypes 测试不同类型的键值对
func TestDifferentTypes(t *testing.T) {
	// 测试整数类型键
	t.Run("整数类型键测试", func(t *testing.T) {
		ht := New[int, string](8)
		ht.Put(1, "one")
		ht.Put(2, "two")

		if val, _ := ht.Get(1); val != "one" {
			t.Errorf("期望值为'one', 实际为 %s", val)
		}
	})

	// 测试浮点数类型键
	t.Run("浮点数类型键测试", func(t *testing.T) {
		ht := New[float64, bool](8)
		ht.Put(1.1, true)
		ht.Put(2.2, false)

		if val, _ := ht.Get(1.1); !val {
			t.Error("期望值为true")
		}
	})

	// 测试自定义结构体值
	t.Run("结构体值测试", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		ht := New[string, Person](8)
		ht.Put("alice", Person{"Alice", 25})

		if val, _ := ht.Get("alice"); val.Name != "Alice" || val.Age != 25 {
			t.Errorf("期望值为 {Alice 25}, 实际为 %v", val)
		}
	})
}

// TestResizing 测试哈希表的自动扩容
func TestResizing(t *testing.T) {
	ht := New[int, int](4) // 从小的初始大小开始

	// 插入足够多的元素触发扩容
	for i := 0; i < 10; i++ {
		ht.Put(i, i*i)
	}

	// 验证所有数据在扩容后仍然完整
	for i := 0; i < 10; i++ {
		if val, exists := ht.Get(i); !exists || val != i*i {
			t.Errorf("扩容后数据不完整: key=%d, expected=%d, actual=%d, exists=%v",
				i, i*i, val, exists)
		}
	}
}

// TestEdgeCases 测试边界条件
func TestEdgeCases(t *testing.T) {
	// 测试创建大小为0的哈希表
	t.Run("创建大小为0测试", func(t *testing.T) {
		ht := New[string, int](0)
		ht.Put("test", 1)
		if val, _ := ht.Get("test"); val != 1 {
			t.Error("即使初始大小为0也应该能正常工作")
		}
	})

	// 测试空字符串键
	t.Run("空字符串键测试", func(t *testing.T) {
		ht := New[string, int](8)
		ht.Put("", 100)
		if val, exists := ht.Get(""); !exists || val != 100 {
			t.Error("空字符串键应该能正常工作")
		}
	})

	// 测试零值
	t.Run("零值测试", func(t *testing.T) {
		ht := New[int, int](8)
		ht.Put(0, 0)
		if val, exists := ht.Get(0); !exists || val != 0 {
			t.Error("零值应该能正常存储和获取")
		}
	})
}

// TestConcurrency 测试并发操作
func TestConcurrency(t *testing.T) {
	ht := New[int, int](16)
	var wg sync.WaitGroup
	n := 1000 // 并发操作数量

	// 并发写入
	t.Run("并发写入测试", func(t *testing.T) {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				ht.Put(val, val)
			}(i)
		}
		wg.Wait()

		// 验证写入结果
		count := 0
		for i := 0; i < n; i++ {
			if _, exists := ht.Get(i); exists {
				count++
			}
		}
		if count != n {
			t.Errorf("期望写入 %d 个元素, 实际写入 %d 个", n, count)
		}
	})

	// 并发读取
	t.Run("并发读取测试", func(t *testing.T) {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				ht.Get(val)
			}(i)
		}
		wg.Wait()
	})

	// 并发删除
	t.Run("并发删除测试", func(t *testing.T) {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				ht.Delete(val)
			}(i)
		}
		wg.Wait()

		// 验证删除结果
		if size := ht.Size(); size != 0 {
			t.Errorf("删除后期望大小为0, 实际为 %d", size)
		}
	})
}

// TestPerformance 性能测试
func BenchmarkHashTable(b *testing.B) {
	ht := New[string, int](16)

	// 测试插入性能
	b.Run("Put性能", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i)
			ht.Put(key, i)
		}
	})

	// 测试查询性能
	b.Run("Get性能", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i)
			ht.Get(key)
		}
	})

	// 测试删除性能
	b.Run("Delete性能", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i)
			ht.Delete(key)
		}
	})
}
