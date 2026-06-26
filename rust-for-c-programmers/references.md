# References and Further Reading

> **What you'll learn.** Where to go next: the official docs, the best books, places
> to practice, the everyday tools, the community, and a short list aimed
> specifically at C programmers.

Every link below is to an official or well-known resource. Where a resource has no
stable official URL, it is named without a link.

## Official documentation

- **The Rust Programming Language** ("the book") — the canonical introduction.
  <https://doc.rust-lang.org/book/>
- **Rust by Example** — runnable examples for each feature.
  <https://doc.rust-lang.org/rust-by-example/>
- **The Standard Library (`std`) docs** — API reference for everything in `std`.
  <https://doc.rust-lang.org/std/>
- **The Rust Reference** — the detailed language reference (syntax and semantics).
  <https://doc.rust-lang.org/reference/>
- **The Rustonomicon** — the dark arts of `unsafe` Rust.
  <https://doc.rust-lang.org/nomicon/>
- **The Rust API Guidelines** — how to design idiomatic, predictable APIs.
  <https://rust-lang.github.io/api-guidelines/>
- **The Cargo Book** — the full guide to Cargo and `Cargo.toml`.
  <https://doc.rust-lang.org/cargo/>
- **Asynchronous Programming in Rust** ("the async book") — `async`/`.await`,
  futures, and executors. <https://rust-lang.github.io/async-book/>
- **The Edition Guide** — what changed in each edition, including 2024.
  <https://doc.rust-lang.org/edition-guide/>

## Books

- **Programming Rust** — Jim Blandy, Jason Orendorff & Leonora Tindall (O'Reilly).
  Thorough and systems-focused; excellent for C and C++ programmers.
- **Rust for Rustaceans** — Jon Gjengset (No Starch Press). Intermediate; deepens
  your understanding once you know the basics.
- **Rust Atomics and Locks** — Mara Bos (O'Reilly). The clearest treatment of
  low-level concurrency, atomics, and memory ordering. Read online:
  <https://marabos.nl/atomics/>
- **Zero To Production In Rust** — Luca Palmieri. A hands-on guide to building real
  web services (APIs, databases, testing, deployment).
  <https://www.zero2prod.com/>

## Practice

- **Rustlings** — small fix-the-code exercises that track the official book.
  <https://github.com/rust-lang/rustlings>
- **Exercism — Rust track** — graded exercises with mentor feedback.
  <https://exercism.org/tracks/rust>
- **Advent of Code** — yearly programming puzzles; a popular, fun way to practice
  Rust on real problems. <https://adventofcode.com/>

## Tools

- **Clippy** — the official linter; catches mistakes and unidiomatic code.
  <https://doc.rust-lang.org/clippy/>
- **rustfmt** — the official code formatter.
  <https://github.com/rust-lang/rustfmt>
- **rust-analyzer** — the language server for editor support (completion, go-to,
  inline errors). <https://rust-analyzer.github.io/>
- **Miri** — an interpreter that detects undefined behavior in `unsafe` code.
  <https://github.com/rust-lang/miri>
- **Tokio** — the most widely used async runtime.
  <https://docs.rs/tokio/> and <https://tokio.rs/>
- **Serde** — the standard serialization/deserialization framework.
  <https://serde.rs/>
- **clap** — the standard command-line argument parser.
  <https://docs.rs/clap/>

## Community

- **The Rust Users Forum** — ask questions, get help.
  <https://users.rust-lang.org/>
- **r/rust** — the Rust subreddit, news and discussion.
  <https://www.reddit.com/r/rust/>
- **This Week in Rust** — a weekly newsletter of news, blog posts, and crates.
  <https://this-week-in-rust.org/>

## Especially for C programmers

- **The Rustonomicon — FFI chapter** — calling C from Rust and back across the
  boundary. <https://doc.rust-lang.org/nomicon/ffi.html>
- **`std::ffi` docs** — `CString`, `CStr`, `OsString`, and C-compatible types.
  <https://doc.rust-lang.org/std/ffi/>
- **The `bindgen` User Guide** — auto-generate Rust bindings from C headers.
  <https://rust-lang.github.io/rust-bindgen/>
- **The `cbindgen` User Guide** — generate C headers from a Rust library so C code
  can call it. <https://github.com/mozilla/cbindgen/blob/master/docs.md>
