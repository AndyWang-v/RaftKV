# Chapter 26 — Gotchas for C Programmers (the checklist)

> **What you'll learn.** The traps that bite C programmers most often when they
> write Go, gathered into one checklist. For each trap you get: what it is, why it
> surprises you coming from C, and the fix with a tiny snippet — plus a pointer to
> the chapter that teaches it in full.

This is the chapter to keep open while you write your first Go. Most of these are
not bugs in Go; they are places where Go's rules differ from C's, so your C
instinct points the wrong way. Read it once now, then come back to it whenever
something behaves in a way you did not expect. Each item is short on purpose.

## Arrays, slices, and strings

These are the single biggest source of surprises. The full story is in Chapter 8 —
Arrays, Slices, and Strings.

### Arrays are value types

In C, an array name *decays* to a pointer to its first element, so passing an array
to a function passes a pointer — the function can change the original. In Go, an
array is a **value**. Assigning it, or passing it to a function, **copies the whole
array**.

```go
a := [3]int{1, 2, 3}
b := a   // copies all three elements
b[0] = 99
// a is still {1, 2, 3}; b is {99, 2, 3}
```

> **Watch out.** A `[1000000]int` passed by value copies a megabyte every call. In
> Go you almost never use arrays directly. **Use a slice (`[]int`)**, which is a
> small header (pointer, length, capacity) and is cheap to pass.

### Slice aliasing: sub-slices share the backing array

A slice is a view into a backing array. Two slices can point at the **same** array.
Slicing (`s[1:3]`) does not copy; it makes another view. So writing through one
slice can change another.

```go
s := []int{1, 2, 3, 4}
t := s[:2]            // t shares s's backing array
t = append(t, 99)    // there is spare capacity, so this writes into s[2]
// s is now {1, 2, 99, 4} — t's append overwrote a neighbor
```

```
s := []int{1,2,3,4}     backing array:  [ 1 | 2 | 3 | 4 ]   len=4 cap=4
t := s[:2]              t ─▶ same array  [ 1 | 2 | 3 | 4 ]   len=2 cap=4
t = append(t, 99)       fits in cap, so it writes index 2:
                                         [ 1 | 2 | 99| 4 ]
                        s[2] changed too — that is aliasing.
```

> **Watch out.** `append` either writes into spare capacity (mutating shared data)
> **or** allocates a new array and copies (leaving the old slice behind). You cannot
> tell which from the call. Two rules keep you safe: **always assign the result**
> (`s = append(s, x)`), and if you need an independent copy, make one explicitly.

```go
indep := make([]int, len(s))
copy(indep, s)        // indep has its own backing array
```

To force `append` to allocate (so it never touches a shared array), cut the
capacity with a *three-index slice*: `t := s[:2:2]`.

### A small slice can pin a big array (a leak)

Because a slice keeps its whole backing array alive, a tiny slice of a huge buffer
keeps the **entire** buffer in memory. The garbage collector cannot free it.

```go
func firstLine(huge []byte) []byte {
	return huge[:80] // keeps all of `huge` alive, maybe megabytes
}
```

> **Watch out.** This is a real memory leak even with a garbage collector (Chapter
> 17 — Memory and the Garbage Collector). The fix is to **copy** the part you need
> so the big array can be freed:

```go
func firstLine(huge []byte) []byte {
	out := make([]byte, 80)
	copy(out, huge[:80]) // out does not reference `huge`
	return out
}
```

### A string index is a byte, not a character

A Go `string` is an immutable, read-only sequence of **bytes**, encoded in UTF-8,
with **no NUL terminator**. Indexing gives one *byte*, and `len` counts *bytes*,
not characters.

```go
s := "héllo"
fmt.Println(len(s))      // 6, not 5 (é is two bytes in UTF-8)
fmt.Println(s[0])        // 104, the byte 'h' (a byte, not a rune)
for i, r := range s {    // range decodes UTF-8: r is a rune (code point)
	_ = i
	_ = r
}
runes := []rune(s)       // convert to runes to index by character
fmt.Println(len(runes))  // 5
```

