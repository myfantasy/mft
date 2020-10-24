package mft

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const rwcmLocked int32 = -1

// RWCMutex - Read Write Context Mutex
type RWCMutex struct {
	state int32
	mx    sync.Mutex
	ch    chan struct{}
}

func (m *RWCMutex) chGet() chan struct{} {
	m.mx.Lock()
	if m.ch == nil {
		m.ch = make(chan struct{}, 1)
	}
	r := m.ch
	m.mx.Unlock()
	return r
}

func (m *RWCMutex) chClose() {
	// it's need only when exists parallel
	// to make faster need add counter to add drop listners of chan

	if m.ch == nil {
		return // it neet to test!!!! theoreticly works when channel get operation is befor atomic operations
	}

	var o chan struct{}
	m.mx.Lock()
	if m.ch != nil {
		o = m.ch
		m.ch = nil
	}
	m.mx.Unlock()
	if o != nil {
		close(o)
	}
}

// Lock - locks mutex
func (m *RWCMutex) Lock() {
	if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

		return
	}

	// Slow way
	m.lockS()
}

// Unlock - unlocks mutex
func (m *RWCMutex) Unlock() {
	if atomic.CompareAndSwapInt32(&m.state, rwcmLocked, 0) {
		m.chClose()
		return
	}

	panic(ErrorCS(-10001, "Unlock fail"))
}

// TryLock - try locks mutex with context
func (m *RWCMutex) TryLock(ctx context.Context) bool {
	if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {
		return true
	}

	// Slow way
	return m.lockSC(ctx)
}

// LockD - try locks mutex with time duration
func (m *RWCMutex) LockD(d time.Duration) bool {
	if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {
		return true
	}

	// Slow way
	return m.lockSD(d)
}

// RLock - read locks mutex
func (m *RWCMutex) RLock() {
	k := atomic.LoadInt32(&m.state)
	if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
		return
	}

	// Slow way
	m.rlockS()
}

// RUnlock - unlocks mutex
func (m *RWCMutex) RUnlock() {
	i := atomic.AddInt32(&m.state, -1)
	if i > 0 {
		return
	} else if i == 0 {
		m.chClose()
		return
	}

	panic(ErrorCS(-10002, "RUnlock fail"))
}

// RTryLock - try read locks mutex with context
func (m *RWCMutex) RTryLock(ctx context.Context) bool {
	k := atomic.LoadInt32(&m.state)
	if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
		return true
	}

	// Slow way
	return m.rlockSC(ctx)
}

// RLockD - try read locks mutex with time duration
func (m *RWCMutex) RLockD(d time.Duration) bool {
	k := atomic.LoadInt32(&m.state)
	if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
		return true
	}

	// Slow way
	return m.rlockSD(d)
}

func (m *RWCMutex) lockS() {
	ch := m.chGet()
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return
		}

		select {
		case <-ch:
			ch = m.chGet()
		}
	}

}

func (m *RWCMutex) lockSC(ctx context.Context) bool {
	ch := m.chGet()
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return true
		}

		if ctx == nil {
			return false
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-ctx.Done():
			return false
		}

	}
}

func (m *RWCMutex) lockSD(d time.Duration) bool {
	// may be use context.WithTimeout(context.Background(), d) however NO it's not fun
	t := time.After(d)
	ch := m.chGet()
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-t:
			return false
		}

	}
}

func (m *RWCMutex) rlockS() {

	ch := m.chGet()
	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return
		}

		select {
		case <-ch:
			ch = m.chGet()
		}

	}

}

func (m *RWCMutex) rlockSC(ctx context.Context) bool {
	ch := m.chGet()
	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return true
		}

		if ctx == nil {
			return false
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-ctx.Done():
			return false
		}

	}

}

func (m *RWCMutex) rlockSD(d time.Duration) bool {
	ch := m.chGet()
	t := time.After(d)
	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-t:
			return false
		}

	}

}
