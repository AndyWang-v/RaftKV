# Appendix B — Glossary

Every important Go term in this book, defined in one to three plain sentences, with a
C comparison where it helps. Terms are grouped by first letter.

## B

**backing array** — The fixed-size array in memory that a slice points into. Several
slices can share one backing array, which is why writing through one slice can change
another. It is like the heap block a C dynamic array lives in, except the garbage
collector owns its lifetime.

**buffered channel** — A channel with room for `n` values, created with
`make(chan T, n)`. A send blocks only when the buffer is full; a receive blocks only
when it is empty. Think of it as a fixed-size, thread-safe queue.

**build cache** — A local cache of compiled packages and test results that makes
repeat builds fast. The `go` tool manages it; clear it with `go clean -cache`. A plain
C/Make workflow has no built-in equivalent.

## C

**capacity** — For a slice, the number of elements that fit in its backing array
starting from the slice's first element; read it with `cap(s)`. `append` uses it to
decide whether to grow in place or allocate a new array.

**channel** — A typed, thread-safe conduit that passes values between goroutines,
created with `make(chan T)`. Send with `ch <- v` and receive with `v := <-ch`. It
replaces the hand-built mutex-plus-condition-variable queue you would write in C.

**closure** — A function value that captures and keeps a reference to variables from
the surrounding scope; those variables live as long as the closure does. C has
function pointers but no closures, so you pass a separate `void *` context argument
instead.

**comparable** — A built-in constraint matching types whose values work with `==` and
`!=` (numbers, strings, booleans, pointers, channels, and structs or arrays of such).
It is used as a generic constraint, for example for map key types.

**constraint** — In generics, an interface that limits which types a type parameter
accepts and which operations are allowed, such as `[T comparable]` or
`[T constraints.Ordered]`. It is the type-safe replacement for C's untyped macros.

**context** — A `context.Context` value that carries a cancellation signal, a
deadline, and request-scoped values across function calls. By convention it is the
first argument, named `ctx`. It is the standard way to cancel a tree of goroutines.

## D

**deadlock** — A state where goroutines are all blocked waiting on one another and none
can proceed. If the entire program is stuck, the runtime detects it and aborts with a
"deadlock" message, which a C pthread deadlock would not do.

**defer** — A statement that schedules a function call to run when the surrounding
function returns, in last-in-first-out order. It is used for cleanup such as closing
files or unlocking mutexes, and it replaces C's `goto cleanup` pattern.

**dynamic dispatch** — Selecting which concrete method to call at run time based on the
value stored in an interface. The interface holds a pointer to a method table, much
like a C++ vtable or a C struct of function pointers.

## E

**embedding** — Placing a type inside a struct or interface without a field name, so
the outer type promotes the inner type's fields and methods. It gives reuse by
composition, not inheritance: there is no base class and no method overriding.

**escape analysis** — The compiler step that decides whether a variable can live on the
stack or must "escape" to the heap. A value escapes if it is still reachable after its
function returns. It is why returning a pointer to a local variable is safe in Go.

**exported / unexported** — An identifier starting with an uppercase letter is exported
(visible outside its package), like a non-`static` symbol in a C header. A lowercase
identifier is unexported (package-private), like `static`.

## F

**fat pointer** — An informal name for a multi-word value carrying a pointer plus extra
data. A slice is a fat pointer (pointer, length, capacity); an interface is a fat
pointer (type, value). A C pointer is a single word.

## G

**garbage collector (GC)** — The runtime component that automatically frees memory the
program can no longer reach, so you never call `free`. Go's is a concurrent, tri-color,
mark-and-sweep collector with very short pauses; Go 1.26 enables the "Green Tea" design
by default. It removes use-after-free and double-free bugs.

**generics** — Type parameters that let one function or type work for many types while
staying type-safe, added in Go 1.18. They replace C's `void *` containers and
preprocessor macros, for example `func Max[T constraints.Ordered](a, b T) T`.

**GMP** — The Go scheduler's model: G is a goroutine, M is a machine (an OS thread), and
P is a processor (a scheduling context with a run queue). An M needs a P to run Gs, and
the number of Ps is `GOMAXPROCS`.

**gofmt** — The standard tool that rewrites Go source into one official format (tabs,
brace placement, spacing). Because everyone runs it, formatting is never debated and
diffs stay small. C has no single official formatter.

**GOMAXPROCS** — The maximum number of OS threads that may run Go code at the same time,
which sets real parallelism. It defaults to the number of logical CPUs; since Go 1.25
the runtime also honors Linux cgroup CPU limits.

**GOPATH** — An older environment variable naming the workspace directory for Go code
and downloads. Since Go modules (Go 1.11+) it no longer controls project layout and
mostly holds the module cache and installed tools.

**GOROOT** — The directory where the Go toolchain and standard library are installed.
You rarely set it by hand; the installer and the `go` tool already know where it is.

**goroutine** — A function running concurrently, scheduled by the Go runtime rather than
the OS, started with `go f()`. It begins with a tiny stack (about 2 KB) that grows on
demand, so you can run hundreds of thousands of them. It is far cheaper than an OS
thread.

## H

**happens-before** — The ordering rule in the Go memory model that defines when one
goroutine's write is guaranteed visible to another goroutine's read. If two accesses
are not ordered this way and at least one writes, you have a data race. The C11 memory
model has the same concept.

## I

**interface** — A type that lists a set of method signatures; any type with those
methods satisfies it automatically, with no `implements` keyword. At run time an
interface value is a (type, value) pair. It is Go's type-safe version of a C struct of
function pointers.

**iota** — A counter usable inside a `const` block that starts at `0` and increases by
one per line, used to build enumerations. `const ( A = iota; B; C )` yields `0`, `1`,
`2`. It replaces a C `enum` or a run of `#define`s.

