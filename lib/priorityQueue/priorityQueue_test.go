package priorityQueue

import (
	"github.com/khorcarol/AgentOfThings/lib/option"
	"testing"
)

func Test_pops_correct_order(t *testing.T) {
	ls := [6]int{3, 4, 1, 6, 2, 5}

	pq := NewPriorityQueue[int]()
	for i, e := range ls {
		pq.Push(e, i)
	}

	for i := 0; i < 6; i++ {
		x := pq.Pop()
		if !x.GetSet() {
			t.Fatalf("Fail on %s: Expected a value, got none", t.Name())
		}
		val := x.GetVal()
		if ls[5-i] != val {
			t.Fatalf("Fail on %s: Expected %d, got %d", t.Name(), ls[i], val)
		}
	}
}

func Test_empty_returns_none(t *testing.T) {
	pq := NewPriorityQueue[int]()

	res := pq.Pop()

	if res != option.OptionNil[int]() {
		t.Fatalf("Fail on %s: Pop operation did not return Nil", t.Name())
	}
}

func Test_only_pops_items(t *testing.T) {
	pq := NewPriorityQueue[int]()

	pq.Push(1, 1)
	pq.Pop()

	v := pq.Pop()
	if v != option.OptionNil[int]() {
		t.Fatalf("Fail on %s: Pop operation did not return Nil", t.Name())
	}
}

func Test_removes_item(t *testing.T) {
	pq := NewPriorityQueue[int]()

	pq.Push(1, 1)
	pq.Push(2, 2)
	pq.Push(3, 3)

	pq.Remove(2)

	v1 := pq.Pop().GetVal()
	v2 := pq.Pop().GetVal()

	if v1 != 3 {
		t.Errorf("Fail on %s: First pop meant to be 3, got: %d", t.Name(), v1)
	}
	if v2 != 1 {
		t.Errorf("Fail on %s: Second pop meant to be 1, got: %d", t.Name(), v2)
	}
}

func Test_updates_item(t *testing.T) {
	pq := NewPriorityQueue[int]()

	for i := 1; i <= 3; i++ {
		pq.Push(i, i)
	}

	pq.Update(2, 5)

	v1 := pq.Pop().GetVal()
	v2 := pq.Pop().GetVal()
	v3 := pq.Pop().GetVal()

	if v1 != 2 {
		t.Errorf("Fail on %s: First pop meant to be 2, got: %d", t.Name(), v1)
	}
	if v2 != 3 {
		t.Errorf("Fail on %s: First pop meant to be 3, got: %d", t.Name(), v1)
	}
	if v3 != 1 {
		t.Errorf("Fail on %s: First pop meant to be 1, got: %d", t.Name(), v1)
	}
}
