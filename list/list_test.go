package list

import (
	"testing"
)

// TestNewList 测试创建新链表
func TestNewList(t *testing.T) {
	list := New[int]()
	if list == nil {
		t.Fatal("New()返回了nil")
	}

	if !list.IsEmpty() {
		t.Error("新创建的链表应该为空")
	}

	if size := list.Size(); size != 0 {
		t.Errorf("新创建的链表大小应为0，实际为%d", size)
	}
}

// TestAppendAndPrepend 测试在链表两端添加节点
func TestAppendAndPrepend(t *testing.T) {
	list := New[int]()

	// 测试Append
	t.Run("Testing Append", func(t *testing.T) {
		values := []int{1, 2, 3}
		for _, v := range values {
			list.Append(v)
		}

		// 验证节点顺序和数量
		if size := list.Size(); size != len(values) {
			t.Errorf("Append后Size()=%d，期望值为%d", size, len(values))
		}

		slice := list.ToSlice()
		for i, v := range values {
			if slice[i] != v {
				t.Errorf("位置%d的值为%d，期望值为%d", i, slice[i], v)
			}
		}
	})

	// 清空链表
	list = New[int]()

	// 测试Prepend
	t.Run("Testing Prepend", func(t *testing.T) {
		values := []int{1, 2, 3}
		for _, v := range values {
			list.Prepend(v)
		}

		// 验证节点顺序和数量
		if size := list.Size(); size != len(values) {
			t.Errorf("Prepend后Size()=%d，期望值为%d", size, len(values))
		}

		// 验证顺序(应该是反序)
		slice := list.ToSlice()
		for i := 0; i < len(values); i++ {
			expected := values[len(values)-1-i]
			if slice[i] != expected {
				t.Errorf("位置%d的值为%d，期望值为%d", i, slice[i], expected)
			}
		}
	})
}

// TestInsert 测试在指定位置插入节点
func TestInsert(t *testing.T) {
	list := New[int]()

	// 测试在空链表开头插入
	list.Insert(0, 1)
	if size := list.Size(); size != 1 {
		t.Errorf("插入后Size()=%d，期望值为1", size)
	}

	// 测试在结尾插入
	list.Insert(1, 3)

	// 测试在中间插入
	list.Insert(1, 2)

	// 验证最终顺序
	expected := []int{1, 2, 3}
	slice := list.ToSlice()
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("位置%d的值为%d，期望值为%d", i, slice[i], v)
		}
	}

	// 测试越界插入
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("插入越界索引应该触发panic")
			}
		}()
		list.Insert(4, 4)
	}()
}

// TestRemove 测试删除节点
func TestRemove(t *testing.T) {
	list := New[int]()

	// 准备测试数据
	values := []int{1, 2, 3, 2, 4}
	for _, v := range values {
		list.Append(v)
	}

	// 测试删除存在的值
	if !list.Remove(2) { // 删除第一个2
		t.Error("删除存在的值应该返回true")
	}

	// 再次删除2，确保能删除重复值
	if !list.Remove(2) {
		t.Error("删除第二个2应该返回true")
	}

	// 测试删除不存在的值
	if list.Remove(5) {
		t.Error("删除不存在的值应该返回false")
	}

	// 验证最终结果
	expected := []int{1, 3, 4}
	slice := list.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("删除后长度为%d，期望值为%d", len(slice), len(expected))
	}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("位置%d的值为%d，期望值为%d", i, slice[i], v)
		}
	}
}

// TestRemoveAt 测试按索引删除节点
func TestRemoveAt(t *testing.T) {
	list := New[int]()

	// 准备测试数据
	values := []int{1, 2, 3, 4}
	for _, v := range values {
		list.Append(v)
	}

	// 测试删除中间节点
	value, ok := list.RemoveAt(1)
	if !ok {
		t.Error("删除有效索引应该返回true")
	}
	if value != 2 {
		t.Errorf("删除的值为%d，期望值为2", value)
	}

	// 测试删除首节点
	value, ok = list.RemoveAt(0)
	if !ok {
		t.Error("删除首节点应该返回true")
	}
	if value != 1 {
		t.Errorf("删除的值为%d，期望值为1", value)
	}

	// 测试删除尾节点
	value, ok = list.RemoveAt(list.Size() - 1)
	if !ok {
		t.Error("删除尾节点应该返回true")
	}
	if value != 4 {
		t.Errorf("删除的值为%d，期望值为4", value)
	}

	// 测试删除无效索引
	_, ok = list.RemoveAt(list.Size())
	if ok {
		t.Error("删除无效索引应该返回false")
	}
}