> **C vs Go.** In C a string is `char*`, NUL-terminated, and `strlen` walks until
> the `\0`. A Go string carries its length, has no terminator, and is bytes of
> UTF-8. To edit a string, convert to `[]byte` or `[]rune`, change that, and convert
> back — the string itself cannot be modified in place.

```c
char s[] = "hi";
s[0] = 'H';   /* legal in C: strings are mutable arrays */
```

```go
s := "hi"
// s[0] = 'H'           // does NOT compile: strings are immutable
b := []byte(s)
b[0] = 'H'
s = string(b)           // "Hi"
```

## Maps

See Chapter 9 — Maps.

### Writing to a nil map panics

A map's zero value is `nil`. You can **read** a nil map (every key returns the zero
value), but **writing** to it panics. You must create it with `make` (or a literal)
first.

```go
var m map[string]int
_ = m["x"]    // fine: returns 0
// m["x"] = 1 // PANIC: assignment to entry in nil map

m = make(map[string]int) // now writable
m["x"] = 1
```

> **Watch out.** This bites when a struct has a map field you forgot to initialize.
> The zero struct compiles and reads fine, then panics the first time you assign to
> the map. Initialize map fields in your `NewX` constructor.

### Map iteration order is random

Ranging over a map visits keys in an **unspecified, deliberately randomized** order
that changes from run to run. C has no built-in map, but programmers expect ordered
or stable containers. Go guarantees the opposite.

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
slices.Sort(keys) // sort if you need a stable order
```

### Map values are not addressable

You cannot take the address of a map element, so you cannot assign to a field of a
struct stored *by value* in a map. C lets you write `table[i].x = 1` freely.

```go
type P struct{ X int }
m := map[string]P{"a": {}}
// m["a"].X = 1        // does NOT compile: m["a"] is not addressable
v := m["a"]            // copy out
v.X = 1
m["a"] = v             // write back
// or store pointers: map[string]*P, then m["a"].X = 1 works
```

## Numbers

See Chapter 4 — Types, Variables, and Constants.

### No implicit numeric conversions

C silently promotes and converts between numeric types. Go does **not**: mixing
types is a compile error. You must convert on purpose.

```go
var i int = 5
var f float64 = 2.0
// x := i * f          // does NOT compile: mismatched types
x := float64(i) * f    // explicit conversion required
_ = x
```

> **C vs Go.** This removes a whole class of silent C bugs (lost precision, sign
> surprises), at the cost of more `int64(x)`-style noise. The compiler will never
> guess what you meant.

### `int` is 64-bit on 64-bit platforms

In C, `int` is almost always 32 bits, even on a 64-bit machine. In Go, `int` is the
machine word: **64 bits on a 64-bit platform**, 32 bits on a 32-bit platform. It is
*not* guaranteed to be 32 bits.

```go
var n int                       // 64 bits on a typical machine today
fmt.Println(unsafe.Sizeof(n))   // 8
```

> **Watch out.** When a size must be fixed (file formats, wire protocols, hashing),
> use a sized type: `int32`, `int64`, `uint8`, and so on. Do not assume `int` is 32
> bits the way you might in C.

### Integer division truncates; signed overflow wraps

Integer division truncates toward zero, as in C. The difference is **overflow**: in
C, signed integer overflow is *undefined behavior*; in Go it is **defined** — it
wraps around using two's complement.

```go
fmt.Println(7 / 2)   // 3 (truncated)
var x int8 = 127
x++                  // defined: wraps to -128 (not UB)
fmt.Println(x)       // -128
```

> **C vs Go.** A C compiler may *assume* signed overflow never happens and optimize
> on that assumption, so overflowing code can do anything. Go promises wraparound,
> so the behavior is portable and predictable (though usually still a bug you want
> to avoid).

### Do not compare floats with `==`

Floating-point math is not exact, so equality rarely holds. This is true in C too,
but it surprises people who expect `0.1 + 0.2` to equal `0.3`.

```go
x, y := 0.1, 0.2
sum := x + y
fmt.Println(sum == 0.3)               // false
fmt.Println(math.Abs(sum-0.3) < 1e-9) // true: compare within a tolerance
```

## Declarations, scope, and visibility

### `:=` can shadow an outer variable (especially `err`)

`:=` always declares a **new** variable in the current block. Inside an `if` or
`for` block, `x, err := f()` creates a *new* `err`, hiding the outer one. The famous
case is checking the inner `err` while the outer stays stale.

```go
data, err := load()        // outer err
if err == nil {
	data, err := parse(data) // BUG: := makes a NEW data and err, shadowing both
	_ = data
	_ = err                  // this err is checked; the outer one is not updated
}
// outer `data` is unchanged here — probably not what you wanted
```

> **Watch out.** Use `=` (not `:=`) when you mean to assign to an existing variable.
> `go vet` does not flag shadowing by default; the optional `shadow` analyzer does
> (Chapter 22 — Tooling). Blocks in `if`, `for`, and `switch` are new scopes
> (Chapter 5 — Control Flow).

### Unused imports and unused local variables are compile errors

In C these are warnings at most. In Go they **stop the build**. This feels harsh on
day one and keeps code clean forever (Chapter 1 — Why Go for a C Programmer).

```go
import "fmt"   // if unused: "imported and not used" — build fails
func f() {
	x := 1     // if unused: "declared and not used" — build fails
}
```

> **Watch out.** The escape hatch is the blank identifier `_`: `_ = x` to silence an
> unused variable on purpose, or `import _ "pkg"` for an import you need only for its
> side effects (Chapter 3 — Program Structure: Packages, Imports, and Visibility).
> Unused *function arguments* and *package-level variables* are allowed.

### Capitalization is visibility, not style

An identifier starting with an **uppercase** letter is exported (visible to other
packages); **lowercase** is package-private. This is a language rule, not a
convention (Chapter 3).

```go
type User struct {
	Name  string // exported: visible everywhere
	email string // unexported: only this package
}
```

> **Watch out.** Renaming `Total` to `total` is not cosmetic — it removes the name
> from every other package and breaks their code. It also hides the field from
> reflection-based libraries (see the JSON gotcha below).

## Control flow

See Chapter 5 — Control Flow.

### `switch` does not fall through

This is the **opposite** of C. Each `case` ends by itself; there is no implicit
fall-through and you do not write `break`. If you actually want C's behavior, say
so with the `fallthrough` keyword.

```go
switch x {
case 1:
	fmt.Println("one") // does NOT continue into case 2
case 2:
	fmt.Println("two")
	fallthrough        // explicit: now continue into the next case
case 3:
	fmt.Println("three")
}
```

> **Watch out.** A forgotten `break` is a classic C bug. In Go the danger is
> reversed: if you *relied* on fall-through, you must add `fallthrough` explicitly.

### `range` gives a copy; the range expression runs once

The loop value from `range` is a **copy** of each element. Writing to it does not
change the slice. Also, the thing you range over is evaluated **once**, before the
loop starts.

```go
type Big struct{ N int }
items := []Big{{1}, {2}}
for _, item := range items {
	item.N = 99        // changes the COPY; items is unchanged
}
for i := range items {
	items[i].N = 99    // correct: index into the slice to modify in place
}
```

```go
s := []int{1, 2, 3}
for range s {
	s = append(s, 0)   // loops exactly 3 times; the length is read once
}
```

### Loop-variable capture (historical — fixed since Go 1.22)

You will see old tutorials warn that a goroutine or closure created in a `for` loop
captures the *same* loop variable, so all of them see the final value. **That bug is
fixed.** Since Go 1.22, each iteration gets its **own** copy of the loop variable,
so capturing it is safe.

```go
// Go 1.22+ (this book targets 1.26): each i is per-iteration, so this is correct.
for i := range 3 {
	go func() { fmt.Println(i) }() // prints 0, 1, 2 (in some order)
}
```

> **Watch out.** This only matters when reading **old code** or tutorials. Before
> Go 1.22 the fix was to add `i := i` at the top of the loop body. You no longer
> need it, and leaving it in is harmless. (Chapter 13 — Goroutines and the
> Scheduler.)

## Functions and defer

See Chapter 6 — Functions.

### `defer`: when arguments and bodies run

`defer` schedules a call to run when the **function** returns. Two surprises: the
deferred call's **arguments are evaluated immediately**, at the `defer` line; and
defers run at *function* exit, not at the end of the enclosing block.

```go
func f() {
	x := 10
	defer fmt.Println("x was", x) // 10 is captured now, even though x changes
	x = 20
	// at return, this prints "x was 10"
}
```

> **Watch out.** A `defer` inside a loop does **not** run at the end of each
> iteration — every deferred call piles up until the function returns. Opening many
> files in a loop with `defer f.Close()` can exhaust file descriptors. Close inside
> the loop, or move the body into its own function.

```go
for _, name := range names {
	func() {
		f, err := os.Open(name)
		if err != nil {
			return
		}
		defer f.Close() // runs at the end of THIS inner function, each iteration
		// ... use f
	}()
}
```

> **Watch out.** `os.Exit` (and a fatal `log.Fatal`, which calls it) terminates the
> program **immediately and skips all deferred calls**. Do not rely on `defer` for
> cleanup right before `os.Exit`.

### Returning a pointer to a local is safe

This is the **opposite** of C, and it is worth saying loudly. In C, returning the
address of a local variable is a bug: the stack frame is gone and the pointer
dangles. In Go it is perfectly safe — the compiler's *escape analysis* sees that the
variable outlives the function and allocates it on the heap, and the garbage
collector frees it later (Chapter 7 — Pointers; Chapter 17 — Memory and the Garbage
Collector).

```c
int *bad(void) {
    int n = 0;
    return &n;   /* BUG in C: dangling pointer to a dead stack frame */
}
```

```go
func good() *int {
	n := 0
	return &n     // SAFE in Go: n escapes to the heap, GC frees it later
}
```

## Pointers and interfaces

### Dereferencing a nil pointer panics (it is defined)

Pointers exist in Go (`*T`, `&x`), but there is no pointer arithmetic. Dereferencing
`nil` does not crash with undefined behavior as in C — it raises a **defined**,
recoverable panic.

```go
var p *int
// _ = *p // PANIC: runtime error: invalid memory address or nil pointer dereference
```

> **C vs Go.** In C, dereferencing `NULL` is undefined behavior — it might
> segfault, might corrupt memory, might appear to work. In Go it always panics with
> a clear message, which you can even `recover` from (Chapter 12 — Errors).

### The typed-nil interface trap

This one fools experienced programmers. An interface value holds **two** parts: a
type and a value (Chapter 11 — Interfaces). An interface is `nil` **only when both
parts are nil**. If you put a nil *pointer* into an interface, the interface is
**not** nil, because it knows the type.

```go
type myErr struct{}
func (*myErr) Error() string { return "boom" }

