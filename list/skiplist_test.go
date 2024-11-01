package list

import (
	"math/rand"
	"testing"
	"time"
)

// 比较函数
func intCmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// TestSkipListNewSkipList 测试创建新跳表
func TestSkipListNewSkipList(t *testing.T) {
	skipList := NewSkipList(intCmp)
	if skipList == nil {
		t.Fatal("NewSkipList返回了nil")
	}

	if skipList.level != 1 {
		t.Errorf("新创建的跳表层级应为1，实际为%d", skipList.level)
	}

	if skipList.header == nil {
		t.Fatal("跳表的头节点不应为nil")
	}

	if len(skipList.header.next) != MaxLevel {
		t.Errorf("头节点的next数组长度应为%d，实际为%d", MaxLevel, len(skipList.header.next))
	}
}

// TestSkipListInsert 测试插入操作
func TestSkipListInsert(t *testing.T) {
	skipList := NewSkipList(intCmp)

	// 测试插入有序数据
	t.Run("Ordered Insert", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		for _, v := range values {
			skipList.Insert(v)
		}

		// 验证所有值都能找到
		for _, v := range values {
			if result := skipList.Search(v); result == nil {
				t.Errorf("未找到已插入的值: %d", v)
			}
		}
	})

	// 测试插入重复数据
	t.Run("Duplicate Insert", func(t *testing.T) {
		skipList := NewSkipList(intCmp)
		skipList.Insert(1)
		skipList.Insert(1)

		// 验证重复值可以正确存储和查找
		if result := skipList.Search(1); result == nil {
			t.Error("未找到重复插入的值")
		}
	})

	// 测试随机插入
	t.Run("Random Insert", func(t *testing.T) {
		skipList := NewSkipList(intCmp)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		nums := make(map[int]bool)

		// 插入100个随机数
		for i := 0; i < 100; i++ {
			value := r.Intn(1000)
			skipList.Insert(value)
			nums[value] = true
		}

		// 验证所有插入的数都能找到
		for num := range nums {
			if result := skipList.Search(num); result == nil {
				t.Errorf("未找到随机插入的值: %d", num)
			}
		}
	})
}

// TestSkipListSearch 测试查找操作
func TestSkipListSearch(t *testing.T) {
	skipList := NewSkipList(intCmp)

	// 测试空跳表查找
	t.Run("Empty List Search", func(t *testing.T) {
		if result := skipList.Search(1); result != nil {
			t.Error("空跳表查找应该返回nil")
		}
	})

	// 测试查找存在的值
	t.Run("Existing Value Search", func(t *testing.T) {
		values := []int{1, 3, 5, 7, 9}
		for _, v := range values {
			skipList.Insert(v)
		}

		for _, v := range values {
			if result := skipList.Search(v); result == nil || *result != v {
				t.Errorf("查找已存在的值 %d 失败", v)
			}
		}
	})

	// 测试查找不存在的值
	t.Run("Non-existing Value Search", func(t *testing.T) {
		nonExisting := []int{2, 4, 6, 8}
		for _, v := range nonExisting {
			if result := skipList.Search(v); result != nil {
				t.Errorf("查找不存在的值 %d 应该返回nil", v)
			}
		}
	})
}

// TestSkipListDelete 测试删除操作
func TestSkipListDelete(t *testing.T) {
	skipList := NewSkipList(intCmp)

	// 测试空跳表删除
	t.Run("Empty List Delete", func(t *testing.T) {
		if skipList.Delete(1) {
			t.Error("从空跳表删除应该返回false")
		}
	})

	// 测试删除存在的值
	t.Run("Existing Value Delete", func(t *testing.T) {
		values := []int{1, 3, 5, 7, 9}
		for _, v := range values {
			skipList.Insert(v)
		}

		// 删除并验证
		for _, v := range values {
			if !skipList.Delete(v) {
				t.Errorf("删除已存在的值 %d 失败", v)
			}
			if result := skipList.Search(v); result != nil {
				t.Errorf("删除后仍能查找到值 %d", v)
			}
		}
	})

	// 测试删除不存在的值
	t.Run("Non-existing Value Delete", func(t *testing.T) {
		if skipList.Delete(100) {
			t.Error("删除不存在的值应该返回false")
		}
	})

	// 测试重复删除
	t.Run("Duplicate Delete", func(t *testing.T) {
		skipList.Insert(1)
		if !skipList.Delete(1) {
			t.Error("首次删除应该返回true")
		}
		if skipList.Delete(1) {
			t.Error("重复删除应该返回false")
		}
	})
}

// TestSkipListLevelDistribution 测试层级分布
func TestSkipListLevelDistribution(t *testing.T) {
	skipList := NewSkipList(intCmp)
	levels := make(map[int]int)

	// 生成大量随机层级并统计分布
	for i := 0; i < 10000; i++ {
		level := skipList.randomLevel()
		levels[level]++

		if level < 1 || level > MaxLevel {
			t.Errorf("生成的层级 %d 超出有效范围 [1, %d]", level, MaxLevel)
		}
	}

	// 验证层级分布符合概率要求
	// P(level = k) = P * P(level = k-1)
	for i := 2; i <= MaxLevel; i++ {
		ratio := float64(levels[i]) / float64(levels[i-1])
		// 允许10%的误差
		if ratio > Probability*1.1 || ratio < Probability*0.9 {
			t.Logf("警告：第 %d 层的比例 %.2f 与期望值 %.2f 相差较大", i, ratio, Probability)
		}
	}
}

// TestSkipListDifferentTypes 测试不同数据类型
func TestSkipListDifferentTypes(t *testing.T) {
	// 测试字符串类型
	t.Run("String Type", func(t *testing.T) {
		stringCmp := func(a, b string) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		}

		skipList := NewSkipList(stringCmp)
		words := []string{"apple", "banana", "cherry"}

		// 插入并验证
		for _, word := range words {
			skipList.Insert(word)
		}

		// 查找验证
		for _, word := range words {
			if result := skipList.Search(word); result == nil || *result != word {
				t.Errorf("未找到已插入的字符串: %s", word)
			}
		}
	})

	// 测试浮点数类型
	t.Run("Float Type", func(t *testing.T) {
		floatCmp := func(a, b float64) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		}

		skipList := NewSkipList(floatCmp)
		nums := []float64{3.14, 2.718, 1.414}

		// 插入并验证
		for _, num := range nums {
			skipList.Insert(num)
		}

		// 查找验证
		for _, num := range nums {
			if result := skipList.Search(num); result == nil || *result != num {
				t.Errorf("未找到已插入的浮点数: %f", num)
			}
		}
	})
}

// TestSkipListOrderMaintenance 测试顺序维护
func TestSkipListOrderMaintenance(t *testing.T) {
	skipList := NewSkipList(intCmp)
	values := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}

	// 乱序插入
	for _, v := range values {
		skipList.Insert(v)
	}

	// 验证最底层的节点是否有序
	current := skipList.header.next[0]
	for current != nil && current.next[0] != nil {
		currentValue := current.value
		nextValue := current.next[0].value
		if intCmp(currentValue, nextValue) >= 0 {
			t.Errorf("顺序错误：%v >= %v", currentValue, nextValue)
		}
		current = current.next[0]
	}
}
