# Chapter 4 — Types, Variables, and Constants

> **What you'll learn.** How Go declares variables and why the type comes *after*
> the name, the exact sizes of Go's basic types, the all-important **zero value**,
> why Go refuses implicit numeric conversions, and how constants and `iota` work.
> Everything is compared directly to C.

C and Go are both statically typed: every variable has a type fixed at compile
time. But Go makes three choices that surprise a C programmer right away. It writes
the type *after* the name. It gives every variable a defined starting value instead
of garbage. And it refuses to convert between number types for you. This chapter
covers all three, plus constants and named types.

## Declaring variables

Go has one keyword for declaring a variable: `var`. There are a few forms.

```go
var x int        // declare x of type int; it starts at 0 (the zero value)
var y int = 5    // declare and initialize
var z = 5        // declare and initialize; type is inferred as int
```

Inside a function you also get the **short variable declaration** `:=`, which
declares and infers the type in one step:

```go
w := 5           // same as: var w int = 5
name := "Ada"    // string
ok := true       // bool
```

> **Watch out.** `:=` works **only inside a function**. At package level (outside
> any function) you must use `var`. Also, `:=` must declare at least one *new*
> variable on its left; otherwise use plain `=` to assign.

You can group several `var` declarations in one block. This is the common style for
package-level variables, and it reads like a small table:

```go
var (
	count   int        // 0
	rate    float64    // 0.0
	label   string     // ""
	enabled bool       // false
)
```

Go also supports **multiple assignment**, which assigns several variables at once.
The right-hand side is fully evaluated first, so swapping needs no temporary:

```go
a, b := 1, 2
a, b = b, a          // swap; no temp variable, no XOR trick
fmt.Println(a, b)    // 2 1
```

This is also how you receive **multiple return values**, which are common in Go
(see Chapter 6 — Functions):

```go
v, err := strconv.Atoi("42")   // v is int, err is error
```

| Concept | C | Go |
|---|---|---|
| Declare, uninitialized | `int x;` (garbage value) | `var x int` (zero value 0) |
| Declare with value | `int x = 5;` | `var x = 5` or `x := 5` |
| Infer the type | not possible (must name it) | `x := 5` or `var x = 5` |
| Declare many | `int a, b, c;` | `var ( a int; b int; c int )` |
| Swap two values | needs a temp | `a, b = b, a` |

### Why the type comes after the name

In C, the type wraps around the name, and complex declarations become hard to read.
The classic example:

```c
int *p[10];        /* p is an array of 10 pointers to int... or is it? */
int (*q)[10];      /* q is a pointer to an array of 10 int */
char *(*f)(int);   /* f is a pointer to a function taking int, returning char* */
```

You must read these "inside out," following C's precedence rules. Go puts the name
first and the type second, so every declaration reads left to right, like a
sentence:

```go
var p [10]*int               // p is an array of 10 pointers to int
var q *[10]int               // q is a pointer to an array of 10 ints
var f func(int) *byte        // f is a function taking int, returning *byte
```

> **Mental model.** Read a Go declaration left to right: "`p` is a `[10]` array of
> `*int`." The type is a phrase that follows the name. There is no "spiral rule"
> like C's.

> **C vs Go.** In C the `*` binds to the variable, so `int* a, b;` declares one
> pointer and one plain `int` — a famous trap. In Go the type is written once and
> applies to the whole declaration, so this confusion cannot happen.

## The basic types

Go's numeric types have **exact, fixed sizes** that are the same on every platform.
This is a real difference from C, where the standard only promises *minimum* sizes
(`int` is at least 16 bits, `long` at least 32) and the real size depends on the
compiler and platform.

