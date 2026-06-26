# Chapter 4 — Variables, Types, and Constants

> **What you'll learn.** How `let` bindings work and why variables are immutable
> by default, the difference between shadowing and mutation, Rust's exact-sized
> scalar and compound types, why there are no implicit conversions, and how
> `const` and `static` differ.

## Variables with `let`

In C you declare a variable by writing its type and name: `int x = 5;`. In Rust
you write `let`, and the type is usually inferred:

```rust
fn main() {
    let x = 5;          // type inferred as i32
    let y: i32 = 5;     // type written explicitly
    println!("{x} {y}");
}
```

The keyword `let` introduces a **binding**: it ties a name to a value. The
compiler figures out the type from the value on the right, so you rarely write
the type. When you want to be explicit, or when inference cannot decide, you add
`: Type` after the name, as in `let y: i32 = 5`.

> **C vs Rust.** C: `int x = 5;` (type first). Rust: `let x = 5;` (type inferred,
> or written after the name as `let x: i32 = 5;`). The name comes before the type
> in Rust, the reverse of C.

## Immutable by default

Here is the first big surprise for a C programmer. In C, every variable is
mutable: you can assign to it again whenever you like. In Rust, a `let` binding
is **immutable** by default — once set, you cannot change it.

```rust
// COMPILE ERROR: cannot assign twice to immutable variable `x`
fn main() {
    let x = 5;
    x = 6;          // error[E0384]: cannot assign twice to immutable variable `x`
    println!("{x}");
}
```

To make a variable changeable, add the keyword `mut` (short for "mutable"):

```rust
fn main() {
    let mut x = 5;
    x = 6;          // ok: x is mutable
    x += 1;
    println!("{x}");    // prints 7
}
```

This is the opposite default from C. In C you write `const` to *forbid* change;
in Rust you write `mut` to *allow* it. The default is the safe one.

> **C vs Rust.** C is mutable by default; you add `const` to lock a value. Rust
> is immutable by default; you add `mut` to unlock it. The defaults are reversed.

> **Rule of thumb.** Start without `mut`. Add it only when the compiler tells you
> that you need it. Code with fewer `mut` bindings is easier to reason about,
> because you know those values never change.

## Shadowing vs mutation

Rust lets you declare a new binding with the **same name** as an old one. This is
called **shadowing**. The new binding hides the old one from that point on:

```rust
fn main() {
    let x = 5;
    let x = x + 1;      // a NEW binding, also named x, value 6
    let x = x * 2;      // another NEW binding, value 12
    println!("{x}");    // prints 12
}
```

Shadowing is not mutation. Each `let` creates a brand-new variable that happens
to reuse the name. The old value is still there underneath; it is just hidden.
Because it is a new variable, it can even have a **different type**:

```rust
fn main() {
    let spaces = "   ";          // spaces is &str (text)
    let spaces = spaces.len();   // spaces is now usize (a number)
    println!("{spaces}");        // prints 3
}
```

Compare that with `mut`, which keeps the *same* variable and so cannot change its
type:

```rust
// COMPILE ERROR: mismatched types
fn main() {
    let mut spaces = "   ";
    spaces = spaces.len();   // error[E0308]: expected `&str`, found `usize`
}
```

> **Mental model.** Mutation changes the value inside one box. Shadowing throws
> the old box away and puts a new box (possibly a different shape) under the same
> label. C has no shadowing in the same scope — a second declaration of the same
> name is an error.

> **Watch out.** Shadowing is easy to confuse with mutation. `let x = ...` again
> is a new variable; `x = ...` (no `let`) is assignment and needs `mut`. Reach for
> shadowing when you transform a value into a new form (often a new type); reach
> for `mut` when one value genuinely changes over time.

## Type inference and annotations

Rust uses **type inference**: the compiler reads how a value is created and used,
then assigns a type. This is more powerful than C's, because it can look at later
uses, not just the initializer. Sometimes that means you must help it:

```rust
fn main() {
    let guess: u32 = "42".parse().expect("not a number");  // annotation needed
    println!("{guess}");
}
```

Here `parse` can produce many number types, so the compiler cannot guess. The
`: u32` tells it which one. Without the annotation you get a compile error asking
you to specify the type.

## Scalar types: integers

Rust's integer types state their size and signedness **in the name**. There is no
guessing about how many bits an `int` has, the way there is in C.

| Length | Signed | Unsigned |
|---|---|---|
| 8-bit | `i8` | `u8` |
| 16-bit | `i16` | `u16` |
| 32-bit | `i32` | `u32` |
| 64-bit | `i64` | `u64` |
| 128-bit | `i128` | `u128` |
| pointer-sized | `isize` | `usize` |

