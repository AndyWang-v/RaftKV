# Authoring guide for "Go for C Programmers"

This file defines the voice, structure, and formatting rules for every chapter, so
the whole book reads as if one person wrote it. Read it fully before writing.

## Audience

- An experienced **C programmer** (6+ years) who is **new to Go** and wants to be
  productive fast.
- **English is not their first language.** Write in **simple, clear English**:
  short sentences, common words, one idea per sentence. Define every technical term
  the first time it appears. Avoid idioms, slang, and humor that depends on
  wordplay.

## Target Go version

- **Go 1.26** (installed on this machine). Write modern, idiomatic Go:
  - Loop variables are **per-iteration** (Go 1.22+); the old "loop variable capture"
    bug does **not** apply to `for` loops anymore. When you mention it (e.g. in the
    gotchas chapter), explain it as *historical* and note it is fixed since 1.22.
  - Use `any`, not `interface{}`.
  - Use the builtins `min`, `max`, `clear` where natural.
  - `for i := range n` (integer range, Go 1.22+) is allowed.
  - Generics (type parameters) are available and stable.

## Read these before writing (match their voice and quality)

- `00-preface.md` — the conventions and callout boxes.
- `01-why-go.md` — **the exemplar chapter.** Match its structure, depth, and tone.
- `README.md` — the full chapter list and titles (use these exact titles when you
  cross-reference other chapters).
- `../book/03-stage0-scaffolding.md` and `../book/appendix-go-for-c.md` — existing
  high-quality writing in this repository; same house style.

## Chapter structure (follow exactly)

1. Title line: `# Chapter N — Title` (use an em dash `—`). For appendices use
   `# Appendix A — Title` etc.
2. A `> **What you'll learn.**` blockquote: 1–3 short sentences.
3. Body: `##` sections, `###` subsections. **Do not** manually number sections
   (write `## Slices`, not `## 8.2 Slices`).
4. End with these sections, in this order:
   - `## Key takeaways` — a bullet list of the chapter's main points.
   - `## Watch out (gotchas for C programmers)` — bullet list of traps/surprises.
   - `## Interview questions` — 3 to 6 items. Format each as a bold question then a
     plain-text answer:

     ```
     **Q: question text?**
     A: a clear, correct, concise answer.
     ```
   - Optional: `## Try it` — one tiny hands-on exercise.

Keep each chapter focused and complete but not bloated: about **200–450 lines** of
Markdown. Clarity and good examples beat length.

## Callout boxes (use these exact labels)

Use Markdown blockquotes. Pick the label that fits:

- `> **Mental model.**` — a quick analogy to build intuition.
- `> **C vs Go.**` — a direct comparison with C.
- `> **Watch out.**` — a trap or surprise (inline, in addition to the end-of-chapter list).
- `> **Rule of thumb.**` — practical advice.
- `> **Deep dive.**` — optional, more advanced aside.
- `> **Try it.**` — a tiny exercise.

## Comparing to C (do this constantly — it is the whole point)

- Explain each new Go idea by relating it to C: "in C you would…; in Go you…".
- Use comparison tables with this header style when helpful:

  ```
  | Concept | C | Go | Note |
  |---|---|---|---|
  ```

- Show C code in ```c fences, Go in ```go fences, shell commands in ```sh fences.

## Code rules

- **All Go code must compile on Go 1.26** and be idiomatic. Prefer small, complete,
  runnable snippets. If a snippet is a fragment, make that obvious (e.g. `// ...`).
- Handle errors in examples the idiomatic way (`if err != nil { ... }`); do not
  ignore them silently unless you are explicitly teaching that point.
- Run `gofmt` mentally: tabs for indentation, braces on the same line.
- Keep examples realistic but minimal. No giant programs except where a chapter is
  explicitly about building one (web/CLI/capstone).

## Diagrams (include where they aid understanding)

The book uses BOTH formats:

- **Mermaid** for flow, sequence, state, and architecture diagrams. Use a fenced
  block tagged `mermaid`. KEEP IT SIMPLE AND VALID:
  - Prefer `flowchart LR` / `flowchart TD`, `sequenceDiagram`, `stateDiagram-v2`.
  - Node labels: keep them plain. Use `\n` for a line break inside a label. Avoid
    parentheses, quotes, and special characters inside labels (they break the
    parser). Example label: `goroutine\n(G)`.
  - Test your mental model of the syntax; invalid Mermaid will fail the PDF build.
- **ASCII/box diagrams** in a plain ``` block for memory layouts, struct/slice
  headers, stack vs heap, byte buffers, etc. These are great for showing bytes and
  pointers — exactly what C programmers like to see.

Aim for at least one diagram in any chapter where structure or flow matters
(slices, strings, interfaces, goroutines/scheduler, channels, GC, HTTP lifecycle).

## Cross-references

Refer to other chapters by number and title, e.g. "(see Chapter 8 — Arrays, Slices,
and Strings)". Do not invent chapter numbers; use `README.md`.

## Do NOT

- Do not include a table of contents or nav links inside a chapter.
- Do not add comments to code that merely restate the code; comments should explain
  *why* or a subtlety.
- Do not use emojis.
- Do not write files outside the paths you are explicitly assigned.

## When done

Reply with a SHORT summary only: the files you wrote and their line counts, plus any
caveats. Do **not** paste chapter contents back.
