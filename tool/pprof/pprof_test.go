package main

import "testing"

func TestAdd(t *testing.T) {
	s := Add("xiaoming")
	if s == "" {
		t.Errorf("Test.Add error!")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add("xiaoming ")
	}
}
