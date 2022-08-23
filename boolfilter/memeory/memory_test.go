package memeory

import (
	"github.com/balance/boolfilter"
	"strconv"
	"testing"
)

func TestMemoryFilter(t *testing.T) {
	// Initial a memory filter
	memFilter := NewMemoryFilter(make([]byte, 10240), boolfilter.DefaultHash...)
	// Push 2000-3000 numbers to the filter.
	// 把2000-3000的数字压入过滤器
	for i := 2000; i <= 3000; i++ {
		memFilter.Push([]byte(strconv.Itoa(i)))
	}
	// Check whether 2500-3000 and 3001-3500 exist in the filter or not.
	// 查看2500-3000，以及3001-3500是否存在于过滤器中
	for i := 2500; i < 3500; i++ {
		r, _ := memFilter.Exists([]byte(strconv.Itoa(i)))
		t.Logf("%d, %t", i, r)
	}
}
