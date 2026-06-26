# Chapter 28 — Gotchas for C Programmers (the checklist)

> **What you'll learn.** The traps that bite C programmers most when they first
> write Rust, gathered in one place. For each one: the trap, why it surprises a C
> programmer, and the fix. Use this as a checklist when the compiler fights you.

This chapter does not teach new ideas. It collects the surprises from the rest of
the book into a single list you can scan fast. Each gotcha points to the chapter
where it is taught in full. If you read only one chapter twice, make it this one.

The format is always the same:

- **The trap** — what you wrote.
- **Why it surprises you** — what C taught you to expect.
- **The fix** — the Rust way, with a tiny snippet.

Intentionally broken examples start with `// COMPILE ERROR:` so you know the
failure is on purpose.

## Variables and types

### Variables are immutable by default

**The trap.** You declare a variable and then try to change it.

```rust
// COMPILE ERROR: cannot assign twice to immutable variable `x`
fn main() {
    let x = 5;
    x = 6; // error[E0384]
}
```

**Why it surprises you.** In C, every variable is mutable unless you write
`const`. In Rust it is the other way around: `let` makes a value that cannot
change. This is the single most common first-day error.

**The fix.** Add `mut` ("mutable").

```rust
fn main() {
    let mut x = 5;
    x = 6; // fine now
    println!("{x}");
}
```

> **C vs Rust.** C: mutable by default, `const` to lock. Rust: immutable by
> default, `mut` to unlock. (Chapter 4 — Variables, Types, and Constants.)

### No implicit numeric conversions

**The trap.** You mix integer and float, or `i32` and `i64`, and expect C's
automatic promotion.

```rust
// COMPILE ERROR: mismatched types (expected i64, found i32)
fn main() {
    let a: i32 = 5;
    let b: i64 = 10;
    let c = a + b; // error[E0308]
}
```

**Why it surprises you.** C silently promotes and converts numbers in arithmetic.
This is convenient but hides bugs (truncation, sign changes). Rust never converts
number types for you.

**The fix.** Convert on purpose with `as` (a cast) or `From`/`into` (a safe
widening conversion).

```rust
fn main() {
    let a: i32 = 5;
    let b: i64 = 10;
    let c = a as i64 + b; // explicit cast
    let d = i64::from(a) + b; // safe widening, no data loss
    println!("{c} {d}");
}
```

> **Rule of thumb.** Use `From`/`into` when the conversion cannot lose data (small
> to big). Use `as` when you accept truncation. (Chapter 4.)

### Integer overflow panics in debug, wraps in release

**The trap.** You add two numbers and the result is too big for the type.

```rust
use std::hint::black_box;

fn main() {
    // `black_box` keeps these out of the compiler's constant folding; in real
    // code the values would come from input. With plain literals the compiler
    // rejects the overflow at compile time instead (even nicer).
    let x: u8 = black_box(250);
    let y: u8 = black_box(10);
    let z = x + y; // debug build: panics "attempt to add with overflow"
    println!("{z}"); // release build: wraps to 4
}
```

**Why it surprises you.** In C, unsigned overflow wraps silently and signed
overflow is undefined behavior. Many C programmers rely on wrapping. Rust treats
accidental overflow as a bug: it panics in debug builds so you find it.

**The fix.** Say what you mean with the explicit methods.

```rust
fn main() {
    let x: u8 = 250;
    let a = x.wrapping_add(10); // wraps on purpose -> 4
    let b = x.checked_add(10); // Option: None on overflow
    let c = x.saturating_add(10); // clamps to u8::MAX (255)
    let (d, overflowed) = x.overflowing_add(10); // value + a bool flag
    println!("{a} {b:?} {c} {d} {overflowed}");
}
```

> **Watch out.** Behavior differs between debug and release builds. Never depend
> on the default wrap. (Chapter 4; performance in Chapter 18 — Memory Without a GC.)

### `usize` is the type for indexing and sizes

**The trap.** You index a collection with an `i32`.

```rust
// COMPILE ERROR: the type `[i32]` cannot be indexed by `i32`
fn main() {
    let v = vec![10, 20, 30];
    let i: i32 = 1;
    let x = v[i]; // error[E0277]: indexing wants usize
    println!("{x}");
}
```

