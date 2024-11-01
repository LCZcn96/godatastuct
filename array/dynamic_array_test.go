package dynamicarray

import (
	"testing"
)

// TestNew 测试创建新的动态数组
func TestNew(t *testing.T) {
	// 创建新的动态数组
	arr := New[int]()

	// 验证初始状态
	if arr.Len() != 0 {
		t.Errorf("期望长度为0, 实际为 %d", arr.Len())
	}
	if arr.Cap() != initialCapacity {
		t.Errorf("期望容量为%d, 实际为 %d", initialCapacity, arr.Cap())
	}
}

// TestAppend 测试添加元素操作
func TestAppend(t *testing.T) {
	arr := New[int]()

	// 测试添加元素
	testCases := []struct {
		value    int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{5, 5}, // 这次添加会触发扩容
	}

	for i, tc := range testCases {
		arr.Append(tc.value)
		if arr.Len() != tc.expected {
			t.Errorf("测试用例 %d: 期望长度为 %d, 实际为 %d", i, tc.expected, arr.Len())
		}

		// 验证最后添加的元素
		if val, err := arr.Get(arr.Len() - 1); err != nil || val != tc.value {
			t.Errorf("测试用例 %d: 期望值为 %d, 实际为 %d", i, tc.value, val)
		}
	}

	// 验证扩容是否正确
	if arr.Cap() != initialCapacity*2 {
		t.Errorf("扩容后期望容量为 %d, 实际为 %d", initialCapacity*2, arr.Cap())
	}
}

// TestInsert 测试插入元素操作
func TestInsert(t *testing.T) {
	arr := New[int]()

	// 准备测试数据
	arr.Append(1)
	arr.Append(3)

	// 测试在中间位置插入
	if err := arr.Insert(1, 2); err != nil {
		t.Errorf("插入元素失败: %v", err)
	}

	// 验证插入后的顺序
	expected := []int{1, 2, 3}
	for i := 0; i < arr.Len(); i++ {
		if val, _ := arr.Get(i); val != expected[i] {
			t.Errorf("位置 %d: 期望值为 %d, 实际为 %d", i, expected[i], val)
		}
	}

	// 测试边界情况
	if err := arr.Insert(-1, 0); err == nil {
		t.Error("期望插入负索引时返回错误")
	}
	if err := arr.Insert(arr.Len()+1, 0); err == nil {
		t.Error("期望插入越界索引时返回错误")
	}
}

// TestRemove 测试删除元素操作
func TestRemove(t *testing.T) {
	arr := New[int]()

	// 准备测试数据
	values := []int{1, 2, 3, 4}
	for _, v := range values {
		arr.Append(v)
	}

	// 测试删除操作
	val, err := arr.Remove(1)
	if err != nil {
		t.Errorf("删除元素失败: %v", err)
	}
	if val != 2 {
		t.Errorf("删除的元素期望为 2, 实际为 %d", val)
	}

	// 验证删除后的顺序
	expected := []int{1, 3, 4}
	for i := 0; i < arr.Len(); i++ {
		if val, _ := arr.Get(i); val != expected[i] {
			t.Errorf("位置 %d: 期望值为 %d, 实际为 %d", i, expected[i], val)
		}
	}

	// 测试边界情况
	if _, err := arr.Remove(-1); err == nil {
		t.Error("期望删除负索引时返回错误")
	}
	if _, err := arr.Remove(arr.Len()); err == nil {
		t.Error("期望删除越界索引时返回错误")
	}
}

// TestGetSet 测试获取和设置元素操作
func TestGetSet(t *testing.T) {
	arr := New[int]()
	arr.Append(1)

	// 测试Get操作
	val, err := arr.Get(0)
	if err != nil || val != 1 {
		t.Errorf("Get操作失败: 期望值为1, 实际为 %d", val)
	}

	// 测试Set操作
	if err := arr.Set(0, 2); err != nil {
		t.Errorf("Set操作失败: %v", err)
	}
	if val, _ := arr.Get(0); val != 2 {
		t.Errorf("Set后Get的值不匹配: 期望值为2, 实际为 %d", val)
	}

	// 测试边界情况
	if _, err := arr.Get(-1); err == nil {
		t.Error("期望获取负索引时返回错误")
	}
	if err := arr.Set(-1, 0); err == nil {
		t.Error("期望设置负索引时返回错误")
	}
}

// TestShrink 测试数组缩容
func TestShrink(t *testing.T) {
	arr := New[int]()

	// 添加足够多的元素触发扩容
	for i := 0; i < 10; i++ {
		arr.Append(i)
	}
	originalCap := arr.Cap()

	// 删除大部分元素触发缩容
	for arr.Len() > 2 {
		arr.Remove(arr.Len() - 1)
	}

	// 验证是否正确缩容
	if arr.Cap() >= originalCap {
		t.Errorf("期望容量减小, 原容量: %d, 现容量: %d", originalCap, arr.Cap())
	}
}
