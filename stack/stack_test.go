package stack

import (
	"testing"
)

// TestNewStack 测试创建新栈
func TestNewStack(t *testing.T) {
	// 创建新的空栈
	s := New[int]()

	// 验证初始状态
	if !s.IsEmpty() {
		t.Error("新创建的栈应该为空")
	}
	if s.Size() != 0 {
		t.Errorf("新栈的大小应该为0，实际为 %d", s.Size())
	}
}

// TestPush 测试入栈操作
func TestPush(t *testing.T) {
	s := New[string]()
	testCases := []struct {
		value           string
		expectedSize    int
		expectedIsEmpty bool
	}{
		{"第一个元素", 1, false},
		{"第二个元素", 2, false},
		{"第三个元素", 3, false},
	}

	for i, tc := range testCases {
		s.Push(tc.value)

		// 验证栈的大小
		if size := s.Size(); size != tc.expectedSize {
			t.Errorf("测试用例 %d: 期望大小为 %d, 实际为 %d", i, tc.expectedSize, size)
		}

		// 验证栈的空状态
		if empty := s.IsEmpty(); empty != tc.expectedIsEmpty {
			t.Errorf("测试用例 %d: 期望空状态为 %v, 实际为 %v", i, tc.expectedIsEmpty, empty)
		}

		// 验证栈顶元素
		if top, err := s.Peek(); err != nil || top != tc.value {
			t.Errorf("测试用例 %d: 期望栈顶元素为 %s, 实际为 %s", i, tc.value, top)
		}
	}
}

// TestPop 测试出栈操作
func TestPop(t *testing.T) {
	s := New[int]()

	// 准备测试数据
	testValues := []int{1, 2, 3, 4, 5}
	for _, v := range testValues {
		s.Push(v)
	}

	// 测试正常出栈
	for i := len(testValues) - 1; i >= 0; i-- {
		val, err := s.Pop()
		if err != nil {
			t.Errorf("第 %d 次出栈操作失败: %v", i, err)
		}
		if val != testValues[i] {
			t.Errorf("第 %d 次出栈: 期望值为 %d, 实际为 %d", i, testValues[i], val)
		}
	}

	// 验证栈已经为空
	if !s.IsEmpty() {
		t.Error("所有元素出栈后，栈应该为空")
	}

	// 测试空栈出栈
	_, err := s.Pop()
	if err == nil {
		t.Error("从空栈出栈应该返回错误")
	}
}

// TestPeek 测试查看栈顶元素操作
func TestPeek(t *testing.T) {
	s := New[float64]()

	// 测试空栈的Peek操作
	_, err := s.Peek()
	if err == nil {
		t.Error("对空栈执行Peek操作应该返回错误")
	}

	// 添加元素并测试Peek
	testValue := 3.14
	s.Push(testValue)

	// 多次Peek，确保不会改变栈的状态
	for i := 0; i < 3; i++ {
		val, err := s.Peek()
		if err != nil {
			t.Errorf("第 %d 次Peek操作失败: %v", i, err)
		}
		if val != testValue {
			t.Errorf("第 %d 次Peek: 期望值为 %f, 实际为 %f", i, testValue, val)
		}
		if s.Size() != 1 {
			t.Errorf("Peek操作不应该改变栈的大小，当前大小: %d", s.Size())
		}
	}
}

// TestMixedOperations 测试混合操作
func TestMixedOperations(t *testing.T) {
	s := New[int]()

	// 测试用例：执行一系列混合操作
	operations := []struct {
		op          string
		pushValue   int
		expectedErr bool
		expectedVal int
	}{
		{"push", 1, false, 1},
		{"push", 2, false, 2},
		{"pop", 0, false, 2},
		{"peek", 0, false, 1},
		{"push", 3, false, 3},
		{"pop", 0, false, 3},
		{"pop", 0, false, 1},
		{"pop", 0, true, 0}, // 应该返回错误
	}

	for i, op := range operations {
		switch op.op {
		case "push":
			s.Push(op.pushValue)
			if val, _ := s.Peek(); val != op.expectedVal {
				t.Errorf("操作 %d (Push): 期望栈顶值为 %d, 实际为 %d", i, op.expectedVal, val)
			}
		case "pop":
			val, err := s.Pop()
			if (err != nil) != op.expectedErr {
				t.Errorf("操作 %d (Pop): 错误状态不符合预期", i)
			}
			if err == nil && val != op.expectedVal {
				t.Errorf("操作 %d (Pop): 期望值为 %d, 实际为 %d", i, op.expectedVal, val)
			}
		case "peek":
			val, err := s.Peek()
			if (err != nil) != op.expectedErr {
				t.Errorf("操作 %d (Peek): 错误状态不符合预期", i)
			}
			if err == nil && val != op.expectedVal {
				t.Errorf("操作 %d (Peek): 期望值为 %d, 实际为 %d", i, op.expectedVal, val)
			}
		}
	}
}

// TestStackWithCustomTypes 测试使用自定义类型
func TestStackWithCustomTypes(t *testing.T) {
	// 定义一个简单的自定义结构体
	type Person struct {
		Name string
		Age  int
	}

	s := New[Person]()

	// 测试数据
	p1 := Person{"张三", 20}
	p2 := Person{"李四", 25}

	// 测试入栈
	s.Push(p1)
	s.Push(p2)

	// 验证出栈顺序和数据完整性
	top, err := s.Pop()
	if err != nil || top != p2 {
		t.Errorf("期望出栈的人员信息为 %v, 实际为 %v", p2, top)
	}

	top, err = s.Pop()
	if err != nil || top != p1 {
		t.Errorf("期望出栈的人员信息为 %v, 实际为 %v", p1, top)
	}
}
