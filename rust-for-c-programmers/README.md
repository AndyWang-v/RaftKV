# Rust for C Programmers

> **Learn Rust fast by building on what you already know.**

This is a complete, beginner-friendly book that teaches the **Rust programming
language** to someone who already knows **C**. Every new idea is explained by
comparing it to the C you already understand. The English is kept simple and the
sentences short on purpose, so the ideas — not the words — are the hard part.

You know pointers, memory, the stack and heap, structs, and `malloc`/`free`. That
is exactly the background Rust was designed around. Rust keeps C's core promise —
**no garbage collector, no hidden runtime, speed and control down to the byte** —
but adds a compiler that proves your memory and threading are safe *before* the
program ever runs. The hard part is not the syntax; it is one new idea called
**ownership**. This book gets you there step by step.

This book targets **Rust 1.96** with **edition 2024** (the version installed here).

---

## Who this is for

- You have written C for a while (the examples assume you are comfortable with
  pointers, memory, the stack/heap, structs, and the compile/link model).
- You are new to Rust, or bounced off the borrow checker once, and you want to be
  productive **soon**.
- English may not be your first language — so this book avoids slang and explains
  every term the first time it appears.

---

## How to read this book (the fast track)

The book is comprehensive, but you do **not** need to read all of it before you
write real Rust. Follow the **fast track** first; it is the shortest path from
"I know C" to "I can fight the borrow checker and win." Then read the rest as you
meet each topic at work.

| Order | Read these first (the fast track) | Why |
|------:|-----------------------------------|-----|
| 1 | Ch. 1–6 (the basics) | Syntax, Cargo, types, control flow, functions. Get a program running. |
| 2 | **Ch. 7–9** (ownership, borrowing, lifetimes) | The heart of Rust. This is what is new. Do not skip. |
| 3 | Ch. 10–13 (slices/strings, structs, enums, errors) | The everyday building blocks of all Rust code. |
| 4 | Ch. 15 (traits) + Ch. 16 (collections & iterators) | How Rust does abstraction and data structures. |
| 5 | **Ch. 28** (gotchas) | The borrow-checker traps that catch C programmers, in one place. |

Everything else — smart pointers, concurrency, async, unsafe/FFI, macros, the
tooling, idioms, interview prep, and the capstone projects — is there when you
need it. The **C → Rust cheat sheet** (Appendix A) is your day-one desk reference.

> If you have one evening: read Chapters 1, 2, 7, and 8. Ownership and borrowing
> are 80% of what makes Rust *feel* different from C. Get those and the rest falls
> into place.

---

## Table of contents

### Part I — Start Here (the fast track)
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 1 | [`01-why-rust.md`](01-why-rust.md) | Why Rust exists, Rust vs C philosophy, safety without a GC, your first program. |
| 2 | [`02-toolchain-and-cargo.md`](02-toolchain-and-cargo.md) | `rustup`, `rustc`, and the full `cargo` tour: `new`, `build`, `run`, `test`, `check`, `add`, `fmt`, `clippy`, `doc`. |
| 3 | [`03-program-structure.md`](03-program-structure.md) | Crates, modules, paths, `use`, and `pub` visibility — no headers. |
| 4 | [`04-variables-types-constants.md`](04-variables-types-constants.md) | `let`, `mut`, shadowing, scalar/compound types with exact sizes, type inference, constants. |
| 5 | [`05-control-flow.md`](05-control-flow.md) | `if`/`loop`/`while`/`for` as expressions, ranges, and an intro to `match` and `if let`. |
| 6 | [`06-functions-and-closures.md`](06-functions-and-closures.md) | Functions, expression-vs-statement, closures, and the `Fn`/`FnMut`/`FnOnce` traits. |

### Part II — Ownership (the heart of Rust)
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 7 | [`07-ownership-and-moves.md`](07-ownership-and-moves.md) | **The big idea.** Owners, moves, `Drop`, and RAII — `malloc`/`free` enforced by the compiler. |
| 8 | [`08-borrowing-and-references.md`](08-borrowing-and-references.md) | `&` and `&mut`, the borrow rules (aliasing XOR mutation), and the borrow checker. |
| 9 | [`09-lifetimes.md`](09-lifetimes.md) | Lifetimes explained simply: what `'a` means, why it exists, and elision. |
| 10 | [`10-slices-and-strings.md`](10-slices-and-strings.md) | `&[T]`/`Vec<T>` and `&str`/`String` (UTF-8), and how they relate to C arrays and `char *`. |

