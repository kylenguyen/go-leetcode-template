package main

import (
	"container/heap"
	"fmt"
)

type ItemHeap struct {
	items []int
	index map[int]int // item -> index in heap
}

func NewItemHeap() *ItemHeap {
	return &ItemHeap{
		items: []int{},
		index: make(map[int]int),
	}
}

func (h *ItemHeap) Len() int           { return len(h.items) }
func (h *ItemHeap) Less(i, j int) bool { return h.items[i] < h.items[j] }
func (h *ItemHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.index[h.items[i]] = i
	h.index[h.items[j]] = j
}

func (h *ItemHeap) Push(x any) {
	item := x.(int)
	h.index[item] = len(h.items)
	h.items = append(h.items, item)
}

func (h *ItemHeap) Pop() any {
	n := len(h.items)
	item := h.items[n-1]
	h.items = h.items[:n-1]
	delete(h.index, item)
	return item
}

func (h *ItemHeap) Init() {
	heap.Init(h)
}

func (h *ItemHeap) Insert(x int) {
	heap.Push(h, x)
}

func (h *ItemHeap) GetMin() int {
	return h.items[0]
}

func (h *ItemHeap) Remove(x int) bool {
	i, ok := h.index[x]
	if !ok {
		return false
	}
	last := len(h.items) - 1
	h.Swap(i, last)
	h.items = h.items[:last]
	delete(h.index, x)
	if i < len(h.items) {
		heap.Fix(h, i)
	}
	return true
}
func main() {
	h := NewItemHeap()
	h.Init()

	h.Insert(5)
	h.Insert(3)
	h.Insert(8)

	fmt.Println("Min:", h.GetMin()) // 3

	h.Remove(3)
	fmt.Println("Min after removing 3:", h.GetMin()) // 5
}
