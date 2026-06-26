# Go for C Programmers

> **Learn Go fast by building on what you already know.**

This is a complete, beginner-friendly book that teaches the **Go programming
language** to someone who already knows **C**. Every new idea is explained by
comparing it to the C you already understand. The English is kept simple and the
sentences short on purpose, so the ideas — not the words — are the hard part.

You know pointers, memory, structs, and how a program is laid out. That is a huge
head start. Go keeps most of that mental model and removes the painful parts
(manual `malloc`/`free`, header files, the preprocessor, pointer arithmetic) while
adding a few powerful new tools (garbage collection, goroutines, interfaces, and
one of the best toolchains in any language). This book shows you exactly what
maps over, what changes, and what to watch out for.

This book targets **Go 1.26** (the version installed here). A few behaviors
changed in recent versions; where that matters, the text says so.

---

## Who this is for

- You have written C for a while (the examples assume you are comfortable with
  pointers, structs, stack vs heap, and the compile/link model).
- You are new to Go, or have only dabbled, and you need to be productive **soon**.
- English may not be your first language — so this book avoids slang and explains
  every term the first time it appears.

---

## How to read this book (the fast track)

The book is comprehensive, but you do **not** need to read all of it before you
start writing real Go. Follow the **fast track** first; it is the shortest path
from "I know C" to "I can contribute to a Go project." Then read the rest as you
hit each topic at work.

| Order | Read these first (the fast track) | Why |
|------:|-----------------------------------|-----|
| 1 | Ch. 1–7 (the basics) | Syntax, the `go` command, types, functions, pointers. Get a program running. |
| 2 | **Ch. 8** (slices & strings) | The single biggest difference from C arrays. Do not skip this. |
| 3 | Ch. 9–12 (maps, structs, interfaces, errors) | The everyday building blocks of all Go code. |
| 4 | Ch. 13–15 (goroutines, channels, sync) | Concurrency is *why* many teams choose Go. |
| 5 | **Ch. 26** (gotchas) + Ch. 21 (testing) | Avoid the common traps; write tests like a local. |

Everything else — the runtime internals, generics, the standard-library tour, web
and CLI building, idioms, interview prep, and the capstone projects — is there
when you need it. The **C → Go cheat sheet** (Appendix A) is your day-one desk
reference.

> If you have one evening: read Chapters 1, 2, 8, and 26. That alone will stop you
> from writing "C in Go" and from falling into the famous slice and goroutine
> traps.

---

## Table of contents

### Part I — Start Here (the fast track)
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 1 | [`01-why-go.md`](01-why-go.md) | Why Go exists, Go vs C philosophy, what's gone and what's new, your first program. |
| 2 | [`02-toolchain-and-commands.md`](02-toolchain-and-commands.md) | Installing Go and the full `go` command tour: `run`, `build`, `test`, `mod`, `fmt`, `vet`, `get`, `install`, `work`, `env`, `doc`. |
| 3 | [`03-program-structure.md`](03-program-structure.md) | Packages vs headers, imports, `main`, exported vs unexported, `init`. |
| 4 | [`04-types-variables-constants.md`](04-types-variables-constants.md) | Types and exact sizes, `var` vs `:=`, constants, `iota`, and Go's all-important **zero values**. |
| 5 | [`05-control-flow.md`](05-control-flow.md) | `if`, the one `for` loop, `switch` (no fallthrough!), `range`, labels. |
| 6 | [`06-functions.md`](06-functions.md) | Multiple return values, named returns, variadics, closures, and `defer`. |
| 7 | [`07-pointers.md`](07-pointers.md) | Pointers without pointer arithmetic, `new`, `nil`, and pass-by-value semantics. |

### Part II — Data and Types
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 8 | [`08-arrays-slices-strings.md`](08-arrays-slices-strings.md) | **The big one.** Arrays vs slices (the fat pointer), `append`, aliasing traps, strings, bytes, and runes (UTF-8). |
| 9 | [`09-maps.md`](09-maps.md) | Hash maps built in: `make`, the comma-ok idiom, random iteration order, nil maps. |
| 10 | [`10-structs-and-methods.md`](10-structs-and-methods.md) | Structs, methods, value vs pointer receivers, and embedding (composition, not inheritance). |
| 11 | [`11-interfaces.md`](11-interfaces.md) | The big leap: implicit interfaces, the `(type, value)` fat pointer, `any`, type switches. |
| 12 | [`12-errors.md`](12-errors.md) | Errors are values (no exceptions): the `error` type, wrapping, `panic`/`recover`. |

