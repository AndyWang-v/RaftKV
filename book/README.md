# The Book

This directory is a living textbook that grows alongside the project. It records
the *journey* (what we built and why), the *knowledge* (Raft theory and Go), the
*code* (every file, explained), and the *trade-offs* we hit at each step.

It is written so you can read it **standalone**, away from the keyboard — print it,
mark it up, and come back to the code.

## Reading order

| File | Contents |
|------|----------|
| `00-preface.md` | What this is, the learning philosophy, how to read it. |
| `01-foundations.md` | The problem: replicated state machines, consensus, roles, terms. |
| `02-leader-election.md` | Election timer, `RequestVote`, split votes, the safety restriction. |
| `03-stage0-scaffolding.md` | Stage 0 code: types, interfaces, testing — with Go-for-C notes. |
| `appendix-go-for-c.md` | A growing Go-vs-C cheat sheet. |
| `appendix-bug-checklist.md` | The canonical Raft bug list to consult when tests fail. |
| `references.md` | Papers, courses, and tools. |

New chapters are added as we complete each build stage.

## Rendering a print-ready copy (optional)

The Markdown is the source of truth; render it whenever you want a printed copy.
`pandoc` is already installed.

```sh
make epub          # e-reader friendly, no extra install
make html          # open in a browser, then Print > Save as PDF
make pdf           # polished PDF; first run: brew install tectonic
```