**Why it surprises you.** In C you index with any integer. Rust uses `usize` (an
unsigned, pointer-sized integer, like C's `size_t`) for all lengths and indices.

**The fix.** Use `usize`, or cast.

```rust
fn main() {
    let v = vec![10, 20, 30];
    let i: usize = 1;
    println!("{}", v[i]);
}
```

## Ownership, moves, and borrows

### Assignment and passing MOVE non-`Copy` values

**The trap.** You use a value after assigning it elsewhere or passing it to a
function.

```rust
// COMPILE ERROR: borrow of moved value: `s`
fn main() {
    let s = String::from("hi");
    let t = s; // ownership MOVES to t
    println!("{s}"); // error[E0382]: s is no longer valid
    println!("{t}");
}
```

**Why it surprises you.** In C, `t = s` copies the pointer; both names work (and
that is exactly how you get double-free bugs). In Rust, assigning a non-`Copy`
value *moves* ownership. The old name becomes invalid. Small `Copy` types (`i32`,
`bool`, `char`, fixed arrays of `Copy`) are duplicated instead, so they keep
working.

**The fix.** Borrow instead of moving, or clone if you truly need a second copy.

```rust
fn main() {
    let s = String::from("hi");
    let t = &s; // borrow: s still owns the data
    println!("{s} {t}"); // both work

    let u = s.clone(); // deep copy if you need an independent String
    println!("{s} {u}");
}
```

> **C vs Rust.** A move is like `memcpy` of the handle plus a promise that the old
> handle will never be touched again — enforced by the compiler. (Chapter 7 —
> Ownership and Moves.)

### Returning a reference to a local fails to compile

**The trap.** You return a pointer to a local variable.

```rust
// COMPILE ERROR: cannot return reference to local variable `s`
fn dangling() -> &String {
    let s = String::from("oops");
    &s // error[E0515]: `s` is dropped at the end of the function
}
```

**Why it surprises you.** This is the *opposite* of C. In C this compiles fine and
gives you a dangling pointer that crashes later. In Rust it is a compile error, so
the dangling pointer can never exist.

**The fix.** Return the owned value (move it out), not a reference to it.

```rust
fn not_dangling() -> String {
    let s = String::from("fine");
    s // move ownership to the caller
}

fn main() {
    println!("{}", not_dangling());
}
```

> **Mental model.** A reference is a borrowed key, not a copy. You cannot hand out
> a key to a room that is about to be demolished. (Chapters 8 and 9.)

### The borrow rules: shared XOR mutable

**The trap.** You hold a shared reference and a mutable reference at the same time.

```rust
// COMPILE ERROR: cannot borrow `v` as mutable because it is also borrowed as immutable
fn main() {
    let mut v = vec![1, 2, 3];
    let first = &v[0]; // shared borrow
    v.push(4); // mutable borrow while `first` is alive -> error[E0502]
    println!("{first}");
}
```

**Why it surprises you.** In C you can read through one pointer and write through
another at the same time. That is how iterator-invalidation bugs happen (`push`
might reallocate and free the buffer `first` points into). Rust forbids it.

**The rule.** At any moment you may have **either** any number of shared
references (`&T`) **or** exactly one mutable reference (`&mut T`), never both. This
is "aliasing XOR mutation."

**The fix.** End the shared borrow before mutating (use the value, then mutate),
or copy out what you need.

```rust
fn main() {
    let mut v = vec![1, 2, 3];
    let first = v[0]; // copy the i32 out; no live borrow
    v.push(4);
    println!("{first} {v:?}");
}
```

> **Watch out.** Only one `&mut` may exist at a time, and no `&mut` while any `&`
> is alive. This is what makes data races impossible. (Chapter 8 — Borrowing and
> References.)

### "Fighting the borrow checker": common patterns

When the borrow checker rejects code that you *know* is logically fine, reach for
one of these patterns, roughly in order of preference:

1. **Shorten the borrow.** Compute and finish reading before you mutate. Thanks to
   non-lexical lifetimes, a borrow ends at its last use, not at the closing brace.
2. **Copy or clone to break the tie.** Pull a small `Copy` value out, or `.clone()`
   an owned value to start. Optimize later if it matters.
3. **Split borrows.** Borrow two different fields of a struct mutably at once
   (allowed), or use `split_at_mut` to get two mutable halves of a slice.
4. **Use indices instead of references.** Store `usize` indices into a `Vec` rather
   than references between elements. Common for graphs and trees.
5. **Restructure.** Move the borrowing into a smaller function, or change who owns
   what.
6. **Reach for shared ownership.** `Rc<RefCell<T>>` (single thread) or
   `Arc<Mutex<T>>` (across threads) when you genuinely need many owners or shared
   mutable state.

```rust
fn main() {
    // Split borrows: two halves of a slice, both mutable, no conflict.
    let mut a = [1, 2, 3, 4];
    let (left, right) = a.split_at_mut(2);
    left[0] += 10;
    right[0] += 20;
    println!("{a:?}"); // [11, 2, 23, 4]
}
```

> **Rule of thumb.** When stuck, clone first to make it compile, then remove the
> clone once you see the data flow. Fighting the checker for an hour is worse than
> one `.clone()`. (Chapters 8, 9, and 17 — Smart Pointers.)

### Drop order is reverse of declaration

**The trap.** You assume destructors run in declaration order, as in C++ would be
the question; in C there are no destructors at all.

**Why it surprises you.** C has no automatic cleanup, so you may not expect any
order. Rust drops (cleans up) values automatically. Within a scope, values are
dropped in **reverse** order of declaration — last declared, first dropped — like
unwinding a stack.

```rust
struct Noisy(&'static str);

impl Drop for Noisy {
    fn drop(&mut self) {
        println!("drop {}", self.0);
    }
}

fn main() {
    let _a = Noisy("a");
    let _b = Noisy("b");
    // prints: drop b, then drop a
}
```

> **C vs Rust.** Rust's `Drop` is RAII: cleanup is automatic and deterministic,
> like calling `free` at the right moment for you. (Chapter 7.)

## Strings and characters

### `char` is 4 bytes, not 1

**The trap.** You treat `char` like a byte.

**Why it surprises you.** In C, `char` is one byte and doubles as "small integer"
and "ASCII letter." In Rust, `char` is a **Unicode scalar value**: it is 4 bytes
and holds any single character, like `'é'` or `'🦀'`. A single byte is `u8`.

```rust
fn main() {
    let c: char = '🦀';
    println!("{}", std::mem::size_of::<char>()); // 4
    let b: u8 = b'A'; // a byte literal is u8, value 65
    println!("{c} {b}");
}
```

### Strings are UTF-8; you cannot index by integer

**The trap.** You try `s[0]` to get the first character, like a C `char` array.

```rust
// COMPILE ERROR: `String` cannot be indexed by `{integer}`
fn main() {
    let s = String::from("héllo");
    let first = s[0]; // error[E0277]: no such indexing
    println!("{first}");
}
```

**Why it surprises you.** A Rust `String` (and `&str`) is UTF-8 bytes. One
character may take several bytes, so `s[0]` has no clear meaning and is forbidden.
There is also **no NUL terminator** — the length is stored, not found by scanning.

**The fix.** Iterate by character or by byte, or take a checked byte range.

```rust
fn main() {
    let s = String::from("héllo");
    let first: char = s.chars().next().unwrap(); // 'h'
    let nth: Option<char> = s.chars().nth(1); // 'é'
    let bytes: &[u8] = s.as_bytes(); // raw UTF-8 bytes
    println!("{first} {nth:?} {}", bytes.len());
}
```

### `.len()` is bytes, not characters

**The trap.** You expect `.len()` to count characters.

```rust
fn main() {
    let s = "héllo";
    println!("{}", s.len()); // 6, not 5 ('é' is two bytes in UTF-8)
    println!("{}", s.chars().count()); // 5 characters
}
```

**Why it surprises you.** C's `strlen` counts bytes up to the NUL, and for ASCII
that equals the character count. In UTF-8 they differ. `.len()` returns the number
of bytes; use `.chars().count()` for characters. (Chapter 10 — Slices and Strings.)

### `&str` vs `String` (like `const char *` vs an owned buffer)

**The trap.** You are not sure which string type to use in a function signature.

**Why it surprises you.** C has just `char *`. Rust splits the idea in two:

| Rust | C analogy | Owns the data? |
|---|---|---|
| `String` | a `malloc`'d, growable buffer | yes |
| `&str` | a `const char *` view (ptr + len) | no, it borrows |

**The fix.** Take `&str` in parameters (it accepts both a `String` and a literal),
return `String` when you build new text.

```rust
fn shout(s: &str) -> String {
    s.to_uppercase()
}

fn main() {
    let owned = String::from("hi");
    println!("{}", shout(&owned)); // &String coerces to &str
    println!("{}", shout("there")); // a literal is already &str
}
```

The same idea applies to data: take `&[T]` (a slice), return `Vec<T>`. A slice is
"pointer + length," exactly like the `(ptr, len)` pair you pass around in C.

## Collections and indexing

### Slices and `Vec` are bounds-checked

**The trap.** You index past the end.

```rust
fn main() {
    let v = vec![1, 2, 3];
    let x = v[10]; // panic: index out of bounds: len is 3 but index is 10
    println!("{x}");
}
```

**Why it surprises you.** In C, `v[10]` reads whatever memory is there — silent
corruption or a crash. Rust checks the bound and panics with a clear message
instead of reading out of bounds.

**The fix.** Use `.get()`, which returns an `Option`, when the index might be
invalid.

```rust
fn main() {
    let v = vec![1, 2, 3];
    match v.get(10) {
        Some(x) => println!("got {x}"),
        None => println!("out of range"),
    }
}
```

> **Deep dive.** The bounds check has a tiny runtime cost, but the optimizer
> removes it in many loops, and iterators avoid it entirely. (Chapter 16.)

## Null, errors, and return values

### There is no null — use `Option`

**The trap.** You look for a null pointer to mean "nothing."

**Why it surprises you.** In C, any pointer can be `NULL`, and forgetting to check
is a top cause of crashes. Rust has no null. "Maybe a value" is the type
`Option<T>`, which is `Some(value)` or `None`, and the compiler forces you to
handle `None`.

```rust
fn first_word(s: &str) -> Option<&str> {
    s.split_whitespace().next()
}

fn main() {
    match first_word("hello world") {
        Some(w) => println!("first: {w}"),
        None => println!("empty"),
    }
}
```

> **C vs Rust.** `Option<&T>` replaces a nullable pointer, and because of "niche
> optimization" it takes the same space as the pointer: `None` reuses the
> all-zeros bit pattern. You get the safety for free. (Chapter 12 — Enums and
> Pattern Matching.)

### `Result` is `#[must_use]`; ignoring it warns

**The trap.** You call a function that returns `Result` and ignore the result.

```rust
use std::fs::File;

fn main() {
    File::create("out.txt"); // warning: unused `Result` that must be used
}
```

**Why it surprises you.** In C, ignoring a return code is silent; nobody stops you
from skipping an error check. Rust marks `Result` as `#[must_use]`, so the
compiler warns when you drop it.

**The fix.** Handle it: `match`, `?`, `.unwrap()`/`.expect()` (panics on error),
or `.ok()` if you truly do not care.

```rust
use std::fs::File;

fn main() -> std::io::Result<()> {
    let _f = File::create("out.txt")?; // `?` returns the error to the caller
    Ok(())
}
```

> **Watch out.** `.unwrap()` and `.expect()` panic on `Err` or `None`. They are
> fine in examples and tests; in real code prefer `?` and proper handling.
> (Chapter 13 — Error Handling.)

### `Result`/`panic` are not C++ exceptions

**The trap.** You expect errors to be thrown and caught.

**Why it surprises you.** Rust has no exceptions. A recoverable error is a normal
return value (`Result`), passed up with `?`. A `panic!` is for bugs that should
never happen; it unwinds and usually ends the program. There is no `try`/`catch`.
(Chapter 13.)

## Expressions and control flow

### `if`, `loop`, and `match` are expressions

**The trap.** You write a C-style assignment from an `if` and get confused, or you
add a stray semicolon that throws away the value.

```rust
fn main() {
    // `if` is an expression: it produces a value.
    let n = 7;
    let label = if n % 2 == 0 { "even" } else { "odd" };
    println!("{label}");

    // `loop` can return a value with `break value`.
    let mut i = 0;
    let doubled = loop {
        i += 1;
        if i == 5 {
            break i * 2;
        }
    };
    println!("{doubled}");
}
```

**Why it surprises you.** In C, `if` and `switch` are statements; you assign inside
them. In Rust they are expressions that evaluate to a value, which replaces C's
ternary `?:` and many temporary variables.

### A trailing semicolon swallows the return value

**The trap.** You add a semicolon to the last line of a function.

```rust
// COMPILE ERROR: mismatched types (expected i32, found ())
fn square(x: i32) -> i32 {
    x * x; // the `;` turns this into a statement returning ()
}
```

**Why it surprises you.** In C, `return` is explicit, so a semicolon is harmless.
In Rust, the **last expression with no semicolon** is the return value. Adding `;`
turns the expression into a statement whose value is `()` (the empty "unit" type,
like C's `void`).

**The fix.** Drop the semicolon (or write `return x * x;`).

```rust
fn square(x: i32) -> i32 {
    x * x // no semicolon: this is the result
}

fn main() {
    println!("{}", square(9));
}
```

### `match` is exhaustive and does not fall through

**The trap.** You expect C `switch` behavior: cases fall through, and missing cases
are allowed.

```rust
// COMPILE ERROR: non-exhaustive patterns: `Blue` not covered
enum Color { Red, Green, Blue }

fn name(c: Color) -> &'static str {
    match c {
        Color::Red => "red",
        Color::Green => "green",
        // forgot Blue -> error[E0004]
    }
}
```

**Why it surprises you.** C's `switch` falls through to the next case unless you
write `break`, and you can omit cases. Rust `match` arms never fall through, and
the compiler forces you to cover every case (or add a `_` catch-all).

**The fix.** Cover all cases, or use `_` for the rest.

```rust
enum Color { Red, Green, Blue }

fn name(c: Color) -> &'static str {
    match c {
        Color::Red => "red",
        Color::Green => "green",
        Color::Blue => "blue",
    }
}

fn main() {
    println!("{}", name(Color::Blue));
}
```

> **Rule of thumb.** Exhaustiveness is a feature: add a variant to an enum and the
> compiler lists every `match` you must update. (Chapter 12.)

## Macros and functions

### `println!` and friends are macros (the `!`)

**The trap.** You write `println(...)` without the `!`.

**Why it surprises you.** They look like function calls but the `!` marks a
**macro** — code that runs at compile time. The macro checks your format string
against the arguments, which a normal function cannot do. Common ones: `println!`,
`format!`, `vec!`, `panic!`, `assert!`, `write!`. (Chapter 26 — Macros.)

```rust
fn main() {
    let x = 42;
    println!("x = {x}"); // captures `x` from scope; checked at compile time
}
```

### No function overloading, no default arguments

**The trap.** You define two functions with the same name and different parameters.

**Why it surprises you.** C does not have overloading either, but you may expect it
from C++ or other languages. Rust has no overloading and no default argument
values.

**The fix.** Use different names, an `Option` parameter, a builder, or a trait
(generics) to accept several types.

```rust
fn connect(host: &str, port: u16) { /* ... */ }
fn connect_default(host: &str) {
    connect(host, 8080); // wrapper supplies the "default"
}

fn main() {
    connect_default("localhost");
}
```

(Chapter 6 — Functions and Closures; builders in Chapter 27 — Idioms and Style.)

### No inheritance — use traits and composition

**The trap.** You look for class inheritance to share behavior.

**Why it surprises you.** Rust has no classes and no inheritance. Shared behavior
comes from **traits** (like interfaces, or a C struct of function pointers), and
data reuse comes from **composition** (put one struct inside another). Traits can
provide default method bodies. (Chapter 15 — Traits.)

## Structs and layout

### Field order is not guaranteed — use `#[repr(C)]` for FFI

**The trap.** You assume struct fields are laid out in memory in declaration order,
then pass the struct to C.

**Why it surprises you.** In C, struct fields are in declaration order (with
padding). Rust may **reorder** fields to save space. That is fine for pure Rust,
but it breaks any code that expects the C layout — including FFI and casting to
bytes.

**The fix.** Add `#[repr(C)]` to force C layout.

```rust
#[repr(C)]
struct Point {
    x: i32,
    y: i32,
}
```

> **Watch out.** Only use `#[repr(C)]` when layout matters (FFI, memory-mapped
> structures). Let Rust optimize layout otherwise. (Chapter 25 — Unsafe and FFI.)

### Deref coercion can surprise you

**The trap.** You call a method and it seems to come from a different type.

**Why it surprises you.** Rust automatically converts `&String` to `&str`,
`&Vec<T>` to `&[T]`, and `&Box<T>` to `&T` when needed. This "deref coercion" is
convenient but can be confusing: a method you call on a `String` may actually be a
`str` method. It also means you usually do not need to write `&*` or `.as_str()`.
(Chapter 17 — Smart Pointers.)

## Smart pointers and shared state

### `Rc` is not thread-safe — use `Arc` across threads

**The trap.** You share an `Rc` between threads.

```rust
// COMPILE ERROR: `Rc<i32>` cannot be sent between threads safely
use std::rc::Rc;
use std::thread;

fn main() {
    let data = Rc::new(5);
    let d = Rc::clone(&data);
    thread::spawn(move || {
        println!("{d}"); // error[E0277]: Rc is not Send
    });
}
```

**Why it surprises you.** `Rc` (reference counting) is fast because its count is
*not* atomic. That makes it unsafe to share across threads. The compiler catches
this through the `Send` trait.

**The fix.** Use `Arc` (atomic reference count) for sharing across threads.

```rust
use std::sync::Arc;
use std::thread;

fn main() {
    let data = Arc::new(5);
    let d = Arc::clone(&data);
    let h = thread::spawn(move || println!("{d}"));
    h.join().unwrap();
}
```

(Chapter 17; threads in Chapter 19 — Threads and Concurrency.)

### `Rc`/`Arc` cycles leak — use `Weak`

**The trap.** You build a cycle of `Rc`s (for example, a parent that points to a
child that points back to the parent).

**Why it surprises you.** Reference counting cannot free a cycle: each node keeps
the other alive, so the count never reaches zero. The memory leaks. This is one of
the few ways safe Rust can leak.

**The fix.** Make one direction a `Weak` reference, which does not keep the value
alive.

```rust
use std::rc::{Rc, Weak};
use std::cell::RefCell;

struct Node {
    parent: RefCell<Weak<Node>>, // weak: does not own the parent
    children: RefCell<Vec<Rc<Node>>>, // strong: owns the children
}
```

(Chapter 17.)

### `RefCell` moves borrow checking to runtime — and panics

**The trap.** You borrow a `RefCell` mutably twice at once.

```rust
use std::cell::RefCell;

fn main() {
    let c = RefCell::new(5);
    let a = c.borrow_mut();
    let b = c.borrow_mut(); // panic: already mutably borrowed: BorrowMutError
    println!("{a} {b}");
}
```

**Why it surprises you.** `RefCell` lets you mutate through a shared reference
("interior mutability"). It still enforces the borrow rules, but at **runtime**, not
compile time. A violation is a panic, not a compile error.

**The fix.** Keep borrows short, and do not hold two `borrow_mut()` at once. Drop
the first borrow before taking the second. (Chapter 17.)

## Concurrency

### Data races are caught at compile time

**The trap.** You move data into a thread, or share it without synchronization, and
expect it to "just work" like C.

**Why it surprises you.** In C, sharing data between threads compiles fine; data
races are your problem to avoid, and they cause the worst bugs. In Rust the
compiler refuses to compile a data race, using two marker traits: `Send` (safe to
move to another thread) and `Sync` (safe to share by reference between threads).

**The fix.** Use a `move` closure to give the thread ownership, and wrap shared
mutable state in `Arc<Mutex<T>>`.

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for _ in 0..4 {
        let c = Arc::clone(&counter);
        handles.push(thread::spawn(move || {
            let mut n = c.lock().unwrap(); // lock; auto-unlocks when `n` drops
            *n += 1;
        }));
    }
    for h in handles {
        h.join().unwrap();
    }
    println!("{}", *counter.lock().unwrap()); // 4
}
```

> **C vs Rust.** A `Mutex` in Rust *wraps* the data it protects, so you cannot read
> the data without locking. In C the mutex and the data are separate, and nothing
> stops you from forgetting the lock. (Chapters 19 and 20.)

## Key takeaways

- **Immutable by default.** Add `mut` to change a variable.
- **Moves, not copies.** Assigning or passing a non-`Copy` value moves it;
  use-after-move is a compile error. Borrow with `&`/`&mut`, or `.clone()`.
- **Aliasing XOR mutation.** Many `&` or one `&mut`, never both. This is the borrow
  checker's core rule.
- **No dangling references.** Returning `&local` is a compile error — the opposite
  of the classic C bug.
- **No null, no implicit conversions.** Use `Option`; convert numbers with `as` or
  `From`.
- **Strings are UTF-8.** `char` is 4 bytes, `s[i]` does not compile, `.len()` is
  bytes, there is no NUL terminator.
- **Overflow panics in debug.** Use `wrapping_`/`checked_`/`saturating_`; index
  with `usize`.
- **Exhaustive `match`, no fall-through. Macros end in `!`. `Result` is
  `#[must_use]`. The last expression (no semicolon) is the return value.**
- **`Rc` is single-thread; use `Arc` across threads.** Cycles leak — use `Weak`.
  `RefCell` panics on a borrow violation at runtime.
- **Data races are compile errors.** Use `move` closures and `Arc<Mutex<T>>`.
- Slices/`Vec` are bounds-checked; use `.get()` for an `Option`. Drop order is
  reverse of declaration. Prefer `&str`/`&[T]` in parameters, return
  `String`/`Vec`.

## Interview questions

**Q: A C programmer writes `let s2 = s1;` for two `String`s and then uses `s1`.
Why does it not compile?**
A: Assigning a non-`Copy` value moves ownership. After the move, `s1` is invalid,
so using it is a compile error (E0382). This guarantees a single owner and so the
buffer is freed exactly once. Borrow (`&s1`) or `.clone()` to keep `s1` usable.

**Q: Why does returning a reference to a local variable fail in Rust but compile in
C?**
A: The local is dropped when the function returns, so the reference would dangle.
Rust's borrow checker proves the reference cannot outlive the value and rejects it
(E0515). C compiles it and gives you undefined behavior at run time. Return the
owned value instead.

**Q: What is the difference between `s.len()` and `s.chars().count()` on a Rust
string?**
A: `.len()` is the number of UTF-8 *bytes*; `.chars().count()` is the number of
characters (Unicode scalar values). They are equal for ASCII but differ when any
character takes more than one byte, such as `'é'`.

**Q: When does integer overflow panic, and how do you opt into wrapping?**
A: In debug builds, overflow panics ("attempt to add with overflow"). In release
builds it wraps silently. To choose behavior explicitly use `wrapping_add`,
`checked_add` (returns `Option`), `saturating_add`, or `overflowing_add`.

**Q: Why can you not share an `Rc<T>` between threads, and what do you use
instead?**
A: `Rc`'s reference count is not atomic, so concurrent updates would race; `Rc` is
not `Send`, and the compiler rejects sending it to another thread. Use `Arc<T>`,
which uses an atomic count and is safe to share. For shared mutation, wrap the data:
`Arc<Mutex<T>>`.

**Q: How does Rust prevent data races at compile time?**
A: Through the `Send` and `Sync` marker traits. `Send` means a type is safe to move
to another thread; `Sync` means `&T` is safe to share between threads. The compiler
only allows thread APIs on types that satisfy these, so code with a data race does
not compile. Shared mutable state must go through a `Mutex`/`RwLock` (often inside
an `Arc`).
