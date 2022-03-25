package reverse

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
		isErr    bool
	}{
		{"Hello, world", "dlrow ,olleH", false},
		{" ", " ", false},
		{"!12345", "54321!", false},
		{"\xfa", "afx\\", true},
	}
	for _, tc := range testcases {
		rev, err := Reverse(tc.in)
		if err != nil {
			if tc.isErr {
				continue
			}
			t.Error(err)
		}
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
