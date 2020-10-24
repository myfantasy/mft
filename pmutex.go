package mft

import (
	"context"
	"sync"
	"time"
)

// PMutex - mutex with identity of lock
// read write mutex with context and change Priority
type PMutex struct {
	mx     sync.Mutex
	th     map[int64]struct{}
	idL    int64
	id     int64
	ch     chan struct{}
	isInit bool
}

func (m *PMutex) chGet() chan struct{} {
	m.mx.Lock()
	if m.ch == nil {
		m.ch = make(chan struct{}, 1)
	}
	r := m.ch
	m.mx.Unlock()
	return r
}

func (m *PMutex) chClose() {
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

func (m *PMutex) init() {

	m.mx.Lock()
	if !m.isInit {
		m.isInit = true
		m.th = make(map[int64]struct{})
	}
	m.mx.Unlock()
}

// Lock - locks mutex
func (m *PMutex) Lock() (i int64) {

	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 && len(m.th) == 0 {
		m.id++
		i = m.id
		m.idL = i

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i
	}

	return m.lockS()
}

func (m *PMutex) lockS() (i int64) {

	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 && len(m.th) == 0 {
			m.id++
			i = m.id
			m.idL = i

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i
		}

		select {
		case <-ch:
			ch = m.chGet()
		}
	}
}

// Unlock - unlocks mutex (also works as RUnlock)
func (m *PMutex) Unlock(i int64) {

	ok := false

	m.mx.Lock()

	_, ok = m.th[i]
	if m.idL == 0 && ok {
		delete(m.th, i)
	}

	if m.idL == i {
		m.idL = 0

		delete(m.th, i)

		ok = true
	}

	m.mx.Unlock()

	if !ok {
		panic(ErrorCS(-10101, "Unlock fail"))
	}

	m.chClose()
}

// TryUnlock - unlocks mutex without panic
func (m *PMutex) TryUnlock(i int64) {

	ok := false

	m.mx.Lock()

	_, ok = m.th[i]
	if m.idL == 0 && ok {
		delete(m.th, i)
	}

	if m.idL == i {
		m.idL = 0

		delete(m.th, i)

		ok = true
	}

	m.mx.Unlock()

	if ok {
		m.chClose()
	}
}

// TryLock - try locks mutex with context
func (m *PMutex) TryLock(ctx context.Context) (i int64, ok bool) {
	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 && len(m.th) == 0 {
		m.id++
		i = m.id
		m.idL = i

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i, true
	}

	// Slow way
	return m.lockSC(ctx)
}

func (m *PMutex) lockSC(ctx context.Context) (i int64, ok bool) {

	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 && len(m.th) == 0 {
			m.id++
			i = m.id
			m.idL = i

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i, true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-ctx.Done():
			return 0, false
		}
	}
}

// LockD - try locks mutex with time duration
func (m *PMutex) LockD(d time.Duration) (i int64, ok bool) {
	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 && len(m.th) == 0 {
		m.id++
		i = m.id
		m.idL = i

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i, true
	}

	// Slow way
	return m.lockSD(d)
}

func (m *PMutex) lockSD(d time.Duration) (i int64, ok bool) {
	t := time.After(d)
	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 && len(m.th) == 0 {
			m.id++
			i = m.id
			m.idL = i

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i, true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-t:
			return 0, false
		}
	}
}

// RLock - read locks mutex
func (m *PMutex) RLock() (i int64) {

	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 {
		m.id++
		i = m.id

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i
	}

	return m.rlockS()
}

func (m *PMutex) rlockS() (i int64) {

	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 {
			m.id++
			i = m.id

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i
		}

		select {
		case <-ch:
			ch = m.chGet()
		}
	}
}

// RLockD - read locks mutex with time duration
func (m *PMutex) RLockD(d time.Duration) (i int64, ok bool) {

	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 {
		m.id++
		i = m.id

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i, true
	}

	return m.rlockSD(d)
}

func (m *PMutex) rlockSD(d time.Duration) (i int64, ok bool) {

	t := time.After(d)
	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 {
			m.id++
			i = m.id

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i, true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-t:
			return 0, false
		}
	}
}

// RTryLock - read locks mutex with context
func (m *PMutex) RTryLock(ctx context.Context) (i int64, ok bool) {

	if !m.isInit {
		m.init()
	}

	m.mx.Lock()

	if m.idL == 0 {
		m.id++
		i = m.id

		m.th[i] = struct{}{}
	}

	m.mx.Unlock()

	if i > 0 {
		return i, true
	}

	return m.rlockSC(ctx)
}

