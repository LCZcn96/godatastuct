package queue

import "errors"

// Deque 双端队列接口
// 支持在队列两端进行插入和删除操作
type Deque[T any] interface {
	PushFront(value T)    // 在队首插入元素
	PushBack(value T)     // 在队尾插入元素
	PopFront() (T, error) // 移除并返回队首元素
	PopBack() (T, error)  // 移除并返回队尾元素
	Front() (T, error)    // 查看队首元素但不移除
	Back() (T, error)     // 查看队尾元素但不移除
	IsEmpty() bool        // 检查双端队列是否为空
	Size() int            // 获取双端队列中元素个数
}

// deque 双端队列的具体实现
type deque[T any] struct {
	elements []T // 使用切片存储元素
}

// NewDeque 创建一个新的空双端队列
// 时间复杂度: O(1)
func NewDeque[T any]() Deque[T] {
	return &deque[T]{elements: []T{}}
}

// PushFront 在队首插入元素
// 时间复杂度: O(n) - 需要移动所有现有元素
func (d *deque[T]) PushFront(value T) {
	d.elements = append([]T{value}, d.elements...)
}

// PushBack 在队尾插入元素
// 时间复杂度: 平均 O(1)，需要扩容时，最坏 O(n)
func (d *deque[T]) PushBack(value T) {
	d.elements = append(d.elements, value)
}

// PopFront 移除并返回队首元素
// 时间复杂度: O(n)，需要移动所有剩余元素
func (d *deque[T]) PopFront() (T, error) {
	if d.IsEmpty() {
		var zero T
		return zero, errors.New("双端队列为空")
	}
	value := d.elements[0]
	d.elements = d.elements[1:]
	return value, nil
}

// PopBack 移除并返回队尾元素
// 时间复杂度: O(1)
func (d *deque[T]) PopBack() (T, error) {
	if d.IsEmpty() {
		var zero T
		return zero, errors.New("双端队列为空")
	}
	index := len(d.elements) - 1
	value := d.elements[index]
	d.elements = d.elements[:index]
	return value, nil
}

// Front 返回队首元素但不移除
// 时间复杂度: O(1)
func (d *deque[T]) Front() (T, error) {
	if d.IsEmpty() {
		var zero T
		return zero, errors.New("双端队列为空")
	}
	return d.elements[0], nil
}

// Back 返回队尾元素但不移除
// 时间复杂度: O(1)
func (d *deque[T]) Back() (T, error) {
	if d.IsEmpty() {
		var zero T
		return zero, errors.New("双端队列为空")
	}
	return d.elements[len(d.elements)-1], nil
}

// IsEmpty 检查双端队列是否为空
// 时间复杂度: O(1)
func (d *deque[T]) IsEmpty() bool {
	return len(d.elements) == 0
}

// Size 返回双端队列中元素的个数
// 时间复杂度: O(1)
func (d *deque[T]) Size() int {
	return len(d.elements)
}
