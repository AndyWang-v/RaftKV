# Chapter 27 — Interview Questions and Answers

This chapter is a bank of real Go interview questions with short, correct answers.
It also doubles as a fast review of the whole book.

> **How to use this.** Cover the answer, say your own answer out loud, then check.
> If an answer surprises you, jump to the chapter named in it and read that section.
> Aim to explain each idea in two or three sentences, the way you would on a call.
> Where it helps, the answer compares Go to C, because that is how you already think.

## Language basics

**Q: What is a "zero value", and why does Go have one?**
A: Every variable in Go is given a defined default value when you declare it without
an initializer. Numbers become `0`, booleans `false`, strings `""`, and pointers,
slices, maps, channels, functions, and interfaces become `nil`. In C an
uninitialized local holds garbage; in Go there is no garbage. Good Go types are
designed so the zero value is already useful (an empty `sync.Mutex` is unlocked; a
`nil` slice is an empty slice you can append to). See Chapter 4 — Types, Variables,
and Constants.

**Q: What is the difference between `:=` and `var`?**
A: `var x int = 5` declares a variable, and the type can be explicit or inferred.
`x := 5` is short variable declaration: it declares *and* initializes and infers the
type, and it only works inside a function. Use `var` for package-level variables, for
the zero value (`var buf []byte`), or when you want to name the type. Use `:=` for
ordinary local variables. They are not interchangeable: `:=` cannot appear at package
scope, and it must introduce at least one new variable on its left.

**Q: How do you make an enum in Go, and what is `iota`?**
A: Go has no `enum` keyword. You use a `const` block with `iota`, a counter that
starts at `0` and increases by one for each line in the block.

```go
type Color int

const (
	Red   Color = iota // 0
	Green              // 1
	Blue               // 2
)
```

This is like C's `enum { Red, Green, Blue }`, but `Color` is a distinct named type,
so the compiler stops you from mixing it with a plain `int` by accident.

**Q: Why does Go have no implicit numeric conversions?**
A: To stop silent bugs. In C, `int` and `double` mix freely, and a narrowing
conversion can lose data quietly. In Go you must convert on purpose:
`float64(i) / float64(n)`. Even `int` and `int32` are different types and will not
mix. This is more typing but removes a whole class of "where did that rounding come
from?" bugs.

**Q: What is the difference between an array and a slice?**
A: An array `[N]T` has a fixed length that is *part of its type*, and it is a value:
assigning or passing it copies all elements. A slice `[]T` is a small header that
points at a backing array and has a length and a capacity; passing a slice copies the
header, not the data. In practice you use slices almost everywhere and rarely write a
fixed array. See Chapter 8 — Arrays, Slices, and Strings.

**Q: Why does Go have only one loop keyword?**
A: Simplicity and one obvious way. `for` covers every case: a three-part C-style
loop, a condition-only loop (`for cond {}` replaces `while`), an infinite loop
(`for {}`), a `for range` over collections, and `for i := range n` to count from
`0` to `n-1`. There is no `while` and no `do`/`while`. See Chapter 5 — Control Flow.

**Q: What does `any` mean, and how is it different from `interface{}`?**
A: `any` is a built-in alias for `interface{}`, the empty interface, which every type
satisfies. They are identical; `any` is just the modern, readable spelling and is
preferred since Go 1.18. It is the closest thing Go has to C's `void *`, but it is
type-safe: you must use a type assertion or type switch to get the concrete value
back out.

## Slices and strings

**Q: What is a slice "header", exactly?**
A: A slice is a three-word struct: a pointer to the first element in a backing array,
an `int` length (`len`), and an `int` capacity (`cap`). This is a "fat pointer". When
you pass a slice to a function, you copy these three words, so the function sees the
same backing array but its own copy of the header.

```
slice []int  ->  [ ptr | len | cap ]
                    |
                    v
backing array     [ 10 | 20 | 30 | 40 ]
```

**Q: What is the difference between length and capacity?**
A: Length is how many elements the slice currently has and can index. Capacity is how
many elements fit in the backing array starting from the slice's first element, before
a new array must be allocated. `append` uses capacity to decide whether it can grow in
place or must allocate. You read them with `len(s)` and `cap(s)`.

