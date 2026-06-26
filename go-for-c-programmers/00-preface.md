# Preface

## What this book is

This is a textbook for one specific reader: an **experienced C programmer** who
needs to learn **Go** quickly and well. It does not assume you know any other
language. It does assume you are comfortable with the things C teaches deeply:
pointers, memory, the stack and the heap, structs, and how source files become a
running program.

That C background is worth a lot. Go was created by people who wrote C and C++ for
decades (some of them helped create Unix and C itself). Go *feels* like a cleaned-up
systems language. So instead of starting from zero, this book starts from C and
shows you the differences.

The writing is deliberately simple. Sentences are short. Every new word is
explained the first time it appears. The goal is that the **ideas** are the only
hard part — never the language the book is written in.

## The one big promise

**You already know most of Go.** Variables, functions, structs, pointers, loops,
integers, and the compile-and-run model are all here and behave much like C.

What is left to learn is a small set of *new ideas* and a small set of *changed
rules*. This book is organized around exactly those. If we do our job, Go should
feel less like a new language and more like "C, with the annoying parts removed
and a few superpowers added."

## The seven biggest shifts from C (your mental map)

Keep these in mind. The whole book is really an expansion of this list.

1. **No manual memory management.** There is no `malloc` or `free`. A *garbage
   collector* frees memory for you. You still care about pointers and layout — you
   just stop tracking ownership and lifetimes by hand. (Chapter 17)

2. **Slices replace raw arrays + length.** The thing you reach for is not `T*` plus
   a separate `int len`. It is a *slice* (`[]T`): a small 3-field header
   (pointer, length, capacity) that the language manages and the GC frees. This is
   the #1 source of surprises for C programmers. (Chapter 8)

3. **No pointer arithmetic.** Pointers exist (`*T`, `&x`), but you cannot do
   `p + 1`. To walk memory, you index a slice. This removes a giant class of bugs.
   (Chapter 7)

4. **Errors are ordinary values, not exceptions.** A function that can fail returns
   an `error` as its last result, and you check it with `if err != nil`. There is
   no `try`/`catch`. (Chapter 12)

5. **Concurrency is built into the language.** `go f()` starts a lightweight thread
   (a *goroutine*). *Channels* pass data between goroutines safely. This is the main
   reason many teams pick Go. (Chapters 13–16)

6. **Interfaces are implicit.** A type satisfies an interface just by having the
   right methods — there is no `implements` keyword. Think of it as a type-safe,
   automatic version of a C struct full of function pointers (a vtable). (Chapter 11)

7. **One official style, enforced by tools.** `gofmt` formats every Go file the same
   way, so all Go code looks alike. Capitalization decides visibility (public vs
   private). The toolchain (build, test, format, vet) is one command: `go`.
   (Chapters 2 and 22)

If those seven sink in, everything else is detail.

## A few things C programmers must *unlearn*

- **Zero values are real and useful.** A freshly declared variable is never
  garbage; it is set to a well-defined "zero" (`0`, `""`, `false`, or `nil`). Go
  code relies on this. (Chapter 4)
- **`==` on a struct compares all fields.** Assigning a struct copies it. There is
  no implicit reference like a C++ reference; if you want to share, pass a pointer.
- **Capitalization is not style — it is meaning.** `Name` is exported (public);
  `name` is not (package-private, like `static`). (Chapter 3)
- **There are no implicit numeric conversions.** You must write `int64(x)`
  explicitly. The compiler will not quietly mix `int` and `int64`. (Chapter 4)

## How the book is organized

The chapters build on each other, but you do not have to read them all before
writing code. The **Start Here** part (Chapters 1–7) plus the slices chapter
(Chapter 8) is enough to be dangerous. The `README.md` lays out a **fast-track**
reading order; follow it, then read the rest as your work demands.

Two reference sections are meant to stay open on your desk:

- **Appendix A — the C → Go cheat sheet:** a translation table for nearly every
  construct, plus the commands you'll type all day.
- **Chapter 26 — gotchas:** the traps that catch C programmers, in one checklist.

## Conventions

You will see short **callout boxes**. Each one always means the same thing:

> **Mental model.** A quick analogy to build intuition.

> **C vs Go.** A direct comparison with C.

> **Watch out.** A common trap or surprise. Read these carefully.

> **Rule of thumb.** Advice you can apply right away.

> **Deep dive.** An optional, deeper aside. Safe to skip on the fast track.

> **Try it.** A tiny exercise.

Code blocks are labeled by language — `go`, `c`, or `sh` (a shell command you type
in a terminal). Most chapters end with **Key takeaways**, a **Watch out**
checklist, and a few **Interview questions** with answers, so the book doubles as
interview preparation.

Let's begin with *why* Go exists, and what it feels like coming from C.
