# Stage 0: Scaffolding

> **Learning objectives.** Set up a Go module and project layout; meet the Go
> language constructs a C programmer needs first (defined types, `iota`, slices,
> structs, exported vs. unexported, zero values); write the core data types and
> interfaces; and learn Go's built-in testing workflow.

## 3.1 The goal of Stage 0

Stage 0 builds everything the Raft logic will *hang on* and be *tested with*: the
data types, the interfaces (the seams), a simulated network we can torture, a
virtual clock, and a test harness. **There is no Raft algorithm in Stage 0.** The
point is that if we cannot reproducibly crash, partition, and delay a cluster, we
cannot trust any green checkmark later. Stage 0 is what makes Stages 1–5 debuggable.

This chapter covers the first half of Stage 0 — module setup, the core types, and
the interfaces. The virtual clock, the simulated network, and the harness follow.

## 3.2 Module and layout

```sh
git init
go mod init raftkv          # module path = raftkv; imports look like raftkv/internal/raft
mkdir -p cmd/kvnode cmd/kvctl internal/raft internal/kv \
         internal/transport internal/storage internal/clock web test
```

The resulting `go.mod` simply records the module path and Go version:

```
module raftkv

go 1.26.4
```

The layout mirrors the architecture: `internal/raft` is the consensus core;
`internal/kv` the state machine; `internal/transport`, `internal/storage`, and
`internal/clock` hold the implementations behind the interfaces; `cmd/` holds the
runnable binaries; `web/` the dashboard; `test/` the cluster harness. Go's `internal`
directory is special: packages under it can only be imported by code rooted at the
same parent, which keeps the API surface honest.

## 3.3 Go for C programmers: the constructs in the first file

Everything in the data-types file has a C analog, but each has a twist worth knowing.

| Go | C analog | The important difference |
|---|---|---|
| `package raft` | translation unit / namespace | Files in a directory share a package and scope — no headers, no `#include`. |
| `type NodeID int` | `typedef int NodeID;` | **Stronger than typedef:** a *distinct* type. You cannot pass a bare `int` where `NodeID` is wanted without an explicit `NodeID(x)` conversion. The compiler catches "I mixed up an id and a count." |
| `const ( Follower Role = iota; ... )` | `enum { ... }` | `iota` auto-increments: 0, 1, 2, ... |
| `struct { ... }` | `struct { ... }` | Nearly identical layout, but methods attach to it and visibility is by capitalization. |
| Capitalized name (`Term`) | non-`static`, in a header | **Exported** — visible outside the package (public). |
| lowercase name (`mu`) | `static` / file-private | **Unexported** — package-private. |
| `[]byte` | `struct { byte *p; size_t len; size_t cap; }` | A **slice**: a fat pointer (data, length, capacity), garbage-collected. *Not* a raw array. This is the biggest mental shift from C. |
| zero values | `calloc` / static init | Go zero-initializes everything. A fresh `Role` is `0` = `Follower`, and that is intentional. |

Two things that bite C programmers specifically:

1. **`[]byte` is a slice, not `char[]`.** It carries its own length, grows with
   `append`, and the garbage collector frees it. We use `[]byte` for commands and
   snapshots because Raft treats them as **opaque blobs** — it never looks inside;
   only the state machine interprets them. That opacity is the clean seam between
   consensus and application.

2. **`NodeID(-1)` as a sentinel.** Where C would `#define NONE -1`, Go uses a typed
   constant `const None NodeID = -1`, so it compares cleanly against other `NodeID`
   values.

## 3.4 The core data types (`internal/raft/types.go`)

This file is pure data — no behavior, no goroutines, no locking. The RPC structs
*are* Figure 2 of the paper.