func (m *PMutex) rlockSC(ctx context.Context) (i int64, ok bool) {

	ch := m.chGet()
	for {
		m.mx.Lock()

		if m.idL == 0 {
			m.id++
			i = m.id

			m.th[i] = struct{}{}
		}

		m.mx.Unlock()

		if i > 0 {
			return i, true
		}

		select {
		case <-ch:
			ch = m.chGet()

		case <-ctx.Done():
			return 0, false
		}
	}
}

// Promote - promotes read lock mutex to lock
func (m *PMutex) Promote(i int64) bool {

	ok := false
	l := false

	m.mx.Lock()

	_, ok = m.th[i]

	if ok {
		if m.idL == 0 && len(m.th) == 1 {
			l = true
			m.idL = i
		}

		if m.idL == i {
			l = true
		}
	}

	m.mx.Unlock()

	if !ok {
		return false
	}

	if l {
		return true
	}

	return m.promoteS(i)
}

func (m *PMutex) promoteS(i int64) bool {

	ok := false
	l := false

	ch := m.chGet()
	for {
		m.mx.Lock()

		_, ok = m.th[i]

		if ok {
			if m.idL == 0 && len(m.th) == 1 {
				l = true
				m.idL = i
			}

			if m.idL == i {
				l = true
			}
		}

		m.mx.Unlock()

		if !ok {
			return false
		}

		if l {
			return true
		}

		select {
		case <-ch:
			ch = m.chGet()
		}
	}
}

// PromoteD - promotes read lock mutex to lock with time duration
func (m *PMutex) PromoteD(d time.Duration, i int64) bool {

	ok := false
	l := false

	m.mx.Lock()

	_, ok = m.th[i]

	if ok {
		if m.idL == 0 && len(m.th) == 1 {
			l = true
			m.idL = i
		}

		if m.idL == i {
			l = true
		}
	}

	m.mx.Unlock()

	if !ok {
		return false
	}

	if l {
		return true
	}

	return m.promoteSD(d, i)
}

func (m *PMutex) promoteSD(d time.Duration, i int64) bool {

	t := time.After(d)

	ok := false
	l := false

	ch := m.chGet()
	for {
		m.mx.Lock()

		_, ok = m.th[i]

		if ok {
			if m.idL == 0 && len(m.th) == 1 {
				l = true
				m.idL = i
			}

			if m.idL == i {
				l = true
			}
		}

		m.mx.Unlock()

		if !ok {
			return false
		}

		if l {
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

// TryPromote - promotes read lock mutex to lock with context
func (m *PMutex) TryPromote(ctx context.Context, i int64) bool {

	ok := false
	l := false

	m.mx.Lock()

	_, ok = m.th[i]

	if ok {
		if m.idL == 0 && len(m.th) == 1 {
			l = true
			m.idL = i
		}

		if m.idL == i {
			l = true
		}
	}

	m.mx.Unlock()

	if !ok {
		return false
	}

	if l {
		return true
	}

	return m.promoteSC(ctx, i)
}

func (m *PMutex) promoteSC(ctx context.Context, i int64) bool {

	ok := false
	l := false

	ch := m.chGet()
	for {
		m.mx.Lock()

		_, ok = m.th[i]

		if ok {
			if m.idL == 0 && len(m.th) == 1 {
				l = true
				m.idL = i
			}

			if m.idL == i {
				l = true
			}
		}

		m.mx.Unlock()

		if !ok {
			return false
		}

		if l {
			return true
		}

		select {
		case <-ch:
			ch = m.chGet()
		case <-ctx.Done():
			return false
		}
	}
}

// Reduce - reduce lock mutex to read lock
func (m *PMutex) Reduce(i int64) bool {

	ok := false
	l := false

	m.mx.Lock()

	_, ok = m.th[i]

	if ok {
		if m.idL == 0 {
			l = true
		}

		if m.idL == i {
			m.idL = 0
			l = true
		}
	}

	m.mx.Unlock()

	if !ok {
		return false
	}

	m.chClose()

	return l
}

// Status - status of lock
// 0 - free (not lock)
// 1 - read lock
// 2 - (write) lock
func (m *PMutex) Status(i int64) int {

	ok := false
	l := false

	m.mx.Lock()

	_, ok = m.th[i]

	if ok {

		if m.idL == i {
			l = true
		}
	}

	m.mx.Unlock()

	if !ok {
		return Free
	}

	if l {
		return Lock
	}

	return ReadLock
}