**Q: What does `append` actually do?**
A: `append(s, x)` adds `x` after the last element. If `len < cap`, it writes into the
existing backing array and returns a slice with `len+1`. If there is no room, it
allocates a bigger backing array, copies the elements, writes `x`, and returns a slice
that points at the new array. Because the return value may point somewhere new, you
must always write `s = append(s, x)`.

**Q: Explain the slice aliasing bug.**
A: Two slices can share one backing array. Writing through one, or appending when
there is spare capacity, can change the other. A re-slice like `b := a[:2]` keeps `a`'s
capacity, so `append(b, x)` can overwrite `a`'s elements. To get an independent copy,
use `copy`, or cut capacity with a full slice expression `a[low:high:max]`.

**Q: What does the built-in `copy` do?**
A: `copy(dst, src)` copies `min(len(dst), len(src))` elements and returns that count.
It is the safe way to duplicate slice data so two slices no longer share memory. It is
Go's `memmove`: it handles overlapping regions correctly.

**Q: What is the difference between a `byte` and a `rune`?**
A: `byte` is an alias for `uint8`: one byte of raw data. `rune` is an alias for
`int32` and holds one Unicode code point (one "character"). A `string` is bytes
encoded as UTF-8, so one rune can be 1 to 4 bytes. Indexing `s[i]` gives a `byte`;
ranging `for _, r := range s` gives runes.

**Q: Are Go strings mutable? How are they stored?**
A: No, a `string` is immutable: a read-only slice of bytes, stored as UTF-8. You
cannot assign to `s[i]`. To edit text, convert to `[]byte` (or `[]rune`), change it,
and convert back; each conversion copies. `len(s)` returns the number of *bytes*, not
characters, so an emoji adds several to `len`. See Chapter 8 — Arrays, Slices, and
Strings.

## Maps

**Q: What is the "comma-ok" idiom for maps?**
A: Reading a missing key returns the value type's zero value, which you cannot tell
apart from a real zero. The two-result form tells you which it is:

```go
v, ok := m["key"] // ok is true only if the key is present
```

The same `, ok` form appears for channel receives and type assertions.

**Q: What happens when you read from or write to a `nil` map?**
A: Reading from a `nil` map is safe and returns the zero value. Writing to a `nil` map
*panics* at run time. The zero value of a map is `nil`, so you must create it with
`make(map[K]V)` or a map literal before you store anything.

**Q: Why is map iteration order random?**
A: Go deliberately randomizes the order each time you range over a map. This stops
programs from accidentally depending on an order that the implementation never
promised. If you need order, copy the keys into a slice and `sort` it.

**Q: Why can't I write `m[k].field = x` when the value is a struct?**
A: Map values are not addressable, so you cannot take their address or assign to a
field in place. You must replace the whole value:

```go
p := m[k]
p.field = x
m[k] = p
```

This is unlike C, where a struct stored in a hash table you wrote could be edited
through a pointer. With a map of *pointers* (`map[K]*T`) you can edit fields directly.

**Q: Are maps safe for concurrent use?**
A: No. Concurrent writes (or a write at the same time as any read) cause a fatal
run-time error, "concurrent map writes", which you cannot recover from. Protect a map
with a `sync.Mutex`/`sync.RWMutex`, or use `sync.Map` for some read-heavy cases. See
Chapter 9 — Maps.

## Structs and interfaces

**Q: When should a method use a value receiver versus a pointer receiver?**
A: Use a pointer receiver (`func (s *S) M()`) when the method must change the receiver,
or when the struct is large and you want to avoid copying, or when some methods already
need a pointer (keep them consistent). Use a value receiver for small, immutable types.
A value receiver gets a copy, so changes to it are lost. See Chapter 10 — Structs and
Methods.

**Q: What is a method set, and why does it matter for interfaces?**
A: The method set of type `T` includes only methods with a value receiver. The method
set of `*T` includes methods with value *and* pointer receivers. So if a method has a
pointer receiver, only `*T` satisfies an interface that lists it; a plain `T` value
does not. This is the usual reason for "cannot use x (type T) as type I" errors.

**Q: What does "implicit interface satisfaction" mean?**
A: A type satisfies an interface just by having the required methods. There is no
`implements` keyword and no declaration linking them, unlike C++ or Java. This is
"structural typing": if it has the methods, it fits. The check happens at compile time
where the concrete value is used as the interface.

