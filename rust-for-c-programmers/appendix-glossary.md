# Appendix B — Glossary

> **What you'll learn.** Every important Rust term used in this book, defined in one
> to three plain sentences, with a C comparison where it helps.

Terms are listed alphabetically. Bold words inside a definition are themselves
defined elsewhere in this glossary.

**adaptor (iterator adaptor)** — A method that takes an iterator and returns a new, lazy iterator, such as `map`, `filter`, or `take`. Nothing runs until a consumer (like `collect` or `sum`) pulls values through the chain.

**Arc** — "Atomically Reference Counted" pointer (`Arc<T>`): shared ownership of a heap value that is safe to send between threads. Like `Rc` but with thread-safe (atomic) counting; use it instead of `Rc` across threads.

**associated type** — A type slot declared inside a trait, written `type Item;`, that each implementer fills in. For example, `Iterator` has `type Item`, so each iterator names the type it yields without an extra generic parameter.

**blanket impl** — A trait implementation written for *all* types that satisfy some bound, such as `impl<T: Display> ToString for T`. One `impl` block thereby covers many types at once.

**borrow** — To temporarily access a value through a **reference** without taking ownership. A borrow ends when the reference goes out of scope; the owner keeps the value. Like passing a pointer in C, but the compiler proves it stays valid.

**borrow checker** — The part of the Rust compiler that enforces the borrowing rules: you may have many shared `&T` borrows or exactly one `&mut T` borrow, never both, and no borrow may outlive its referent. It prevents dangling pointers and data races at compile time.

**Box** — `Box<T>` is the simplest smart pointer: a single owned value on the heap, freed automatically when the `Box` drops. It is Rust's `malloc` + guaranteed `free` for one value.

**Cargo** — Rust's build tool and package manager. It compiles, tests, formats, lints, fetches dependencies, and builds docs — replacing Makefiles plus a package manager.

**Cargo.lock** — A generated file that records the exact versions of every dependency used in a build, so builds are reproducible. Commit it for binaries; libraries usually do not.

**Cargo.toml** — The manifest file that describes a **package**: its name, version, edition, and dependencies. The analog of a Makefile plus dependency list.

**Cell** — `Cell<T>` provides **interior mutability** for `Copy` types by moving values in and out, with no references handed out. It allows mutation through a shared `&` without runtime borrow tracking.

**Clone** — The trait for explicit, possibly expensive deep copies, via `.clone()`. Unlike **Copy**, cloning is never implicit; you ask for it.

**closure** — An anonymous function that can capture variables from its surrounding scope, written `|args| body`. Like a C function pointer plus a captured environment, but type-safe and without manual context structs.

**Copy** — A marker trait for small types (like integers) that are duplicated bit-for-bit on assignment instead of **moved**. If a type is `Copy`, `let b = a;` leaves `a` still usable.

**crate** — The unit of compilation in Rust: either a binary (has `main`) or a library. Roughly one tree of modules compiled together; a `package` contains one or more crates.

**deref coercion** — The compiler's automatic conversion from `&T` to `&U` when `T` implements `Deref<Target = U>`, e.g. `&String` to `&str`, or `&Box<T>` to `&T`. It lets methods and arguments work across these pointer-like types without manual conversion.

**Drop** — The trait whose `drop` method runs automatically when a value goes out of scope, releasing resources. It is Rust's destructor and the basis of **RAII** — the compiler-inserted equivalent of `free`/`close`.

**dyn** — The keyword marking a **trait object**, e.g. `&dyn Trait` or `Box<dyn Trait>`. It selects dynamic dispatch through a **vtable** instead of static dispatch.

**edition** — A named, opt-in language revision (2015, 2018, 2021, 2024) set in `Cargo.toml`. Editions can change syntax and idioms without breaking older crates, which keep their own edition. This book targets edition 2024.

**elision (lifetime elision)** — A set of compiler rules that infer common **lifetime** annotations so you do not have to write them. For example, a method taking `&self` usually returns a reference tied to `self` automatically.

**enum** — A type that is exactly one of several named **variants**, each able to carry data. It is a tagged union (a C `union` plus a tag) but type-safe, with the tag checked by `match`.

**exhaustiveness** — The rule that a `match` must handle every possible value, or include a `_` wildcard. The compiler rejects an incomplete match, so adding an enum variant forces you to update all matches.

**FFI** — Foreign Function Interface: calling C from Rust (or Rust from C) across the language boundary. It requires `extern "C"`, `unsafe`, and often `repr(C)` types.

**Fn / FnMut / FnOnce** — The three traits a **closure** can implement, by how it uses captures: `Fn` reads them (`&`), `FnMut` mutates them (`&mut`), and `FnOnce` consumes them (by value, callable once). Function arguments accept whichever they need.