```go
// Package raft implements the Raft consensus algorithm from scratch.
package raft

// NodeID identifies one server in the cluster. A distinct type (not a bare int)
// so the compiler stops us mixing a node id with a log index or a count.
type NodeID int

// None is the sentinel meaning "no node" (no vote yet; unknown leader).
const None NodeID = -1

// Role is the server's current role. Zero value is Follower (servers boot as one).
type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// String makes Role satisfy fmt.Stringer, so logs print "Follower" not "0".
func (r Role) String() string {
	switch r {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// LogEntry is one command in the replicated log. Term keys all consistency checks
// (Log Matching). Command is opaque to Raft; only the StateMachine decodes it.
type LogEntry struct {
	Term    int
	Command []byte
}

// --- RPC payloads: this section IS Figure 2 ---

type RequestVoteArgs struct {
	Term         int
	CandidateID  NodeID
	LastLogIndex int
	LastLogTerm  int
}
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     NodeID
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}
type AppendEntriesReply struct {
	Term    int
	Success bool
	// Fast-backup hints (not in Figure 2; a standard optimization). Set only on
	// failure, so the leader can skip many indices at once.
	ConflictTerm  int
	ConflictIndex int
}

type InstallSnapshotArgs struct {
	Term              int
	LeaderID          NodeID
	LastIncludedIndex int
	LastIncludedTerm  int
	Data              []byte
}
type InstallSnapshotReply struct{ Term int }

// ApplyMsg carries committed work from Raft up to the application, in log order,
// exactly once per entry. Either CommandValid or SnapshotValid is true, never both.
type ApplyMsg struct {
	CommandValid bool
	Command      []byte
	CommandIndex int

	SnapshotValid bool
	Snapshot      []byte
	SnapshotIndex int
}

// RaftState is the durable state that must survive a crash (persisted together
// with the snapshot so the two can never disagree).
type RaftState struct {
	CurrentTerm       int
	VotedFor          NodeID
	Log               []LogEntry
	LastIncludedIndex int
	LastIncludedTerm  int
}
```

## 3.5 Go's testing workflow

C has no standard test framework; Go bakes one in:

- Test files end in `_test.go` and sit next to the code they test (same package).
- Test functions are `func TestXxx(t *testing.T)`; `go test` finds them by signature.
- The idiomatic pattern is **table-driven**: a slice of `{input, want}` cases, looped.

A minimal example for `Role.String()` — small on purpose, so the *pattern* is the
lesson:

```go
package raft

import "testing"

func TestRoleString(t *testing.T) {
	cases := []struct {
		name string
		role Role
		want string
	}{
		{"follower", Follower, "Follower"},
		{"candidate", Candidate, "Candidate"},
		{"leader", Leader, "Leader"},
		{"unknown", Role(99), "Unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.role.String(); got != tc.want {
				t.Errorf("Role(%d).String() = %q, want %q", tc.role, got, tc.want)
			}
		})
	}
}
```

Run it:

```sh
go test -v ./internal/raft/     # verbose: shows each subtest
go test -race ./internal/raft/  # with the race detector — our gold standard
```

The commands we will use constantly: `go build ./...`, `go vet ./...` (static
checks), `gofmt -l .` (lists files not in canonical format — Go has exactly one
official style), and `go test -race ./...`.

## 3.6 The interfaces — the seams (`internal/raft/interfaces.go`)

This is the single most important design idea in the project. Raft talks to the
outside world through four plugs: a `Transport` (reach peers), a `Persister` (save
to disk), a `Clock` (know the time), and a `StateMachine` (the app it replicates).
**All are interfaces**, so tests inject fakes and the demo injects real
implementations — with no change to Raft.

```
   Raft depends only      Tests inject:        Demo injects:
   on the INTERFACES ->   simnet (fake net)    gRPC (real net)
   never the concrete     sim clock (virtual)  real clock
   types                  in-mem persister     file persister
```

### Go interfaces vs. C

In C you fake polymorphism with a struct of function pointers plus a `void *self`
(think `struct file_operations`). A Go interface is that pattern, type-safe and
first-class, with two special properties:

1. **Implicit satisfaction.** There is no `implements` keyword. A type satisfies an
   interface simply by *having the right methods*. The check happens **at compile
   time, at the point of use** — when you assign or pass the concrete type where the
   interface is expected. If a method is missing or mis-typed, the compiler stops you
   right there.

2. **An interface value is a fat pointer** — a `(type, data)` pair, essentially
   `(vtable, self)`.

Because satisfaction is implicit, idiomatic Go **defines an interface in the package
that consumes it** (here, `raft`), not next to the implementations. And we document
intent with a free, zero-cost compile-time assertion:

```go
var _ Transport = (*SimNet)(nil)  // fails to compile if *SimNet isn't a Transport
```

That declares a throwaway variable (`_` discards it) of the interface type and
assigns a nil concrete pointer; it emits no runtime code but forces the compiler to
verify satisfaction exactly where the type is defined.

### The file

```go
package raft

import "context"

// Transport delivers Raft RPCs to peers.
//
// Contract (non-negotiable):
//   - All calls are synchronous request/response.
//   - The CALLER owns timeouts (via ctx) and retries.
//   - The caller MUST NOT hold the Raft mutex while calling. Doing so deadlocks
//     the cluster. Snapshot under the lock, release it, call, re-acquire, re-validate.
type Transport interface {
	RequestVote(ctx context.Context, peer NodeID, args *RequestVoteArgs) (*RequestVoteReply, error)
	AppendEntries(ctx context.Context, peer NodeID, args *AppendEntriesArgs) (*AppendEntriesReply, error)
	InstallSnapshot(ctx context.Context, peer NodeID, args *InstallSnapshotArgs) (*InstallSnapshotReply, error)
}

// Persister stores durable Raft state.
//
// Contract:
//   - Save MUST be durable (fsync) before the triggering RPC handler replies.
//   - State and snapshot are saved together so a crash can't leave them disagreeing.
type Persister interface {
	Save(state RaftState, snapshot []byte) error
	Load() (RaftState, []byte, error)
	RaftStateSize() int
}

// StateMachine is the application Raft replicates (our KV map).
//
// Contract:
//   - Apply runs strictly in log order, exactly once per committed entry.
//   - Apply MUST be deterministic.
type StateMachine interface {
	Apply(cmd []byte) (result []byte)
	Snapshot() []byte
	Restore(snapshot []byte)
}
```

