package binarytree

import (
	"testing"
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

// TestNewBinaryTree 测试创建新二叉树
func TestNewBinaryTree(t *testing.T) {
	tree := New(intCmp)
	if tree == nil {
		t.Fatal("New()返回了nil")
	}
}

// TestInsert 测试插入操作
func TestInsert(t *testing.T) {
	tree := New(intCmp)

	// 测试基本插入
	t.Run("Basic Insert", func(t *testing.T) {
		values := []int{5, 3, 7, 1, 4, 6, 8}
		for _, v := range values {
			tree.Insert(v)
		}

		// 验证插入的值都能找到
		for _, v := range values {
			if node := tree.Search(v); node == nil || node.Value != v {
				t.Errorf("未找到已插入的值: %d", v)
			}
		}
	})

	// 测试重复插入
	t.Run("Duplicate Insert", func(t *testing.T) {
		tree := New(intCmp)
		tree.Insert(1)
		tree.Insert(1)

		// 验证值存在
		if node := tree.Search(1); node == nil {
			t.Error("未找到插入的值1")
		}
	})
}

// TestSearch 测试查找操作
func TestSearch(t *testing.T) {
	tree := New(intCmp)
	values := []int{5, 3, 7, 1, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// 测试查找存在的值
	t.Run("Existing Values", func(t *testing.T) {
		for _, v := range values {
			if node := tree.Search(v); node == nil || node.Value != v {
				t.Errorf("未找到已存在的值: %d", v)
			}
		}
	})

	// 测试查找不存在的值
	t.Run("Non-existing Values", func(t *testing.T) {
		nonExisting := []int{0, 2, 9, 10}
		for _, v := range nonExisting {
			if node := tree.Search(v); node != nil {
				t.Errorf("找到了不应存在的值: %d", v)
			}
		}
	})
}

// TestRemove 测试删除操作
func TestRemove(t *testing.T) {
	// 初始化树
	tree := New(intCmp)
	values := []int{5, 3, 7, 1, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// 测试删除叶子节点
	t.Run("Remove Leaf", func(t *testing.T) {
		if !tree.Remove(1) {
			t.Error("删除叶子节点1失败")
		}
		if node := tree.Search(1); node != nil {
			t.Error("删除后仍能找到节点1")
		}
	})

	// 测试删除有一个子节点的节点
	t.Run("Remove Node with One Child", func(t *testing.T) {
		tree := New(intCmp)
		tree.Insert(2)
		tree.Insert(1)
		if !tree.Remove(2) {
			t.Error("删除带一个子节点的节点2失败")
		}
		if node := tree.Search(2); node != nil {
			t.Error("删除后仍能找到节点2")
		}
	})

	// 测试删除有两个子节点的节点
	t.Run("Remove Node with Two Children", func(t *testing.T) {
		if !tree.Remove(7) {
			t.Error("删除带两个子节点的节点7失败")
		}
		if node := tree.Search(7); node != nil {
			t.Error("删除后仍能找到节点7")
		}
		// 验证子树结构完整
		if node := tree.Search(6); node == nil {
			t.Error("节点7的左子节点6丢失")
		}
		if node := tree.Search(8); node == nil {
			t.Error("节点7的右子节点8丢失")
		}
	})

	// 测试删除不存在的节点
	t.Run("Remove Non-existing Node", func(t *testing.T) {
		if tree.Remove(100) {
			t.Error("删除不存在的节点应该返回false")
		}
	})
}

// TestTraversals 测试遍历操作
func TestTraversals(t *testing.T) {
	tree := New(intCmp)
	values := []int{5, 3, 7, 1, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	// 测试前序遍历
	t.Run("PreOrder Traversal", func(t *testing.T) {
		expected := []int{5, 3, 1, 4, 7, 6, 8}
		result := make([]int, 0)
		tree.PreOrderTraversal(func(v int) {
			result = append(result, v)
		})

		if !sliceEqual(result, expected) {
			t.Errorf("前序遍历结果错误，期望 %v，得到 %v", expected, result)
		}
	})

	// 测试中序遍历
	t.Run("InOrder Traversal", func(t *testing.T) {
		expected := []int{1, 3, 4, 5, 6, 7, 8}
		result := make([]int, 0)
		tree.InOrderTraversal(func(v int) {
			result = append(result, v)
		})

		if !sliceEqual(result, expected) {
			t.Errorf("中序遍历结果错误，期望 %v，得到 %v", expected, result)
		}
	})

	// 测试后序遍历
	t.Run("PostOrder Traversal", func(t *testing.T) {
		expected := []int{1, 4, 3, 6, 8, 7, 5}
		result := make([]int, 0)
		tree.PostOrderTraversal(func(v int) {
			result = append(result, v)
		})

		if !sliceEqual(result, expected) {
			t.Errorf("后序遍历结果错误，期望 %v，得到 %v", expected, result)
		}
	})
}

// TestDifferentTypes 测试不同数据类型
func TestDifferentTypes(t *testing.T) {
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

		tree := New(stringCmp)
		words := []string{"banana", "apple", "cherry"}
		for _, word := range words {
			tree.Insert(word)
		}

		// 验证插入和查找
		for _, word := range words {
			if node := tree.Search(word); node == nil || node.Value != word {
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

		tree := New(floatCmp)
		nums := []float64{3.14, 2.718, 1.414}
		for _, num := range nums {
			tree.Insert(num)
		}

		// 验证插入和查找
		for _, num := range nums {
			if node := tree.Search(num); node == nil || node.Value != num {
				t.Errorf("未找到已插入的浮点数: %f", num)
			}
		}
	})
}

// TestEmptyTree 测试空树操作
func TestEmptyTree(t *testing.T) {
	tree := New(intCmp)

	// 测试空树搜索
	if node := tree.Search(1); node != nil {
		t.Error("空树搜索应该返回nil")
	}

	// 测试空树删除
	if tree.Remove(1) {
		t.Error("空树删除应该返回false")
	}

	// 测试空树遍历
	count := 0
	tree.InOrderTraversal(func(v int) {
		count++
	})
	if count != 0 {
		t.Error("空树遍历不应该有任何回调")
	}
}

// sliceEqual 辅助函数：比较两个切片是否相等
func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
