package sort

import "testing"

func TestQuick(t *testing.T) {
	r := []int{1, 2, 2, 3, 3, 4, 5, 5, 6}
	b := QuickSort([]int{3, 2, 3, 1, 2, 4, 5, 5, 6})
	for i := 0; i < len(b); i++ {
		if b[i] != r[i] {
			t.Errorf("the error sort:%v", b)
		}
	}
}
