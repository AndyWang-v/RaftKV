# Appendix A — C to Rust Cheat Sheet

> **What you'll learn.** A dense, scannable lookup table from C to Rust: types,
> declarations, memory, strings, pointers, control flow, errors, concurrency, the
> standard library headers, common idioms, and a `cargo` command reference.

This is a desk reference, not a tutorial. Each row maps something you do in C to the
idiomatic Rust way to do it. Chapters are named so you can read the full story.

## Types

| C | Rust | Note |
|---|---|---|
| `char` | `i8` / `u8` | C `char` signedness is platform-defined; Rust is explicit. |
| `signed char` / `unsigned char` | `i8` / `u8` | A raw byte is `u8`. |
| `short` / `unsigned short` | `i16` / `u16` | Fixed width in Rust. |
| `int` / `unsigned int` | `i32` / `u32` | `int` is usually 32-bit; Rust pins it. |
| `long` | `i32` or `i64` | C `long` varies by platform; pick the exact Rust type. |
| `long long` / `unsigned long long` | `i64` / `u64` | Fixed width. |
| `__int128` | `i128` / `u128` | Built in, no extension needed. |
| `size_t` | `usize` | Index/length type; pointer-sized. |
| `ptrdiff_t`, `intptr_t` | `isize` | Pointer-sized signed integer. |
| `float` / `double` | `f32` / `f64` | IEEE-754, same as C. |
| `_Bool` / `bool` | `bool` | `true` / `false`; size 1 byte. |
| `void` (no value) | `()` (unit) | The empty tuple; size 0. |
| `void` (return: never) | `!` (never) | For functions that never return. |
| `enum { ... }` | `enum` | Rust enums are tagged unions, far more powerful (Ch. 12). |
| `struct { ... }` | `struct` | Same idea; see declarations below (Ch. 11). |
| `union { ... }` | `enum` (safe) or `union` (`unsafe`) | Prefer `enum`; raw `union` needs `unsafe`. |
| `char *` (text) | `&str` / `String` | UTF-8; not null-terminated (Ch. 10). |
| `T *` | `&T` / `&mut T` / `Box<T>` / `*const T` / `*mut T` | Pick by ownership and mutability. |
| `T[N]` | `[T; N]` | Fixed-size array, length known at compile time. |
| `T *` + length | `&[T]` / `Vec<T>` | Slice (borrowed) or owned growable array. |

## Declarations

| Thing | C | Rust |
|---|---|---|
| Variable | `int x = 5;` | `let x: i32 = 5;` (type optional) |
| Mutable variable | `int x = 5;` | `let mut x: i32 = 5;` |
| Constant | `#define N 10` / `const int N = 10;` | `const N: i32 = 10;` |
| Global mutable | `static int c = 0;` | `static mut C: i32 = 0;` (needs `unsafe`; prefer atomics) |
| Immutable global | `static const int K = 1;` | `static K: i32 = 1;` |
| Array | `int a[3] = {1,2,3};` | `let a: [i32; 3] = [1, 2, 3];` |
| Repeated array | `int a[100] = {0};` | `let a = [0; 100];` |
| Pointer | `int *p = &x;` | `let p: &i32 = &x;` (or `&mut`, `*const`, `*mut`) |
| Struct | `struct P { int x, y; };` | `struct P { x: i32, y: i32 }` |
| Tuple struct | (none) | `struct P(i32, i32);` |
| Enum | `enum C { Red, Green };` | `enum C { Red, Green }` |
| Typedef | `typedef uint32_t Id;` | `type Id = u32;` (alias) or `struct Id(u32);` (newtype) |
| Function | `int add(int a, int b)` | `fn add(a: i32, b: i32) -> i32` |
| Function pointer | `int (*f)(int)` | `fn(i32) -> i32` or `Box<dyn Fn(i32) -> i32>` |
| Extern function | `extern int f(void);` | `extern "C" { fn f() -> i32; }` (Ch. 25) |

> **C vs Rust.** Variables are immutable by default in Rust. Where C has mutable
> everything, Rust makes you write `mut`. Types come *after* the name, not before.

## Memory management