**Q: How is an interface value represented in memory?**
A: As a two-word pair: a pointer to type information (which concrete type and its
method table) and a pointer to the value (or the value itself if it fits). This is like
a C struct that holds both a `vtable` pointer and a `void *data`, but the compiler
builds and type-checks it for you. The zero value of an interface is `nil`, meaning
both words are nil.

**Q: What is the "nil interface" trap?**
A: An interface is `nil` only when *both* its type word and value word are nil. If you
put a typed nil pointer into an interface, the interface holds (type=`*T`, value=nil),
which is *not* equal to `nil`. So a function that returns a typed nil error will fail
the `if err != nil` check. Return a literal `nil`, not a nil typed pointer.

**Q: What is the empty interface used for?**
A: `any` (the empty interface) can hold a value of any type, so it is used for
general containers, printing (`fmt.Println(args ...any)`), and decoding unknown JSON.
You retrieve the real value with a type assertion `v, ok := x.(int)` or a type switch.
Prefer concrete types or generics when you can; `any` loses compile-time type checks.

**Q: How is embedding different from inheritance?**
A: Embedding puts one type inside another without a field name, and the outer type
"promotes" the inner type's fields and methods so you can call them directly. It is
*composition*, not inheritance: there is no base class, no virtual override, and no
"is-a" relationship for type checking. You get code reuse; polymorphism comes from
interfaces instead. See Chapter 11 — Interfaces.

## Errors

**Q: How does Go handle errors, and why not exceptions?**
A: Errors are ordinary values of the built-in `error` interface, returned as the last
result and checked with `if err != nil`. There are no exceptions for normal failures.
The reason is clarity: the error path is visible in the code, not hidden in invisible
unwinding, so you cannot forget a failure case silently. See Chapter 12 — Errors.

**Q: What do `%w`, `errors.Is`, and `errors.As` do?**
A: `fmt.Errorf("...: %w", err)` wraps an error, keeping the original inside for later
inspection. `errors.Is(err, target)` reports whether `target` is anywhere in the wrap
chain (use it to compare against sentinel errors like `io.EOF`). `errors.As(err,
&target)` finds the first error in the chain of a given concrete type and assigns it,
so you can read its fields.

**Q: What is the difference between `panic` and returning an `error`?**
A: An `error` is for expected, recoverable problems (file not found, bad input); the
caller decides what to do. `panic` is for programmer mistakes or impossible states
(index out of range, nil dereference); it unwinds the stack and crashes unless a
deferred `recover` stops it. Do not use `panic` for ordinary control flow.

**Q: When is it acceptable to call `panic`?**
A: At program startup when a required condition is impossible to continue from (a
missing config that you cannot run without), in truly "this should never happen"
internal checks, and inside a package when you `recover` at the boundary and turn it
back into an error. Library code that callers depend on should return errors, not
panic.

**Q: What does `recover` do, and where must it be called?**
A: `recover` stops a panic and returns the panic value, but it only works when called
directly inside a deferred function. Outside `defer` it returns `nil` and does nothing.
The common use is at an API boundary or per-request handler to keep one bad request
from crashing the whole server.

**Q: How do you define a custom error type?**
A: Define a type with an `Error() string` method, which satisfies the `error`
interface.

```go
type NotFoundError struct{ Key string }

func (e *NotFoundError) Error() string {
	return "not found: " + e.Key
}
```

Callers can then use `errors.As` to detect it and read `Key`.

## Concurrency — goroutines and channels

**Q: How is a goroutine different from an OS thread?**
A: A goroutine is a function scheduled by the Go runtime, not the OS. It starts with a
tiny stack (about 2 KB) that grows and shrinks on demand, while an OS thread reserves
megabytes. You can run hundreds of thousands of goroutines; the runtime multiplexes
them onto a small number of OS threads. Starting one is `go f()`. See Chapter 13 —
Goroutines and the Scheduler.

**Q: What is the G-M-P scheduler model?**
A: `G` is a goroutine, `M` is a machine (an OS thread), and `P` is a processor (a
scheduling context that owns a queue of runnable goroutines). An `M` must hold a `P` to
run Gs. The number of `P`s is `GOMAXPROCS`, which caps how many goroutines run in
parallel. Idle `P`s use "work stealing" to take Gs from busy `P`s, which balances load.