## M

**method set** — The set of methods callable on a type, which determines the interfaces
it satisfies. The method set of `T` has value-receiver methods only; the method set of
`*T` has both value- and pointer-receiver methods. This is why a pointer-receiver
method requires `*T` to satisfy an interface.

**module** — A collection of related packages versioned and released together, defined
by a `go.mod` file at its root. It is the unit of dependency management, roughly like a
whole versioned library you depend on.

**module cache** — The local directory where downloaded module versions are stored and
shared across projects, verified against `go.sum`. Clear it with `go clean -modcache`.
C has no built-in dependency cache.

**mutex** — A mutual-exclusion lock from the `sync` package that lets only one goroutine
touch shared data at a time. Its zero value is an unlocked, ready-to-use mutex, with no
init call unlike `pthread_mutex_init`. Lock with `mu.Lock()` and `defer mu.Unlock()`.

**MVS** — Minimal Version Selection, the algorithm Go uses to pick dependency versions.
It chooses the lowest version that satisfies all requirements, which keeps builds
reproducible. Many other package managers instead pick the newest allowed version.

## N

**nil** — The zero value of pointers, slices, maps, channels, functions, and
interfaces, meaning "no value" or "not set". It is like C's `NULL`, but typed and
broader: a `nil` slice is a valid empty slice you can append to, and a `nil` map is
safe to read.

## P

**package** — One directory of `.go` files compiled together, and the unit of imports
and visibility; every file begins with `package name`. It replaces C's
header-plus-translation-unit model, so there are no `#include`s.

**panic** — A run-time error that unwinds the stack and crashes the program unless a
deferred `recover` stops it. It is for programmer mistakes and impossible states (index
out of range, nil dereference), not for ordinary errors. It resembles an unhandled
exception or `abort`.

**pointer** — A value holding the address of another value, declared `*T`, taken with
`&x`, and read with `*p`. Go has pointers but **no pointer arithmetic**, and the GC
manages what they point to, so dangling pointers and manual frees do not happen.

## R

**race condition** — A bug where two goroutines access the same memory without
synchronization and at least one writes, giving undefined results. Find races by
building or testing with the `-race` flag. It is the same hazard as unsynchronized
pthreads in C.

**receiver** — The special parameter of a method, written before the method name, that
the method acts on, as in `func (r Rectangle) Area() float64`. A value receiver gets a
copy; a pointer receiver (`*Rectangle`) can modify the original. It is the explicit
`self` or `this`.

**recover** — A built-in that stops a panic and returns its value, but only when called
directly inside a deferred function. It turns a crash into a handled error at an API
boundary; outside `defer` it does nothing.

**rune** — An alias for `int32` that holds one Unicode code point (one character).
Ranging over a string yields runes, decoding the UTF-8 bytes. A C `char` is one byte,
which is not one character for non-ASCII text.

## S

**select** — A statement that waits on several channel operations and runs whichever is
ready first, choosing at random among ties. A `default` case makes it non-blocking. It
is like `select`/`poll` on file descriptors, but for channels.

**semantic import versioning** — Go's rule that a module's major version 2 or higher
must appear in its import path, such as `example.com/lib/v2`. This lets two major
versions of a library coexist in one build without conflict.

**slice** — Go's dynamic array: a header (pointer, length, capacity) that views a
backing array, written `[]T`. It grows with `append` and is passed by copying its small
header. It is the everyday replacement for C's manual `malloc`/`realloc` arrays.

**slice header** — The three-word value that represents a slice: a pointer to the
backing array, an `int` length, and an `int` capacity. Passing a slice copies this
header, not the elements it points to.

**stop-the-world** — A brief phase where the garbage collector pauses every goroutine to
do work that needs a consistent view of memory. Modern Go keeps these pauses well under
a millisecond by doing most GC work concurrently.

**struct tag** — A string literal attached to a struct field, written in backticks after
the field, that carries metadata libraries read by reflection, such as a JSON field name
like `json:"name"`. C has no equivalent; you would hand-write serialization code.

## T

**type assertion** — Extracting the concrete value from an interface, as in
`s, ok := x.(string)`, which succeeds if `x` holds a `string`. The `, ok` form avoids a
panic when the type does not match. It is like a checked downcast.

**type parameter** — A placeholder type in a generic function or type, written in square
brackets with a constraint, as in `func Map[T, U any](...)`. It is filled in with a real
type at the call site.

**type switch** — A `switch` that branches on the concrete type inside an interface:
`switch v := x.(type) { case int: ...; case string: ... }`. It is the idiomatic way to
handle a value that could be one of several types.

## V

**vendoring** — Copying a project's dependencies into a local `vendor/` directory so
builds use those copies instead of the module cache or the network. It helps with
reproducible or offline builds, like checking third-party source into your repository.

**vet** — `go vet`, a static analysis tool that reports likely mistakes that still
compile, such as wrong `Printf` verbs or misused struct tags. Run it in CI alongside the
tests.

## W

**WaitGroup** — A `sync.WaitGroup` counter that waits for a group of goroutines to
finish: `Add` before starting them, `Done` (via `defer`) inside each, and `Wait` to
block until the count reaches zero. It is like joining a group of threads.

**work stealing** — A scheduler technique where an idle processor (P) takes runnable
goroutines from another P's queue, balancing load across CPUs without central
coordination.

## Z

**zero value** — The defined default a variable receives when declared without an
initializer: `0` for numbers, `false` for booleans, `""` for strings, and `nil` for
pointers, slices, maps, channels, and interfaces. There is no uninitialized garbage as
in C, and good Go types make the zero value useful.
