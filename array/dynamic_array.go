package dynamicarray

import (
	"errors"
)

// 常量定义
const (
	initialCapacity = 4    // 初始容量大小
	shrinkFactor    = 0.25 // 缩容触发因子：当数组大小/容量小于此值时触发缩容
)

// DynamicArray 动态数组接口
// 支持泛型类型 T
// 实现了自动扩容和缩容的动态数组数据结构
type DynamicArray[T any] interface {
	Append(value T)                  // 在数组末尾添加元素
	Insert(index int, value T) error // 在指定位置插入元素
	Remove(index int) (T, error)     // 删除指定位置的元素并返回
	Get(index int) (T, error)        // 获取指定位置的元素
	Set(index int, value T) error    // 设置指定位置的元素
	Len() int                        // 获取数组当前长度
	Cap() int                        // 获取数组当前容量
}

// dynamicArray 动态数组实现
type dynamicArray[T any] struct {
	data     []T // 底层切片
	size     int // 当前元素数量
	capacity int // 当前容量
}

// New 创建新的动态数组
// 时间复杂度: O(1)
func New[T any]() DynamicArray[T] {
	return &dynamicArray[T]{
		data:     make([]T, initialCapacity),
		size:     0,
		capacity: initialCapacity,
	}
}

// Append 在数组末尾添加元素
// 时间复杂度: 平均 O(1)，需要扩容时，最坏 O(n)
func (da *dynamicArray[T]) Append(value T) {
	// 如果size达到容量上限,需要扩容
	if da.size == da.capacity {
		da.resize(da.capacity * 2)
	}
	da.data[da.size] = value
	da.size++
}

// Insert 在指定索引位置插入元素
// 参数:
//   - index: 插入位置
//   - value: 待插入的值
//
// 返回值:
//   - error: 索引越界时返回错误
//
// 时间复杂度: O(n) - 需要移动插入位置后的所有元素
func (da *dynamicArray[T]) Insert(index int, value T) error {
	if index < 0 || index > da.size {
		return errors.New("索引越界")
	}
	// 容量检查
	if da.size == da.capacity {
		da.resize(da.capacity * 2)
	}
	// 移动元素，为插入腾出空间
	copy(da.data[index+1:], da.data[index:da.size])
	da.data[index] = value
	da.size++
	return nil
}

// Remove 删除并返回指定索引位置的元素
// 参数:
//   - index: 要删除元素的索引
//
// 返回值:
//   - T: 被删除的元素
//   - error: 索引越界时返回错误
//
// 时间复杂度: O(n)，需要移动删除位置后的所有元素
func (da *dynamicArray[T]) Remove(index int) (T, error) {
	if index < 0 || index >= da.size {
		var zero T
		return zero, errors.New("索引越界")
	}
	value := da.data[index]
	// 移动元素填补空缺
	copy(da.data[index:], da.data[index+1:da.size])
	da.size--
	var zero T
	da.data[da.size] = zero // 清理最后一个元素

	// 缩容检查
	if da.size > 0 && float64(da.size)/float64(da.capacity) <= shrinkFactor {
		da.resize(da.capacity / 2)
	}

	return value, nil
}

// resize 调整数组容量
// 参数:
//   - newCapacity: 新的容量大小
//
// 时间复杂度: O(n) - 需要复制所有元素到新的底层数组
func (da *dynamicArray[T]) resize(newCapacity int) {
	newData := make([]T, newCapacity)
	copy(newData, da.data[:da.size])
	da.data = newData
	da.capacity = newCapacity
}

// Get 获取指定索引位置的元素
// 时间复杂度: O(1)
func (da *dynamicArray[T]) Get(index int) (T, error) {
	if index < 0 || index >= da.size {
		var zero T
		return zero, errors.New("索引越界")
	}
	return da.data[index], nil
}

// Set 设置指定索引位置的元素值
// 时间复杂度: O(1)
func (da *dynamicArray[T]) Set(index int, value T) error {
	if index < 0 || index >= da.size {
		return errors.New("索引越界")
	}
	da.data[index] = value
	return nil
}

// Len 返回数组中元素的个数
// 时间复杂度: O(1)
func (da *dynamicArray[T]) Len() int {
	return da.size
}

// Cap 返回数组的容量
// 时间复杂度: O(1)
func (da *dynamicArray[T]) Cap() int {
	return da.capacity
}
