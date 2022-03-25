package main

import (
	"fmt"

	"fuzzing/reverse"
)

func main() {
	input := "The quick brown fox jumped over the lazy собакка"
	rev, _ := reverse.Reverse(input)
	doubleRev, _ := reverse.Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
