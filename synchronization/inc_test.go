package main

import (
	"testing"
)

func Benchmark_sync(t *testing.B) {
	x := NewSync()
	expectedCount := t.N

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x.Inc()
		}
	})
	if x.Sum() != int64(expectedCount) {
		t.Fatal("the sum is not equal to the number of increments")
	}
}
func Benchmark_atomic(t *testing.B) {
	x := NewAtomic()
	expectedCount := t.N

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x.Inc()
		}
	})
	if x.Sum() != int64(expectedCount) {
		t.Fatal("the sum is not equal to the number of increments")
	}
}

func Benchmark_channels(t *testing.B) {
	x := NewChannel()
	res := ChannelWorker(*x, t.N)
	expectedCount := t.N

	if res != int64(expectedCount) {
		t.Fatal("the sum is not equal to the number of increments")
	}
}
