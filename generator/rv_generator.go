package generator

import (
	"sync"
	"syscall"
	"time"
)

// G Generator
type G struct {
	rvc  int64
	rvgl int64
	rvmx sync.Mutex

	AddValue int64
}

var g G

// RvGet2 - Generate Next RV (sync)
func (g *G) RvGet2() int64 {
	t := time.Now().UnixNano() / 10000 * 10000

	g.rvmx.Lock()

	if g.rvgl == t {
		g.rvc = g.rvc + 1
		t = t + g.rvc
	} else {
		g.rvc = 0
		g.rvgl = t
	}

	g.rvmx.Unlock()

	return t + g.AddValue
}

// RvGet - Generate Next RV (sync)
func (g *G) RvGet() int64 {
	t := GoTime() / 10000 * 10000

	g.rvmx.Lock()

	if g.rvgl == t {
		g.rvc = g.rvc + 1
		t = t + g.rvc
	} else {
		g.rvc = 0
		g.rvgl = t
	}

	g.rvmx.Unlock()

	return t + g.AddValue
}

// GoTime - fast get time
func GoTime() int64 {
	a := syscall.Timeval{}
	syscall.Gettimeofday(&a)
	return syscall.TimevalToNsec(a)
}

// RvGet2 - Generate Next RV (sync)
func RvGet2() int64 {
	return g.RvGet2()
}

// RvGet - Generate Next RV (sync)
func RvGet() int64 {
	return g.RvGet()
}
