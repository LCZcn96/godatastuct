package list

import (
	"math/rand"
	"time"
)

// 跳表实现
const (
	MaxLevel    = 16  // 最大层数
	Probability = 0.5 // 向上提升的概率
)

// node 跳表节点
type node[T any] struct {
	value T          // 节点值
	next  []*node[T] // 不同层级的下一个节点指针数组
}

// SkipList 跳表结构
type SkipList[T any] struct {
	header *node[T]         // 头节点（哨兵节点）
	level  int              // 当前最大层数
	cmp    func(a, b T) int // 比较函数
	rand   *rand.Rand       // 随机数生成器
}

func NewSkipList[T any](cmp func(a, b T) int) *SkipList[T] {
	return &SkipList[T]{
		header: &node[T]{next: make([]*node[T], MaxLevel)},
		level:  1,
		cmp:    cmp,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *SkipList[T]) randomLevel() int {
	level := 1
	for s.rand.Float64() < Probability && level < MaxLevel {
		level++
	}
	return level
}

func (s *SkipList[T]) Insert(value T) {
	update := make([]*node[T], MaxLevel)
	current := s.header

	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && s.cmp(current.next[i].value, value) < 0 {
			current = current.next[i]
		}
		update[i] = current
	}

	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.header
		}
		s.level = level
	}

	newNode := &node[T]{value: value, next: make([]*node[T], level)}
	for i := 0; i < level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

func (s *SkipList[T]) Search(value T) *T {
	current := s.header
	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && s.cmp(current.next[i].value, value) < 0 {
			current = current.next[i]
		}
	}
	current = current.next[0]
	if current != nil && s.cmp(current.value, value) == 0 {
		return &current.value
	}
	return nil
}

func (s *SkipList[T]) Delete(value T) bool {
	update := make([]*node[T], MaxLevel)
	current := s.header
	found := false

	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && s.cmp(current.next[i].value, value) < 0 {
			current = current.next[i]
		}
		update[i] = current
	}

	current = current.next[0]
	if current != nil && s.cmp(current.value, value) == 0 {
		found = true
		for i := 0; i < s.level; i++ {
			if update[i].next[i] != current {
				break
			}
			update[i].next[i] = current.next[i]
		}
		for s.level > 1 && s.header.next[s.level-1] == nil {
			s.level--
		}
	}
	return found
}
