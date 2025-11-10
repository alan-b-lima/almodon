package heap_test

import (
	"math/rand/v2"
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/heap"
)

type Int int

func (i0 Int) Less(i1 Int) bool {
	return i0 < i1
}

func TestExpectedBehavior(t *testing.T) {
	var heap Heap[Int]

	for range 10000 {
		op := rand.IntN(3)

		switch op {
		case 0:
			heap.Push(Int(rand.IntN(1000)))

		case 1:
			if heap.Len() > 0 {
				heap.Pop()
			}

		case 2:
			if heap.Len() == 0 {
				break
			}

			str := heap.String()
			if e, o := heap.Peek(), heap.Pop(); e != o {
				t.Errorf("%d should have been equal to %d in %s", e, o, str)
			}
		}
	}
}
