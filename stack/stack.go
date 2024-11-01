package stack

import "errors"

// Stack 栈接口
// 支持泛型类型T
type Stack[T any] interface {
	Push(value T)     // 将元素压入栈顶
	Pop() (T, error)  // 弹出栈顶元素
	Peek() (T, error) // 查看栈顶元素但不移除
	IsEmpty() bool    // 检查栈是否为空
	Size() int        // 获取栈中元素个数
}

// stack 栈的结构体
// 使用切片作为底层存储结构
type stack[T any] struct {
	elements []T // 存储元素的切片
}

// New 创建一个新的空栈
// 时间复杂度: O(1)
func New[T any]() Stack[T] {
	return &stack[T]{elements: []T{}}
}

// Push 将元素压入栈顶
// 时间复杂度: 平均O(1)，当需要扩容时，最坏O(n)
func (s *stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

// Pop 弹出并返回栈顶元素
// 如果栈为空，返回错误
// 时间复杂度: O(1)
func (s *stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, errors.New("栈为空")
	}
	index := len(s.elements) - 1
	value := s.elements[index]
	s.elements = s.elements[:index]
	return value, nil
}

// Peek 返回栈顶元素但不移除
// 如果栈为空，返回错误
// 时间复杂度: O(1)
func (s *stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, errors.New("栈为空")
	}
	return s.elements[len(s.elements)-1], nil
}

// IsEmpty 检查栈是否为空
// 时间复杂度: O(1)
func (s *stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// Size 返回栈中元素的个数
// 时间复杂度: O(1)
func (s *stack[T]) Size() int {
	return len(s.elements)
}
