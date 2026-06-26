// Package clock provides an injectable time abstraction. Raft is full of
// timeouts (election timers, heartbeats); routing every time operation through
// this interface lets production use real wall-clock time while tests use a
// deterministic virtual clock that advances only when the test says so — no real
// sleeps, no flakiness, thousands of fault-injection seeds in seconds.
package clock

import "time"

// Clock is the seam between Raft and time. Production wires in RealClock; tests
// wire in *SimClock.
type Clock interface {
	// Now returns the current time (real or virtual).
	Now() time.Time
	// NewTimer returns a one-shot Timer that fires once after d has elapsed.
	NewTimer(d time.Duration) Timer
	// After is shorthand for NewTimer(d).C(): a receive-only channel that gets
	// one value after d elapses. Receive-only so callers can wait but not send.
	After(d time.Duration) <-chan time.Time
}

// Timer mirrors the essential surface of the standard library's time.Timer, so
// the simulated clock can supply a drop-in fake.
type Timer interface {
	// C is the channel on which the time is delivered when the timer fires.
	C() <-chan time.Time
	// Stop prevents the timer from firing. It reports whether the timer was still
	// active (true) or had already fired/been stopped (false).
	Stop() bool
	// Reset restarts the timer to fire after d from now. It reports whether the
	// timer was active before the reset.
	Reset(d time.Duration) bool
}