```
 i8  : [#]                      1 byte
 i16 : [##]                     2 bytes
 i32 : [####]                   4 bytes   <-- default
 i64 : [########]               8 bytes
 i128: [################]      16 bytes
 isize/usize: pointer-sized   (8 bytes on a 64-bit machine)
```

The `i` means signed (it can be negative), the `u` means unsigned (zero or
positive). The number is the width in bits. The two odd ones are `isize` and
`usize`: they are **pointer-sized**, meaning they match the address width of the
machine (4 bytes on a 32-bit target, 8 bytes on a 64-bit target). This is like
`intptr_t`/`uintptr_t` and `ptrdiff_t`/`size_t` in C.

`usize` is special: Rust uses it for **all indexing and lengths**. The length of
an array or a `Vec`, and any value you use to index into one, is a `usize`. This
is the same reason C uses `size_t` for `sizeof` and array sizes.

> **Watch out.** If you do not annotate an integer and nothing forces a type, Rust
> defaults to **`i32`**. In C the default integer type is also `int`, usually
> 32-bit, so this will feel familiar — but in Rust it is a precise, fixed-width
> `i32`, not an "at least 16 bits" `int`.

### Integer literals

You can write numbers in several bases and add a `_` as a digit separator to make
big numbers readable. You can also glue the type onto the end as a **suffix**:

```rust
fn main() {
    let a = 1_000_000;     // decimal, underscores ignored by the compiler
    let b = 0xff;          // hex
    let c = 0o77;          // octal
    let d = 0b1010_0001;   // binary
    let e = b'A';          // a byte literal: the u8 value 65
    let f = 1_000u64;      // suffix: this literal is a u64
    println!("{a} {b} {c} {d} {e} {f}");
}
```

The suffix (`1_000u64`) is a quick way to pick a type without a separate
annotation. It is handy when a literal alone would be ambiguous.

### Integer overflow

In C, signed integer overflow is **undefined behavior**, and unsigned overflow
wraps silently. Rust is stricter and more predictable:

- In a **debug** build, an overflowing `+`, `-`, or `*` **panics** (the program
  stops with an error message) instead of silently giving a wrong answer.
- In a **release** build (`--release`), it **wraps around** using two's
  complement, like C's unsigned arithmetic — fast, but you opted into it.

Relying on either default is usually a mistake. When you actually want a specific
overflow behavior, say so with a method:

```rust
fn main() {
    let big: u8 = 250;

    println!("{}", big.wrapping_add(10));        // 4   (wraps around 256)
    println!("{:?}", big.checked_add(10));       // None (would overflow)
    println!("{:?}", big.checked_add(3));        // Some(253)
    println!("{}", big.saturating_add(10));      // 255 (clamps at the max)
    let (v, overflowed) = big.overflowing_add(10);
    println!("{v} {overflowed}");                // 4 true
}
```

- `wrapping_*` — wrap around on overflow (the C unsigned behavior).
- `checked_*` — return `Option`: `Some(result)` or `None` if it overflowed.
- `saturating_*` — stop at the type's minimum or maximum instead of wrapping.
- `overflowing_*` — return the wrapped result plus a `bool` flag.

> **C vs Rust.** C: signed overflow is undefined behavior, a real source of bugs.
> Rust: overflow panics in debug builds so you find it, and you choose explicit
> behavior with `wrapping_*`, `checked_*`, or `saturating_*` when you need it.

## Scalar types: floating point, bool, char

### Floating point

Rust has exactly two floating-point types: `f32` (single precision) and `f64`
(double precision). The default is `f64`. They follow the IEEE 754 standard, the
same as C's `float` and `double`.

```rust
fn main() {
    let x = 2.0;          // f64 by default
    let y: f32 = 3.0;     // f32
    println!("{x} {y}");
}
```

### Booleans

`bool` has two values, `true` and `false`, and is one byte. The key difference
from C: a `bool` is **not** a number. You cannot use an integer where a `bool` is
expected, and you cannot do arithmetic on a `bool`.

```rust
// COMPILE ERROR: mismatched types, expected `bool`, found integer
fn main() {
    let flag = 1;
    if flag {            // error[E0308]: expected `bool`, found integer
        println!("yes");
    }
}
```

In C, `if (n)` and `if (ptr)` are common: any nonzero value is "true." Rust has
no such rule. A condition must be a real `bool`. We cover this fully in Chapter 5
— Control Flow.

### Characters