**generic** — A type or function parameterized over types, written with `<T>`. Rust generics are real and type-checked, unlike C macros or `void *`, and are compiled by **monomorphization**.

**interior mutability** — The ability to mutate data through a shared `&` reference, provided by types like **Cell**, **RefCell**, **Mutex**, and atomics. It moves the aliasing-vs-mutation check from compile time to runtime (or to atomic hardware).

**iterator** — A value that produces a sequence on demand via its `next` method (the `Iterator` trait). Iterators are lazy and compile to tight loops — a **zero-cost abstraction**.

**lifetime** — A named region of code, written `'a`, over which a **reference** is valid. The compiler uses lifetimes to prove no reference outlives its data; they are erased at compile time and have no runtime cost.

**'static** — The longest **lifetime**: data that lives for the entire program, such as string literals (`&'static str`) and `const`/`static` items. As a bound (`T: 'static`), it means the type holds no shorter-lived references.

**macro** — Code that generates code at compile time, invoked with a `!` such as `println!` or `vec!`. Far more powerful and hygienic than C's text-substitution `#define`.

**macro_rules** — The built-in system for *declarative* macros, which match on syntax patterns and expand to code. The everyday way to write your own macros.

**match** — Rust's pattern-matching control structure, like a `switch` on steroids: it matches structure and binds values, and it must be **exhaustive**. It is also an expression that yields a value.

**module** — A namespace inside a crate, declared with `mod`, that groups items and controls visibility. The replacement for C's header-and-file organization, without separate `.h` files.

**monomorphization** — The compiler's strategy of generating a separate specialized copy of a **generic** function or type for each concrete type used. This makes generics zero-cost (no runtime dispatch), at the price of larger code.

**move** — Transfer of ownership from one binding to another, e.g. `let b = a;` for a non-**Copy** type, after which `a` is invalid. It guarantees a single owner, so each value is freed exactly once.

**&mut** — An exclusive, mutable **reference** (`&mut T`). While it exists, no other reference to the same value may exist, which is how Rust prevents data races and aliasing bugs.

**Mutex** — `Mutex<T>` wraps data and grants access only while locked, via `.lock()`. The lock guard unlocks automatically on drop, so you cannot forget to unlock; combine with `Arc` to share across threads.

**newtype** — A single-field tuple struct that wraps another type to give it a distinct identity, e.g. `struct Meters(f64);`. It prevents mixing up values, adds methods, and sidesteps the **orphan rule**.

**niche optimization** — A layout trick where the compiler stores an enum's tag inside an unused bit pattern of its fields, so `Option<&T>` is the same size as `&T` (using null for `None`). It makes many `Option` and enum types free of extra space.

**Option** — The enum `Option<T>` with variants `Some(T)` and `None`, representing a value that may be absent. It replaces null pointers and sentinel values; the compiler forces you to handle `None`.

**orphan rule** — The coherence rule that you may implement a trait for a type only if you define the trait or the type in your crate. It prevents conflicting implementations across crates; the **newtype** pattern works around it.

**ownership** — The core rule that every value has exactly one owner, and the value is dropped (freed) when the owner goes out of scope. It replaces manual `malloc`/`free` and makes use-after-free and double-free impossible.

**package** — A bundle managed by **Cargo**, defined by one `Cargo.toml`, containing one or more **crates** (often one library and/or one binary).

**panic** — An unrecoverable error that stops the current thread, by default **unwinding** the stack and running destructors. Triggered by `panic!`, failed `assert!`, out-of-bounds indexing, or `unwrap` on `None`/`Err`.

**pattern** — The structure on the left of a `match` arm or in a `let`, which both tests a value's shape and binds parts of it (e.g. `Some(x)`, `Point { x, .. }`). Destructuring and matching in one step.

**procedural macro** — A macro written as Rust code that transforms a token stream into new tokens, used for custom `#[derive]`, attributes, and function-like macros. More flexible than **macro_rules** but lives in its own crate.

**pub** — The visibility keyword that exports an item from its **module** (e.g. `pub fn`, `pub struct`). Without it, items are private to their module; the analog of choosing what goes in a C header.

**RAII** — "Resource Acquisition Is Initialization": tying a resource's lifetime to a value's scope, so cleanup happens automatically via **Drop**. The same idiom C++ uses; in Rust it is built into ownership.

**raw pointer** — `*const T` or `*mut T`: an unchecked C-style pointer with no borrow or lifetime guarantees. Dereferencing one requires **unsafe**; used for **FFI** and low-level code.