### Part III — Data and Abstraction
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 11 | [`11-structs-and-methods.md`](11-structs-and-methods.md) | Structs, `impl`, methods (`&self`/`&mut self`/`self`), and associated functions. |
| 12 | [`12-enums-and-pattern-matching.md`](12-enums-and-pattern-matching.md) | Sum types, `Option`, `Result`, and exhaustive `match` — Rust's killer feature. |
| 13 | [`13-error-handling.md`](13-error-handling.md) | `Result`, the `?` operator, `panic!`, and error crates — no exceptions. |
| 14 | [`14-generics.md`](14-generics.md) | Generics by monomorphization (vs C macros and `void *`). |
| 15 | [`15-traits.md`](15-traits.md) | Traits: shared behavior, trait objects (`dyn`), and the standard traits. |
| 16 | [`16-collections-and-iterators.md`](16-collections-and-iterators.md) | `Vec`, `String`, `HashMap`, and the lazy iterator pipeline. |

### Part IV — Memory and Smart Pointers
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 17 | [`17-smart-pointers.md`](17-smart-pointers.md) | `Box`, `Rc`, `RefCell`, `Deref`, and interior mutability. |
| 18 | [`18-memory-without-gc.md`](18-memory-without-gc.md) | The stack, the heap, allocation, leaks, and performance — with no GC. |

### Part V — Concurrency
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 19 | [`19-threads-and-concurrency.md`](19-threads-and-concurrency.md) | `std::thread`, `move` closures, `Send`/`Sync`, and "fearless concurrency." |
| 20 | [`20-channels-and-shared-state.md`](20-channels-and-shared-state.md) | Channels (`mpsc`) and shared state with `Arc<Mutex<T>>`. |
| 21 | [`21-async-await.md`](21-async-await.md) | `async`/`.await`, futures, and the Tokio runtime. |

### Part VI — Systems and Tooling
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 22 | [`22-cargo-crates-and-workspaces.md`](22-cargo-crates-and-workspaces.md) | Dependencies, `Cargo.toml`, `Cargo.lock`, semver, and workspaces. |
| 23 | [`23-testing.md`](23-testing.md) | Unit tests, integration tests, and documentation tests with `cargo test`. |
| 24 | [`24-tooling.md`](24-tooling.md) | `rustfmt`, `clippy`, `rust-analyzer`, `miri`, benchmarks, and `cargo expand`. |
| 25 | [`25-unsafe-and-ffi.md`](25-unsafe-and-ffi.md) | `unsafe`, raw pointers, and calling C (and being called by C) via FFI. |
| 26 | [`26-macros.md`](26-macros.md) | `macro_rules!`, `derive`, and a tour of procedural macros. |

### Part VII — Mastery and Reference
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 27 | [`27-idioms-and-style.md`](27-idioms-and-style.md) | How Rust *should* be written: newtypes, builders, `From`/`Into`, "make illegal states unrepresentable." |
| 28 | [`28-gotchas-for-c-programmers.md`](28-gotchas-for-c-programmers.md) | A single checklist of the traps that bite C programmers most. |
| 29 | [`29-interview-questions.md`](29-interview-questions.md) | Real Rust interview questions with clear, correct answers. |
| 30 | [`30-capstone-projects.md`](30-capstone-projects.md) | Three guided projects: a CLI tool, a multithreaded program, and a small data structure. |
| A | [`appendix-c-to-rust-cheatsheet.md`](appendix-c-to-rust-cheatsheet.md) | One-page C↔Rust translation table + command cheat sheet. |
| B | [`appendix-glossary.md`](appendix-glossary.md) | Every Rust term, defined in plain English. |
| — | [`references.md`](references.md) | Official docs, books, courses, and tools. |

---

## How to run the examples

Every code sample is real Rust. To try one, make a scratch project and paste it in:

```sh
cargo new play && cd play      # creates a Cargo project (explained in Chapter 2)
# paste a sample into src/main.rs, then:
cargo run
```

The `examples/` folder in this book holds the larger, runnable programs from the
chapters as its own Cargo project.

## How to make a PDF (optional — the Markdown is the source of truth)

```sh
make html          # open in a browser, then Print > Save as PDF
make epub          # e-reader friendly
make pdf           # polished PDF; needs: brew install tectonic
make mermaid-deps  # optional: render the Mermaid diagrams as pictures in the PDF
```

---

## Conventions used in this book

Throughout the chapters you will see short **callout boxes**. They always mean the
same thing:

> **Mental model.** A quick analogy to lock in the intuition.

> **C vs Rust.** A direct, side-by-side comparison with C.

> **Watch out.** A common trap or surprise — read these carefully.

> **Rule of thumb.** Practical advice you can apply immediately.

> **Deep dive.** An optional, more advanced aside. Safe to skip on the fast track.

> **Try it.** A tiny exercise to make the idea stick.

Code is shown in fenced blocks and labeled by language: `rust`, `c`, or `sh`
(shell). Some Rust examples **intentionally do not compile** — they show a mistake
the compiler catches. Those always start with a `// COMPILE ERROR:` comment that
explains what the compiler will say. Most chapters end with **Key takeaways**, a
**Watch out** checklist, and a few **Interview questions** with answers.

Let's begin.
