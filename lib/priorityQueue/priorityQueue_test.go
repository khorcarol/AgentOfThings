package priorityQueue 

import (
	"testing"
	"github.com/khorcarol/AgentOfThings/lib/option"
)

func Test_pops_correct_order(t *testing.T){
	ls := [6]int{3, 4, 1, 6, 2, 5}

	pq := NewPriorityQueue[int]()
	for i, e := range ls {
		pq.Push(e, i)
	}

	for i:=0; i<6; i++{
		x := pq.Pop()
		if !x.GetSet(){
			t.Fatalf("Fail on %s: Expected a value, got none", t.Name())
		}
		val := x.GetVal()
		if ls[5-i] != val {
			t.Fatalf("Fail on %s: Expected %d, got %d", t.Name(), ls[i], val)
		}
	}
}

func Test_empty_returns_none(t *testing.T){
	pq := NewPriorityQueue[int]()

	res := pq.Pop()

	if res != option.OptionNil[int](){
		t.Fatalf("Fail on %s: Pop operation did not return Nil", t.Name())
	}
}

func Test_only_pops_items(t *testing.T){
	pq := NewPriorityQueue[int]()

	pq.Push(1, 1)
	pq.Pop()

	v := pq.Pop()
	if v != option.OptionNil[int]() {
		t.Fatalf("Fail on %s: Pop operation did not return Nil", t.Name())
	}
}

