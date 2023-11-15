package main

import (
	"fmt"
	"testing"
)

func TestStripSlice(t *testing.T) {
	x := StripSlice([]string{"a", "b", "c"}, "b")
	fmt.Println(x)
}