A Rust `char` is a **Unicode scalar value**: it is **4 bytes** and can hold any
single Unicode character, not just ASCII. This is very different from a C `char`,
which is one byte and is really just a small integer.

```rust
fn main() {
    let letter = 'A';        // char, 4 bytes
    let pi = 'π';            // also one char
    let crab = '🦀';         // a single char, even though it is multi-byte as text
    println!("{letter} {pi} {crab}");
}
```

Write a `char` with single quotes (`'A'`) and text with double quotes (`"A"`).
A `char` is one character; a string is a sequence of bytes encoded as UTF-8. If
you want C's "one byte" meaning, use `u8` (and the byte literal `b'A'`).

> **C vs Rust.** C `char` is one byte and doubles as a tiny integer. Rust `char`
> is a 4-byte Unicode scalar value and is **not** an integer. For a raw byte, use
> `u8`. We cover text and UTF-8 in Chapter 10 — Slices and Strings.

> **Watch out.** Single quotes and double quotes mean different types in Rust.
> `'a'` is a `char`; `"a"` is a string slice (`&str`). They are not
> interchangeable, unlike the loose char/string handling habits from C.

## Compound types: tuples

A **tuple** groups a fixed number of values, which may have different types, into
one value. It is like an anonymous, lightweight struct.

```rust
fn main() {
    let point: (i32, f64, char) = (1, 2.5, 'z');

    // Destructuring: unpack the tuple into separate names.
    let (a, b, c) = point;
    println!("{a} {b} {c}");

    // Or access fields by index with a dot.
    println!("{} {} {}", point.0, point.1, point.2);
}
```

You read a tuple field with `.0`, `.1`, `.2`, and so on. **Destructuring** lets
you unpack all the fields at once into separate names. The size and the type of
each position are fixed when you write the tuple.

### The unit type

The empty tuple `()` is called the **unit type**. It has exactly one value, also
written `()`, and carries no information. It is what a function returns when it
returns "nothing" — the rough equivalent of C's `void`. You will see it in
Chapter 6 — Functions and Closures.

## Compound types: arrays

An **array** in Rust is a fixed-length sequence of values that all have the same
type. Its type is written `[T; N]`: element type `T`, length `N`. Crucially, the
**length is part of the type**. A `[i32; 3]` and a `[i32; 4]` are different,
incompatible types.

```rust
fn main() {
    let nums: [i32; 3] = [10, 20, 30];
    let zeros = [0u8; 16];           // sixteen bytes, all 0  (like memset to 0)

    println!("{}", nums[0]);         // indexing with []
    println!("{}", nums.len());      // 3
    println!("{}", zeros.len());     // 16
}
```

Rust arrays are **stack-allocated** and **bounds-checked**. If you index past the
end, the program panics instead of reading stray memory:

```rust
use std::hint::black_box;

// This compiles, but PANICS at runtime: index out of bounds.
fn main() {
    let nums = [10, 20, 30];
    // `black_box` hides the value from the compiler so it cannot reject this
    // at compile time; in real code the index usually comes from input.
    let i = black_box(5);
    println!("{}", nums[i]); // thread panicked: index out of bounds: len is 3
}
```

> **Watch out.** If you index with a *constant* the compiler already knows is out
> of range (e.g. `nums[5]` written literally), Rust rejects it at **compile time**
> instead — an even nicer outcome. Indexes computed at runtime are the ones checked
> at runtime.

```
let nums = [10, 20, 30];   type [i32; 3]

  index:    0     1     2
          +-----+-----+-----+
  stack:  | 10  | 20  | 30  |
          +-----+-----+-----+
            i32   i32   i32      <- 3 * 4 = 12 bytes, on the stack

  nums[5]  ->  PANIC: index out of bounds (checked at runtime)
```

> **C vs Rust.** A C array decays to a pointer and carries no length, so
> out-of-bounds access is silent and dangerous. A Rust array knows its length
> (it is in the type), stays on the stack, and is bounds-checked at runtime.

> **Watch out.** Because the length is part of the type, a function that takes
> `[i32; 3]` will not accept `[i32; 4]`. For a view into a sequence of any length,
> you use a **slice** (`&[T]`), and for a growable, heap-allocated array you use
> `Vec<T>`. Both are covered in Chapter 10 — Slices and Strings.

## No implicit conversions

C quietly converts between numeric types all the time: assign an `int` to a
`long`, mix `int` and `double` in an expression, pass a `char` where an `int` is
expected. This is convenient and a frequent source of bugs.

Rust does **none** of this automatically. Different numeric types never mix:

```rust
// COMPILE ERROR: mismatched types
fn main() {
    let a: i32 = 5;
    let b: i64 = 10;
    let c = a + b;       // error[E0308]: expected `i32`, found `i64`
    println!("{c}");
}
```

To convert, you ask for it explicitly with the `as` keyword:

```rust
fn main() {
    let a: i32 = 5;
    let b: i64 = 10;
    let c = a as i64 + b;        // convert a to i64, then add
    println!("{c}");

    let big: i32 = 300;
    let small = big as u8;       // truncates: 300 mod 256 = 44
    println!("{small}");         // prints 44
}
```

The `as` cast is like a C cast: it can **truncate** or change the value. Casting
`300i32` to `u8` keeps only the low 8 bits, giving 44. `as` never panics and
never reports an error, so use it knowing it can silently change the value, just
like in C.

For safe, lossless conversions, idiomatic Rust prefers the `From` and `Into`
traits, such as `i64::from(a)`, which only compile when the conversion cannot
lose data. We cover those in Chapter 27 — Idioms and Style.

> **Rule of thumb.** Use `as` for deliberate, possibly-lossy casts (and for
> casting to and from raw pointers later). Prefer `From`/`Into` when you want the
> compiler to guarantee the conversion is lossless.

## Constants and statics

Rust has two ways to declare a fixed, named value: `const` and `static`. Both
replace C's `#define` constants and `const`/`static` globals, but they are not
the same as each other.

### `const`

A `const` is a compile-time constant. It must have a type, its value must be
computable at compile time, and the compiler **inlines** it wherever it is used —
there is no single storage location, much like a C `#define` but type-checked.

```rust
const MAX_USERS: u32 = 100_000;
const SECONDS_PER_DAY: u32 = 60 * 60 * 24;   // computed at compile time

fn main() {
    println!("{MAX_USERS} {SECONDS_PER_DAY}");
}
```

Names use `SCREAMING_SNAKE_CASE` by convention. A `const` can be declared in any
scope, including inside a function or a module, and it has no fixed address.

> **C vs Rust.** A Rust `const` is closest to a C `#define` of a literal or to an
> `enum` constant, but it is fully typed and scoped. It is *not* a variable: it
> has no address and never changes.

### `static`

A `static` is a value with a **single, fixed memory address** that lives for the
entire run of the program. This is like a C global variable. Because it has one
real location, you can take a reference to it, and it is not inlined.

```rust
static GREETING: &str = "hello";

fn main() {
    println!("{GREETING}");
}
```

An immutable `static` is safe. A **mutable** global is exactly the kind of
shared mutable state that causes data races in C, so Rust makes it hard on
purpose. The old `static mut` form is now strongly discouraged — in edition 2024
the compiler even refuses to let you take a reference to one. The safe, modern way
to have a mutable global counter is an **atomic**:

```rust
use std::sync::atomic::{AtomicU32, Ordering};

static COUNTER: AtomicU32 = AtomicU32::new(0);

fn main() {
    COUNTER.fetch_add(1, Ordering::Relaxed); // safe: no `unsafe`, no data race
    println!("{}", COUNTER.load(Ordering::Relaxed));
}
```

> **Rule of thumb.** Reach for `const` almost always. Use `static` only when you
> truly need a single fixed address (for example, a large lookup table). For a
> mutable global, use an `Atomic` or a `Mutex` (Chapter 20 — Channels and Shared
> State), not `static mut`.

| Feature | `const` | `static` |
|---|---|---|
| Has a fixed address | No (inlined) | Yes (one location) |
| Can be mutable | No | Only `static mut`, and that is `unsafe` |
| Lives for the whole program | Effectively yes (inlined) | Yes (`'static`) |
| Type required | Yes | Yes |
| Closest C idea | typed `#define` | global variable |

## C type to Rust type

| C type | Rust type | Note |
|---|---|---|
| `int8_t` / `signed char` | `i8` | exact width |
| `uint8_t` / `unsigned char` | `u8` | also Rust's "raw byte" |
| `int16_t` | `i16` | |
| `uint16_t` | `u16` | |
| `int` / `int32_t` | `i32` | Rust's default integer |
| `unsigned` / `uint32_t` | `u32` | |
| `long long` / `int64_t` | `i64` | |
| `uint64_t` | `u64` | |
| `size_t` | `usize` | indexing and lengths |
| `ptrdiff_t` / `intptr_t` | `isize` | pointer-sized signed |
| `float` | `f32` | IEEE 754 single |
| `double` | `f64` | IEEE 754 double, the default |
| `_Bool` / `bool` | `bool` | not a number in Rust |
| `char` (one byte) | `u8` | a byte, not text |
| a Unicode character | `char` | 4 bytes, Unicode scalar value |
| `void` (no value) | `()` | the unit type |
| `T arr[N]` | `[T; N]` | length is part of the type |
| `struct { ... }` (anonymous group) | `(A, B, C)` | tuple |

