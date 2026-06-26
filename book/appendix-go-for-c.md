# Appendix A: Go for C Programmers (a growing cheat sheet)

This appendix accumulates the Go concepts we meet, framed against C. It grows as
the project introduces new language features (goroutines, channels, `select`,
`sync.Mutex`, `context`, generics, ...).

## Types and values

| Concept | C | Go | Note |
|---|---|---|---|
| Type alias with teeth | `typedef int NodeID;` | `type NodeID int` | Go's is a *distinct* type; no implicit conversion to/from `int`. |
| Enum | `enum { A, B }` | `const ( A = iota; B )` | `iota` auto-increments within a `const` block. |
| Struct | `struct {...}` | `struct {...}` | Similar memory layout; Go adds methods and capitalization-based visibility. |
| Dynamic array | manual `malloc`/`realloc` | `[]T` (slice) | Fat pointer (ptr,len,cap), GC-managed, grows via `append`. |
| Fixed array | `T a[N]` | `[N]T` | Value type in Go (copied on assignment!); slices are the usual choice. |
| String | `char *` + len | `string` (immutable) / `[]byte` (mutable) | `string` is immutable bytes; convert with `[]byte(s)` / `string(b)`. |
| Null/sentinel | `#define NONE -1` | `const None NodeID = -1` | Typed constant. |
| Uninitialized memory | indeterminate | **zero value** | Go zero-initializes everything (0, "", nil, false). |

## Visibility and packages

- **Capitalized identifier** = exported (public, visible outside the package), like
  a non-`static` symbol declared in a header.
- **lowercase identifier** = unexported (package-private), like `static`.
- A directory is a package; files in it share scope. No headers, no `#include`; you
  `import "module/path/pkg"` and refer to `pkg.Name`.
- `internal/` packages are importable only within the same module subtree.

## Methods and interfaces

- A **method** is a function with a receiver: `func (r Role) String() string`. Think
  of it as a function whose first argument (`self`) is written before the name.
- An **interface** is a set of method signatures — like a struct of function
  pointers (a vtable) in C, but type-safe and first-class.
- **Implicit satisfaction:** a type implements an interface just by having the
  methods; there is no `implements` keyword. Checked at compile time, at the point
  where the concrete type is used as the interface.
- **Compile-time assertion idiom:** `var _ Iface = (*Concrete)(nil)` verifies
  satisfaction with zero runtime cost.
- An interface value is a `(type, value)` pair; the zero value is `nil`.

## Concurrency and time

| Concept | C analog | Go | Note |
|---|---|---|---|
| Instant | `struct timespec` | `time.Time` | A value type (copied on assignment). |
| Span | `long` nanoseconds | `time.Duration` | An `int64` of nanoseconds; `100 * time.Millisecond`. |
| Add/sub time | manual | `t.Add(d)`, `t2.Sub(t1)` | Instant ± span; instant − instant = span. |
| Thread | `pthread_create` | `go f()` | A goroutine: a cheap user-space thread (KBs, not MBs). |
| Mutex | `pthread_mutex_t` | `sync.Mutex` | Zero value is ready to use; no init call. |
| Cleanup on return | `goto cleanup` / `__attribute__((cleanup))` | `defer f()` | Runs at function return, LIFO order. |
| Thread-safe queue | hand-built (mutex+condvar+ring) | channel `chan T` | Built-in. Send `ch <- v`; receive `<-ch`. |
| Channel direction | — | `<-chan T` (recv-only), `chan<- T` (send-only) | Enforced by the type system. |
| Buffered channel | bounded queue | `make(chan T, n)` | `n=0` (unbuffered) = rendezvous; send blocks until receive. |
| Wait on many | `select()` / `poll()` | `select { case ... }` | Cases are channel ops; a `default` makes it non-blocking. |

Two rules of thumb that recur throughout the project:

- **Never send on a channel (or make a blocking call) while holding a mutex.** Snapshot
  what you need under the lock, release it, then do the blocking work.
- **The zero value should be useful.** A `sync.Mutex{}` is an unlocked mutex; a `nil`
  slice is an empty appendable slice. Design types so their zero value just works.

## Tooling

| Task | Command |
|---|---|
| Build everything | `go build ./...` |
| Static checks | `go vet ./...` |
| Canonical formatting (list offenders) | `gofmt -l .` |
| Run tests (verbose) | `go test -v ./...` |
| Run tests with race detector | `go test -race ./...` |
| Create a module | `go mod init <path>` |

## Testing

- Test files end in `_test.go`, live beside the code, same package.
- Test functions: `func TestXxx(t *testing.T)`.
- Idiomatic pattern is **table-driven**: a slice of cases looped with `t.Run(name, ...)`
  so a failure names the exact row.
- `t.Errorf` records a failure and continues; `t.Fatalf` records and stops the test.
