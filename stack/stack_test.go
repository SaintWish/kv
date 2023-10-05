package stack

import (
	"fmt"
	"testing"
)

const maxEntries = 500

func TestMoveToBack_KeyString(t *testing.T) {
	m := New[int](2048)

	m.Push(1337)
	m.Push(1338)
	m.Push(1339)
	m.Push(1400)

	fmt.Println(m.Stack())

	m.MoveToBack(0)

	fmt.Println(m.Stack())
}

func Benchmark_KV1_SZ2048(b *testing.B) {
	m := New[int](2048)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := 1; e <= maxEntries; e++ {
			m.Push(e)
		}

		for e := 1; e <= maxEntries; e++ {
			m.Pop()
		}
	}
}