| Go type | Size | Meaning | Closest C type |
|---|---|---|---|
| `bool` | 1 byte | `true` or `false` | `_Bool` / `bool` |
| `int8` / `uint8` | 1 byte | 8-bit signed / unsigned | `int8_t` / `uint8_t` |
| `int16` / `uint16` | 2 bytes | 16-bit | `int16_t` / `uint16_t` |
| `int32` / `uint32` | 4 bytes | 32-bit | `int32_t` / `uint32_t` |
| `int64` / `uint64` | 8 bytes | 64-bit | `int64_t` / `uint64_t` |
| `int` / `uint` | **word size** | 32-bit on 32-bit CPUs, **64-bit on 64-bit CPUs** | `int` / `size_t` (roughly) |
| `uintptr` | word size | unsigned int big enough to hold a pointer's bits | `uintptr_t` |
| `byte` | 1 byte | **alias for `uint8`**; raw data | `unsigned char` |
| `rune` | 4 bytes | **alias for `int32`**; one Unicode code point | `wchar_t` (loosely) |
| `float32` | 4 bytes | IEEE-754 single precision | `float` |
| `float64` | 8 bytes | IEEE-754 double precision | `double` |
| `complex64` / `complex128` | 8 / 16 bytes | complex numbers (built in) | none |
| `string` | 16-byte header | immutable UTF-8 bytes | `char *` + length |

A few points a C programmer must internalize:

- **`int` is the machine word, not 32 bits.** On any modern 64-bit machine, Go's
  `int` is **64 bits**. Do not assume it is 32. When you need an exact width (file
  formats, wire protocols, hardware registers), use a sized type like `int32` or
  `uint64`, exactly as you would use `int32_t` in C.
- **`byte` and `rune` are aliases**, not new types. `byte` is literally `uint8` and
  `rune` is literally `int32`. Use `byte` for raw bytes and `rune` for a Unicode
  code point; the names document intent (see Chapter 8 — Arrays, Slices, and
  Strings).
- **`string` is its own type**, not a `char *`. It carries a length and its bytes
  are immutable. Strings get full treatment in Chapter 8.
- Complex numbers exist (`complex128`, the builtins `complex`, `real`, `imag`) but
  are rarely used outside scientific code; we will not dwell on them.

Here is the byte layout, the kind of picture a C programmer likes to see. Each box
is one byte:

```
  int8  / uint8   [#]                            1 byte
  int16 / uint16  [#][#]                         2 bytes
  int32 / uint32  [#][#][#][#]                    4 bytes   (rune is int32)
  int64 / uint64  [#][#][#][#][#][#][#][#]        8 bytes
  int   / uint    [#][#][#][#][#][#][#][#]        word size: 8 on 64-bit CPUs
                                                             4 on 32-bit CPUs

  The sized types are identical on every platform.
  Only int, uint, and uintptr change with the machine word.
```

You can ask for a type's size at runtime with `unsafe.Sizeof`, the rough equal of
C's `sizeof`. It returns a `uintptr` (a count of bytes):

```go
import "unsafe"

var n int
fmt.Println(unsafe.Sizeof(n)) // 8 on a 64-bit machine
fmt.Println(unsafe.Sizeof(int32(0))) // always 4
```

> **Watch out.** Despite the package name, `unsafe.Sizeof` is safe and computed at
> compile time. The `unsafe` package is only "unsafe" when you use it to bypass the
> type system (pointer tricks); just asking for a size is fine.

## Zero values

This is one of Go's most important rules, and it removes a whole category of C bugs.

**Every variable is initialized to its type's *zero value* when declared, even if
you do not assign one.** There is no uninitialized garbage in Go.

| Type | Zero value |
|---|---|
| numbers (`int`, `float64`, ...) | `0` / `0.0` |
| `bool` | `false` |
| `string` | `""` (empty string) |
| pointers, slices, maps, channels, functions, interfaces | `nil` |
| struct | each field set to *its* zero value |

```go
var i int        // 0
var f float64    // 0.0
var s string     // "" (empty, not nil — a string is never nil)
var p *int       // nil
var b bool       // false
```

```
C:  int x;          // x holds whatever bytes were on the stack (garbage)
Go: var x int       // x is exactly 0, guaranteed
```

> **C vs Go.** In C, a local `int x;` has an indeterminate value; reading it before
> assignment is undefined behavior, and forgetting to initialize is a classic bug.
> In Go, the same declaration always yields `0`. Static and global variables in C
> *are* zero-initialized, but locals are not — Go makes the rule uniform.

Go's standard library is designed so the zero value is **useful**: a freshly
declared variable is often ready to use without any setup.

```go
var buf bytes.Buffer    // ready to use; no New(), no init
buf.WriteString("hi")   // works immediately

var mu sync.Mutex       // an unlocked mutex, ready to Lock()
mu.Lock()
```

