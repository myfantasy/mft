package generator

import (
	"testing"
)

func BenchmarkRvGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RvGet()
	}
}

func BenchmarkRvGet2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RvGet2()
	}
}

func TestRvGet(t *testing.T) {
	RvGet()
}

func TestRvGet2(t *testing.T) {
	RvGet2()
}
