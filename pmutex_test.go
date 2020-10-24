package mft

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkPMutexLockUnlock(b *testing.B) {
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		k := mx.Lock()
		mx.Unlock(k)
	}
}

func BenchmarkPMutexRLockRUnlock(b *testing.B) {
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		k := mx.Lock()
		mx.Unlock(k)
	}
}

func TestPMutex(t *testing.T) {

	var mx PMutex

	k := mx.Lock()
	mx.Unlock(k)

	k = mx.Lock()
	k1, t1 := mx.RLockD(time.Millisecond)
	if t1 {
		t.Fatal("TestPMutex t1 fail R lock duration")
	}

	if mx.Status(k) != Lock {
		t.Fatal("TestPMutex t0 fail status " + strconv.Itoa(mx.Status(k)))
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock(k)
	}()

	k2, t2 := mx.RLockD(10 * time.Millisecond)

	if mx.Status(k) != Free {
		t.Fatal("TestPMutex t0 fail status unlock " + strconv.Itoa(mx.Status(k)))
	}

	k3, t3 := mx.RLockD(10 * time.Millisecond)

	if !t2 {
		t.Fatal("TestPMutex t2 fail R lock duration")
	}
	if !t3 {
		t.Fatal("TestPMutex t3 fail R lock duration")
	}

	if mx.Status(k3) != ReadLock {
		t.Fatal("TestPMutex t3 fail status " + strconv.Itoa(mx.Status(k3)))
	}

	mx.TryUnlock(k2)
	mx.TryUnlock(k3)

	if mx.Status(k3) != Free {
		t.Fatal("TestPMutex t3 fail status unlock " + strconv.Itoa(mx.Status(k3)))
	}

	if mx.Status(k1) != Free {
		t.Fatal("TestPMutex t1 fail status " + strconv.Itoa(mx.Status(k1)))
	}

	mx.TryUnlock(k1)

}

func TestPMutexPromote(t *testing.T) {

	var mx PMutex

	k := mx.RLock()
	k1 := mx.RLock()

	if mx.Status(k1) != ReadLock {
		t.Fatal("TestPMutex t1 fail status " + strconv.Itoa(mx.Status(k1)))
	}

	ok := mx.PromoteD(time.Millisecond, k1)

	if ok {
		t.Fatal("TestPMutex t0 must be not promote")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock(k)
	}()

	ok = mx.PromoteD(10*time.Millisecond, k1)

	if !ok {
		t.Fatal("TestPMutex t0 must be promote")
	}

	if mx.Status(k1) != Lock {
		t.Fatal("TestPMutex t1 fail status promote " + strconv.Itoa(mx.Status(k1)))
	}

	_, ok2 := mx.RLockD(1 * time.Millisecond)

	if ok2 {
		t.Fatal("TestPMutex t0 must be fail by timeout")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Reduce(k1)
	}()

	_, ok3 := mx.RLockD(10 * time.Millisecond)

	if !ok3 {
		t.Fatal("TestPMutex t0 must be fail by timeout")
	}

}

// TODO: make normal test
