package rbtree

import (
	"golang.org/x/exp/constraints"
)

// Color 节点颜色
type Color bool

const (
	RED   Color = true  // 红色节点
	BLACK Color = false // 黑色节点
)

// Node 红黑树节点
// T 必须是可比较的类型(constraints.Ordered)
type Node[T constraints.Ordered] struct {
	Value  T        // 节点值
	Color  Color    // 节点颜色
	Left   *Node[T] // 左子节点
	Right  *Node[T] // 右子节点
	Parent *Node[T] // 父节点
}

// Tree 红黑树结构
type Tree[T constraints.Ordered] struct {
	Root *Node[T] // 根节点
	size int      // 树中节点数量
}

// NewTree 创建新的红黑树
// 时间复杂度: O(1)
func NewTree[T constraints.Ordered]() *Tree[T] {
	return &Tree[T]{
		Root: nil,
		size: 0,
	}
}

// Insert 插入新节点
// 红黑树的五个性质:
// 1. 每个节点要么是红色，要么是黑色
// 2. 根节点是黑色
// 3. 所有叶子节点都是黑色
// 4. 如果一个节点是红色，则它的子节点必须是黑色
// 5. 从任一节点到其每个叶子的所有路径都包含相同数目的黑色节点
// 时间复杂度: O(log n)
func (t *Tree[T]) Insert(value T) {
	// 创建新节点，初始为红色
	newNode := &Node[T]{
		Value:  value,
		Color:  RED, // 新节点默认为红色
		Left:   nil,
		Right:  nil,
		Parent: nil,
	}

	// 如果是空树，直接作为根节点
	if t.Root == nil {
		t.Root = newNode
		t.fixInsert(newNode) // 修复可能违反的红黑树性质
		return
	}

	// 找到合适的插入位置
	current := t.Root
	var parent *Node[T]
	for current != nil {
		parent = current
		if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}

	// 连接新节点
	newNode.Parent = parent
	if value < parent.Value {
		parent.Left = newNode
	} else {
		parent.Right = newNode
	}

	// 修复红黑树性质
	t.fixInsert(newNode)
	t.size++
}

// fixInsert 修复插入后可能违反的红黑树性质
// 时间复杂度: O(log n)，最多需要旋转O(log n)次
func (t *Tree[T]) fixInsert(node *Node[T]) {
	// 情况1：节点是根节点
	if node.Parent == nil {
		node.Color = BLACK
		return
	}

	// 如果父节点是黑色，不需要修复
	if node.Parent.Color == BLACK {
		return
	}

	// 获取父节点、叔叔节点和祖父节点
	parent := node.Parent
	grandparent := parent.Parent
	var uncle *Node[T]

	if grandparent.Left == parent {
		uncle = grandparent.Right
	} else {
		uncle = grandparent.Left
	}

	// 情况2：叔叔节点是红色
	// 解决方案：父节点和叔叔节点变黑，祖父节点变红，然后对祖父节点递归处理
	if uncle != nil && uncle.Color == RED {
		parent.Color = BLACK
		uncle.Color = BLACK
		grandparent.Color = RED
		t.fixInsert(grandparent)
		return
	}

	// 情况3：叔叔节点是黑色（或NIL），当前节点是“内侧子节点”
	// 解决方案：先对父节点进行一次旋转，转化为情况4
	if parent == grandparent.Left && node == parent.Right {
		t.rotateLeft(parent)
		node = parent
		parent = node.Parent
	} else if parent == grandparent.Right && node == parent.Left {
		t.rotateRight(parent)
		node = parent
		parent = node.Parent
	}

	// 情况4：叔叔节点是黑色（或NIL），当前节点是“外侧子节点”
	// 解决方案：对祖父节点进行一次旋转，并重新着色
	parent.Color = BLACK
	grandparent.Color = RED
	if node == parent.Left {
		t.rotateRight(grandparent)
	} else {
		t.rotateLeft(grandparent)
	}
}

// rotateLeft 左旋操作
// 时间复杂度: O(1)
func (t *Tree[T]) rotateLeft(node *Node[T]) {
	rightChild := node.Right
	node.Right = rightChild.Left

	if rightChild.Left != nil {
		rightChild.Left.Parent = node
	}

	rightChild.Parent = node.Parent
	if node.Parent == nil {
		t.Root = rightChild
	} else if node == node.Parent.Left {
		node.Parent.Left = rightChild
	} else {
		node.Parent.Right = rightChild
	}

	rightChild.Left = node
	node.Parent = rightChild
}

// rotateRight 右旋操作
// 时间复杂度: O(1)
func (t *Tree[T]) rotateRight(node *Node[T]) {
	leftChild := node.Left
	node.Left = leftChild.Right

	if leftChild.Right != nil {
		leftChild.Right.Parent = node
	}

	leftChild.Parent = node.Parent
	if node.Parent == nil {
		t.Root = leftChild
	} else if node == node.Parent.Right {
		node.Parent.Right = leftChild
	} else {
		node.Parent.Left = leftChild
	}

	leftChild.Right = node
	node.Parent = leftChild
}

// Search 查找节点
// 时间复杂度: O(log n)
func (t *Tree[T]) Search(value T) bool {
	current := t.Root
	for current != nil {
		if current.Value == value {
			return true
		}
		if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return false
}

// Size 返回树中节点数量
// 时间复杂度: O(1)
func (t *Tree[T]) Size() int {
	return t.size
}