| C | Rust | Note |
|---|---|---|
| `malloc(sizeof(T))` | `Box::new(value)` | One heap value, freed when the `Box` drops (Ch. 17). |
| `free(p)` | (automatic) | `Drop` runs at end of scope; no manual `free` (Ch. 7). |
| `calloc(n, sizeof(T))` | `vec![T::default(); n]` / `vec![0; n]` | Zero-initialized buffer. |
| `malloc(n * sizeof(T))` (array) | `Vec::with_capacity(n)` then `push`, or `vec![x; n]` | Growable array (Ch. 16). |
| `realloc(p, n)` (grow) | `vec.resize(n, value)` / `vec.push(x)` | `Vec` reallocates for you. |
| `realloc` (shrink) | `vec.truncate(n)` / `vec.shrink_to_fit()` | |
| `sizeof(T)` | `std::mem::size_of::<T>()` | Compile-time size in bytes. |
| `sizeof(expr)` | `std::mem::size_of_val(&expr)` | Size of a value. |
| `_Alignof(T)` | `std::mem::align_of::<T>()` | Alignment in bytes. |
| `memcpy(dst, src, n)` | `dst.copy_from_slice(src)` / `slice::copy_from_slice` | Safe, checked. |
| `memset(p, 0, n)` | `slice.fill(0)` | |
| `memmove` | `slice.copy_within(range, dest)` | Overlap-safe. |
| Move semantics | (manual) | `let b = a;` *moves* ownership; `a` is invalid after (Ch. 7). |
| Reference count | (manual) | `Rc<T>` (single thread) / `Arc<T>` (threads) (Ch. 17, 20). |

> **Mental model.** A `Box<T>` is `malloc` + a guaranteed `free` at scope end. You
> never call `free`; the compiler inserts it exactly once.

## Strings

| C | Rust | Note |
|---|---|---|
| `char *` literal | `&str` (`"hi"`) | UTF-8, length known, not null-terminated. |
| `char buf[N]` (owned text) | `String` | Heap-allocated, growable (Ch. 10). |
| `strlen(s)` | `s.len()` | Byte length; `s.chars().count()` for characters. |
| `strcmp(a, b) == 0` | `a == b` | `==` compares contents. |
| `strcmp` ordering | `a.cmp(b)` | Returns `Ordering::{Less,Equal,Greater}`. |
| `strcpy(d, s)` / `strdup(s)` | `s.to_owned()` / `s.to_string()` / `String::from(s)` | Owned copy. |
| copy a `String` | `s.clone()` | Deep copy. |
| `strcat(d, s)` | `d.push_str(s)` | In place. |
| `a + b` (concat) | `a + &b` (a is `String`, b is `&str`) | Or `format!("{a}{b}")`. |
| `strchr` / `strstr` | `s.find(c)` / `s.find("sub")` | Returns `Option<usize>`. |
| `strtok` | `s.split(sep)` | Returns an iterator. |
| `printf("%d\n", x)` | `println!("{x}")` | Type-checked format string (Ch. 5). |
| `fprintf(stderr, ...)` | `eprintln!(...)` | |
| `sprintf(buf, ...)` | `let s = format!(...);` | Returns a `String`; no buffer overflow. |
| `snprintf` | `format!(...)` | Always bounded; allocates as needed. |
| `puts(s)` | `println!("{s}")` | |
| `atoi` / `strtol` | `s.parse::<i32>()?` | Returns `Result`, no silent failure. |
| C string for FFI | `std::ffi::CString` / `CStr` | Null-terminated bridge to C (Ch. 25). |

> **Watch out.** A Rust `&str`/`String` is **not** null-terminated and is UTF-8.
> Indexing by byte (`s[0]`) is not allowed; use `s.as_bytes()[0]` or `s.chars()`.

## Pointers and references

| C | Rust | Note |
|---|---|---|
| `const T *p` (read) | `&T` | Shared reference; many allowed at once. |
| `T *p` (write) | `&mut T` | Exclusive reference; only one at a time (Ch. 8). |
| Owning pointer | `Box<T>` | Single owner on the heap. |
| `T *` (raw, unchecked) | `*const T` / `*mut T` | Only dereferenced in `unsafe` (Ch. 25). |
| `NULL` | `Option::None` | A nullable `&T` becomes `Option<&T>` (Ch. 12). |
| `p == NULL` check | `match opt { Some(x) => .., None => .. }` | Compiler forces the check. |
| `p->field` | `p.field` | Rust auto-dereferences; no `->`. |
| `(*p).field` | `(*p).field` or `p.field` | `.` works through references. |
| `&x` | `&x` | Take a reference. |
| `*p` | `*p` | Dereference. |
| `p + 1` (pointer math) | `slice[i]` / iterators / `ptr.add(1)` (`unsafe`) | No arithmetic on safe refs. |

