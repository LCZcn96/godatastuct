package queue

import (
	"errors"
	"testing"
)

// TestNewQueue 测试创建新队列
func TestNewQueue(t *testing.T) {
	// 测试创建有效容量的队列
	q, err := NewQueue[int](5)
	if err != nil {
		t.Fatalf("使用有效容量创建队列失败: %v", err)
	}
	if q == nil {
		t.Fatal("创建的队列不应为nil")
	}

	// 测试创建无效容量的队列
	q, err = NewQueue[int](0)
	if err == nil {
		t.Fatal("使用无效容量创建队列应该返回错误")
	}
	if q != nil {
		t.Fatal("使用无效容量创建队列时应返回nil")
	}

	// 测试创建默认队列
	q = NewDefaultQueue[int]()
	if q == nil {
		t.Fatal("创建默认队列失败")
	}
}

// TestQueueOperations 测试队列的基本操作
func TestQueueOperations(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	// 测试Add方法
	tests := []struct {
		name    string
		value   int
		wantErr error
	}{
		{"添加第一个元素", 1, nil},
		{"添加第二个元素", 2, nil},
		{"添加第三个元素", 3, nil},
		{"队列已满", 4, ErrQueueFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := q.Add(tt.value)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// 测试Size方法
	if size := q.Size(); size != 3 {
		t.Errorf("Size() = %v, want %v", size, 3)
	}

	// 测试Remove方法
	value, err := q.Remove()
	if err != nil {
		t.Fatalf("Remove()失败: %v", err)
	}
	if value != 1 {
		t.Errorf("Remove() = %v, want %v", value, 1)
	}
}

// TestQueueEmptyOperations 测试空队列操作
func TestQueueEmptyOperations(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	// 测试空队列的Remove操作
	_, err = q.Remove()
	if !errors.Is(err, ErrQueueEmpty) {
		t.Errorf("从空队列Remove应返回ErrQueueEmpty, got %v", err)
	}

	// 测试空队列的Element操作
	_, err = q.Element()
	if !errors.Is(err, ErrQueueEmpty) {
		t.Errorf("从空队列Element应返回ErrQueueEmpty, got %v", err)
	}

	// 测试空队列的Poll操作
	_, ok := q.Poll()
	if ok {
		t.Error("从空队列Poll应返回false")
	}

	// 测试空队列的Peek操作
	_, ok = q.Peek()
	if ok {
		t.Error("从空队列Peek应返回false")
	}
}

// TestQueueFullOperations 测试满队列操作
func TestQueueFullOperations(t *testing.T) {
	q, err := NewQueue[int](2)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	// 填充队列
	if err := q.Add(1); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}
	if err := q.Add(2); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}

	// 测试向满队列添加元素
	err = q.Add(3)
	if !errors.Is(err, ErrQueueFull) {
		t.Errorf("向满队列Add应返回ErrQueueFull, got %v", err)
	}

	// 测试向满队列offer元素
	if ok := q.Offer(3); ok {
		t.Error("向满队列Offer应返回false")
	}
}

// TestQueueClear 测试清空队列操作
func TestQueueClear(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	// 添加元素
	for _, v := range []int{1, 2, 3} {
		if err := q.Add(v); err != nil {
			t.Fatalf("添加元素失败: %v", err)
		}
	}

	// 清空队列
	q.Clear()

	// 验证队列为空
	if !q.IsEmpty() {
		t.Error("Clear()后队列应为空")
	}

	if size := q.Size(); size != 0 {
		t.Errorf("Clear()后Size() = %v, want 0", size)
	}
}

// TestCircularBehavior 测试循环队列行为
func TestCircularBehavior(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	// 添加和移除元素以测试循环行为
	if err := q.Add(1); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}
	if err := q.Add(2); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}

	// 移除1
	if _, err := q.Remove(); err != nil {
		t.Fatalf("移除元素失败: %v", err)
	}

	if err := q.Add(3); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}

	// 移除2
	if _, err := q.Remove(); err != nil {
		t.Fatalf("移除元素失败: %v", err)
	}

	if err := q.Add(4); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}

	// 验证元素
	value, err := q.Element()
	if err != nil {
		t.Fatalf("Element()失败: %v", err)
	}
	if value != 3 {
		t.Errorf("Element() = %v, want 3", value)
	}

	// 测试队列转换为切片
	circularQueue, ok := q.(*CircularQueue[int])
	if !ok {
		t.Fatal("无法将队列转换为CircularQueue类型")
	}

	slice := circularQueue.ToSlice()
	expected := []int{3, 4}
	if len(slice) != len(expected) {
		t.Errorf("ToSlice()长度 = %v, want %v", len(slice), len(expected))
	}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("ToSlice()[%d] = %v, want %v", i, slice[i], v)
		}
	}
}

// TestQueueString 测试队列的字符串表示
func TestQueueString(t *testing.T) {
	q, err := NewQueue[int](3)
	if err != nil {
		t.Fatalf("创建队列失败: %v", err)
	}

	circularQueue, ok := q.(*CircularQueue[int])
	if !ok {
		t.Fatal("无法将队列转换为CircularQueue类型")
	}

	// 测试空队列的字符串表示
	if s := circularQueue.String(); s != "[]" {
		t.Errorf("空队列String() = %v, want []", s)
	}

	// 添加元素后测试字符串表示
	if err := q.Add(1); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}
	if err := q.Add(2); err != nil {
		t.Fatalf("添加元素失败: %v", err)
	}

	expected := "[1 2]"
	if s := circularQueue.String(); s != expected {
		t.Errorf("String() = %v, want %v", s, expected)
	}
}