**Q: What is `GOMAXPROCS` and what is its default?**
A: It is the maximum number of OS threads that may execute Go code at the same time,
which sets the real parallelism. By default it is the number of logical CPUs available.
Since Go 1.25 the runtime is container-aware: on Linux it also honors the cgroup CPU
limit and lowers `GOMAXPROCS` to match, and it adjusts if the limit changes.

**Q: What is the difference between a buffered and an unbuffered channel?**
A: An unbuffered channel (`make(chan T)`) has no storage: a send blocks until a
receiver is ready, so it is a rendezvous and also a synchronization point. A buffered
channel (`make(chan T, n)`) holds up to `n` values; a send blocks only when the buffer
is full, and a receive blocks only when it is empty. See Chapter 14 — Channels and
Select.

**Q: What does `select` do?**
A: `select` waits on several channel operations at once and runs the one that is ready
first; if several are ready it picks one at random. A `default` case makes it
non-blocking (do something else if no channel is ready). It is Go's version of
`select`/`poll` on file descriptors, but for channels.

**Q: What are the rules for closing a channel?**
A: Only the sender should close a channel, and only once. Sending on a closed channel
panics; closing an already-closed channel panics; closing a `nil` channel panics.
Receiving from a closed channel never blocks: it returns the zero value immediately,
and `v, ok := <-ch` sets `ok` to `false`. You do not have to close every channel; close
to signal "no more values" to receivers.

**Q: What are the common causes of a deadlock?**
A: All goroutines are blocked with no one able to make progress: sending on an
unbuffered channel with no receiver, receiving from a channel that no one will send on
or close, every goroutine waiting on a `WaitGroup` that never reaches zero, or two
goroutines each holding one lock and waiting for the other. If the *whole* program is
blocked, the runtime detects it and crashes with "all goroutines are asleep -
deadlock!".

## Concurrency — synchronization and patterns

**Q: When should I use a mutex versus a channel?**
A: Use a `sync.Mutex` to protect shared state that several goroutines read and write
(a cache, a counter). Use a channel to *transfer ownership* of data or to coordinate
steps between goroutines (a pipeline, a worker pool). The slogan is "share memory by
communicating", but a mutex is simpler and faster for plain shared state. See Chapter
15 — Sync and Context.

**Q: What is a `sync.WaitGroup` for?**
A: It waits for a set of goroutines to finish. Call `wg.Add(n)` before starting them,
`wg.Done()` (usually via `defer`) inside each one, and `wg.Wait()` to block until the
counter reaches zero. It is like `pthread_join` for a group, without holding each
thread handle.

```go
var wg sync.WaitGroup
for _, item := range items {
	wg.Add(1)
	go func() {
		defer wg.Done()
		process(item)
	}()
}
wg.Wait()
```

**Q: What is `context`, and how is it used for cancellation?**
A: A `context.Context` carries a deadline, a cancellation signal, and request-scoped
values across API boundaries. A parent creates one with `context.WithCancel` or
`context.WithTimeout` and passes it as the first argument. Worker code selects on
`ctx.Done()` and stops when it is closed. It is the standard way to cancel a tree of
goroutines (for example, when an HTTP client disconnects).

**Q: What is the race detector and how do you use it?**
A: It is a tool built into the toolchain that finds data races: unsynchronized access
to the same memory from two goroutines where at least one is a write. Run your tests or
program with the `-race` flag (`go test -race ./...`, `go run -race .`). It adds
overhead, so use it in testing and CI, not production. See Chapter 16 — Concurrency
Patterns.

**Q: What is a worker pool and why use one?**
A: It is a fixed number of goroutines that pull jobs from a shared channel. It limits
how much work runs at once, which bounds memory and CPU use, instead of spawning one
goroutine per job without limit. Producers send jobs and close the channel; workers
range over it and exit when it is closed.

**Q: What are fan-out and fan-in?**
A: Fan-out: start several goroutines that read from the same input channel, to spread
work across CPUs. Fan-in: merge several input channels into one output channel so a
single consumer can read all results. Together they form a pipeline stage that scales.

**Q: What is a goroutine leak, and how do you avoid it?**
A: A goroutine leak is a goroutine that blocks forever and is never cleaned up, so it
holds memory and resources for the life of the program. The usual cause is a send or
receive with no partner (for example, writing to a channel whose reader has gone away).
Avoid it by giving every goroutine a clear way to exit: close the channel, or pass a
`context` and select on `ctx.Done()`.

## Runtime and memory

