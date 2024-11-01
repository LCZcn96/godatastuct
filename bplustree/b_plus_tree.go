package bplustree

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

// TreeNode B+ 树节点结构
// K: 键类型，必须是可比较的
// V: 值类型，可以是任意类型
type TreeNode[K constraints.Ordered, V any] struct {
	isLeaf   bool              // 是否为叶子节点
	keys     []K               // 键数组
	children []*TreeNode[K, V] // 子节点指针数组（仅对非叶子节点有效）
	values   []V               // 值数组（仅对叶子节点有效）
	next     *TreeNode[K, V]   // 指向下一个叶子节点的指针（用于范围查询）
	parent   *TreeNode[K, V]   // 父节点指针
}

// BPlusTree B+ 树结构
type BPlusTree[K constraints.Ordered, V any] struct {
	root  *TreeNode[K, V] // 根节点
	order int             // 树的阶数（每个节点最多可以有order个子节点）
}

// NewBPlusTree 创建新的 B+ 树
// 参数：
//   - order: 树的阶数，必须大于等于3
//
// 返回：
//   - *BPlusTree[K, V]: 新创建的 B+ 树指针
func NewBPlusTree[K constraints.Ordered, V any](order int) *BPlusTree[K, V] {
	if order < 3 {
		panic("阶数必须至少为3")
	}
	return &BPlusTree[K, V]{
		root: &TreeNode[K, V]{
			isLeaf: true,
			keys:   make([]K, 0),
			values: make([]V, 0),
		},
		order: order,
	}
}

// Insert 向 B+ 树中插入键值对
// 参数：
//   - key: 要插入的键
//   - value: 要插入的值
func (tree *BPlusTree[K, V]) Insert(key K, value V) {
	// 处理空树的情况
	if len(tree.root.keys) == 0 {
		tree.root.keys = append(tree.root.keys, key)
		tree.root.values = append(tree.root.values, value)
		return
	}

	// 查找要插入的叶子节点
	targetLeaf := tree.findLeaf(key)

	// 在叶子节点中查找插入位置
	insertPos := 0
	for insertPos < len(targetLeaf.keys) && targetLeaf.keys[insertPos] < key {
		insertPos++
	}

	// 如果键已存在，更新值
	if insertPos < len(targetLeaf.keys) && targetLeaf.keys[insertPos] == key {
		targetLeaf.values[insertPos] = value
		return
	}

	// 插入新的键值对
	targetLeaf.keys = append(targetLeaf.keys, key)
	targetLeaf.values = append(targetLeaf.values, value)

	// 将新插入的键值对移动到正确的位置
	for i := len(targetLeaf.keys) - 1; i > insertPos; i-- {
		targetLeaf.keys[i] = targetLeaf.keys[i-1]
		targetLeaf.values[i] = targetLeaf.values[i-1]
	}
	targetLeaf.keys[insertPos] = key
	targetLeaf.values[insertPos] = value

	// 检查是否需要分裂
	if len(targetLeaf.keys) >= tree.order {
		tree.splitLeafNode(targetLeaf)
	}
}

// findLeaf 查找要插入的叶子节点
// 参数：
//   - key: 要查找的键
//
// 返回：
//   - *TreeNode[K, V]: 包含给定键的叶子节点
func (tree *BPlusTree[K, V]) findLeaf(key K) *TreeNode[K, V] {
	currentNode := tree.root
	for !currentNode.isLeaf {
		pos := 0
		// 找到第一个大于或等于key的位置
		for pos < len(currentNode.keys) && currentNode.keys[pos] <= key {
			pos++
		}
		currentNode = currentNode.children[pos]
	}
	return currentNode
}

// splitLeafNode 分裂叶子节点
// 参数：
//   - leafNode: 需要分裂的叶子节点
func (tree *BPlusTree[K, V]) splitLeafNode(leafNode *TreeNode[K, V]) {
	midIndex := (len(leafNode.keys) + 1) / 2

	// 创建新的右侧节点
	newRightNode := &TreeNode[K, V]{
		isLeaf: true,
		keys:   make([]K, len(leafNode.keys[midIndex:])),
		values: make([]V, len(leafNode.values[midIndex:])),
		next:   leafNode.next,
		parent: leafNode.parent,
	}

	// 复制数据到新节点
	copy(newRightNode.keys, leafNode.keys[midIndex:])
	copy(newRightNode.values, leafNode.values[midIndex:])

	// 更新原节点
	leafNode.keys = leafNode.keys[:midIndex]
	leafNode.values = leafNode.values[:midIndex]
	leafNode.next = newRightNode

	// 获取用于父节点的键
	separatorKey := newRightNode.keys[0]

	// 处理父节点
	if leafNode == tree.root {
		// 创建新的根节点
		newRoot := &TreeNode[K, V]{
			isLeaf:   false,
			keys:     []K{separatorKey},
			children: []*TreeNode[K, V]{leafNode, newRightNode},
		}
		tree.root = newRoot
		leafNode.parent = newRoot
		newRightNode.parent = newRoot
	} else {
		tree.insertIntoParent(leafNode, separatorKey, newRightNode)
	}
}

