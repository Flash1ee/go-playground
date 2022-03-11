package main

import (
	"testing"
)

func bench(t *testing.B, c Incrementable) {
	expectedCount := t.N

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
	if c.Sum() != int64(expectedCount) {
		t.Fatal("the sum is not equal to the number of increments")
	}
}
func Benchmark_sync(t *testing.B) {
	x := NewSync()

	bench(t, x)
}
func Benchmark_atomic(t *testing.B) {
	x := NewAtomic()

	bench(t, x)
}

func Benchmark_channels(t *testing.B) {
	x := NewChannel()

	bench(t, x)
}

func Benchmark_buf_channels(t *testing.B) {
	x := NewBufChannel()

	bench(t, x)
}
