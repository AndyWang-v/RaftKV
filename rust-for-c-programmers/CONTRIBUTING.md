# Authoring guide for "Rust for C Programmers"

This file defines the voice, structure, and formatting rules for every chapter, so
the whole book reads as if one person wrote it. Read it fully before writing.

## Audience

- An experienced **C programmer** (6+ years) who is **new to Rust** and wants to be
  productive fast.
- **English is not their first language.** Write in **simple, clear English**:
  short sentences, common words, one idea per sentence. Define every technical term
  the first time it appears. Avoid idioms, slang, and wordplay.

## Target Rust version

- **Rust 1.96, edition 2024** (installed on this machine). Write modern, idiomatic
  Rust:
  - Use `?` for error propagation; return `Result<T, E>` from fallible functions.
  - Prefer iterators and pattern matching over index loops where it is clearer.
  - Use `derive` (`#[derive(Debug, Clone, PartialEq)]`) where natural.
  - snake_case for functions/variables, CamelCase for types/traits, SCREAMING_CASE
    for consts.
  - In small examples, `.unwrap()`/`.expect()` is acceptable to keep focus, but say
    in a sentence that real code would handle the error (and show `?` in the error
    chapter properly).
  - Edition-2024 details: the prelude includes `Future`/`IntoFuture`; `gen` is a
    reserved keyword; closures capture disjointly. Do not rely on older-edition
    behavior.

## Read these before writing (match their voice and quality)

- `00-preface.md` — the conventions and callout boxes.
- `01-why-rust.md` — **the exemplar chapter.** Match its structure, depth, and tone.
- `README.md` — the full chapter list and titles (use these exact titles when you
  cross-reference other chapters).

## Chapter structure (follow exactly)

1. Title line: `# Chapter N — Title` (use an em dash `—`). For appendices use
   `# Appendix A — Title` etc.
2. A `> **What you'll learn.**` blockquote: 1–3 short sentences.
3. Body: `##` sections, `###` subsections. **Do not** manually number sections
   (write `## Ownership`, not `## 7.2 Ownership`).
4. End with these sections, in this order:
   - `## Key takeaways` — a bullet list of the chapter's main points.
   - `## Watch out (gotchas for C programmers)` — bullet list of traps/surprises.
   - `## Interview questions` — 3 to 6 items, each formatted as:

     ```
     **Q: question text?**
     A: a clear, correct, concise answer.
     ```
   - Optional: `## Try it` — one tiny hands-on exercise.

Keep each chapter focused and complete but not bloated: about **250–500 lines** of
Markdown (the ownership, traits, and enum chapters may run longer — depth and good
examples beat brevity). Clarity wins.

## Callout boxes (use these exact labels)

Use Markdown blockquotes. Pick the label that fits:

- `> **Mental model.**` — a quick analogy to build intuition.
- `> **C vs Rust.**` — a direct comparison with C.
- `> **Watch out.**` — a trap or surprise (inline, in addition to the end list).
- `> **Rule of thumb.**` — practical advice.
- `> **Deep dive.**` — optional, more advanced aside.
- `> **Try it.**` — a tiny exercise.

## Comparing to C (do this constantly — it is the whole point)

- Explain each new Rust idea by relating it to C: "in C you would…; in Rust you…".
- Map Rust concepts onto C ones: ownership/`Drop` ↔ `malloc`/`free` and RAII;
  `&`/`&mut` ↔ pointers but checked; `Vec<T>` ↔ `malloc`'d array + length;
  `&str`/`String` ↔ `char *`; traits ↔ structs of function pointers / interfaces;
  generics ↔ macros/`void *`; `Result`/`?` ↔ return codes/`errno`; `unsafe` ↔
  ordinary C.
- Use comparison tables with this header style when helpful:

  ```
  | Concept | C | Rust | Note |
  |---|---|---|---|
  ```

- Show C code in ```c fences, Rust in ```rust fences, shell commands in ```sh.

## Code rules (CRITICAL — the code must be correct)

- **All Rust code must compile on Rust 1.96, edition 2024**, UNLESS it is a
  teaching example of a compiler error.
- **Mark every intentionally-non-compiling example.** Begin such a block with a
  comment on the first line: `// COMPILE ERROR: <short reason>`, and where useful
  show the gist of the compiler's message in a comment or following sentence. This
  is essential — it tells the reader (and the verification step) that the failure
  is on purpose. Do NOT leave a broken example unmarked.
- Complete, runnable examples should include `fn main() { ... }` so they can be
  compiled and run. Fragments (a lone function or struct) are fine for
  illustration; make it obvious they are fragments.
- Prefer small, focused examples. Handle errors idiomatically in the error chapter;
  elsewhere `unwrap()` is acceptable with a one-line caveat.
- Run `rustfmt` style in your head: 4-space indentation, braces on the same line.

## Diagrams (include where they aid understanding)

The book uses BOTH formats:

- **Mermaid** for flow, sequence, state, and ownership/borrow diagrams. Use a
  fenced block tagged `mermaid`. KEEP IT SIMPLE AND VALID:
  - Prefer `flowchart LR` / `flowchart TD`, `sequenceDiagram`, `stateDiagram-v2`.
  - Node labels must be plain: NO parentheses, quotes, `&`, `<`, `>`, or commas
    inside labels (they break the parser). Use `\n` for a line break. For example,
    write a label like `borrow mut\nexclusive` rather than `&mut (exclusive)`.
- **ASCII/box diagrams** in a plain ``` block for memory layouts: the stack vs the
  heap, a `Vec` header (ptr/len/cap) pointing at a heap buffer, ownership moves,
  `Box`/`Rc`/`RefCell`, a trait object's (data, vtable) pair. C programmers love
  these byte/pointer pictures — use them often.

Aim for at least one diagram in any chapter where structure or flow matters
(ownership, borrowing, lifetimes, slices/strings, smart pointers, threads, async,
trait objects).

## Cross-references

Refer to other chapters by number and title, e.g. "(see Chapter 8 — Borrowing and
References)". Do not invent chapter numbers; use `README.md`.

## Do NOT

- Do not include a table of contents or nav links inside a chapter.
- Do not add comments to code that merely restate the code; comments should explain
  *why* or a subtlety (the `// COMPILE ERROR:` marker is required, though).
- Do not use emojis.
- Do not write files outside the paths you are explicitly assigned.

## When done

Reply with a SHORT summary only: the files you wrote and their line counts, plus any
caveats (especially any snippet you are unsure compiles). Do **not** paste chapter
contents back.