> **C vs Rust.** A `&T`/`&mut T` can never be null and never dangle — the borrow
> checker proves the referent outlives the reference (Ch. 9). For "maybe absent,"
> use `Option<&T>`, which is the same size as a pointer thanks to niche
> optimization.

## Control flow

| C | Rust |
|---|---|
| `if (c) { } else { }` | `if c { } else { }` (no parentheses; braces required) |
| ternary `c ? a : b` | `if c { a } else { b }` (an expression) |
| `switch (x) { case ...: }` | `match x { pattern => expr, _ => expr }` (Ch. 12) |
| `while (c) { }` | `while c { }` |
| `do { } while (c);` | `loop { ...; if !c { break; } }` |
| `for (i=0; i<n; i++)` | `for i in 0..n { }` |
| `for` over array | `for x in &arr { }` / `for x in arr.iter() { }` |
| infinite loop | `loop { }` |
| `break` / `continue` | `break` / `continue` (and labeled `break 'outer`) |
| `goto` | (none) | use loops, functions, `match`, or `?`. |
| fallthrough `case` | `match x { 1 \| 2 => ... }` | Combine patterns with `\|`. |

> **Watch out.** `match` is **exhaustive**: you must cover every case or add `_`.
> There is no implicit fallthrough between arms.

## Error handling

| C | Rust | Note |
|---|---|---|
| return code `int` | `Result<T, E>` | Errors are values you cannot ignore (Ch. 13). |
| `errno` | the `E` in `Result<T, E>` | Carry rich error data, not a global int. |
| check-and-return on error | `let x = f()?;` | `?` returns the error early. |
| `NULL` for "no value" | `Option<T>` | `Some`/`None` (Ch. 12). |
| `assert(cond)` | `assert!(cond)` / `debug_assert!(cond)` | |
| `abort()` | `std::process::abort()` | |
| `exit(code)` | `std::process::exit(code)` | |
| fatal error | `panic!("msg")` | Unwinds the stack, then aborts (Ch. 13). |
| `setjmp`/`longjmp` (exceptions) | `Result` + `?`; or `panic!`/`catch_unwind` | Rust has no exceptions. |
| `goto cleanup;` (RAII) | `Drop` runs automatically | No manual cleanup label needed. |

```rust
use std::num::ParseIntError;

fn double(s: &str) -> Result<i32, ParseIntError> {
    let n = s.parse::<i32>()?; // `?` returns Err early, like checking a return code
    Ok(n * 2)
}

fn main() {
    println!("{:?}", double("21")); // Ok(42)
    println!("{:?}", double("x"));  // Err(ParseIntError { .. })
}
```

## Concurrency

| C (POSIX) | Rust | Note |
|---|---|---|
| `pthread_create` | `std::thread::spawn(\|\| { .. })` | Returns a `JoinHandle` (Ch. 19). |
| `pthread_join` | `handle.join().unwrap()` | |
| `pthread_mutex_t` + lock/unlock | `std::sync::Mutex<T>` + `.lock()` | Lock guards the data; auto-unlocks on drop. |
| `pthread_rwlock_t` | `std::sync::RwLock<T>` | |
| `pthread_cond_t` | `std::sync::Condvar` | |
| shared owned data | `Arc<T>` | Atomic reference count (Ch. 20). |
| shared mutable data | `Arc<Mutex<T>>` | The standard pattern. |
| `_Atomic int` | `std::sync::atomic::AtomicI32` | With explicit `Ordering`. |
| `atomic_fetch_add` | `a.fetch_add(1, Ordering::SeqCst)` | |
| thread-local `__thread` | `thread_local!` macro | |
| message passing | `std::sync::mpsc::channel()` | `tx.send` / `rx.recv` (Ch. 20). |

> **C vs Rust.** The compiler checks thread safety with the `Send` and `Sync`
> marker traits. Sharing non-synchronized mutable data across threads simply does
> not compile — data races are caught before the program runs.

## libc headers → Rust

