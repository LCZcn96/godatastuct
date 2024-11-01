package bplustree

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"testing"
)

// 辅助函数：验证B+树的基本属性
func validateBPlusTree[K constraints.Ordered, V any](t *testing.T, tree *BPlusTree[K, V]) {
	if tree.root == nil {
		t.Error("根节点不应为空")
		return
	}

	// 验证根节点属性
	if !tree.root.isLeaf && len(tree.root.children) == 0 {
		t.Error("非叶子根节点必须有子节点")
	}

	// 如果根节点不是叶子节点，验证所有内部节点
	if !tree.root.isLeaf {
		validateInternalNode(t, tree.root, tree.order)
	}
}

// 辅助函数：验证内部节点
func validateInternalNode[K constraints.Ordered, V any](t *testing.T, node *TreeNode[K, V], order int) {
	// 验证键的数量与子节点数量的关系
	if len(node.children) != len(node.keys)+1 {
		t.Errorf("内部节点的子节点数量应该等于键数量+1，当前键数量：%d，子节点数量：%d",
			len(node.keys), len(node.children))
	}

	// 验证子节点的父指针
	for _, child := range node.children {
		if child.parent != node {
			t.Error("子节点的父指针不正确")
		}
	}
}

func TestBPlusTreeBasicOperations(t *testing.T) {
	tree := NewBPlusTree[int, string](3)

	t.Run("空树操作", func(t *testing.T) {
		// 验证空树查找
		_, found := tree.Search(1)
		if found {
			t.Error("空树不应该找到任何值")
		}

		// 验证空树插入
		tree.Insert(1, "一")
		value, found := tree.Search(1)
		if !found || value != "一" {
			t.Errorf("插入后未找到值，got (%v, %v), want (一, true)", value, found)
		}

		validateBPlusTree(t, tree)
	})

	t.Run("基本插入和查找", func(t *testing.T) {
		testData := map[int]string{
			2: "二",
			3: "三",
			4: "四",
		}

		for k, v := range testData {
			tree.Insert(k, v)
			// 每次插入后立即验证是否可以查找到
			value, found := tree.Search(k)
			if !found || value != v {
				t.Errorf("插入后未找到键 %d，got (%v, %v), want (%v, true)",
					k, value, found, v)
			}
		}

		validateBPlusTree(t, tree)
	})

	t.Run("更新已存在的值", func(t *testing.T) {
		tree.Insert(1, "一一")
		value, found := tree.Search(1)
		if !found || value != "一一" {
			t.Errorf("更新后的值不正确，got (%v, %v), want (一一, true)",
				value, found)
		}

		validateBPlusTree(t, tree)
	})
}

func TestBPlusTreeNodeSplit(t *testing.T) {
	tree := NewBPlusTree[int, string](3) // 3阶B+树，方便测试分裂

	t.Run("叶子节点分裂", func(t *testing.T) {
		// 插入数据直到触发分裂
		data := []struct {
			key   int
			value string
		}{
			{1, "一"},
			{2, "二"},
			{3, "三"}, // 这次插入应该触发分裂
		}

		for _, d := range data {
			t.Logf("插入键: %d, 值: %s", d.key, d.value)
			tree.Insert(d.key, d.value)

			// 验证插入后的查找
			value, found := tree.Search(d.key)
			if !found || value != d.value {
				t.Errorf("分裂后查找失败 - 键 %d: got (%v, %v), want (%v, true)",
					d.key, value, found, d.value)
			}
		}

		// 验证树的结构
		validateBPlusTree(t, tree)

		// 打印树的结构以便调试
		t.Logf("叶子节点分裂后的树结构:\n%s", tree)
	})

	t.Run("内部节点分裂", func(t *testing.T) {
		// 继续插入数据触发内部节点分裂
		data := []struct {
			key   int
			value string
		}{
			{4, "四"},
			{5, "五"},
			{6, "六"},
			{7, "七"}, // 这些插入应该触发内部节点分裂
		}

		for _, d := range data {
			t.Logf("插入键: %d, 值: %s", d.key, d.value)
			tree.Insert(d.key, d.value)

			// 验证所有已插入的值是否都能找到
			value, found := tree.Search(d.key)
			if !found || value != d.value {
				t.Errorf("内部节点分裂后查找失败 - 键 %d: got (%v, %v), want (%v, true)",
					d.key, value, found, d.value)
			}
		}

		validateBPlusTree(t, tree)
		t.Logf("内部节点分裂后的树结构:\n%s", tree)
	})
}

