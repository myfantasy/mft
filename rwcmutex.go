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
		return // it neet to test!!!!
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

// LockC - try locks mutex with context
func (m *RWCMutex) LockC(ctx context.Context) bool {
	if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {
		return true
	}

	// Slow way
	return m.lockSC(ctx)
}

// LockD - try locks mutex with time duretion
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

// RLockC - try read locks mutex with context
func (m *RWCMutex) RLockC(ctx context.Context) bool {
	k := atomic.LoadInt32(&m.state)
	if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
		return true
	}

	// Slow way
	return m.rlockSC(ctx)
}

// RLockD - try read locks mutex with time duretion
func (m *RWCMutex) RLockD(d time.Duration) bool {
	k := atomic.LoadInt32(&m.state)
	if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
		return true
	}

	// Slow way
	return m.rlockSD(d)
}

func (m *RWCMutex) lockS() {
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return
		}

		ch := m.chGet()

		select {
		case <-ch:
		}

	}
}

func (m *RWCMutex) lockSC(ctx context.Context) bool {
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return true
		}

		if ctx == nil {
			return false
		}

		ch := m.chGet()

		select {
		case <-ch:
		case <-ctx.Done():
			return false
		}

	}
}

func (m *RWCMutex) lockSD(d time.Duration) bool {
	// may be use context.WithTimeout(context.Background(), d) however NO it's not fun
	t := time.After(d)
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, rwcmLocked) {

			return true
		}

		ch := m.chGet()

		select {
		case <-ch:
		case <-t:
			return false
		}

	}
}

func (m *RWCMutex) rlockS() {

	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return
		}

		ch := m.chGet()

		select {
		case <-ch:
		}

	}

}

func (m *RWCMutex) rlockSC(ctx context.Context) bool {

	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return true
		}

		if ctx == nil {
			return false
		}

		ch := m.chGet()

		select {
		case <-ch:
		case <-ctx.Done():
			return false
		}

	}

}

func (m *RWCMutex) rlockSD(d time.Duration) bool {

	t := time.After(d)
	var k int32
	for {
		k = atomic.LoadInt32(&m.state)
		if k >= 0 && atomic.CompareAndSwapInt32(&m.state, k, k+1) {
			return true
		}

		ch := m.chGet()

		select {
		case <-ch:
		case <-t:
			return false
		}

	}

}
