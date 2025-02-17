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

type PriorityQueue[T any] struct {
	h      tHeap
	length int
}

type Pair[T any] struct {
	First  T
	Second T
	Score  int
}

// Heap functions
func (h *tHeap) Len() int { return len(*h) }

func (h *tHeap) Less(i int, j int) bool {
	return (*h)[i].priority < (*h)[j].priority
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
	return item.value
}

func (h tHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Priority Queue Functions

func NewPriorityQueue[T any](size int) PriorityQueue[T] {

	h := make(tHeap, size)
	pq := PriorityQueue[T]{h, 0}

	heap.Init(&pq.h)

	return pq
}

func (np PriorityQueue[T]) Len() int {
	return np.length
}

func (pq PriorityQueue[T]) Push(val T, priority int) {
	item := new(Item)
	item.value = val
	item.priority = priority
	heap.Push(&pq.h, item)
}

func (pq PriorityQueue[T]) Pop() option.Option[T] {
	if pq.length == 0 {
		return option.OptionNil[T]()
	} else {
		pq.length--
		return option.OptionVal(pq.h.Pop().(T))
	}
}

func (pq PriorityQueue[T]) To_list() []T {
	res := make([]T, pq.Len())
	cp := pq

	for i := 0; i < cp.Len(); i++ {
		val := cp.Pop()
		if val.GetSet() {
			res[i] = val.GetVal()
		}
	}
	return res
}

func (pq PriorityQueue[T]) To_pairs() []T {
	return pq.To_list()
}