// insertIntoParent 将分裂后的节点插入到父节点
func (tree *BPlusTree[K, V]) insertIntoParent(leftNode *TreeNode[K, V], key K, rightNode *TreeNode[K, V]) {
	parent := leftNode.parent
	insertPos := 0

	// 查找插入位置
	for insertPos < len(parent.keys) && parent.keys[insertPos] < key {
		insertPos++
	}

	// 插入键和子节点
	parent.keys = append(parent.keys, key)
	parent.children = append(parent.children, nil)

	// 移动现有的键和子节点
	for i := len(parent.keys) - 1; i > insertPos; i-- {
		parent.keys[i] = parent.keys[i-1]
		parent.children[i+1] = parent.children[i]
	}
	parent.keys[insertPos] = key
	parent.children[insertPos+1] = rightNode
	rightNode.parent = parent

	// 检查是否需要分裂父节点
	if len(parent.keys) >= tree.order {
		tree.splitInternalNode(parent)
	}
}

// splitInternalNode 分裂内部节点
func (tree *BPlusTree[K, V]) splitInternalNode(internalNode *TreeNode[K, V]) {
	midIndex := len(internalNode.keys) / 2
	promoteKey := internalNode.keys[midIndex]

	// 创建新的右侧节点
	newRightNode := &TreeNode[K, V]{
		isLeaf:   false,
		keys:     make([]K, len(internalNode.keys[midIndex+1:])),
		children: make([]*TreeNode[K, V], len(internalNode.children[midIndex+1:])),
	}

	// 复制键和子节点到新节点
	copy(newRightNode.keys, internalNode.keys[midIndex+1:])
	copy(newRightNode.children, internalNode.children[midIndex+1:])

	// 更新子节点的父指针
	for _, child := range newRightNode.children {
		child.parent = newRightNode
	}

	// 更新原节点
	internalNode.keys = internalNode.keys[:midIndex]
	internalNode.children = internalNode.children[:midIndex+1]

	// 处理父节点
	if internalNode == tree.root {
		newRoot := &TreeNode[K, V]{
			isLeaf:   false,
			keys:     []K{promoteKey},
			children: []*TreeNode[K, V]{internalNode, newRightNode},
		}
		tree.root = newRoot
		internalNode.parent = newRoot
		newRightNode.parent = newRoot
	} else {
		newRightNode.parent = internalNode.parent
		tree.insertIntoParent(internalNode, promoteKey, newRightNode)
	}
}

// Search 在 B+ 树中查找指定键对应的值
// 参数：
//   - key: 要查找的键
//
// 返回：
//   - V: 找到的值
//   - bool: 是否找到该键
func (tree *BPlusTree[K, V]) Search(key K) (V, bool) {
	currentNode := tree.root

	// 找到包含目标键的叶子节点
	for !currentNode.isLeaf {
		pos := 0
		for pos < len(currentNode.keys) && key >= currentNode.keys[pos] {
			pos++
		}
		currentNode = currentNode.children[pos]
	}

	// 在叶子节点中查找键
	for i := 0; i < len(currentNode.keys); i++ {
		if currentNode.keys[i] == key {
			return currentNode.values[i], true
		}
	}

	var zero V
	return zero, false
}

// String 返回树的字符串表示，用于调试
func (tree *BPlusTree[K, V]) String() string {
	if tree.root == nil {
		return "空树"
	}
	return tree.printTree(tree.root, 0)
}

// printTree 递归打印树的结构
func (tree *BPlusTree[K, V]) printTree(node *TreeNode[K, V], level int) string {
	var sb strings.Builder
	indent := strings.Repeat("  ", level)

	if node.isLeaf {
		sb.WriteString(fmt.Sprintf("%s叶子节点: keys=%v values=%v\n",
			indent, node.keys, node.values))
	} else {
		sb.WriteString(fmt.Sprintf("%s内部节点: keys=%v\n", indent, node.keys))
		for _, child := range node.children {
			sb.WriteString(tree.printTree(child, level+1))
		}
	}
	return sb.String()
}
