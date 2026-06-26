# Chapter 29 — Interview Questions and Answers

This chapter is a bank of real Rust interview questions with short, correct
answers, grouped by topic, plus code puzzles and design questions at the end.

> **How to use this.** Read the question, try to answer out loud before you read
> the answer, and write the code snippets by hand. The questions double as a
> review of the whole book; each group maps to chapters you have already read. The
> "Predict / explain" puzzles are the kind interviewers love, because they reveal
> whether you really understand ownership, not just the syntax.

## Basics

**Q: What is the difference between `let` and `let mut`?**
A: `let` binds an immutable value; you cannot reassign it. `let mut` binds a
mutable value you can change. Rust is immutable by default — the opposite of C,
where everything is mutable unless you write `const`.

**Q: What is shadowing, and how does it differ from mutation?**
A: Shadowing means declaring a new variable with the same name, which hides the
old one. It can change the type, and the old value still existed independently.
Mutation changes the same variable in place and keeps its type.

```rust
fn main() {
    let x = "42"; // &str
    let x: i32 = x.parse().unwrap(); // new variable, new type
    println!("{x}");
}
```

**Q: Why is immutability the default?**
A: It makes code easier to reason about, helps the optimizer, and supports the
borrow rules (you can share an immutable value freely between many readers). You
opt into mutation with `mut`, which documents intent.