func TestBPlusTreeEdgeCases(t *testing.T) {
	t.Run("无效的阶数", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("使用无效的阶数应该导致panic")
			}
		}()
		NewBPlusTree[int, string](2)
	})

	t.Run("查找不存在的键", func(t *testing.T) {
		tree := NewBPlusTree[int, string](3)
		tree.Insert(1, "一")

		_, found := tree.Search(999)
		if found {
			t.Error("不应该找到不存在的键")
		}
	})

	t.Run("重复键的处理", func(t *testing.T) {
		tree := NewBPlusTree[int, string](3)

		// 多次插入相同的键
		values := []string{"一", "一一", "一二"}
		for _, v := range values {
			tree.Insert(1, v)

			// 每次插入后验证值是否被正确更新
			got, found := tree.Search(1)
			if !found || got != v {
				t.Errorf("更新后的值不正确，got (%v, %v), want (%v, true)",
					got, found, v)
			}
		}

		validateBPlusTree(t, tree)
	})
}

func TestBPlusTreeLargeDataset(t *testing.T) {
	tree := NewBPlusTree[int, string](4)
	const dataSize = 100

	t.Run("大量数据插入和查询", func(t *testing.T) {
		// 按顺序插入数据
		t.Log("开始插入数据...")
		for i := 0; i < dataSize; i++ {
			value := fmt.Sprintf("值_%d", i)
			tree.Insert(i, value)

			// 验证插入是否成功
			got, found := tree.Search(i)
			if !found || got != value {
				t.Logf("树的当前状态:\n%s", tree)
				t.Fatalf("插入验证失败 - 键 %d: got (%v, %v), want (%v, true)",
					i, got, found, value)
			}
		}

		// 验证所有数据
		t.Log("验证所有已插入的数据...")
		for i := 0; i < dataSize; i++ {
			expectedValue := fmt.Sprintf("值_%d", i)
			value, found := tree.Search(i)
			if !found || value != expectedValue {
				t.Errorf("最终验证失败 - 键 %d: got (%v, %v), want (%v, true)",
					i, value, found, expectedValue)
			}
		}

		validateBPlusTree(t, tree)
		t.Log("大数据测试完成")
	})
}

func TestBPlusTreeRandomOperations(t *testing.T) {
	tree := NewBPlusTree[int, string](4)

	t.Run("随机顺序插入", func(t *testing.T) {
		data := []struct {
			key   int
			value string
		}{
			{5, "五"},
			{3, "三"},
			{7, "七"},
			{1, "一"},
			{6, "六"},
			{2, "二"},
			{4, "四"},
			{8, "八"},
		}

		// 记录已经插入的键
		inserted := make(map[int]string)

		for _, d := range data {
			t.Logf("插入键: %d, 值: %s", d.key, d.value)
			tree.Insert(d.key, d.value)

			// 记录这个键值对
			inserted[d.key] = d.value

			// 验证所有已插入的数据
			for k, v := range inserted {
				value, found := tree.Search(k)
				if !found || value != v {
					t.Errorf("随机插入后验证失败 - 键 %d: got (%v, %v), want (%v, true)",
						k, value, found, v)
				}
			}
		}

		validateBPlusTree(t, tree)
		t.Logf("最终树的结构:\n%s", tree)
	})
}

// 性能测试
func BenchmarkBPlusTreeOperations(b *testing.B) {
	tree := NewBPlusTree[int, string](4)

	b.Run("顺序插入", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Insert(i, fmt.Sprintf("值_%d", i))
		}
	})

	b.Run("随机查找", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Search(i % 100)
		}
	})
}
