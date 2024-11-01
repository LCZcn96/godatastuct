package queue

import (
	"errors"
	"fmt"
)

// 定义队列操作可能遇到的错误
var (
	// ErrQueueEmpty 当队列为空时进行出队等操作会返回此错误
	ErrQueueEmpty = errors.New("队列为空")
	// ErrQueueFull 当队列已满时进行入队操作会返回此错误
	ErrQueueFull = errors.New("队列已满")
)

// Queue 队列接口
// 参考Java队列接口设计，提供两组操作方法：
// 1. 抛出错误的方法：Add/Remove/Element
// 2. 返回状态的方法：Offer/Poll/Peek
type Queue[T any] interface {
	// Add 将指定元素添加到队列尾部
	// 当队列已满时，返回 ErrQueueFull
	// 时间复杂度: O(1)
	Add(value T) error

	// Offer 将指定元素添加到队列尾部
	// 返回添加是否成功
	// 时间复杂度: O(1)
	Offer(value T) bool

	// Remove 移除并返回队首元素
	// 当队列为空时，返回零值和 ErrQueueEmpty
	// 时间复杂度: O(1)
	Remove() (T, error)

	// Poll 移除并返回队首元素
	// 当队列为空时，返回零值和 false
	// 时间复杂度: O(1)
	Poll() (T, bool)

	// Element 获取但不移除队首元素
	// 当队列为空时，返回零值和 ErrQueueEmpty
	// 时间复杂度: O(1)
	Element() (T, error)

	// Peek 获取但不移除队首元素
	// 当队列为空时，返回零值和 false
	// 时间复杂度: O(1)
	Peek() (T, bool)

	// IsEmpty 判断队列是否为空
	// 返回 true 表示队列为空
	// 时间复杂度: O(1)
	IsEmpty() bool

	// IsFull 判断队列是否已满
	// 返回 true 表示队列已满
	// 时间复杂度: O(1)
	IsFull() bool

	// Size 获取队列中元素的数量
	// 时间复杂度: O(1)
	Size() int

	// Clear 清空队列中的所有元素
	// 时间复杂度: O(n)
	Clear()
}

// CircularQueue 循环队列的具体实现
// 使用循环数组实现，提供高效的队列操作
type CircularQueue[T any] struct {
	elements []T // 存储元素的循环数组
	front    int // 队首元素的索引
	rear     int // 队尾元素的下一个位置的索引
	size     int // 当前队列中的元素数量
	capacity int // 队列的最大容量
}

// NewQueue 创建一个指定容量的新队列
// 参数：
//   - initialCapacity: 初始容量，必须大于0
//
// 返回值：
//   - Queue[T]: 队列接口实例
//   - error: 如果初始容量小于等于0，返回错误
func NewQueue[T any](initialCapacity int) (Queue[T], error) {
	if initialCapacity <= 0 {
		return nil, errors.New("初始容量必须大于0")
	}
	return &CircularQueue[T]{
		elements: make([]T, initialCapacity),
		front:    0,
		rear:     0,
		size:     0,
		capacity: initialCapacity,
	}, nil
}

// NewDefaultQueue 创建一个默认容量（16）的新队列
// 返回值：
//   - Queue[T]: 队列接口实例
func NewDefaultQueue[T any]() Queue[T] {
	q, _ := NewQueue[T](16)
	return q
}

// Add 将指定元素添加到队列尾部
// 参数：
//   - value: 要添加的元素
//
// 返回值：
//   - error: 队列已满时返回 ErrQueueFull，添加成功时返回 nil
func (q *CircularQueue[T]) Add(value T) error {
	if q.IsFull() {
		return ErrQueueFull
	}
	q.elements[q.rear] = value
	q.rear = (q.rear + 1) % q.capacity
	q.size++
	return nil
}

// Offer 尝试将指定元素添加到队列尾部
// 参数：
//   - value: 要添加的元素
//
// 返回值：
//   - bool: true表示添加成功，false表示队列已满
func (q *CircularQueue[T]) Offer(value T) bool {
	if q.IsFull() {
		return false
	}
	q.elements[q.rear] = value
	q.rear = (q.rear + 1) % q.capacity
	q.size++
	return true
}

// Remove 移除并返回队首元素
// 返回值：
//   - T: 队首元素，如果队列为空则返回零值
//   - error: 队列为空时返回 ErrQueueEmpty，否则返回 nil
func (q *CircularQueue[T]) Remove() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}
	value := q.elements[q.front]
	var zero T
	q.elements[q.front] = zero // 清除引用，帮助垃圾回收
	q.front = (q.front + 1) % q.capacity
	q.size--
	return value, nil
}

// Poll 尝试移除并返回队首元素
// 返回值：
//   - T: 队首元素，如果队列为空则返回零值
//   - bool: true表示成功移除元素，false表示队列为空
func (q *CircularQueue[T]) Poll() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	value := q.elements[q.front]
	var zero T
	q.elements[q.front] = zero
	q.front = (q.front + 1) % q.capacity
	q.size--
	return value, true
}

// Element 获取但不移除队首元素
// 返回值：
//   - T: 队首元素，如果队列为空则返回零值
//   - error: 队列为空时返回 ErrQueueEmpty，否则返回 nil
func (q *CircularQueue[T]) Element() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}
	return q.elements[q.front], nil
}

// Peek 尝试获取但不移除队首元素
// 返回值：
//   - T: 队首元素，如果队列为空则返回零值
//   - bool: true表示成功获取元素，false表示队列为空
func (q *CircularQueue[T]) Peek() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	return q.elements[q.front], true
}

// IsEmpty 判断队列是否为空
// 返回值：
//   - bool: true表示队列为空，false表示队列非空
func (q *CircularQueue[T]) IsEmpty() bool {
	return q.size == 0
}

// IsFull 判断队列是否已满
// 返回值：
//   - bool: true表示队列已满，false表示队列未满
func (q *CircularQueue[T]) IsFull() bool {
	return q.size == q.capacity
}

// Size 获取队列中元素的数量
// 返回值：
//   - int: 队列中的元素个数
func (q *CircularQueue[T]) Size() int {
	return q.size
}

// Clear 清空队列中的所有元素
// 该方法会清除所有元素的引用，帮助垃圾回收
func (q *CircularQueue[T]) Clear() {
	var zero T
	for i := range q.elements {
		q.elements[i] = zero
	}
	q.front = 0
	q.rear = 0
	q.size = 0
}

// String 返回队列的字符串表示
// 实现 fmt.Stringer 接口
// 返回值：
//   - string: 队列内容的字符串表示
func (q *CircularQueue[T]) String() string {
	if q.IsEmpty() {
		return "[]"
	}

	var result []T
	idx := q.front
	for i := 0; i < q.size; i++ {
		result = append(result, q.elements[idx])
		idx = (idx + 1) % q.capacity
	}
	return fmt.Sprintf("%v", result)
}

// ToSlice 将队列转换为切片
// 返回值：
//   - []T: 包含队列所有元素的切片副本
func (q *CircularQueue[T]) ToSlice() []T {
	if q.IsEmpty() {
		return []T{}
	}

	result := make([]T, q.size)
	idx := q.front
	for i := 0; i < q.size; i++ {
		result[i] = q.elements[idx]
		idx = (idx + 1) % q.capacity
	}
	return result
}