(The fourth interface, `Clock`, lives in its own `internal/clock` package because it
has no dependency on Raft's types; it is covered in the next section of this chapter.)

> **Trade-off: where do interfaces live?** The design doc sketches a `transport`
> package "owning" the `Transport` interface. But the RPC payload types
> (`RequestVoteArgs`, etc.) belong to Raft, and the interface references them, so
> defining the interface in `transport` and the types in `raft` would create an
> import cycle (Go forbids these). The idiomatic resolution — define interfaces in
> the consumer — also happens to avoid the cycle. We follow it: interfaces that
> reference Raft types live in `raft`; the `Clock`, which references none, stands
> alone.

## 3.7 The virtual clock (`internal/clock`)

Raft is full of timeouts: election timers (300–600 ms), heartbeats (100 ms). If
tests used real `time.Sleep`, a single election test would take half a second, and
a thousand-seed chaos suite would take *hours* — and still be flaky, because real
timing is nondeterministic. The fix is to make **time itself injectable**: a fake
clock the test advances by hand, so ten simulated minutes pass in microseconds,
deterministically.

### The concepts this introduces

- **`time.Time`** is an instant (a point on the timeline); **`time.Duration`** is a
  span stored as an `int64` of nanoseconds. Add a span to an instant with
  `t.Add(d)`; subtract two instants with `t2.Sub(t1)`.
- **Channels** are typed, thread-safe pipes — the closest C analog is a hand-built
  thread-safe queue (mutex + condvar + ring buffer), here a language primitive.
  Send with `ch <- v`, receive with `<-ch`. A type like `<-chan time.Time` is
  **receive-only**, so a caller can wait on it but not send into it. `make(chan T)`
  is unbuffered (send blocks until a receive); `make(chan T, 1)` buffers one value.
- **`select`** is a `switch` whose cases are channel operations; a `default` case
  makes it **non-blocking** ("if nothing is ready right now, do default").
- **`sync.Mutex`** is mutual exclusion, like `pthread_mutex_t`, but its zero value
  is ready to use (no init call).
- **`defer`** schedules a call for when the enclosing function returns (LIFO). The
  idiom `mu.Lock(); defer mu.Unlock()` guarantees unlock on every return path — the
  clean version of C's `goto cleanup`.

### The interfaces

```go
type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
	After(d time.Duration) <-chan time.Time
}

type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}
```

### RealClock — boring on purpose

```go
type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (RealClock) NewTimer(d time.Duration) Timer         { return &realTimer{t: time.NewTimer(d)} }

type realTimer struct{ t *time.Timer }

func (r *realTimer) C() <-chan time.Time        { return r.t.C }
func (r *realTimer) Stop() bool                 { return r.t.Stop() }
func (r *realTimer) Reset(d time.Duration) bool { return r.t.Reset(d) }
```

(`realTimer` exists because the standard `time.Timer` exposes its channel as a
*field* `.C`, but an interface can only require *methods* — so we wrap it in a
`C()` method.)

### SimClock — the interesting one

`SimClock` holds a virtual `now` and a list of pending timers under a mutex.
`Advance(d)` pushes `now` forward and fires every timer that has come due.

```go
func (c *SimClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	now := c.now

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
	sort.Slice(due, func(i, j int) bool { return due[i].when.Before(due[j].when) })
	for _, t := range due {
		t.fire(now)
	}
}
```

The timer's `fire` uses a **non-blocking send** so `Advance` can never stall:

```go
func (t *simTimer) fire(now time.Time) {
	select {
	case t.ch <- now: // channel is buffered (cap 1)
	default:          // already has an undrained tick; drop rather than block
	}
}
```

> **Discipline learned here, used everywhere later.** `Advance` collects the due
> timers *under* the lock, releases the lock, and only *then* sends on channels.
> Never send on a channel (or make a blocking call) while holding a mutex. This is
> the exact rule the Raft layer must obey (doc §6.2) — we rehearse it on the clock
> first, where the stakes are low.

### What the tests prove (and how fast)

```
TestSimClockNowAdvance ......... 0.00s   advanced 5 virtual seconds, instantly
TestSimClockReset .............. 0.00s   advanced 100 virtual ms, instantly
TestRealClockAfterFires ........ 0.01s   the only test that waits real time
```

Every `SimClock` test moves seconds or minutes of virtual time yet finishes in
0.00 s. That ratio is the entire justification for the abstraction. A representative
test, showing the non-blocking "did it fire yet?" check:

```go
func TestSimClockAfterFiresOnlyAfterDeadline(t *testing.T) {
	c := NewSimClock()
	ch := c.After(100 * time.Millisecond)

	c.Advance(100*time.Millisecond - 1) // one ns short
	select {
	case <-ch:
		t.Fatal("timer fired before its deadline")
	default: // good: nothing ready
	}

	c.Advance(1) // crossing the deadline
	select {
	case got := <-ch:
		if want := c.Now(); !got.Equal(want) {
			t.Fatalf("fired at %v, want %v", got, want)
		}
	default:
		t.Fatal("timer did not fire at its deadline")
	}
}
```

## 3.8 Where Stage 0 continues

Still to build in Stage 0:

- **The simulated network** — an in-memory `Transport` that can delay, drop,
  duplicate, reorder, and partition messages on command.
- **The cluster harness** — start/stop an N-node cluster in virtual time, inject
  faults, and continuously assert Raft's safety invariants.