## Key takeaways

- `let` binds a name to a value. Bindings are **immutable by default**; add `mut`
  to allow change. This is the reverse of C.
- **Shadowing** (`let x = ...` again) makes a new variable with the same name and
  may change its type; **mutation** (`x = ...` with `mut`) changes one variable's
  value and keeps its type.
- Integers state their size and signedness: `i8`..`i128`, `u8`..`u128`, plus
  pointer-sized `isize`/`usize`. The default integer is `i32`; `usize` is used for
  indexing and lengths.
- Overflow **panics in debug builds** and wraps in release; use `wrapping_*`,
  `checked_*`, or `saturating_*` to choose behavior explicitly.
- `bool` is not a number, and `char` is a **4-byte Unicode scalar value**, not a
  C-style byte (use `u8` for a byte).
- Tuples group mixed types; arrays `[T; N]` are stack-allocated, bounds-checked,
  and the **length is part of the type**. `()` is the unit type (like `void`).
- There are **no implicit numeric conversions**: cast with `as` (may truncate) or
  use `From`/`Into` for lossless conversions.
- `const` is an inlined, typed compile-time constant; `static` has a single fixed
  address; `static mut` is `unsafe`.

## Watch out (gotchas for C programmers)

- **Immutable by default.** You will need `let mut` more often than you expect.
- **Shadowing is not mutation.** A second `let` makes a new variable; assignment
  without `let` needs `mut`.
- **`char` is 4 bytes**, not one. Use `u8` for a raw byte. Single quotes are a
  `char`; double quotes are a string.
- **Use `usize` for indexing and lengths.** Mixing index types causes compile
  errors; `as usize` is common when you have, say, an `i32` index.
- **Overflow panics in debug builds.** Do not rely on silent wraparound; choose
  `wrapping_*`/`checked_*`/`saturating_*` when you mean it.
- **No implicit conversions.** `i32 + i64` will not compile. Write `as` or use
  `From`/`Into`.
- **Array length is part of the type.** `[i32; 3]` and `[i32; 4]` are different
  types; use a slice `&[T]` or `Vec<T>` for variable length (Chapter 10).

## Interview questions

**Q: What is the difference between shadowing and mutation in Rust?**
A: Mutation (`x = ...`, requires `let mut`) changes the value stored in one
variable and cannot change its type. Shadowing (`let x = ...` again) creates a
brand-new variable that reuses the name; the old one is hidden, and the new one
may have a different type. Shadowing works even on immutable bindings.

**Q: How big is a Rust `char`, and how does it differ from a C `char`?**
A: A Rust `char` is 4 bytes and holds a Unicode scalar value, so it can represent
any single Unicode character. A C `char` is one byte and is really a small
integer. For a single raw byte in Rust you use `u8`, not `char`.

**Q: What happens on integer overflow in Rust?**
A: In a debug build, an overflowing arithmetic operation panics, so the bug is
caught. In a release build it wraps around using two's complement. To control the
behavior explicitly you use `wrapping_*` (wrap), `checked_*` (returns `Option`),
`saturating_*` (clamp to min/max), or `overflowing_*` (result plus a flag).

**Q: Why does Rust require `as` for numeric conversions when C does them
implicitly?**
A: Rust has no implicit numeric conversions, to prevent the silent value changes
and bugs that implicit conversions cause in C. You write `as` for a deliberate,
possibly-lossy cast, or use the `From`/`Into` traits for conversions the compiler
can prove are lossless.

**Q: When would you use `const` versus `static`?**
A: Use `const` for a typed compile-time constant; it has no fixed address and is
inlined where used. Use `static` when you need a single value at one fixed memory
address that lives for the whole program (like a C global). A mutable `static mut`
requires `unsafe`, so prefer safe shared-state types like `Mutex` or atomics.

## Try it

1. Write `let x = 5;` then `x = 6;` and read the error. Add `mut` to fix it.
2. Shadow a value to change its type: `let s = "hello"; let s = s.len();` and
   print it. Then try the same with `mut` and watch it fail.
3. Set `let n: u8 = 255;` and print `n + 1` in a debug build (`cargo run`) to see
   the panic, then replace it with `n.wrapping_add(1)` and `n.checked_add(1)`.
