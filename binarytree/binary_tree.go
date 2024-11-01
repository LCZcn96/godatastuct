package binarytree

// TreeNode 定义了二叉树的节点
type TreeNode[T any] struct {
	Value T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

// BinaryTree 定义了二叉树的接口
type BinaryTree[T any] interface {
	Insert(value T)
	Search(value T) *TreeNode[T]
	Remove(value T) bool
	PreOrderTraversal(func(T))
	InOrderTraversal(func(T))
	PostOrderTraversal(func(T))
}

// binaryTree 实现了 BinaryTree 接口
type binaryTree[T any] struct {
	root *TreeNode[T]
	cmp  func(a, b T) int // 比较函数，用于比较节点值
}

// New 创建一个新的二叉树，需要传入一个比较函数
func New[T any](cmp func(a, b T) int) BinaryTree[T] {
	return &binaryTree[T]{cmp: cmp}
}

func (t *binaryTree[T]) Insert(value T) {
	t.root = t.insertRec(t.root, value)
}

func (t *binaryTree[T]) insertRec(node *TreeNode[T], value T) *TreeNode[T] {
	if node == nil {
		return &TreeNode[T]{Value: value}
	}
	if t.cmp(value, node.Value) < 0 {
		node.Left = t.insertRec(node.Left, value)
	} else {
		node.Right = t.insertRec(node.Right, value)
	}
	return node
}

func (t *binaryTree[T]) Search(value T) *TreeNode[T] {
	return t.searchRec(t.root, value)
}

func (t *binaryTree[T]) searchRec(node *TreeNode[T], value T) *TreeNode[T] {
	if node == nil || t.cmp(value, node.Value) == 0 {
		return node
	}
	if t.cmp(value, node.Value) < 0 {
		return t.searchRec(node.Left, value)
	}
	return t.searchRec(node.Right, value)
}

func (t *binaryTree[T]) Remove(value T) bool {
	var removed bool
	t.root, removed = t.removeRec(t.root, value)
	return removed
}

func (t *binaryTree[T]) removeRec(node *TreeNode[T], value T) (*TreeNode[T], bool) {
	if node == nil {
		return nil, false
	}
	var removed bool
	if t.cmp(value, node.Value) < 0 {
		node.Left, removed = t.removeRec(node.Left, value)
	} else if t.cmp(value, node.Value) > 0 {
		node.Right, removed = t.removeRec(node.Right, value)
	} else {
		removed = true
		if node.Left == nil {
			return node.Right, true
		} else if node.Right == nil {
			return node.Left, true
		} else {
			// 找到右子树中最小的节点替换当前节点
			minNode := t.findMin(node.Right)
			node.Value = minNode.Value
			node.Right, _ = t.removeRec(node.Right, minNode.Value)
		}
	}
	return node, removed
}

func (t *binaryTree[T]) findMin(node *TreeNode[T]) *TreeNode[T] {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}
func (t *binaryTree[T]) PreOrderTraversal(f func(T)) {
	t.preOrderRec(t.root, f)
}

func (t *binaryTree[T]) preOrderRec(node *TreeNode[T], f func(T)) {
	if node != nil {
		f(node.Value)
		t.preOrderRec(node.Left, f)
		t.preOrderRec(node.Right, f)
	}
}
func (t *binaryTree[T]) InOrderTraversal(f func(T)) {
	t.inOrderRec(t.root, f)
}

func (t *binaryTree[T]) inOrderRec(node *TreeNode[T], f func(T)) {
	if node != nil {
		t.inOrderRec(node.Left, f)
		f(node.Value)
		t.inOrderRec(node.Right, f)
	}
}

func (t *binaryTree[T]) PostOrderTraversal(f func(T)) {
	t.postOrderRec(t.root, f)
}

func (t *binaryTree[T]) postOrderRec(node *TreeNode[T], f func(T)) {
	if node != nil {
		t.postOrderRec(node.Left, f)
		t.postOrderRec(node.Right, f)
		f(node.Value)
	}
}
