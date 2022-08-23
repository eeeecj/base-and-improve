package hash

import "testing"

func BenchmarkHashCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashCode([]byte("xiaoming"))
	}
}
