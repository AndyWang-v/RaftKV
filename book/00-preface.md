# Preface

## What this is

This is the development journal and textbook for a single, ambitious project:
**building a linearizable, fault-tolerant key-value store, backed by a
from-scratch implementation of the Raft consensus algorithm, in Go**, with a
React dashboard that lets you *watch* the cluster elect leaders, replicate a log,
and recover from partitions in real time.

It is written for one specific reader: an experienced **C programmer** who is new
to Go and to distributed systems, and who wants to *actually understand* what they
build — not paste in code that happens to pass a test. If that's you, this book
walks the whole path with you.

## Why build Raft from scratch?

Almost anyone can wire together an existing consensus library. Very few engineers
can explain *why* a leader must not commit an entry from a previous term by replica
count alone, or *why* the election "up-to-date" check is a safety property rather
than an optimization. The gap between those two groups is exactly the gap this
project closes. The key-value store on top is deliberately trivial — an in-memory
map — because its only job is to make the consensus layer *observable*. The Raft
layer is the part that teaches you something durable.

## The learning philosophy (read this twice)

This project is built **stage by stage**, and the method matters as much as the
code. The rules we follow:

1. **One stage at a time.** We do not generate the whole system at once. That path
   produces plausible code that fails under partition and teaches nothing. Each
   stage has a goal and a set of acceptance tests; a stage is "done" only when its
   tests pass under Go's race detector (`go test -race`).

2. **Explain before writing.** For every Raft rule, we state *why* it exists before
   we implement it — especially the three subtle ones: the election restriction,
   the "previous-term commit" rule (the famous Figure 8), and the election-timer
   reset rule. If the explanation is wrong, the code is wrong.

3. **Own the test harness.** The deterministic simulator and harness are where
   understanding compounds. We read every line; we trust the green checkmark only
   because we know what it checks.

4. **Forbid shortcuts.** No holding a lock across a network call. No skipping
   persistence before replying. No dropping the term clause in commit advancement.
   These are not style preferences; they are correctness.

5. **Debug before asking for a fix.** When a chaos test fails, we reproduce it with
   a fixed seed, replay it, and read the log *before* changing code. The canonical
   bugs (see the appendix) are the curriculum.

6. **Checkpoint the learning.** After each stage we write down what it taught and
   one trade-off we hit. This book *is* that checkpoint, expanded.

## What we are building (and not)

**Goals**

- Correctness under adversarial conditions: crashes, partitions, message loss,
  reordering, and restarts.
- Idiomatic, race-free Go that survives `go test -race`.
- Clean seams between Raft, the state machine, the network, and storage — so each
  is independently testable.

**Non-goals** (scope discipline)

- No sharding / multi-Raft. **One** Raft group.
- No Byzantine fault tolerance (we assume crash-stop, not malicious, failures).
- No production storage engine; file-backed persistence is sufficient (with an
  optional embedded-engine upgrade later).
- Membership changes, if we do them, are single-server-at-a-time, not full joint
  consensus.

## The toolchain

At the time of writing the project uses:

- **Go 1.26.4** — the consensus core, the KV store, the simulator, all tests.
- **Node.js 24 / npm 11** — the React dashboard (a later stage).
- **pandoc 3.9** — to render this book.

## How to read this book

The chapters follow the order we actually built things, which is also a sensible
teaching order:

- **Foundations** and **Leader Election** are pure theory — the mental model you
  need before any code makes sense.
- **Stage 0** begins the code: the data types, the interfaces (the seams), and the
  Go-for-C concepts you meet along the way.
- Later chapters add one build stage each.

Two appendices run throughout: a **Go-for-C cheat sheet** that grows as we meet new
language features, and the **canonical Raft bug checklist** — the list of mistakes
that everyone hits, which doubles as a debugging guide.

Sprinkled through the chapters are two recurring boxes:

> **Checkpoint** boxes pose a question and give the ideal answer. Try to answer
> before reading on; that's where the learning is.

> **Trade-off** boxes flag a design decision and its alternatives. Distributed
> systems are made of trade-offs; noticing them is the skill.

Let's begin.
