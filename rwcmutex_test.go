package mft

import (
	"testing"
	"time"
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

func TestRWCMutex(t *testing.T) {

	var mx RWCMutex

	mx.Lock()
	mx.Unlock()

	mx.Lock()
	t1 := mx.RLockD(time.Millisecond)
	if t1 {
		t.Fatal("TestRWCMutex t1 fail R lock duration")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock()
	}()

	t2 := mx.RLockD(10 * time.Millisecond)
	t3 := mx.RLockD(10 * time.Millisecond)

	if !t2 {
		t.Fatal("TestRWCMutex t2 fail R lock duration")
	}
	if !t3 {
		t.Fatal("TestRWCMutex t2 fail R lock duration")
	}

}

// TODO: make normal test