func bad() error {
	var p *myErr = nil
	return p          // returns an interface holding (type=*myErr, value=nil)
}

// bad() == nil is FALSE, even though the pointer inside is nil!
```

> **Watch out.** Never declare a typed error variable and return it; return the
> literal `nil` for "no error", or return a real error value. This trap most often
> appears as "my function clearly returns nil but `err != nil` is true."

### What you can compare with `==`

You can compare with `==` (and use as map keys): booleans, numbers, strings,
pointers, channels, **structs and arrays whose fields/elements are all comparable**,
and interfaces. You **cannot** compare slices, maps, or functions — except to `nil`.

```go
type Point struct{ X, Y int }
fmt.Println(Point{1, 2} == Point{1, 2}) // true: structs compare field by field

var s1, s2 []int
// s1 == s2          // does NOT compile (slices are not comparable)
fmt.Println(s1 == nil) // the only comparison allowed: against nil
fmt.Println(slices.Equal(s1, s2)) // use slices.Equal for element-wise comparison
```

> **C vs Go.** In C, `==` on structs is not allowed at all and comparing arrays
> compares pointers. In Go, `==` on a comparable struct compares every field — handy,
> but it *panics at runtime* if the struct contains an interface holding an
> uncomparable value (like a slice).

## Concurrency

See Chapters 13–15.

### Goroutines: `main` exiting, leaks, and no external kill

`go f()` starts a lightweight thread called a *goroutine*. Three surprises:

- When `main` returns, the program exits and **all goroutines are killed**
  instantly, even mid-work. They are not joined for you.
- A goroutine blocked forever (on a channel that never receives, say) is a **leak**:
  it sits in memory until the program ends.
- You **cannot kill a goroutine** from outside. It must return on its own. The idiom
  is to signal it to stop with a `context.Context` or a channel (Chapter 15 —
  Synchronization and context).

```go
func main() {
	go fmt.Println("might not print")
	// main returns immediately; the goroutine may never run.
	// To wait, use a sync.WaitGroup or a channel.
}
```

> **Watch out.** Sharing memory between goroutines without synchronization is a
> **data race** and is undefined. Always test with the race detector: `go test
> -race ./...` or `go run -race .` (Chapter 16 — Concurrency Patterns).

### Channels: closed, nil, and unbuffered

Channels are typed pipes between goroutines (Chapter 14 — Channels and select). The
rules trip up newcomers:

- **Sending on a closed channel panics.** (Receiving from one is fine and returns
  the zero value with `ok == false`.)
- **A nil channel blocks forever** — both send and receive. A `var ch chan int` you
  forgot to `make` will hang the goroutine.
- **An unbuffered send blocks until someone receives.** If no one ever does, and no
  other goroutine can run, the runtime detects the deadlock and crashes.
- **Only the sender should close a channel**, and only once. Closing twice, or
  closing as the receiver, panics.

```go
ch := make(chan int) // unbuffered
// ch <- 1           // deadlock if no other goroutine receives

