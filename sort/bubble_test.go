package sort

import "testing"

func TestBubble(t *testing.T) {
	r := []int{1, 4, 5, 6, 7}
	b := Bubble([]int{7, 5, 6, 4, 1})
	for i := 0; i < len(b); i++ {
		if b[i] != r[i] {
			t.Errorf("the error sort:%v", b)
		}
	}
}
