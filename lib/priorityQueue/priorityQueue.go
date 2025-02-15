package lib

import (
	"container/heap"
)

type Item struct {
	value any
	priority int
	index int
}

type THeap []*Item

type PriorityQueue[T any] struct {
	h THeap
	length int
}

// Heap functions
func (h *THeap) Len() int { return len(*h) }

func (h *THeap) Less(i int, j int) bool {
	return (*h)[i].priority < (*h)[j].priority
}

func (h *THeap) Push(x any){
	n := len(*h)
	item := x.(*Item) 
	item.index = n
	*h = append(*h, item)
}

func (h *THeap) Pop() any{
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.priority= -1 // for safety
	*h = old[0 : n-1]
	return item.value
}

func (h THeap) Swap(i int, j int){
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}


// Priority Queue Functions

func newPriorityQueue[T any](size int) PriorityQueue[T]{
	
	h := make(THeap, size)
	pq := PriorityQueue[T]{h, 0}

	heap.Init(&pq.h)

	return pq
}

func (np PriorityQueue[T]) Len() int {
	return np.length
}

func (pq PriorityQueue[T]) Push(val T, priority int){
	item := new(Item)
	item.value = val
	item.priority = priority
	heap.Push(&pq.h, item)
}

func (pq PriorityQueue[T]) Pop() Option[T]{
	if pq.length == 0{
		return optionNil[T]()
	} else{
		pq.length --
		return optionVal(pq.h.Pop().(T))
	}
}

func (pq PriorityQueue[T]) to_list() []T{
	res := make([]T, pq.Len())
	cp := pq
	
	for i:=0; i<cp.Len(); i++ {
		val := cp.Pop()
		if val.set {
			res[i] = val.val
		}
	}
	return res
}

func test() int{
	return 1;
}