**Rc** — "Reference Counted" pointer (`Rc<T>`): shared ownership of a heap value within a single thread, freed when the last `Rc` drops. Not thread-safe — use **Arc** across threads.

**RefCell** — `RefCell<T>` provides **interior mutability** by checking the borrow rules at *runtime* instead of compile time, handing out `Ref`/`RefMut` guards. Violating the rules panics rather than failing to compile.

**reference** — A non-owning pointer to a value, `&T` (shared) or **&mut** `T` (exclusive), guaranteed by the compiler to be non-null and never dangling. Like a C pointer, but checked.

**repr(C)** — An attribute (`#[repr(C)]`) that lays out a struct or enum using C's field order and alignment rules. Required for types shared across **FFI**.

**Result** — The enum `Result<T, E>` with variants `Ok(T)` and `Err(E)`, representing success or failure. It replaces return codes and `errno`; the **? operator** propagates the `Err` case.

**rustfmt** — The official code formatter, run via `cargo fmt`, that enforces a standard style. It removes formatting debates and keeps diffs clean.

**rustup** — The official toolchain installer and manager. It installs `rustc`, **Cargo**, and components, and switches between stable, beta, nightly, and cross-compilation targets.

**Send** — A marker trait meaning a value can be safely *moved* to another thread. Most types are `Send`; **Rc** and raw pointers are not. The compiler uses it to prevent unsafe thread transfers.

**shadowing** — Declaring a new variable with the same name as an earlier one using `let`, which hides the old binding (and may change its type). Different from mutation: it creates a fresh variable.

**slice** — A view into a contiguous sequence, `&[T]` (or `&str` for text), carrying a pointer and a length. It is a C pointer-plus-length pair, but bounds-checked and lifetime-tracked.

**String** — A heap-allocated, growable, owned UTF-8 string. The owned counterpart to **&str**; like a managed `char *` buffer that always knows its length and capacity.

**&str** — A borrowed string slice: a view into UTF-8 text as a pointer and length. Comparable to a `const char *` plus length, but not null-terminated and guaranteed valid UTF-8.

**supertrait** — A trait that another trait requires, written `trait B: A`, so every type implementing `B` must also implement `A`. It lets `B`'s methods rely on `A`'s.

**Sync** — A marker trait meaning a value is safe to *share* between threads by reference (`&T` is **Send**). Types guarded by **Mutex** or atomics are `Sync`; **RefCell** and **Cell** are not.

**trait** — A named set of methods a type can implement, defining shared behavior. Like a C++ interface or a C struct of function pointers, but checked at compile time; the basis of generics' bounds.

**trait object** — A value accessed through `dyn Trait` (behind a pointer like `&dyn Trait` or `Box<dyn Trait>`), enabling runtime polymorphism via a **vtable**. The Rust equivalent of a pointer to a struct of function pointers.

**turbofish** — The `::<Type>` syntax used to specify a generic type argument explicitly, as in `"42".parse::<i32>()` or `vec.iter().collect::<Vec<_>>()`. Use it when inference cannot decide the type.

**unsafe** — A keyword marking code where you, not the compiler, uphold the safety rules: dereferencing **raw pointers**, calling `extern`/`unsafe` functions, or accessing `static mut`. It is "ordinary C," scoped and searchable.

**unwinding** — The default response to a **panic**: walking back up the stack, running each value's **Drop**, until the thread stops. It can be caught with `catch_unwind` or replaced by `abort` via a build setting.

**usize** — The pointer-sized unsigned integer used for indexing and lengths (Rust's `size_t`). `isize` is its signed counterpart (`ptrdiff_t`).

**variant** — One of the named alternatives of an **enum**, possibly carrying data (e.g. `Some(T)` and `None` are the variants of `Option`). A `match` selects on which variant a value is.

**Vec** — `Vec<T>`, a heap-allocated, growable array storing a pointer, length, and capacity. The Rust equivalent of a `malloc`'d array tracked with its length, freed automatically.

**vtable** — A table of function pointers (plus size and alignment) that a **trait object** uses to find the right method at runtime. Equivalent to a hand-built C struct of function pointers used for dynamic dispatch.

**where clause** — A clause that lists trait bounds for generics in a readable block after the signature, e.g. `where T: Clone + Debug`. Useful when bounds are long or involve associated types.

**zero-cost abstraction** — The principle that high-level features (iterators, generics, traits) compile to machine code as efficient as hand-written low-level code. You do not pay at runtime for the convenience.

**? operator** — The postfix `?` that, applied to a **Result** or **Option**, returns the value on success or returns early with the error/`None` on failure. It replaces repetitive check-and-return code around C return codes.
