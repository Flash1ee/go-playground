package main

import (
	"testing"
)

type TestTable[T any] struct {
	name  string
	cases []TestCase[T]
}

type TestCase[T any] struct {
	name    string
	src     []T
	wantDst []T
}

func GetTestDataForUniq_Int(t *testing.T) TestTable[int] {
	t.Helper()
	x := TestTable[int]{
		name: "int tests",
	}
	x.cases = []TestCase[int]{
		{
			name:    "int_uniq_slice",
			src:     []int{1, 2, 3, 4, 5},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name:    "int_need_delete_nums",
			src:     []int{1, 2, 3, 4, 5, 5},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name:    "int_need_delete_nums_from_the_start_end,_center",
			src:     []int{1, 1, 1, 2, 3, 3, 3, 4, 5, 5},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name:    "int_empty",
			src:     []int{},
			wantDst: []int{},
		},
		{
			name:    "int_nil_value",
			src:     nil,
			wantDst: []int{},
		},
	}
	return x
}
func GetTestDataForUniq_String(t *testing.T) TestTable[string] {
	t.Helper()
	return TestTable[string]{
		name: "test strings",
		cases: []TestCase[string]{
			{
				name:    "string_uniq_slice",
				src:     []string{"a", "one", "b", "two"},
				wantDst: []string{"a", "one", "b", "two"},
			},
			{
				name:    "string_need_delete_values_at_the_end",
				src:     []string{"a", "one", "b", "two", "two"},
				wantDst: []string{"a", "one", "b", "two"},
			},
			{
				name:    "string_need_delete_from_the_start_end_center",
				src:     []string{"a", "a", "one", "b", "b", "b", "two", "three", "three"},
				wantDst: []string{"a", "one", "b", "two", "three"},
			},
			{
				name:    "string_empty",
				src:     []string{},
				wantDst: []string{},
			},
			{
				name:    "string_nil_value",
				src:     nil,
				wantDst: []string{},
			},
		},
	}
}
