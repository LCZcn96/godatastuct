package rbtree

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"testing"
)

// validateRedBlackProperties 验证红黑树的所有性质
func validateRedBlackProperties[T constraints.Ordered](t *testing.T, tree *Tree[T]) {
	if tree.Root == nil {
		return // 空树是有效的红黑树
	}

	// 验证根节点是黑色（性质2）
	if tree.Root.Color != BLACK {
		t.Error("根节点必须是黑色")
	}

	// 验证从根节点开始的所有性质
	blackHeight, err := validateNode(tree.Root, BLACK)
	if err != nil {
		t.Errorf("红黑树性质验证失败: %v", err)
	}

	t.Logf("红黑树黑高度为: %d", blackHeight)
}

// validateNode 验证节点及其子树的红黑树性质
func validateNode[T constraints.Ordered](node *Node[T], parentColor Color) (int, error) {
	if node == nil {
		return 1, nil // NIL节点被视为黑色
	}

	// 检查红色节点的子节点是否为黑色（性质4）
	if node.Color == RED && parentColor == RED {
		return 0, fmt.Errorf("发现连续的红色节点")
	}

	// 递归验证左子树
	leftBlackHeight, err := validateNode(node.Left, node.Color)
	if err != nil {
		return 0, err
	}

	// 递归验证右子树
	rightBlackHeight, err := validateNode(node.Right, node.Color)
	if err != nil {
		return 0, err
	}

	// 验证左右子树的黑高度相同（性质5）
	if leftBlackHeight != rightBlackHeight {
		return 0, fmt.Errorf("左右子树的黑高度不相等：左 %d, 右 %d",
			leftBlackHeight, rightBlackHeight)
	}

	// 计算当前节点的黑高度
	blackHeight := leftBlackHeight
	if node.Color == BLACK {
		blackHeight++
	}

	return blackHeight, nil
}

func TestRedBlackTreeBasicOperations(t *testing.T) {
	tree := NewTree[int]()

	t.Run("空树操作", func(t *testing.T) {
		if !tree.Search(1) {
			t.Log("空树查找测试通过")
		}
		validateRedBlackProperties(t, tree)
	})

	t.Run("基本插入操作", func(t *testing.T) {
		values := []int{7, 3, 18, 10, 22, 8, 11, 26, 2, 6}
		for _, v := range values {
			t.Logf("插入值: %d", v)
			tree.Insert(v)
			validateRedBlackProperties(t, tree)

			// 验证插入后能够找到该值
			if !tree.Search(v) {
				t.Errorf("未找到已插入的值: %d", v)
			}
		}
	})
}

func TestRedBlackTreeBalancing(t *testing.T) {
	tree := NewTree[int]()

	t.Run("左旋平衡", func(t *testing.T) {
		// 构造需要左旋的情况
		values := []int{10, 20, 30}
		for _, v := range values {
			tree.Insert(v)
			validateRedBlackProperties(t, tree)
		}
	})

	t.Run("右旋平衡", func(t *testing.T) {
		tree = NewTree[int]() // 重置树
		// 构造需要右旋的情况
		values := []int{30, 20, 10}
		for _, v := range values {
			tree.Insert(v)
			validateRedBlackProperties(t, tree)
		}
	})

	t.Run("双旋转平衡", func(t *testing.T) {
		tree = NewTree[int]() // 重置树
		// 构造需要双旋转的情况
		values := []int{30, 10, 20}
		for _, v := range values {
			tree.Insert(v)
			validateRedBlackProperties(t, tree)
		}
	})
}

func TestRedBlackTreeProperties(t *testing.T) {
	tree := NewTree[int]()

	t.Run("连续插入升序值", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			tree.Insert(i)
			validateRedBlackProperties(t, tree)
		}
	})

	t.Run("连续插入降序值", func(t *testing.T) {
		tree = NewTree[int]()
		for i := 10; i >= 1; i-- {
			tree.Insert(i)
			validateRedBlackProperties(t, tree)
		}
	})
}

// 添加性能测试
func BenchmarkRedBlackTree(b *testing.B) {
	tree := NewTree[int]()

	b.Run("顺序插入", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Insert(i)
		}
	})

	b.Run("查找操作", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Search(i % 100) // 循环查找前100个数
		}
	})
}
