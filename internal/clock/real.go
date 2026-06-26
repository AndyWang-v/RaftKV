package clock

import "time"

// RealClock is the production Clock: a thin pass-through to the standard library
// time package. It carries no state, so the empty struct value RealClock{} is a
// complete, usable clock.
type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (RealClock) NewTimer(d time.Duration) Timer         { return &realTimer{t: time.NewTimer(d)} }

// realTimer wraps *time.Timer to satisfy our Timer interface. We need a method
// C() because the standard timer exposes its channel as a struct field (.C),
// and interfaces can only require methods.
type realTimer struct{ t *time.Timer }

func (r *realTimer) C() <-chan time.Time        { return r.t.C }
func (r *realTimer) Stop() bool                 { return r.t.Stop() }
func (r *realTimer) Reset(d time.Duration) bool { return r.t.Reset(d) }

// Compile-time assertions: these fail to build if the types ever stop satisfying
// the interfaces. Zero runtime cost; they document and enforce intent.
var (
	_ Clock = RealClock{}
	_ Timer = (*realTimer)(nil)
)