| C header | Rust equivalent | Note |
|---|---|---|
| `<stdio.h>` | `std::io`, `print!`/`println!`, `std::fs` | I/O and files. |
| `<stdlib.h>` | `std` (alloc via `Box`/`Vec`), `std::process`, `std::env` | `exit`, `getenv`, `rand` (use the `rand` crate). |
| `<string.h>` | `&str`/`String`/`[T]` methods | `len`, `find`, `split`, `copy_from_slice`. |
| `<stdint.h>` | built-in `i8`..`i128`, `u8`..`u128`, `usize`/`isize` | No header needed. |
| `<stdbool.h>` | `bool` | Built in. |
| `<stddef.h>` | `usize`, `()`, `core::ptr` | |
| `<math.h>` | `f64`/`f32` methods, `std::f64::consts` | `x.sqrt()`, `x.sin()`, `PI`. |
| `<float.h>` / `<limits.h>` | `f64::MAX`, `i32::MAX`, `u8::MIN`, etc. | Associated constants. |
| `<time.h>` | `std::time::{Instant, Duration, SystemTime}` | Monotonic and wall clock. |
| `<assert.h>` | `assert!`, `debug_assert!` | |
| `<errno.h>` | `std::io::Error` / `Result` | No global errno. |
| `<pthread.h>` | `std::thread`, `std::sync` | Threads, mutexes, channels. |
| `<ctype.h>` | `char` methods | `c.is_alphabetic()`, `c.to_ascii_uppercase()`. |
| `<stdarg.h>` (varargs) | macros / slices / generics | Use `&[T]` or a macro instead. |

```rust
use std::time::Instant;

fn main() {
    let start = Instant::now();
    let r = (1.0_f64).sin() + std::f64::consts::PI.sqrt(); // math.h methods
    println!("{r:.4} in {:?}", start.elapsed());           // time.h
}
```

## Idioms (the Rust way)

| Idea | C habit | Rust idiom | Chapter |
|---|---|---|---|
| Cleanup | `goto cleanup; free(...)` | RAII: `Drop` frees on scope exit | 7, 17 |
| Ownership | "who frees this?" comments | One owner, enforced by compiler | 7 |
| Distinct types | `typedef uint32_t UserId;` | Newtype: `struct UserId(u32);` | 27 |
| Conversions | implicit casts | `From`/`Into`: `let s: String = x.into();` | 15, 27 |
| Errors | return codes / `errno` | `Result<T, E>` + `?` | 13 |
| "No value" | `NULL` / sentinel | `Option<T>` | 12 |
| Polymorphism | structs of function pointers | traits + `dyn` | 15 |
| Generic code | macros / `void *` | generics (monomorphized) | 14 |
| Invalid states | runtime checks | "make illegal states unrepresentable" with enums | 27 |
| Builders | many init functions | builder pattern returning `Self` | 27 |

## `cargo` command cheat sheet

| Command | What it does |
|---|---|
| `cargo new my_app` | Create a new binary project (`--lib` for a library). |
| `cargo init` | Turn the current directory into a Cargo project. |
| `cargo build` | Compile (debug). `--release` for an optimized build. |
| `cargo run` | Build and run the binary. Args after `--`: `cargo run -- foo`. |
| `cargo check` | Type-check without producing a binary (fast). |
| `cargo test` | Run unit, integration, and doc tests (Ch. 23). |
| `cargo bench` | Run benchmarks. |
| `cargo add serde` | Add a dependency to `Cargo.toml` (`--features ...`). |
| `cargo remove serde` | Remove a dependency. |
| `cargo update` | Update dependencies within semver, rewriting `Cargo.lock`. |
| `cargo fmt` | Format code with `rustfmt`. |
| `cargo clippy` | Lint with Clippy; `--fix` to auto-apply (Ch. 24). |
| `cargo doc --open` | Build and open HTML docs. |
| `cargo clean` | Delete the `target/` build directory. |
| `cargo tree` | Show the dependency graph. |
| `cargo install ripgrep` | Install a binary crate globally. |

## Common one-liners

```rust
// Read all of stdin into a String
let mut input = String::new();
std::io::Read::read_to_string(&mut std::io::stdin(), &mut input).unwrap();

// Parse and sum whitespace-separated integers
let total: i64 = "1 2 3".split_whitespace().map(|s| s.parse::<i64>().unwrap()).sum();

// Command-line arguments (skip(1) drops the program name)
let args: Vec<String> = std::env::args().skip(1).collect();

// Read a file to a String (returns Result)
let text = std::fs::read_to_string("data.txt")?;

// Sort and remove duplicates
let mut v = vec![3, 1, 2, 1]; v.sort(); v.dedup();

// Filter, transform, collect
let evens: Vec<i32> = (0..10).filter(|n| n % 2 == 0).map(|n| n * n).collect();
```

```sh
rustc hello.rs && ./hello      # compile and run a single file (no Cargo)
cargo run --release            # optimized build + run
RUST_BACKTRACE=1 cargo run     # show a backtrace on panic
```