> **Rule of thumb.** When you design your own types, try to make the zero value do
> something sensible. A struct that works "empty" is easier to use than one that
> needs a constructor call first.

## No implicit conversions

C converts between number types automatically, sometimes silently losing data:

```c
int    i = 300;
char   c = i;        /* silently truncates to 44 on most systems */
double d = i;        /* int promoted to double, no warning */
double r = 7 / 2;    /* r is 3.0: integer division happened first! */
```

Go refuses all of this. **You must convert types explicitly**, with a conversion
that looks like a function call: `T(value)`.

```go
var i int = 300
var c byte = byte(i)      // explicit; you chose to truncate (c is 44)
var d float64 = float64(i) // explicit widening
var r float64 = float64(7) / float64(2) // 3.5; convert before dividing
```

If you forget the conversion, it is a **compile error**, not a warning:

```go
var i int = 10
var f float64 = i   // compile error:
// cannot use i (variable of type int) as float64 value in variable declaration
```

This even applies between `int` and `int64` on a 64-bit machine, *where they have
the same size*. They are still different types, so you must write `int64(x)`. The
compiler is checking *types*, not bit widths.

> **C vs Go.** C's "usual arithmetic conversions" are convenient but cause silent
> truncation, signed/unsigned surprises, and accidental integer division. Go trades
> that convenience for safety: every conversion is visible in the code, so you can
> see exactly where a value changes type.

### The one exception: untyped constants

There is a single, deliberate exception, and it is what makes Go pleasant to write
despite the strictness. **Constant expressions are *untyped* until they are used.**
An untyped constant adapts to whatever type the context needs, with no explicit
conversion.

```go
var x float64 = 3        // fine: the untyped constant 3 becomes a float64
var y int = 3            // fine: the same 3 becomes an int
const big = 1 << 40
var z int64 = big        // fine
f := 2.5
fmt.Println(f * 2)       // fine: 2 is untyped, becomes float64 here
```

Compare that to a *variable*, which has a fixed type and will not auto-convert:

```go
i := 3                   // i is int (a typed variable now)
var w float64 = i        // ERROR: i is int, not an untyped constant
var w2 float64 = 3       // OK: 3 is an untyped constant
```

We cover constants next; this exception is the reason they are special.

## Constants

A `const` is a value fixed at compile time. Unlike a C `#define` (which is just
text substitution by the preprocessor) or a C `const` variable (which still lives
in memory), a Go constant is a true compile-time value with no storage and no
address.

```go
const greeting = "hello"
const pi = 3.14159
const (
	maxRetries = 3
	timeoutSec = 30
)
```

### Typed vs untyped constants

A constant can be **untyped** (the default) or **typed**.

```go
const a = 100          // untyped: adapts to context, very flexible
const b int = 100      // typed: b is exactly int, behaves like a variable's type
```

Untyped constants have two powers:

1. **They adapt to the needed type** (the exception from the previous section).
2. **They are computed with very high precision** — far beyond `int64` or
   `float64` — and only rounded to a real type when assigned to one. This means
   constant math does not overflow during the calculation.

```go
const Huge = 1 << 100      // fine as an untyped constant (arbitrary precision)
fmt.Println(Huge >> 97)    // 8 — the expression result fits in int, so it prints
// var n int = Huge        // ERROR: constant overflows int (caught at compile time)
```

> **Watch out.** Assigning an out-of-range constant is a **compile error**, caught
> before your program ever runs. C might warn (or not) and then silently wrap or
> truncate at runtime. Go stops the build:
> `1 << 40 (untyped int constant ...) as int32 value ... (overflows)`.

A typed constant behaves like a variable of that type: it does *not* auto-convert.

```go
const k int = 3
var f float64 = k    // ERROR: k is typed int; needs float64(k)
```

> **Rule of thumb.** Leave constants untyped unless you have a specific reason to
> pin a type. Untyped constants compose more easily in expressions.

### `iota`: Go's enum generator

Go has no `enum` keyword. Instead it uses `iota`, a counter that resets to `0` at
the start of each `const` block and increases by `1` for each line.

```go
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1 (the expression "Weekday = iota" repeats)
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)
```

