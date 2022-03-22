package main

import (
	"reflect"
	"testing"
)

func Test_getUniqNumSlice(t *testing.T) {
	type args struct {
		src []int
	}
	tests := []struct {
		name    string
		args    args
		wantDst []int
	}{
		{
			name: "simple, uniq slice",
			args: args{
				src: []int{1, 2, 3, 4, 5},
			},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name: "simple, need delete nums",
			args: args{
				src: []int{1, 2, 3, 4, 5, 5},
			},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name: "simple, need delete nums from the start, end, center",
			args: args{
				src: []int{1, 1, 1, 2, 3, 3, 3, 4, 5, 5},
			},
			wantDst: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty",
			args: args{
				src: []int{},
			},
			wantDst: []int{},
		},
		{
			name: "nil value",
			args: args{
				src: nil,
			},
			wantDst: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDst := getUniqNumSlice(tt.args.src); !reflect.DeepEqual(gotDst, tt.wantDst) {
				t.Errorf("getUniqNumSlice() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

func Test_getTemplateUniq_Int(t *testing.T) {
	tests := GetTestDataForUniq_Int(t)

	for _, tCase := range tests.cases {
		t.Run(tCase.name, func(t *testing.T) {
			if gotDst := getTemplateUniq(tCase.src); !reflect.DeepEqual(gotDst, tCase.wantDst) {
				t.Errorf("getUniqNumSlice() - %v, test - %v, res = %v, want %v", tCase.name, tCase.name, gotDst, tCase.wantDst)
			}
		})
	}
}
func Test_getTemplateUniq_String(t *testing.T) {
	tests := GetTestDataForUniq_String(t)

	for _, tCase := range tests.cases {
		t.Run(tCase.name, func(t *testing.T) {
			if gotDst := getTemplateUniq(tCase.src); !reflect.DeepEqual(gotDst, tCase.wantDst) {
				t.Errorf("getUniqNumSlice() - %v, test - %v, res = %v, want %v", tCase.name, tCase.name, gotDst, tCase.wantDst)
			}
		})
	}
}