### Part III — Concurrency (Go's superpower)
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 13 | [`13-goroutines-scheduler.md`](13-goroutines-scheduler.md) | Goroutines vs threads, the G-M-P scheduler, stack growth, leaks. |
| 14 | [`14-channels-select.md`](14-channels-select.md) | Channels, buffered vs unbuffered, `select`, closing, deadlocks. |
| 15 | [`15-sync-and-context.md`](15-sync-and-context.md) | `sync.Mutex`/`WaitGroup`/`Once`, atomics, the memory model, and `context`. |
| 16 | [`16-concurrency-patterns.md`](16-concurrency-patterns.md) | Pipelines, fan-in/fan-out, worker pools, `errgroup`, and the race detector. |

### Part IV — Under the Hood & Organizing Code
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 17 | [`17-memory-and-gc.md`](17-memory-and-gc.md) | Stack vs heap, escape analysis, the garbage collector, reducing allocations. |
| 18 | [`18-packages-and-modules.md`](18-packages-and-modules.md) | Modules, `go.mod`/`go.sum`, versioning, `internal/`, and project layout. |
| 19 | [`19-generics.md`](19-generics.md) | Type parameters and constraints (vs C macros and `void *`). |
| 20 | [`20-standard-library-tour.md`](20-standard-library-tour.md) | The batteries included: `fmt`, `strings`, `strconv`, `io`, `os`, `time`, `encoding/json`, and more. |

### Part V — Doing Real Work
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 21 | [`21-testing.md`](21-testing.md) | `go test`, table-driven tests, subtests, benchmarks, fuzzing, coverage. |
| 22 | [`22-tooling.md`](22-tooling.md) | `gofmt`, `go vet`, `staticcheck`, `golangci-lint`, the race detector, `pprof`, and `delve`. |
| 23 | [`23-web-services.md`](23-web-services.md) | Build an HTTP/JSON API with `net/http`: handlers, routing, middleware, graceful shutdown. |
| 24 | [`24-cli-tools.md`](24-cli-tools.md) | Build a command-line tool: `flag`, subcommands, stdin/stdout, exit codes, config. |

### Part VI — Mastery & Reference
| Ch. | File | What you'll learn |
|----:|------|-------------------|
| 25 | [`25-idioms-and-style.md`](25-idioms-and-style.md) | How Go *should* be written: naming, error style, "accept interfaces, return structs," functional options. |
| 26 | [`26-gotchas-for-c-programmers.md`](26-gotchas-for-c-programmers.md) | A single checklist of the traps that bite C programmers most. |
| 27 | [`27-interview-questions.md`](27-interview-questions.md) | Real Go interview questions with clear, correct answers. |
| 28 | [`28-capstone-projects.md`](28-capstone-projects.md) | Three guided projects: a CLI tool, a JSON REST API, and a concurrent fetcher. |
| A | [`appendix-c-to-go-cheatsheet.md`](appendix-c-to-go-cheatsheet.md) | One-page C↔Go translation table + command cheat sheet. |
| B | [`appendix-glossary.md`](appendix-glossary.md) | Every Go term, defined in plain English. |
| — | [`references.md`](references.md) | Official docs, books, courses, and tools. |

---

## How to run the examples

Every code sample is real Go. To try one, put it in a file and run it — no build
step or `Makefile` needed:

```sh
# 1. make a scratch folder and a module (explained in Chapter 2)
mkdir play && cd play && go mod init play

# 2. save a sample as main.go, then:
go run .
```

The `examples/` folder in this book holds the larger, runnable programs from the
chapters as its own Go module.

## How to make a PDF (optional — the Markdown is the source of truth)

```sh
make html          # open in a browser, then Print > Save as PDF
make epub          # e-reader friendly
make pdf           # polished PDF; first run: brew install tectonic
make mermaid-deps  # optional: render the Mermaid diagrams as pictures in the PDF
```

---

## Conventions used in this book

Throughout the chapters you will see short **callout boxes**. They always mean the
same thing:

> **Mental model.** A quick analogy to lock in the intuition.

> **C vs Go.** A direct, side-by-side comparison with C.

> **Watch out.** A common trap or surprise — read these carefully.

> **Rule of thumb.** Practical advice you can apply immediately.

> **Deep dive.** An optional, more advanced aside. Safe to skip on the fast track.

> **Try it.** A tiny exercise to make the idea stick.

Code is shown in fenced blocks and labeled by language: `go`, `c`, or `sh` (shell).
Most chapters end with **Key takeaways**, a **Watch out** checklist, and a few
**Interview questions** with answers.

Let's begin.