// TestFind 测试查找节点
func TestFind(t *testing.T) {
	list := New[int]()

	// 准备测试数据
	values := []int{1, 2, 3, 4}
	for _, v := range values {
		list.Append(v)
	}

	// 测试查找存在的值
	node := list.Find(3)
	if node == nil {
		t.Error("查找存在的值不应返回nil")
	}
	if node.Value != 3 {
		t.Errorf("找到的节点值为%d，期望值为3", node.Value)
	}

	// 测试查找不存在的值
	node = list.Find(5)
	if node != nil {
		t.Error("查找不存在的值应返回nil")
	}
}

// TestGetAndSet 测试获取和设置节点值
func TestGetAndSet(t *testing.T) {
	list := New[int]()

	// 准备测试数据
	values := []int{1, 2, 3}
	for _, v := range values {
		list.Append(v)
	}

	// 测试Get
	t.Run("Testing Get", func(t *testing.T) {
		value, ok := list.Get(1)
		if !ok {
			t.Error("获取有效索引应该返回true")
		}
		if value != 2 {
			t.Errorf("Get(1)=%d，期望值为2", value)
		}

		// 测试无效索引
		_, ok = list.Get(-1)
		if ok {
			t.Error("获取无效索引应该返回false")
		}
	})

	// 测试Set
	t.Run("Testing Set", func(t *testing.T) {
		if !list.Set(1, 5) {
			t.Error("设置有效索引应该返回true")
		}

		value, _ := list.Get(1)
		if value != 5 {
			t.Errorf("设置后的值为%d，期望值为5", value)
		}

		// 测试无效索引
		if list.Set(-1, 1) {
			t.Error("设置无效索引应该返回false")
		}
	})
}

// TestClear 测试清空链表
func TestClear(t *testing.T) {
	list := New[int]()

	// 添加一些节点
	for i := 1; i <= 3; i++ {
		list.Append(i)
	}

	// 清空链表
	list.Clear()

	if !list.IsEmpty() {
		t.Error("Clear()后链表应该为空")
	}

	if size := list.Size(); size != 0 {
		t.Errorf("Clear()后Size()=%d，期望值为0", size)
	}

	// 验证可以在清空后继续添加元素
	list.Append(1)
	if size := list.Size(); size != 1 {
		t.Errorf("Clear()后追加元素，Size()=%d，期望值为1", size)
	}
}

// TestToSlice 测试转换为切片
func TestToSlice(t *testing.T) {
	list := New[int]()

	// 测试空链表
	slice := list.ToSlice()
	if len(slice) != 0 {
		t.Error("空链表转换为切片应该得到空切片")
	}

	// 添加元素后测试
	values := []int{1, 2, 3}
	for _, v := range values {
		list.Append(v)
	}

	slice = list.ToSlice()
	if len(slice) != len(values) {
		t.Errorf("切片长度为%d，期望值为%d", len(slice), len(values))
	}

	for i, v := range values {
		if slice[i] != v {
			t.Errorf("位置%d的值为%d，期望值为%d", i, slice[i], v)
		}
	}
}

// TestDifferentTypes 测试不同数据类型
func TestDifferentTypes(t *testing.T) {
	// 测试字符串类型
	t.Run("String Type", func(t *testing.T) {
		list := New[string]()
		list.Append("hello")
		list.Append("world")

		if size := list.Size(); size != 2 {
			t.Errorf("Size()=%d，期望值为2", size)
		}

		value, ok := list.Get(0)
		if !ok || value != "hello" {
			t.Errorf("Get(0)=%s，期望值为hello", value)
		}
	})

	// 测试浮点数类型
	t.Run("Float Type", func(t *testing.T) {
		list := New[float64]()
		list.Append(3.14)
		list.Append(2.718)

		if size := list.Size(); size != 2 {
			t.Errorf("Size()=%d，期望值为2", size)
		}

		value, ok := list.Get(1)
		if !ok || value != 2.718 {
			t.Errorf("Get(1)=%f，期望值为2.718", value)
		}
	})
}
