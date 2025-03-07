package priorityQueue

import (
	"container/heap"

	"github.com/khorcarol/AgentOfThings/lib/option"
)

type Item struct {
	value    any
	priority int
	index    int
}

type tHeap []*Item

type PriorityQueue[T comparable] struct {
	h      tHeap
	length int
}

// Heap functions
func (h tHeap) Len() int { return len(h) }

func (h tHeap) Less(i int, j int) bool {
	return h[i].priority > h[j].priority
}

func (h *tHeap) Push(x any) {
	n := len(*h)
	item := x.(*Item)
	item.index = n
	*h = append(*h, item)
}

func (h *tHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil     // don't stop the GC from reclaiming the item eventually
	item.priority = -1 // for safety
	*h = old[0 : n-1]
	return item
}

func (h tHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Priority Queue Functions

func NewPriorityQueue[T comparable]() PriorityQueue[T] {

	h := make(tHeap, 0)
	heap.Init(&h)

	return PriorityQueue[T]{h, 0}

}

func (np PriorityQueue[T]) Len() int {
	return np.length
}

func (pq *PriorityQueue[T]) Push(val T, priority int) {
	item := new(Item)
	item.value = val
	item.priority = priority
	heap.Push(&pq.h, item)
	pq.length++
}

func (pq *PriorityQueue[T]) Pop() option.Option[T] {
	if pq.length == 0 {
		return option.OptionNil[T]()
	}
	v := heap.Pop(&pq.h).(*Item)

	pq.length--
	return option.OptionVal(v.value.(T))

}

// Updates the first item in the priority queue with item.value=val
func (pq *PriorityQueue[T]) Update(val T, priority int) {
	var item *Item = nil

	for _, i := range pq.h {
		if i.value.(T) == val {
			item = i
			break
		}
	}

	if item != nil {
		item.value = val
		item.priority = priority
		heap.Fix(&pq.h, item.index)
	}
}

// Removes the first item with value val
func (pq *PriorityQueue[T]) Remove(val T) {
	n := NewPriorityQueue[T]()
	for _, i := range pq.h {
		if i.value.(T) != val {
			n.Push(i.value.(T), i.priority)
		}
	}

	*pq = n
}

func (pq *PriorityQueue[T]) To_list() []T {
	res := make([]T, pq.Len())

	for i := 0; i < pq.Len(); i++ {
		val := pq.Pop()
		if val.GetSet() {
			res[i] = val.GetVal()
		}
	}
	return res
}