**Q: What is the difference between the stack and the heap in Go?**
A: Each goroutine has its own stack for local variables and call frames; it is fast and
freed automatically when a function returns. The heap holds values that must outlive the
function that created them; it is managed by the garbage collector. Unlike C, *you* do
not choose: the compiler decides, and there is no `malloc`/`free`. See Chapter 17 —
Memory and the Garbage Collector.

**Q: What is escape analysis?**
A: It is the compiler step that decides whether a variable can live on the stack or
must "escape" to the heap. A value escapes if it is still reachable after the function
returns: it is returned by pointer, stored in a global, captured by a goroutine, or
placed in an interface that outlives the call. You can see the decisions with
`go build -gcflags=-m`.

**Q: Is it safe to return a pointer to a local variable in Go?**
A: Yes. This is undefined behavior in C, but in Go it is safe and common. Escape
analysis sees that the local outlives the function, so it allocates that variable on
the heap, and the garbage collector frees it once nothing points to it.

```go
func newCounter() *int {
	n := 0     // escapes to the heap
	return &n  // safe in Go; would be a dangling pointer in C
}
```

**Q: How does Go's garbage collector work?**
A: It is a concurrent, tri-color, mark-and-sweep collector. It runs alongside your
program and uses a write barrier so it can mark live objects while goroutines keep
running; the stop-the-world pauses are tiny (usually well under a millisecond). It does
not move objects (no compaction) and is not generational. Go 1.26 ships the new "Green
Tea" collector by default, which scans small objects with better cache locality and
CPU scaling.

**Q: What do `GOGC` and `GOMEMLIMIT` control?**
A: `GOGC` (default `100`) sets how much the heap may grow before the next collection:
`100` means collect when live data has doubled. Raising it trades memory for fewer GC
cycles; `GOGC=off` disables the GC. `GOMEMLIMIT` sets a *soft* memory limit (default
effectively off); the runtime collects more aggressively as memory approaches it, which
helps avoid out-of-memory kills in containers.

**Q: How do you reduce allocations in a hot path?**
A: Preallocate slices and maps with a capacity (`make([]T, 0, n)`); reuse buffers with
`sync.Pool`; pass and store values instead of pointers when they are small so they stay
on the stack; use `strings.Builder` instead of `+` in a loop; and avoid putting values
into `any`/interfaces needlessly. Measure first with `go test -bench` and `pprof`; do
not guess. See Chapter 17 — Memory and the Garbage Collector.

## Modules and tooling

**Q: What is the difference between a module and a package?**
A: A package is one directory of `.go` files compiled together and is the unit of
*imports*. A module is a collection of packages versioned and released together; it is
the unit of *dependencies* and is defined by a `go.mod` file at its root. Roughly: a
package is like a C library's translation units; a module is like the whole versioned
library you depend on. See Chapter 18 — Packages and Modules.

**Q: What are `go.mod` and `go.sum`?**
A: `go.mod` declares the module path, the Go version, and the required dependencies and
their versions. `go.sum` records cryptographic checksums of those dependencies so the
build can verify it downloaded exactly the expected code. Commit both. You rarely edit
them by hand; `go get` and `go mod tidy` maintain them.

**Q: What do `gofmt` and `go vet` do?**
A: `gofmt` rewrites source into the one canonical format (tabs, brace placement,
spacing), so style is never debated and diffs stay small. `go vet` runs static checks
for likely mistakes that compile but are probably wrong, such as bad `Printf` format
verbs, lost struct tags, or unreachable code. Run both in CI.

**Q: What is a table-driven test?**
A: A test where the cases are a slice of structs and a loop runs each one, usually with
`t.Run` so each row is a named subtest. It is the idiomatic Go pattern because adding a
case is one line and a failure names the exact row.

```go
func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"positive", 3, 3},
		{"negative", -3, 3},
		{"zero", 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := abs(tc.in); got != tc.want {
				t.Errorf("abs(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}
```

**Q: How do you write and read a benchmark?**
A: Write `func BenchmarkXxx(b *testing.B)` and loop `for b.Loop() { ... }` (or the
older `for range b.N`). Run it with `go test -bench=. -benchmem`. The output reports
nanoseconds per operation and, with `-benchmem`, bytes and allocations per operation,
which is what you optimize against. See Chapter 21 — Testing.

## Predict the output

