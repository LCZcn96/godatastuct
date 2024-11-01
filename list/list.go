package list

// Node 链表节点定义
// 类型参数 T 必须是可比较的类型
type Node[T comparable] struct {
	Value T        // 节点存储的值
	Next  *Node[T] // 指向下一个节点的指针
}

// LinkedList 链表接口
// 定义了单链表支持的所有操作
type LinkedList[T comparable] interface {
	Append(value T)               // 在链表末尾添加节点
	Prepend(value T)              // 在链表头部添加节点
	Insert(index int, value T)    // 在指定位置插入节点
	Remove(value T) bool          // 删除指定值的节点
	RemoveAt(index int) (T, bool) // 删除指定位置的节点
	Find(value T) *Node[T]        // 查找指定值的节点
	Get(index int) (T, bool)      // 获取指定位置的值
	Set(index int, value T) bool  // 设置指定位置的值
	IsEmpty() bool                // 检查链表是否为空
	Size() int                    // 获取链表长度
	Clear()                       // 清空链表
	ToSlice() []T                 // 将链表转换为切片
}

// linkedList 链表实现
type linkedList[T comparable] struct {
	head *Node[T] // 头节点指针
	tail *Node[T] // 尾节点指针
	size int      // 链表大小
}

// New 创建新的链表
// 时间复杂度: O(1)
func New[T comparable]() LinkedList[T] {
	return &linkedList[T]{}
}

// Append 在链表末尾添加节点
// 时间复杂度: O(1) - 由于维护了tail指针
func (l *linkedList[T]) Append(value T) {
	newNode := &Node[T]{Value: value}
	if l.head == nil {
		// 空链表情况
		l.head = newNode
		l.tail = newNode
	} else {
		// 非空链表，直接在尾部添加
		l.tail.Next = newNode
		l.tail = newNode
	}
	l.size++
}

// Prepend 在链表头部添加节点
// 时间复杂度: O(1)
func (l *linkedList[T]) Prepend(value T) {
	newNode := &Node[T]{Value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.Next = l.head
		l.head = newNode
	}
	l.size++
}

// Insert 在指定位置插入节点
// 参数:
//   - index: 插入位置
//   - value: 待插入的值
//
// 时间复杂度: O(n) - 需要遍历到指定位置
func (l *linkedList[T]) Insert(index int, value T) {
	if index < 0 || index > l.size {
		panic("索引越界")
	}
	if index == 0 {
		l.Prepend(value)
		return
	}
	if index == l.size {
		l.Append(value)
		return
	}

	newNode := &Node[T]{Value: value}
	prevNode := l.getNodeAt(index - 1)
	newNode.Next = prevNode.Next
	prevNode.Next = newNode
	l.size++
}

// Remove 删除指定值的节点
// 返回是否成功删除
// 时间复杂度: O(n) - 需要遍历查找值
func (l *linkedList[T]) Remove(value T) bool {
	if l.head == nil {
		return false
	}

	// 处理头节点的特殊情况
	if l.head.Value == value {
		l.head = l.head.Next
		if l.head == nil {
			l.tail = nil
		}
		l.size--
		return true
	}

	// 遍历查找要删除的节点
	prev := l.head
	current := l.head.Next
	for current != nil {
		if current.Value == value {
			prev.Next = current.Next
			if current == l.tail {
				l.tail = prev
			}
			l.size--
			return true
		}
		prev = current
		current = current.Next
	}
	return false
}

func (l *linkedList[T]) RemoveAt(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}
	var removedNode *Node[T]
	if index == 0 {
		removedNode = l.head
		l.head = l.head.Next
		if l.head == nil {
			l.tail = nil
		}
	} else {
		prevNode := l.getNodeAt(index - 1)
		removedNode = prevNode.Next
		prevNode.Next = removedNode.Next
		if removedNode == l.tail {
			l.tail = prevNode
		}
	}
	l.size--
	return removedNode.Value, true
}
func (l *linkedList[T]) Find(value T) *Node[T] {
	current := l.head
	for current != nil {
		if current.Value == value {
			return current
		}
		current = current.Next
	}
	return nil
}
func (l *linkedList[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}
	node := l.getNodeAt(index)
	return node.Value, true
}
func (l *linkedList[T]) Set(index int, value T) bool {
	if index < 0 || index >= l.size {
		return false
	}
	node := l.getNodeAt(index)
	node.Value = value
	return true
}
func (l *linkedList[T]) getNodeAt(index int) *Node[T] {
	current := l.head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	return current
}
func (l *linkedList[T]) IsEmpty() bool {
	return l.size == 0
}
func (l *linkedList[T]) Size() int {
	return l.size
}
func (l *linkedList[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}
func (l *linkedList[T]) ToSlice() []T {
	slice := make([]T, 0, l.size)
	current := l.head
	for current != nil {
		slice = append(slice, current.Value)
		current = current.Next
	}
	return slice
}
