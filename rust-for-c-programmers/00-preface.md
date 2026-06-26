# Preface

## What this book is

This is a textbook for one specific reader: an **experienced C programmer** who
needs to learn **Rust** quickly and well. It does not assume you know any other
language. It does assume you are comfortable with the things C teaches deeply:
pointers, memory, the stack and the heap, structs, `malloc`/`free`, and how source
files become a running program.

That C background is the perfect starting point. Rust was built for exactly the
jobs C is used for — operating systems, browsers, embedded devices, network
services, command-line tools — and it makes the same core promise as C: **no
garbage collector, no hidden runtime, and full control over memory and speed.**
The difference is that Rust's compiler *checks your memory and threading for you*,
and refuses to build the program until it can prove there are no dangling
pointers, no use-after-free, no double-free, and no data races.

The writing is deliberately simple. Sentences are short. Every new word is
explained the first time it appears. The goal is that the **ideas** are the only
hard part — never the language the book is written in.

## The one big promise (and the one hard idea)

**Most of Rust will feel familiar.** Integers, structs, functions, control flow,
the stack and heap, and compiling to a fast native binary are all here and behave
much like C.

There is really just **one big new idea** to learn, and the whole language grows
out of it: **ownership**. In C, *you* keep track of who allocated a buffer and who
must free it, and the bugs that follow (use-after-free, double-free, leaks, data
races) are some of the worst in our industry. In Rust, the compiler tracks this
for you using a few simple rules. Learning those rules is the real work of this
book. Once they click, Rust stops feeling like a fight and starts feeling like a
very strict, very helpful colleague.

> **Mental model.** Rust is "C where the compiler is your `malloc`/`free`
> accountant and your thread-safety reviewer." You still control memory; you just
> write down the ownership rules the compiler then enforces.

## The biggest shifts from C (your mental map)

Keep these in mind. The whole book is an expansion of this list.

1. **Ownership replaces `malloc`/`free`.** Every value has exactly one *owner*.
   When the owner goes out of scope, the value is freed automatically (this is
   RAII — the same idea as a C++ destructor). You never call `free`, and you can
   never free twice. (Chapter 7)

2. **Borrowing has rules the compiler enforces.** You can lend out a value by
   reference. The rule: **either many shared `&` readers, or exactly one `&mut`
   writer — never both at once.** This single rule removes data races and dangling
   pointers. (Chapter 8)

3. **No `NULL`, no uninitialized memory.** There is no null pointer. "Maybe a
   value, maybe nothing" is expressed with the `Option<T>` type, and the compiler
   forces you to handle the "nothing" case. Every variable must be initialized
   before use. (Chapter 12)

4. **Errors are values, not exceptions.** A function that can fail returns
   `Result<T, E>`, and the `?` operator makes checking it short and clean. There
   are no exceptions and no `errno`. (Chapter 13)

5. **Values move by default.** Assigning or passing a value *moves* its ownership;
   the old variable can no longer be used. This is how the compiler guarantees a
   single owner. Cheap "plain data" types are `Copy` instead and behave like C.
   (Chapter 7)

6. **Immutable by default.** `let x = 5;` is read-only. You must write `let mut x`
   to allow changes. This is the reverse of C, and it is a big part of why Rust
   code is easy to reason about. (Chapter 4)

7. **Traits instead of inheritance; generics are zero-cost.** Shared behavior is
   described by *traits* (like a type-checked interface), and generics are compiled
   into specialized machine code (monomorphization), so abstraction costs nothing
   at runtime. (Chapters 14–15)

8. **`unsafe` is a small, explicit escape hatch.** Raw pointers, pointer
   arithmetic, and C interop live inside `unsafe` blocks. The other 99% of your
   code is checked. This is how Rust talks to C and the hardware. (Chapter 25)

If those ideas sink in, everything else is detail.

## A few things C programmers must *unlearn*

- **You do not free memory.** Do not look for the `free`. The owner's scope ending
  *is* the free. (Chapter 7)
- **No pointer arithmetic in normal code.** To walk memory you use slices and
  iterators, which are bounds-checked. Raw pointer math exists only in `unsafe`.
- **No implicit numeric conversions.** You must write `x as i64` (or use
  `From`/`Into`). The compiler will not quietly mix `i32` and `i64`. (Chapter 4)
- **Variables are immutable until you say otherwise.** Reach for `mut` only when
  you truly need to change a value.
- **Strings are UTF-8, not NUL-terminated `char` arrays.** A `String` knows its
  length, is not null-terminated, and indexing by integer is deliberately not
  allowed. (Chapter 10)

## How the book is organized

The chapters build on each other, but you do not have to read them all before
writing code. The **Start Here** part (Chapters 1–6) plus the **Ownership** part
(Chapters 7–10) is the real core. The `README.md` lays out a **fast-track** reading
order; follow it, then read the rest as your work demands.

Two reference sections are meant to stay open on your desk:

- **Appendix A — the C → Rust cheat sheet:** a translation table for nearly every
  construct, plus the commands you'll type all day.
- **Chapter 28 — gotchas:** the borrow-checker traps that catch C programmers, in
  one checklist.

## Conventions

You will see short **callout boxes**. Each one always means the same thing:

> **Mental model.** A quick analogy to build intuition.

> **C vs Rust.** A direct comparison with C.

> **Watch out.** A common trap or surprise. Read these carefully.

> **Rule of thumb.** Advice you can apply right away.

> **Deep dive.** An optional, deeper aside. Safe to skip on the fast track.

> **Try it.** A tiny exercise.

Code blocks are labeled by language — `rust`, `c`, or `sh` (a shell command you
type in a terminal). Some Rust examples **intentionally do not compile**, to show a
mistake the compiler catches; those begin with a `// COMPILE ERROR:` comment.
Most chapters end with **Key takeaways**, a **Watch out** checklist, and a few
**Interview questions** with answers, so the book doubles as interview preparation.

Let's begin with *why* Rust exists, and what it feels like coming from C.
