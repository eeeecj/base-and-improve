package heap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeap(t *testing.T) {
	h := NewHeap()
	h.AddSlice([]int{12, 4, 6, 1, 123, 65, 78})
	t.Log(h.BuildMaxHeap(h.Len() - 1))
	assert.Equal(t, []int{1, 4, 6, 12, 65, 78, 123}, h.Sort())
	assert.Equal(t, []int{123, 78, 65, 12, 6, 4, 1}, h.SortDesc())
}
