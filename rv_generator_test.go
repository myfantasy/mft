package mft

import (
	"testing"
	"time"
)

func BenchmarkRvGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RvGet()
	}
}

// func BenchmarkRvGet2(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		RvGet2()
// 	}
// }

func BenchmarkRvGetPart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RvGetPart()
	}
}

func BenchmarkSleepD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Microsecond)
	}
}

func TestRvGet(t *testing.T) {
	RvGet()
}

// func TestRvGet2(t *testing.T) {
// 	RvGet2()
// }

func TestRvGetPart(t *testing.T) {
	RvGetPart()
}