close(ch)
// ch <- 1           // PANIC: send on closed channel
v, ok := <-ch        // fine: v == 0, ok == false (closed and drained)
_ = v
_ = ok
```

### Do not copy a struct that contains a `sync.Mutex`

A `sync.Mutex` must not be copied after first use — a copy has its own lock state, so
two goroutines can each "hold" their own copy and the mutex protects nothing. The
trap is a **value receiver** on a method, which copies the whole struct (including
the mutex) on every call.

```go
type Counter struct {
	mu sync.Mutex
	n  int
}

// func (c Counter) Inc() { ... } // BUG: value receiver copies the mutex
func (c *Counter) Inc() {          // fix: pointer receiver, no copy
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}
```

> **Watch out.** `go vet` catches this with its `copylocks` check and reports
> something like *"Inc passes lock by value: Counter contains sync.Mutex"*. Run
> `go vet ./...` and never copy a value that contains a lock (Chapter 22 — Tooling;
> Chapter 15 — Synchronization and context).

## Serialization

### `encoding/json` only sees exported fields

The standard `encoding/json` package uses reflection, which can only read
**exported** (capitalized) fields. An unexported field is silently skipped on both
marshal and unmarshal (Chapter 20 — A Tour of the Standard Library; Chapter 3 for
visibility).

```go
type User struct {
	Name  string `json:"name"` // included; tag renames the key
	email string                // NOT serialized: unexported
}
b, _ := json.Marshal(User{Name: "Ann", email: "secret"})
fmt.Println(string(b)) // {"name":"Ann"}
```

> **Watch out.** If a field "disappears" from your JSON, check its capitalization
> first. Export every field you want serialized, and use a struct tag to control the
> key name. The same rule applies to `encoding/xml`, `gob`, and most other
> reflection-based encoders.

## Key takeaways

- Arrays are values (copied); strings are immutable UTF-8 bytes. Use slices for
  dynamic data, but remember sub-slices **share** a backing array and `append` may
  mutate a neighbor or reallocate.
- A nil map is read-only — writing panics; map order is random; map values are not
  addressable.
- Go has no implicit numeric conversions; `int` is word-sized (64-bit on 64-bit);
  signed overflow **wraps** (defined), unlike C's undefined behavior; never compare
  floats with `==`.
- `:=` can shadow (watch `err`); unused imports and locals are **compile errors**;
  capitalization is visibility.
- `switch` does not fall through; `range` yields a copy and evaluates its expression
  once; loop-variable capture is **fixed** since Go 1.22.
- `defer` evaluates arguments immediately and runs at function return; defers in a
  loop accumulate; `os.Exit` skips them. Returning `&local` is **safe** in Go.
- A nil dereference panics (defined); a nil pointer inside an interface makes the
  interface **non-nil** (the typed-nil trap); you can `==` structs/arrays but not
  slices/maps/functions.
- Goroutines die when `main` exits, can leak, and cannot be killed externally;
  channel misuse (send-on-closed, nil, unbuffered with no receiver) panics or
  deadlocks; do not copy a `sync.Mutex`; `encoding/json` sees only exported fields.

## Watch out (the one-line checklist)

Scan this list before you ship; each line is a trap explained above.

- Passing an array copies it — use a slice.
- Sub-slices alias; always `s = append(s, x)`; `copy` to get an independent slice.
- A small slice pins its whole backing array — `copy` out to release the rest.
- `len(string)` is bytes; `string[i]` is a byte; strings are immutable UTF-8.
- Writing a nil map panics — `make` it first; map order is random.
- No implicit numeric conversions; `int` is 64-bit here; signed overflow wraps.
- Do not compare floats with `==`.
- `:=` can shadow `err`; unused imports/locals fail the build.
- Capitalization decides exported vs unexported.
- `switch` has no fall-through; `range` value is a copy.
- Loop capture is fixed since Go 1.22 (only old code needs `i := i`).
- `defer` args evaluate now; defers run at function return; loops accumulate them;
  `os.Exit` skips them.
- Returning `&local` is safe (escape analysis).
- nil dereference panics; typed-nil makes an interface non-nil.
- `==` works on structs/arrays, not slices/maps/functions.
- `main` returning kills all goroutines; they can leak; you cannot kill one
  externally; use `-race`.
- Send-on-closed panics; nil channel blocks forever; unbuffered send needs a
  receiver; only the sender closes.
- Never copy a `sync.Mutex` (value receivers do this — `go vet` catches it).
- `encoding/json` serializes only exported fields.

## Interview questions

**Q: A colleague slices a request body to keep its first 16 bytes and the server's
memory keeps growing. What is happening and how do you fix it?**
A: A slice keeps its entire backing array alive, so `body[:16]` pins the whole
(possibly large) body in memory and the garbage collector cannot free it. Fix it by
copying the bytes you need into a fresh slice — `out := make([]byte, 16);
copy(out, body[:16])` — so the original backing array can be collected.

**Q: Why might a function that "returns nil" still fail an `err != nil` check?**
A: The typed-nil interface trap. An interface value stores a type and a value, and
is nil only when both are nil. If the function declares a typed pointer (for
example `var p *MyError`) and returns it, the returned `error` interface holds a
non-nil type with a nil value, so it is not equal to nil. Return the literal `nil`
for success instead of a typed nil pointer.

**Q: How does Go's `switch` differ from C's, and how do `for`/`range` loop variables
behave in modern Go?**
A: Go's `switch` does not fall through; each case ends on its own, and you opt into
C-style fall-through with the `fallthrough` keyword. As of Go 1.22, each loop
iteration gets its own copy of the loop variable, so capturing it in a goroutine or
closure is safe — the old "all closures see the last value" bug is gone and the
`i := i` workaround is no longer needed.

**Q: Why does writing to a map sometimes panic, and reading never does?**
A: A map's zero value is nil. Reading a nil map is defined and returns the element
type's zero value, but writing to a nil map panics because there is no underlying
hash table to store into. You must create the map first with `make` or a composite
literal. This commonly bites when a struct has a map field that was never
initialized.

**Q: Give three ways concurrency in Go surprises a C programmer.**
A: First, when `main` returns the whole program exits and all goroutines are killed
immediately — they are not joined for you. Second, you cannot kill a goroutine from
outside; it must return on its own, so you signal cancellation with a context or
channel, and a goroutine blocked forever leaks. Third, channel rules cause panics or
deadlocks: sending on a closed channel panics, a nil channel blocks forever, and an
unbuffered send with no receiver deadlocks. Shared memory without synchronization is
a data race, so test with `-race`.

**Q: Why is returning the address of a local variable safe in Go but a bug in C?**
A: In C, a local lives in the function's stack frame, which is reclaimed on return,
so a returned pointer dangles. In Go, the compiler performs escape analysis: if a
local's address outlives the function, the compiler allocates it on the heap
instead of the stack, and the garbage collector frees it when nothing references it.
So `return &n` is a normal, safe idiom in Go.
