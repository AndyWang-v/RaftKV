package clock

import (
	"sort"
	"sync"
	"time"
)

// simEpoch is the fixed instant a SimClock starts at, so virtual time is
// identical across runs — a prerequisite for reproducible tests.
var simEpoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// SimClock is a deterministic, manually-advanced clock for tests. Virtual time
// moves only when a test calls Advance; no real time elapses. Advancing by an
// hour returns instantly. It is safe for concurrent use by multiple goroutines.
type SimClock struct {
	mu     sync.Mutex
	now    time.Time
	timers []*simTimer
}

// NewSimClock returns a SimClock positioned at the fixed epoch.
func NewSimClock() *SimClock {
	return &SimClock{now: simEpoch}
}

func (c *SimClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *SimClock) After(d time.Duration) <-chan time.Time {
	return c.NewTimer(d).C()
}

func (c *SimClock) NewTimer(d time.Duration) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &simTimer{
		clock: c,
		when:  c.now.Add(d),
		ch:    make(chan time.Time, 1), // buffered so a fire never blocks Advance
	}
	c.timers = append(c.timers, t)
	return t
}

// Advance moves virtual time forward by d and fires every timer whose deadline
// has now passed, in deadline order.
func (c *SimClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	now := c.now

	// Partition the pending timers into "due now" and "still pending".
	var due, kept []*simTimer
	for _, t := range c.timers {
		if t.when.After(now) {
			kept = append(kept, t)
		} else {
			due = append(due, t)
		}
	}
	c.timers = kept
	c.mu.Unlock()

	// Fire OUTSIDE the lock — never send on a channel while holding a mutex.
	// Deadline order keeps causality sane when several timers come due at once.
	sort.Slice(due, func(i, j int) bool { return due[i].when.Before(due[j].when) })
	for _, t := range due {
		t.fire(now)
	}
}

// removeLocked drops t from the pending set and reports whether it was present.
// The caller MUST hold c.mu.
func (c *SimClock) removeLocked(t *simTimer) bool {
	for i, x := range c.timers {
		if x == t {
			c.timers = append(c.timers[:i], c.timers[i+1:]...)
			return true
		}
	}
	return false
}

// simTimer is a SimClock-backed Timer.
type simTimer struct {
	clock *SimClock
	when  time.Time
	ch    chan time.Time
}

func (t *simTimer) C() <-chan time.Time { return t.ch }

// fire delivers the time without blocking. The channel is buffered (cap 1); if a
// previous fire is still undrained we drop this one rather than stall Advance.
func (t *simTimer) fire(now time.Time) {
	select {
	case t.ch <- now:
	default:
	}
}

func (t *simTimer) Stop() bool {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	return t.clock.removeLocked(t)
}

func (t *simTimer) Reset(d time.Duration) bool {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	active := t.clock.removeLocked(t)
	t.when = t.clock.now.Add(d)
	t.clock.timers = append(t.clock.timers, t)
	return active
}

var (
	_ Clock = (*SimClock)(nil)
	_ Timer = (*simTimer)(nil)
)
