package queue

import (
	"testing"
)

// TestNewDeque 测试创建新的双端队列
func TestNewDeque(t *testing.T) {
	deque := NewDeque[int]()
	if deque == nil {
		t.Fatal("NewDeque返回了nil")
	}

	if !deque.IsEmpty() {
		t.Error("新创建的双端队列应该为空")
	}

	if size := deque.Size(); size != 0 {
		t.Errorf("新创建的双端队列大小应为0，实际为%d", size)
	}
}

// TestPushFrontAndBack 测试从前后添加元素
func TestPushFrontAndBack(t *testing.T) {
	deque := NewDeque[int]()

	// 测试PushFront
	testValues := []int{1, 2, 3}
	for _, v := range testValues {
		deque.PushFront(v)

		// 验证元素是否正确添加到队首
		front, err := deque.Front()
		if err != nil {
			t.Errorf("获取队首元素失败: %v", err)
		}
		if front != v {
			t.Errorf("PushFront(%d)后，Front()=%d，期望值为%d", v, front, v)
		}
	}

	// 清空队列
	deque = NewDeque[int]()

	// 测试PushBack
	for _, v := range testValues {
		deque.PushBack(v)

		// 验证元素是否正确添加到队尾
		back, err := deque.Back()
		if err != nil {
			t.Errorf("获取队尾元素失败: %v", err)
		}
		if back != v {
			t.Errorf("PushBack(%d)后，Back()=%d，期望值为%d", v, back, v)
		}
	}
}

// TestPopFrontAndBack 测试从前后移除元素
func TestPopFrontAndBack(t *testing.T) {
	deque := NewDeque[int]()

	// 准备测试数据
	testValues := []int{1, 2, 3}
	for _, v := range testValues {
		deque.PushBack(v)
	}

	// 测试PopFront
	for i, expected := range testValues {
		value, err := deque.PopFront()
		if err != nil {
			t.Errorf("PopFront()失败(第%d次): %v", i+1, err)
		}
		if value != expected {
			t.Errorf("PopFront()=%d，期望值为%d", value, expected)
		}
	}

	// 重新填充数据
	for _, v := range testValues {
		deque.PushBack(v)
	}

	// 测试PopBack
	for i := len(testValues) - 1; i >= 0; i-- {
		value, err := deque.PopBack()
		if err != nil {
			t.Errorf("PopBack()失败(第%d次): %v", i+1, err)
		}
		if value != testValues[i] {
			t.Errorf("PopBack()=%d，期望值为%d", value, testValues[i])
		}
	}
}

// TestEmptyDequeOperations 测试空队列操作
func TestEmptyDequeOperations(t *testing.T) {
	deque := NewDeque[int]()

	// 测试空队列的Front操作
	_, err := deque.Front()
	if err == nil {
		t.Error("从空队列获取Front应该返回错误")
	}

	// 测试空队列的Back操作
	_, err = deque.Back()
	if err == nil {
		t.Error("从空队列获取Back应该返回错误")
	}

	// 测试空队列的PopFront操作
	_, err = deque.PopFront()
	if err == nil {
		t.Error("从空队列PopFront应该返回错误")
	}

	// 测试空队列的PopBack操作
	_, err = deque.PopBack()
	if err == nil {
		t.Error("从空队列PopBack应该返回错误")
	}
}

// TestSizeAndEmpty 测试大小和空状态检查
func TestSizeAndEmpty(t *testing.T) {
	deque := NewDeque[int]()

	// 初始状态检查
	if !deque.IsEmpty() {
		t.Error("新创建的队列应该为空")
	}
	if size := deque.Size(); size != 0 {
		t.Errorf("空队列的Size()=%d，期望值为0", size)
	}

	// 添加元素后检查
	deque.PushBack(1)
	if deque.IsEmpty() {
		t.Error("添加元素后队列不应该为空")
	}
	if size := deque.Size(); size != 1 {
		t.Errorf("添加一个元素后Size()=%d，期望值为1", size)
	}

	// 移除元素后检查
	_, _ = deque.PopBack()
	if !deque.IsEmpty() {
		t.Error("移除所有元素后队列应该为空")
	}
	if size := deque.Size(); size != 0 {
		t.Errorf("移除所有元素后Size()=%d，期望值为0", size)
	}
}

// TestMixedOperations 测试混合操作
func TestMixedOperations(t *testing.T) {
	deque := NewDeque[int]()

	// 执行一系列混合操作
	operations := []struct {
		operation string
		value     int
		expected  int
	}{
		{"PushFront", 1, 1},
		{"PushBack", 2, 2},
		{"PushFront", 3, 3},
		{"PopBack", 0, 2},
		{"PopFront", 0, 3},
		{"PushBack", 4, 4},
	}

	for i, op := range operations {
		switch op.operation {
		case "PushFront":
			deque.PushFront(op.value)
			front, err := deque.Front()
			if err != nil {
				t.Errorf("步骤%d: Front()失败: %v", i, err)
			}
			if front != op.expected {
				t.Errorf("步骤%d: Front()=%d，期望值为%d", i, front, op.expected)
			}
		case "PushBack":
			deque.PushBack(op.value)
			back, err := deque.Back()
			if err != nil {
				t.Errorf("步骤%d: Back()失败: %v", i, err)
			}
			if back != op.expected {
				t.Errorf("步骤%d: Back()=%d，期望值为%d", i, back, op.expected)
			}
		case "PopFront":
			value, err := deque.PopFront()
			if err != nil {
				t.Errorf("步骤%d: PopFront()失败: %v", i, err)
			}
			if value != op.expected {
				t.Errorf("步骤%d: PopFront()=%d，期望值为%d", i, value, op.expected)
			}
		case "PopBack":
			value, err := deque.PopBack()
			if err != nil {
				t.Errorf("步骤%d: PopBack()失败: %v", i, err)
			}
			if value != op.expected {
				t.Errorf("步骤%d: PopBack()=%d，期望值为%d", i, value, op.expected)
			}
		}
	}
}

// TestDifferentTypes 测试不同类型的数据
func TestDifferentTypes(t *testing.T) {
	// 测试字符串类型
	t.Run("String Type", func(t *testing.T) {
		deque := NewDeque[string]()
		deque.PushBack("hello")
		deque.PushBack("world")

		value, err := deque.PopFront()
		if err != nil {
			t.Errorf("PopFront()失败: %v", err)
		}
		if value != "hello" {
			t.Errorf("PopFront()=%s，期望值为hello", value)
		}
	})

	// 测试浮点数类型
	t.Run("Float Type", func(t *testing.T) {
		deque := NewDeque[float64]()
		deque.PushBack(3.14)
		deque.PushBack(2.718)

		value, err := deque.PopBack()
		if err != nil {
			t.Errorf("PopBack()失败: %v", err)
		}
		if value != 2.718 {
			t.Errorf("PopBack()=%f，期望值为2.718", value)
		}
	})

	// 测试布尔类型
	t.Run("Boolean Type", func(t *testing.T) {
		deque := NewDeque[bool]()
		deque.PushFront(true)
		deque.PushFront(false)

		value, err := deque.Front()
		if err != nil {
			t.Errorf("Front()失败: %v", err)
		}
		if value != false {
			t.Error("Front()应该返回false")
		}
	})
}
