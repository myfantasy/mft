package mft

import (
	"sync"
	"time"
)

// G Generator
type G struct {
	rvc  int64
	rvgl int64
	rvmx sync.Mutex
	//rvmx RWCMutex

	AddValue int64 `json:"generator_add_value,omitempty"`
}

var GlobalGenerator *G = &G{}

// RvGet2 - Generate Next RV (sync)
func (g *G) RvGet() int64 {
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

// RvGetPart - Generate Next RV (sync) partitioned by (x%10000)/10
func (g *G) RvGetPart() int64 {
	t := time.Now().UnixNano() / 10000 * 10000

	g.rvmx.Lock()

	if g.rvgl == t {
		g.rvc = g.rvc + 1
		t = t + g.rvc
	} else {
		g.rvc = 0
		g.rvgl = t
	}

	k := g.rvc
	g.rvmx.Unlock()

	if k >= 10 {
		time.Sleep(time.Microsecond)
		return g.RvGetPart()
	}

	return t + g.AddValue
}

// RvGet2 - Generate Next RV (sync)
// global is used if g == nil
func (g *G) RvGetOrGlobal() int64 {
	if g != nil {
		return g.RvGet()
	}
	return RvGet()
}

// RvGetPartOrGlobal - Generate Next RV (sync) partitioned by (x%10000)/10
// global is used if g == nil
func (g *G) RvGetPartOrGlobal() int64 {
	if g != nil {
		return g.RvGetPart()
	}
	return RvGetPart()
}

// RvGet - Generate Next RV (sync)
func RvGet() int64 {
	return GlobalGenerator.RvGet()
}

// RvGetPart - Generate Next RV (sync)
func RvGetPart() int64 {
	return GlobalGenerator.RvGetPart()
}
