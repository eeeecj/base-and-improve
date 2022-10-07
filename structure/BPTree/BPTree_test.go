package BPTree

import (
	"fmt"
	"testing"
)

func TestBPTree_Get(t *testing.T) {
	tt := NewBPTree(3)
	tt.Insert(1, 2)
	tt.Insert(2, 3)
	tt.Insert(3, 4)
	fmt.Println(tt.Get(2))
	tt.Delete(3)
	fmt.Println(tt.Get(3))
}
