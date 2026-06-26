package clock

import (
	"testing"
	"time"
)

// --- SimClock: the deterministic clock our Raft tests will rely on ---

func TestSimClockNowAdvance(t *testing.T) {
	c := NewSimClock()
	start := c.Now()
	c.Advance(5 * time.Second)
	if got := c.Now().Sub(start); got != 5*time.Second {
		t.Fatalf("after Advance(5s), elapsed = %v, want 5s", got)
	}
}

func TestSimClockAfterFiresOnlyAfterDeadline(t *testing.T) {
	c := NewSimClock()
	ch := c.After(100 * time.Millisecond)

	// One nanosecond short of the deadline: must NOT have fired.
	c.Advance(100*time.Millisecond - 1)
	select {
	case <-ch:
		t.Fatal("timer fired before its deadline")
	default: // good: nothing ready
	}

	// Crossing the deadline fires it, carrying the (virtual) current time.
	c.Advance(1)
	select {
	case got := <-ch:
		if want := c.Now(); !got.Equal(want) {
			t.Fatalf("fired at %v, want %v", got, want)
		}
	default:
		t.Fatal("timer did not fire at its deadline")
	}
}

func TestSimClockMultipleTimersAllFire(t *testing.T) {
	c := NewSimClock()
	chans := []<-chan time.Time{
		c.After(30 * time.Millisecond),
		c.After(10 * time.Millisecond),
		c.After(20 * time.Millisecond),
	}
	c.Advance(time.Second) // all three become due at once
	for i, ch := range chans {
		select {
		case <-ch:
		default:
			t.Fatalf("timer %d did not fire", i)
		}
	}
}

func TestSimClockStopPreventsFire(t *testing.T) {
	c := NewSimClock()
	tm := c.NewTimer(50 * time.Millisecond)
	if active := tm.Stop(); !active {
		t.Fatal("Stop on a pending timer should report it was active")
	}
	c.Advance(time.Second)
	select {
	case <-tm.C():
		t.Fatal("a stopped timer must not fire")
	default:
	}
}

func TestSimClockReset(t *testing.T) {
	c := NewSimClock()
	tm := c.NewTimer(50 * time.Millisecond) // due at t=50ms

	c.Advance(40 * time.Millisecond) // now=40ms, not yet due
	tm.Reset(50 * time.Millisecond)  // new deadline = 40ms + 50ms = 90ms

	c.Advance(40 * time.Millisecond) // now=80ms, still before 90ms
	select {
	case <-tm.C():
		t.Fatal("reset timer fired before its new deadline")
	default:
	}

	c.Advance(20 * time.Millisecond) // now=100ms, past 90ms
	select {
	case <-tm.C():
	default:
		t.Fatal("reset timer should have fired by its new deadline")
	}
}

// --- RealClock: a thin wrapper over the time package (small real waits) ---

func TestRealClockNow(t *testing.T) {
	c := RealClock{}
	before := time.Now()
	got := c.Now()
	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v, want within [%v, %v]", got, before, after)
	}
}

func TestRealClockAfterFires(t *testing.T) {
	c := RealClock{}
	select {
	case <-c.After(5 * time.Millisecond):
		// fired as expected
	case <-time.After(2 * time.Second):
		t.Fatal("RealClock.After(5ms) did not fire within 2s")
	}
}
