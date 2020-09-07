package mft

import (
	"testing"
)

func BenchmarkRWCMutexLockUnlock(b *testing.B) {
	mx := RWCMutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkRWCMutexRLockRUnlock(b *testing.B) {
	mx := RWCMutex{}

	for i := 0; i < b.N; i++ {
		mx.RLock()
		mx.RUnlock()
	}
}

// TODO: make normal test