These are classic puzzles. Read the code, decide the output, then check the answer
and the reason. Each program is complete and runnable.

**1. Defer order**

```go
package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		defer fmt.Println(i)
	}
}
```

Output:

```
2
1
0
```

Why: `defer` pushes each call onto a stack. The argument `i` is evaluated *now* (so
0, 1, 2 are captured), but the calls run in last-in-first-out order when `main`
returns.

**2. Append aliasing**

```go
package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := a[:2]
	b = append(b, 99)
	fmt.Println(a, b)
}
```

Output:

```
[1 2 99] [1 2 99]
```

Why: `b` has `len 2` but `cap 3`, sharing `a`'s backing array. `append` has spare
capacity, so it writes `99` into index 2 of that shared array, overwriting `a[2]`.
Both slices now show `99`.

**3. Goroutines and the loop variable**

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Print(i)
		}()
	}
	wg.Wait()
}
```

Output: some ordering of the digits `0`, `1`, and `2`, for example `012` or `201`.

Why: Since Go 1.22 each iteration has its own `i`, so each goroutine prints its own
value. Before Go 1.22 all goroutines shared one `i` and this famously printed `333`.
The order is not fixed because goroutines run concurrently.

**4. The nil interface that is not nil**

```go
package main

import "fmt"

type MyError struct{}

func (*MyError) Error() string { return "boom" }

func find() error {
	var e *MyError // a nil pointer
	return e       // returned as an interface
}

func main() {
	fmt.Println(find() == nil)
}
```

Output:

```
false
```

Why: `find` returns an `error` interface holding (type=`*MyError`, value=nil). An
interface equals `nil` only when both its type and value are nil; here the type is
set, so the result is not `nil`. Fix it by returning a literal `nil`.

**5. Map iteration order**

```go
package main

import "fmt"

func main() {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	for k := range m {
		fmt.Print(k)
	}
}
```

Output: a random ordering of `1`, `2`, `3`, such as `231`; it changes between runs.

Why: Go randomizes map iteration order on purpose so code cannot depend on it. To
print in order, copy the keys into a slice and `sort.Ints` it first.

**6. Named return changed by defer**

```go
package main

import "fmt"

func double() (result int) {
	defer func() {
		result *= 2
	}()
	return 5
}

func main() {
	fmt.Println(double())
}
```

Output:

```
10
```

Why: `return 5` first assigns `5` to the named return value `result`. Then the
deferred function runs *before* the function truly returns, doubling `result` to
`10`. A deferred closure can read and modify named return values.

## Design questions

These ask you to sketch a small system out loud. Give the data structures and the
flow; you do not need full code.

**Sketch a worker pool.**

```go
func workerPool(jobs <-chan int, results chan<- int, workers int) {
	var wg sync.WaitGroup
	for range workers { // start a fixed number of workers
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs { // ends when jobs is closed
				results <- j * j
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results) // close results once every worker is done
	}()
}
```

Key points to say: a fixed number of goroutines bounds concurrency; the producer
closes `jobs` to signal "no more work"; each worker `range`s over `jobs`; a
`WaitGroup` tracks the workers; a separate goroutine closes `results` after `Wait` so
the consumer's `range results` can end cleanly.

**Sketch a rate limiter.**

```go
// Allow at most one event every 200 milliseconds.
ticker := time.NewTicker(200 * time.Millisecond)
defer ticker.Stop()
for req := range requests {
	<-ticker.C // wait for the next tick (a "token")
	go handle(req)
}
```

Key points: a `time.Ticker` produces tokens at a steady rate, and you take one before
each request. To allow short bursts, use a buffered "token bucket" channel that a
ticker refills up to its capacity. In production, prefer the standard
`golang.org/x/time/rate.Limiter`, which implements a token bucket with burst support.

**Sketch an LRU cache.**

- Goal: O(1) `Get` and `Put`, and evict the least-recently-used entry when full.
- Structures: a `map[Key]*list.Element` for O(1) lookup, plus a doubly linked list
  (`container/list`) that orders entries from most-recently used at the front to
  least-recently used at the back.
- `Get`: find the element in the map, move it to the front, return its value.
- `Put`: if the key exists, update its value and move it to the front; otherwise
  insert at the front. If the size now exceeds the capacity, remove the back element
  and delete its key from the map.
- For concurrent use, guard both structures with one `sync.Mutex`.
