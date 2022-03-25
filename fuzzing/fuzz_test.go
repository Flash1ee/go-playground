package main

import (
	"reflect"
	"sort"
	"testing"

	"fuzzing/quick_sort"
)

func TestQuickSort(t *testing.T) {
	testcases := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{3, 2, 1}, []int{1, 2, 3}},
		{[]int{-1, -3, -2, 0}, []int{-3, -2, -1, 0}},
		{[]int{-1, -2, 1}, []int{-2, -1, 1}},
		{[]int{1}, []int{1}},
	}
	for _, tc := range testcases {
		res := quick_sort.QuickSort(tc.in)
		if !reflect.DeepEqual(tc.want, res) {
			t.Errorf("Receive quick_sort res: %q, want %q", res, tc.want)
		}
	}
}

func FuzzQuickSort(f *testing.F) {
	testcases := [][]int{
		{1, 2, 3},
		{3, 2, 1},
		{1, 3, 2, 4},
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig []int) {
		exp := make([]int, len(orig))
		copy(exp, orig)
		res := quick_sort.QuickSort(orig)

		sort.Ints(exp)
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("Before: %q, after: %q Expect: %q", orig, res, exp)
		}
	})
}