**Q: List Rust's scalar types and their C analogues.**
A: Integers `i8`..`i128`/`u8`..`u128` plus `isize`/`usize` (like `size_t`); floats
`f32`/`f64`; `bool`; and `char` (a 4-byte Unicode scalar, not C's 1-byte `char`).
Sizes are fixed by the language, unlike C where `int` width is
platform-dependent.

**Q: How big is a `char`, and why?**
A: Four bytes. A Rust `char` is a Unicode scalar value, so it can hold any single
character (`'é'`, `'🦀'`), not just ASCII. A single byte is `u8`. Strings store
UTF-8 bytes, so a `char` and the bytes used to encode it in a string are not the
same size.

**Q: What does `as` do, and when should you avoid it?**
A: `as` is an explicit primitive cast (for example `x as u8`). It can truncate or
change sign silently. Prefer `From`/`into` for lossless conversions; use `as` only
when you accept the possible data loss, or `try_into` to get a checked `Result`.

**Q: Why is `if`/`match`/`loop` being an "expression" useful?**
A: They evaluate to a value, so you can write `let y = if c { a } else { b };` with
no temporary mutable variable, and `break value` out of a `loop`. The last
expression in a block (without a semicolon) is its value. A stray semicolon turns
it into a statement of type `()` and discards the value.

## Ownership and borrowing

**Q: What is ownership?**
A: Every value has exactly one owner (a variable). When the owner goes out of
scope, the value is dropped (freed) automatically. This replaces manual
`malloc`/`free` and guarantees each value is freed exactly once.

**Q: What is the difference between a move and a copy?**
A: For types that are not `Copy` (like `String`, `Vec`), assignment or passing
*moves* ownership and invalidates the source. For `Copy` types (small, fixed-size:
integers, `bool`, `char`, references, arrays of `Copy`), the value is duplicated
bit-for-bit and the source stays valid.

**Q: State the borrow rules and explain why they exist.**
A: At any time you may have either any number of shared references (`&T`) or
exactly one mutable reference (`&mut T`), never both, and no reference may outlive
the value. "Aliasing XOR mutation" prevents data races and iterator-invalidation
bugs, and it lets the compiler reason about and optimize the code.

**Q: Why can there be only one `&mut` at a time?**
A: If two mutable references existed, two pieces of code could write the same
memory at once (a data race) or invalidate each other's assumptions. Exclusivity
makes mutation safe and is what makes `&mut T` analogous to a `restrict` pointer in
C — but enforced.

**Q: When do you need `.clone()`?**
A: When you need an independent, owned copy of a non-`Copy` value (for example to
keep using the original after giving one away). Cloning is explicit because it may
be expensive (a deep copy / allocation). It is a reasonable first move to satisfy
the borrow checker, then optimize away if it matters.

**Q: What is `Drop`, and how does it relate to RAII and C?**
A: `Drop` is the trait whose `drop` method runs automatically when a value goes out
of scope, like a C++ destructor. It implements RAII: resources (memory, files,
locks) are released deterministically. In C you would call `free`/`fclose`/`unlock`
by hand; in Rust the compiler inserts the cleanup, in reverse order of declaration.

**Q: Can you call `drop` yourself?**
A: Not the `Drop::drop` method directly. You call `std::mem::drop(value)`, which
takes ownership and lets the value go out of scope early. This is how you release a
lock or free memory before the end of the scope.

## Lifetimes

**Q: What is a lifetime?**
A: A lifetime is the span of code during which a reference is valid. It is a
compile-time concept with no runtime cost. Annotations like `'a` describe
relationships between the lifetimes of inputs and outputs so the compiler can prove
no reference dangles.

**Q: What is lifetime elision?**
A: A set of rules that let you omit lifetime annotations in common cases. For
example, a method that takes `&self` and returns a reference gets the lifetime of
`self` automatically. Most function signatures need no explicit lifetimes thanks to
elision.

**Q: What does `'static` mean?**
A: The reference is valid for the entire program. String literals are `&'static
str` because they live in the binary. As a bound (`T: 'static`) it means the type
contains no shorter-lived references. It does *not* mean "leaked forever" — an
owned `String` also satisfies `T: 'static`.

**Q: When does a struct need a lifetime parameter?**
A: When it stores a reference. The lifetime ties the struct's validity to the
borrowed data, so the struct cannot outlive what it points to.

```rust
struct Parser<'a> {
    input: &'a str, // the struct must not outlive `input`
}
```

**Q: Why do lifetimes exist if they have no runtime cost?**
A: They give the borrow checker the information it needs to reject dangling
references at compile time. They are documentation that the compiler also verifies.

## Enums and pattern matching

**Q: How does `Option<T>` replace null?**
A: `Option<T>` is an enum with `Some(T)` or `None`. "No value" is encoded in the
type, and the compiler forces you to handle `None`, so you cannot dereference a
null by accident. There is no null pointer in safe Rust.

**Q: What is `Result<T, E>`?**
A: An enum with `Ok(T)` for success and `Err(E)` for failure. It is how fallible
functions report errors — Rust has no exceptions. It is `#[must_use]`, so ignoring
it warns.

**Q: What does it mean that `match` is exhaustive?**
A: You must cover every possible case, or the code does not compile. If you add a
variant to an enum, the compiler points you at every `match` you must update.
Arms do not fall through, unlike C's `switch`.

**Q: What is niche optimization?**
A: The compiler reuses an unused bit pattern ("niche") to encode an enum's tag for
free. The classic example: `Option<&T>` is the same size as `&T` because a
reference is never null, so the all-zeros pattern represents `None`. You get
null-safety with no extra space.

**Q: How do you handle just one variant without a full `match`?**
A: Use `if let` (or `let ... else`) to match one pattern and ignore the rest.

```rust
fn main() {
    let maybe = Some(7);
    if let Some(n) = maybe {
        println!("got {n}");
    }
}
```

## Error handling

**Q: When do you return a `Result` versus calling `panic!`?**
A: Return `Result` for expected, recoverable errors the caller should handle (file
not found, bad input). Use `panic!` for bugs and broken invariants that should
never happen (a failed assertion, an index that must be valid). A library should
almost always return `Result`, not panic.

**Q: What does the `?` operator do?**
A: It unwraps an `Ok`/`Some` value, or returns the `Err`/`None` early from the
current function. It is concise error propagation, replacing C's repetitive
"check the return code and return it" pattern.

```rust
use std::num::ParseIntError;

fn double(s: &str) -> Result<i32, ParseIntError> {
    let n: i32 = s.parse()?; // returns the error on failure
    Ok(n * 2)
}
```

**Q: How does `?` convert error types?**
A: It calls `From::from` on the error, so any error type that implements
`From<SourceError>` is converted to the function's declared error type. This lets
one function return one error type while calling code that fails in different ways.

**Q: What are `thiserror` and `anyhow` for?**
A: `thiserror` is a derive macro to build custom error enums (typed errors for
libraries). `anyhow` provides a single boxed error type (`anyhow::Error`) plus
context, convenient for applications where you do not need to match on the error.

**Q: What does `#[must_use]` mean on `Result`?**
A: The compiler warns if you ignore the returned value. It nudges you to handle the
error instead of silently dropping it, unlike C where ignoring a return code is
invisible.

**Q: What does `.unwrap()` do, and when is it acceptable?**
A: It returns the inner value or panics on `Err`/`None`. It is fine in examples,
prototypes, tests, and cases where the value is provably present. In production
paths prefer `?` or explicit handling, and use `.expect("why")` when you do unwrap
so the panic message explains the assumption.

## Traits and generics

**Q: What is the difference between static and dynamic dispatch?**
A: Static dispatch uses generics: the compiler generates a specialized copy of the
function per concrete type (monomorphization), so calls are direct and inlinable —
zero cost. Dynamic dispatch uses a trait object (`dyn Trait`): the concrete type is
chosen at run time through a vtable, like a C struct of function pointers. It adds
an indirection but allows mixing types at run time.

**Q: What is a trait object, and how is it laid out?**
A: `dyn Trait` is a value of some type that implements `Trait`, accessed behind a
pointer. The pointer is "fat": it carries two words — a data pointer and a pointer
to a vtable (the method table). It is the Rust version of a hand-rolled C interface
made of function pointers.

```
&dyn Draw  =  [ data ptr | vtable ptr ]
                            |
                            v
                 [ drop | size | align | draw fn | ... ]
```

**Q: What is monomorphization?**
A: For each concrete type used with a generic, the compiler stamps out a dedicated
copy of the code. This gives C++-template-like performance (no runtime dispatch) at
the cost of larger binaries and longer compile times.

**Q: What is the orphan rule?**
A: You may implement a trait for a type only if you own the trait or the type. This
prevents two crates from giving conflicting implementations for the same pair. To
work around it for a foreign trait on a foreign type, wrap the type in a newtype you
own.

**Q: What is object safety?**
A: A trait is "object safe" (can be made into `dyn Trait`) only if its methods can
be dispatched through a vtable — for example, no generic methods and no methods
that return `Self` by value. If a trait is not object safe, you cannot use it as a
trait object; use generics instead.

**Q: Name common standard traits and what they give you.**
A: `Debug` (the `{:?}` format), `Clone`/`Copy` (duplication), `PartialEq`/`Eq`
(equality), `PartialOrd`/`Ord` (ordering), `Default` (a default value),
`Hash` (use as a `HashMap` key), `From`/`Into` (conversions), `Iterator`,
`Display` (user-facing `{}`), and `Drop` (cleanup). Many can be `#[derive]`d.

**Q: What is the difference between `impl Trait` and `dyn Trait` in a return type?**
A: `-> impl Trait` returns one concrete (but unnamed) type chosen at compile time;
it is static dispatch and zero cost, but every path must return the *same* type.
`-> Box<dyn Trait>` returns a heap-allocated trait object via dynamic dispatch and
can return different concrete types from different branches.

## Smart pointers

**Q: When do you use `Box<T>`?**
A: To put a value on the heap: for large values you do not want on the stack, for
recursive types (a linked list or tree node that contains itself), and to store a
trait object (`Box<dyn Trait>`). It is a single owner, like a `unique_ptr`.

**Q: What is the difference between `Rc` and `Arc`?**
A: Both give shared ownership through reference counting; the value is freed when
the last owner drops. `Rc` uses a non-atomic counter (fast, single-thread only).
`Arc` uses an atomic counter (safe to share across threads, slightly slower). The
compiler stops you from sending `Rc` across threads.

**Q: What is interior mutability, and what is `RefCell`?**
A: Interior mutability lets you mutate data through a shared reference (`&T`).
`RefCell<T>` provides it by enforcing the borrow rules at *run time* instead of
compile time: `borrow()` and `borrow_mut()` track borrows, and a violation panics.
It is for single-threaded code; the thread-safe equivalents are `Mutex`/`RwLock`.

**Q: How do reference-counting cycles cause leaks, and how do you fix them?**
A: If two `Rc`s point at each other, neither count reaches zero, so neither is
freed — a leak. Break the cycle with `Weak<T>`, a non-owning reference you upgrade
to `Rc` only when needed. Typically children own via `Rc` and point back to parents
via `Weak`.

**Q: What is `Rc<RefCell<T>>` for?**
A: Multiple owners (`Rc`) of a value that any of them can mutate (`RefCell`). It is
the common single-threaded pattern for shared mutable state, such as nodes in a
graph. The thread-safe version is `Arc<Mutex<T>>`.

## Concurrency

**Q: What are `Send` and `Sync`?**
A: Marker traits the compiler uses for thread safety. `Send` means a value can be
moved to another thread. `Sync` means `&T` can be shared between threads (i.e. `T`
is `Send + Sync`-friendly to reference). Most types are both automatically; types
like `Rc` and `RefCell` are not, which is how unsafe sharing is rejected.

**Q: What is "fearless concurrency"?**
A: The idea that the same ownership and `Send`/`Sync` rules that ensure
single-threaded safety also rule out data races at compile time. You can write
threaded code aggressively because the compiler refuses to build a race.

**Q: How do you share mutable state between threads?**
A: Wrap it in `Arc<Mutex<T>>` (or `Arc<RwLock<T>>`). `Arc` gives shared ownership
across threads; `Mutex` guarantees one writer at a time. The data lives *inside*
the mutex, so you cannot access it without locking.

**Q: What are channels?**
A: Message-passing pipes (`std::sync::mpsc`): one or more senders push values, a
receiver pulls them. They let threads communicate by transferring ownership of data
rather than sharing memory — "do not communicate by sharing memory; share memory by
communicating."

**Q: How exactly does Rust prevent a data race?**
A: A data race needs two threads, at least one writing, the same memory, with no
synchronization. Rust breaks this: `&mut T` is exclusive (no other access while it
exists), and shared types that allow mutation must be `Sync` (a `Mutex` makes access
exclusive). Types that are not thread-safe are not `Send`/`Sync`, so they cannot
cross threads. The check is at compile time.

**Q: What is the difference between async and threads?**
A: Threads are OS-scheduled and good for CPU-bound or blocking work; each has its
own stack. Async tasks are lightweight, cooperatively scheduled by a runtime (like
Tokio), and excel at many concurrent I/O operations on few threads. Async avoids
the per-thread memory cost when you have thousands of waiting connections.

**Q: What does "futures are lazy" mean?**
A: An `async` function returns a `Future` that does nothing until it is `.await`ed
or handed to an executor. Creating a future does not start the work, unlike, say,
spawning a thread. No runtime polls it, nothing happens.

## Memory and runtime

**Q: What is the difference between the stack and the heap in Rust?**
A: The same as in C. The stack holds fixed-size values with automatic, LIFO
lifetime (function locals). The heap holds dynamically sized or long-lived data,
reached through owners like `Box`, `Vec`, and `String`. Rust frees heap data
automatically when its owner drops.

**Q: Rust has no garbage collector — how is memory freed?**
A: Through ownership and `Drop`. When an owner goes out of scope, its `drop` runs
and frees any owned resources, deterministically and at a known point. There are no
GC pauses; cleanup cost is predictable, exactly as with manual `free` in C but
automatic and correct.

**Q: Can safe Rust leak memory?**
A: Yes. Leaking memory is safe (it does not cause undefined behavior), so it is
allowed. You can leak with reference-count cycles (`Rc`), or on purpose with
`Box::leak` or `std::mem::forget`. Leaks are a logic bug, not a safety bug.

**Q: What is a "zero-cost abstraction"?**
A: A high-level feature that compiles to code as efficient as the hand-written
low-level version. Iterators, generics, `Option`, and `async`/await are examples:
you pay no runtime penalty for the convenience. "What you do not use, you do not pay
for; what you do use, you could not hand-code better."

**Q: Where do `Vec<T>` and `String` keep their data?**
A: The header (pointer, length, capacity) lives on the stack or inside its owner;
the elements live on the heap. This is exactly a `malloc`'d buffer plus a length and
capacity, managed for you.

## Unsafe and FFI

**Q: What does `unsafe` actually allow?**
A: Five extra powers: dereference a raw pointer, call an `unsafe` function (incl.
FFI), access or modify a mutable static, implement an `unsafe` trait, and access
union fields. That is all. It does not turn off the type system.

**Q: Does `unsafe` disable the borrow checker?**
A: No. The borrow checker and type checker still run inside `unsafe` blocks.
`unsafe` only unlocks the five operations above. You are promising the compiler you
upheld the invariants it cannot check (for example, that a raw pointer is valid).

**Q: How do you call a C function from Rust?**
A: Declare it in an `extern "C"` block with its signature and link the library.
Calls are `unsafe` because Rust cannot verify the C side.

```rust
unsafe extern "C" {
    fn abs(input: i32) -> i32; // from the C standard library
}

fn main() {
    let x = unsafe { abs(-5) };
    println!("{x}");
}
```

**Q: Why is `#[repr(C)]` needed for FFI?**
A: Rust may reorder struct fields for efficiency, so the default layout is not
guaranteed. `#[repr(C)]` forces C-compatible layout (declaration order with C
padding rules) so the struct matches what the C code expects.

**Q: What is the goal of an `unsafe` block in good Rust code?**
A: To isolate the few operations the compiler cannot verify, wrap them in a safe API
that upholds the invariants, and keep the unsafe region small and well-documented.
The rest of the program stays safe and relies on that wrapper.

## Cargo and tooling

**Q: What is the difference between a crate and a package?**
A: A crate is a single compilation unit — one library or one binary. A package is a
bundle described by one `Cargo.toml`; it can contain at most one library crate and
any number of binary crates.

**Q: What is `Cargo.lock`, and should you commit it?**
A: It records the exact dependency versions resolved for a build, so builds are
reproducible. Commit it for applications (binaries). For libraries it is
conventionally not committed, so downstream users resolve their own versions.

**Q: What is `clippy`?**
A: The official linter (`cargo clippy`). It catches common mistakes and suggests
more idiomatic code beyond what the compiler warns about — for example, "use
`if let` instead of `match`," or flagging a needless `.clone()`.

**Q: What are Rust editions?**
A: Named language epochs (2015, 2018, 2021, 2024) that let Rust make
otherwise-breaking syntax changes without splitting the ecosystem. Each crate picks
its edition in `Cargo.toml`; crates of different editions interoperate. This book
targets edition 2024.

**Q: How do you add a dependency?**
A: `cargo add <crate>` edits `Cargo.toml` for you (or you add it by hand under
`[dependencies]`). `cargo build` then fetches and compiles it. Versions follow
semantic versioning.

**Q: What does `cargo check` do that `cargo build` does not?**
A: `cargo check` type-checks and borrow-checks without producing a binary, so it is
much faster. It is what you run repeatedly while editing; `cargo build` when you
actually need the executable.

## Predict / explain (small puzzles)

For each puzzle, decide what happens and why before reading the answer.

**1. Use after move.**

```rust
// COMPILE ERROR: borrow of moved value: `s`
fn main() {
    let s = String::from("hi");
    let t = s; // move
    println!("{s} {t}");
}
```

Answer: Does not compile (E0382). `String` is not `Copy`, so `let t = s` moves
ownership; `s` is invalid afterward. Fix: borrow (`let t = &s;`) or clone
(`let t = s.clone();`).

**2. Borrow conflict.**

```rust
// COMPILE ERROR: cannot borrow `v` as mutable because it is also borrowed as immutable
fn main() {
    let mut v = vec![1, 2, 3];
    let first = &v[0]; // shared borrow, still alive below
    v.push(4); // needs a mutable borrow -> error[E0502]
    println!("{first}");
}
```

Answer: Does not compile. `push` may reallocate the buffer and invalidate `first`,
so the borrow checker rejects mutating `v` while `first` lives. Fix: copy the value
out (`let first = v[0];`) before pushing, or finish using `first` first.

**3. Option handling.**

```rust
fn main() {
    let v = vec![10, 20, 30];
    match v.get(5) {
        Some(x) => println!("got {x}"),
        None => println!("nothing"),
    }
}
```

Answer: Compiles and prints `nothing`. `.get()` returns `Option<&T>`, and index 5
is out of range, so it is `None`. Compare with `v[5]`, which would panic. This is
how you index safely.

**4. Iterator laziness.**

```rust
fn main() {
    let nums = vec![1, 2, 3];
    nums.iter().map(|n| println!("{n}")); // no output!
    println!("done");
}
```

Answer: Prints only `done`. `map` is lazy: it builds an iterator but does nothing
until consumed. Nothing drives it, so the closure never runs (the compiler also
warns about an unused iterator). Fix: use `for n in &nums { println!("{n}"); }`, or
add `.for_each(...)`, or `.collect()`.

**5. Integer overflow in debug.**

```rust
fn main() {
    let x: u8 = std::hint::black_box(200);
    let y: u8 = std::hint::black_box(100);
    println!("{}", x + y); // ?
}
```

Answer: In a debug build this panics ("attempt to add with overflow") because
300 > 255. In a release build it wraps to 44. Do not rely on either: use
`x.wrapping_add(y)` to wrap on purpose or `x.checked_add(y)` to get an `Option`.
(The `black_box` calls stop the compiler from constant-folding the values; if you
write `200u8 + 100u8` as plain literals, the compiler rejects it at compile time.)

**6. Rc cycle.**

```rust
use std::rc::Rc;
use std::cell::RefCell;

struct Node {
    next: RefCell<Option<Rc<Node>>>,
}

fn main() {
    let a = Rc::new(Node { next: RefCell::new(None) });
    let b = Rc::new(Node { next: RefCell::new(None) });
    *a.next.borrow_mut() = Some(Rc::clone(&b));
    *b.next.borrow_mut() = Some(Rc::clone(&a)); // cycle!
    // a and b keep each other alive; memory leaks at end of main.
}
```

Answer: Compiles and runs, but leaks. `a` and `b` reference each other, so neither
strong count reaches zero and neither `Node` is dropped. This is a safe leak (no
undefined behavior). Fix: make one link a `Weak<Node>` to break ownership of the
cycle.

## Design questions

These open-ended questions check whether you reach for idiomatic Rust shapes.

**Q: Model a traffic-light state machine.**
A: Use an enum for the states and a method that returns the next state. Enums make
illegal states unrepresentable, and exhaustive `match` guarantees every transition
is handled.

```rust
#[derive(Debug, Clone, Copy, PartialEq)]
enum Light {
    Red,
    Green,
    Yellow,
}

impl Light {
    fn next(self) -> Light {
        match self {
            Light::Red => Light::Green,
            Light::Green => Light::Yellow,
            Light::Yellow => Light::Red,
        }
    }
}

fn main() {
    let mut s = Light::Red;
    for _ in 0..4 {
        println!("{s:?}");
        s = s.next();
    }
}
```

**Q: Build a thread-safe counter shared by many threads.**
A: Wrap the count in `Arc<Mutex<T>>`. `Arc` gives each thread shared ownership;
`Mutex` makes increments mutually exclusive. The lock auto-releases when the guard
drops.

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let counter = Arc::new(Mutex::new(0u64));
    let mut handles = vec![];
    for _ in 0..8 {
        let c = Arc::clone(&counter);
        handles.push(thread::spawn(move || {
            for _ in 0..1000 {
                *c.lock().unwrap() += 1;
            }
        }));
    }
    for h in handles {
        h.join().unwrap();
    }
    println!("{}", *counter.lock().unwrap()); // 8000
}
```

For a counter specifically, you could also use `AtomicU64` with no lock at all,
which is faster for a single number.

**Q: Sketch the design of an LRU (least-recently-used) cache.**
A: You need O(1) lookup and O(1) "move to most-recently-used." The classic design
combines a hash map (key to node) with a doubly linked list ordered by recency:

- On `get`, find the node via the map and move it to the front of the list.
- On `put`, insert at the front; if over capacity, evict the node at the back and
  remove it from the map.

In Rust, a doubly linked list with shared mutable nodes fights the borrow checker.
Three practical options:

1. Use `std::collections::HashMap` plus a `VecDeque` or an index-based list, storing
   `usize` indices into a `Vec<Node>` instead of references — the most common safe
   approach.
2. Use `Rc<RefCell<Node>>` for the nodes and `Weak` for back-pointers to avoid
   cycles.
3. Use the `lru` crate, which has already solved this carefully (sometimes with a
   small, well-audited `unsafe` core).

The interview point is to recognize that intrusive doubly linked lists are the one
data structure that genuinely fights Rust's aliasing rules, and to name the
index-based or `Rc<RefCell>` workarounds rather than reaching for `unsafe` first.