After the first line sets the pattern `Weekday = iota`, the following lines repeat
it automatically, with `iota` increasing each time. This is the standard way to
build an enumeration in Go (compare to C's `enum { SUNDAY, MONDAY, ... }`).

Because `iota` can appear inside an expression, it is perfect for **bit flags**
(`1 << iota`) and for sizes:

```go
type Permission uint

const (
	Read    Permission = 1 << iota // 1  (0b001)
	Write                          // 2  (0b010)
	Execute                        // 4  (0b100)
)

const (
	_  = iota             // skip 0 with the blank identifier
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20
	GB                    // 1 << 30
)
```

You combine and test flags with the same bitwise operators you know from C:

```go
p := Read | Write          // turn two bits on
if p&Write != 0 {          // test a bit
	fmt.Println("writable")
}
```

| Concept | C | Go |
|---|---|---|
| Named constant | `#define MAX 3` or `const int MAX = 3;` | `const Max = 3` |
| Enumeration | `enum { RED, GREEN };` | `const ( Red = iota; Green )` |
| Bit flags | `#define R (1<<0)` | `R = 1 << iota` |
| Computed at | preprocessor / runtime | compile time, high precision |
| Has an address | `const` variable does | no (pure value) |

## Named types vs type aliases

You can give a type a new name with `type`. There are two forms, and the difference
matters.

A **named type** (also called a *defined type*) creates a brand-new, distinct type
with the same underlying representation:

```go
type Celsius float64
type Fahrenheit float64
```

`Celsius` and `Fahrenheit` are both `float64` underneath, but they are **different
types**. You cannot mix them by accident, which prevents unit-confusion bugs:

```go
var c Celsius = 100
var f Fahrenheit = 100
// c = f          // ERROR: cannot use f (Fahrenheit) as Celsius
c = Celsius(f)    // OK: explicit conversion when you really mean it
```

A **type alias** (note the `=`) creates a second *name for the exact same type*.
There is no new type at all:

```go
type Float = float64   // Float and float64 are interchangeable
var x Float = 3
var y float64 = x      // OK: same type, no conversion needed
```

This is precisely how `byte` and `rune` are defined in the language:
`type byte = uint8` and `type rune = int32`.

| | Named type `type T U` | Alias `type T = U` |
|---|---|---|
| Creates a new type? | Yes, distinct from `U` | No, same type as `U` |
| Needs conversion to mix with `U`? | Yes | No |
| Can have its own methods? | Yes | No (it is just `U`) |
| Main use | new domain types, enums | renaming during refactors, compatibility |

The big payoff of a named type: **you can attach methods to it.** This is how Go
gives behavior to a plain number or string without C's tricks. We cover methods in
Chapter 10 — Structs and Methods, but here is a taste:

```go
func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}
```

## Printing and inspecting types

The `fmt` package has verbs that are invaluable while learning. Use them to see a
value *and* its type.

```go
type Point struct{ X, Y int }
p := Point{1, 2}

fmt.Printf("%v\n", p)   // {1 2}            default format
fmt.Printf("%+v\n", p)  // {X:1 Y:2}        adds field names
fmt.Printf("%#v\n", p)  // main.Point{X:1, Y:2}   Go syntax
fmt.Printf("%T\n", p)   // main.Point       the type
fmt.Printf("%T\n", 3.0) // float64          default type of an untyped float
fmt.Printf("%T\n", 'A') // int32            a rune literal is int32
```

| Verb | Shows | Like C's |
|---|---|---|
| `%v` | value in a default format | depends (`%d`, `%f`, ...) |
| `%+v` | value with struct field names | no equivalent |
| `%#v` | value as Go source code | no equivalent |
| `%T` | the dynamic type of the value | no equivalent |

> **Rule of thumb.** When something behaves oddly, print it with `%T` and `%#v`
> first. Nine times out of ten the surprise is "this value is not the type I
> thought it was."

## Key takeaways

- Declare with `var name type`, `var name = value`, or (inside functions) `name :=
  value`. Group with `var ( ... )`; assign many at once with `a, b = b, a`.
- The **type comes after the name**, so declarations read left to right and avoid
  C's confusing pointer/array syntax.
- Sized integer types (`int8`..`int64`, `uint8`..`uint64`) are **exact and the same
  everywhere**. `int`/`uint`/`uintptr` are the **machine word** — 64 bits on 64-bit
  machines, *not* always 32.
- Every variable starts at its **zero value** (`0`, `0.0`, `false`, `""`, `nil`).
  There is no uninitialized garbage. Useful zero values are a design goal.
- Go has **no implicit numeric conversions**; write `float64(i)`, `int64(x)`, etc.
  The single exception is **untyped constants**, which adapt to context.
- Constants are compile-time values. Untyped ones have high precision and a default
  type; out-of-range assignments fail at **compile time**. `iota` builds enums and
  bit flags.
- `type T U` makes a **distinct** type (you can add methods); `type T = U` is a mere
  **alias**. `byte` and `rune` are aliases for `uint8` and `int32`.

## Watch out (gotchas for C programmers)

- **`int` is not 32 bits.** It is the word size (64 bits on 64-bit CPUs). Use a
  sized type when the width must be exact.
- **No implicit conversion**, even between same-size types like `int` and `int64`.
  The compiler checks types, not bit widths.
- **`:=` opens a new scope and can shadow.** Inside an `if` or `for`, `x := ...`
  may create a *new* `x` that hides an outer one. We cover this trap in Chapter 5 —
  Control Flow.
- **Unused local variables are a compile error**, not a warning (so are unused
  imports). Declare it, use it, or remove it.
- **Integer division truncates toward zero**, just like C: `7 / 2 == 3`. Convert to
  `float64` first if you want `3.5`.
- **Signed overflow is defined, not undefined.** Go uses two's-complement wraparound
  for integers: `int8(127) + 1 == -128`. In C, signed overflow is *undefined
  behavior*; in Go it is specified and predictable. (Constant overflow, by contrast,
  is caught at compile time.)

## Interview questions

**Q: Why does Go put the type after the variable name?**
A: So declarations read left to right like a sentence and avoid C's "inside-out"
pointer and array syntax. `var p [10]*int` reads "p is an array of 10 pointers to
int," with no precedence puzzles. It also means the type is written once per
declaration, so the `int* a, b;` trap (only `a` is a pointer) cannot happen.

**Q: What is a zero value, and how does it differ from C?**
A: Every Go variable is automatically set to its type's zero value when declared:
`0` for numbers, `false` for `bool`, `""` for strings, and `nil` for pointers,
slices, maps, channels, functions, and interfaces. In C, local variables are
uninitialized (garbage) unless you assign them; only statics and globals are
zeroed. Go makes the rule uniform and reliable.

**Q: How big is a Go `int`? How is that different from C?**
A: Go's `int` is the machine word size — 32 bits on 32-bit platforms and 64 bits on
64-bit platforms — and it is its own type. The sized types (`int8`..`int64`) are
exact on every platform. In C, `int`, `long`, etc. have only minimum sizes
guaranteed and vary by compiler/platform; you reach for `<stdint.h>` types like
`int32_t` for exact widths, just as you use `int32` in Go.

**Q: What are untyped constants and why are they useful?**
A: Constant expressions in Go are untyped by default. They have arbitrary precision
and adopt whatever type the surrounding context requires, so `var x float64 = 3` and
`var y int = 3` both work without a conversion. This is the one exception to Go's
"no implicit conversion" rule, and it lets constant math stay exact until it is
assigned to a concrete type.

**Q: What is the difference between `type Celsius float64` and `type Alias =
float64`?**
A: The first defines a new, distinct type whose underlying type is `float64`; it is
not interchangeable with `float64` (you must convert), and you can attach methods to
it. The second is a type alias: `Alias` is just another name for `float64`, fully
interchangeable, and cannot have its own methods. `byte` and `rune` are aliases for
`uint8` and `int32`.

## Try it

1. Declare `var x int8 = 127`, then print `x + 1`. Confirm it wraps to `-128`
   (defined behavior), unlike C's undefined signed overflow.
2. Define `type Meters float64` and `type Feet float64`, declare one of each, and
   try to add them. Read the compile error, then fix it with an explicit
   conversion.
3. Build a `const` block of permission flags using `1 << iota`, combine two with
   `|`, and test one with `&`. Print the result with `%b` to see the bits.
