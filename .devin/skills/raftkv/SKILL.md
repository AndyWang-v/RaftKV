---
name: raftkv
description: >-
  RaftKV personal learning project — a from-scratch Raft consensus algorithm +
  linearizable key-value store in Go, with a React dashboard later. Load this
  before working on the Raft/KV code, the deterministic test harness, or the
  project's books. Covers project goals, the stage-by-stage LEARNING method,
  repo layout, conventions, build/test commands, and current status.
triggers:
  - user
  - model
---

# RaftKV — project context and working agreement

This is a **personal learning project**, not a company project. The owner is an
experienced **C programmer** who is newer to Go (and Rust) and wants to *deeply
understand* everything we build — not just get working code. Optimize for
learning and correctness over speed.

## What we're building

A linearizable, fault-tolerant key-value store backed by a **from-scratch**
implementation of the Raft consensus algorithm in Go. The KV map is deliberately
trivial; it exists to make the consensus layer observable. A React dashboard (a
later stage) will visualize the live cluster (leader, log replication, inject a
partition and watch re-election).

There is a detailed design document driving this; the build proceeds in stages
0–7. **Do not generate the whole system at once** — that produces plausible code
that fails under partition and teaches nothing.

## How we work (the learning method — follow this)

1. **One stage at a time.** A stage is "done" only when its tests pass under
   `go test -race`. Move on only then.
2. **Explain before writing.** For every Raft rule, state *why* it exists before
   implementing it — especially the election restriction, the Figure-8
   previous-term commit clause, and the election-timer reset rule. If the
   explanation is wrong, the code is wrong.
3. **Teach Go/Rust with C analogies.** The owner knows C well; relate new
   language constructs (slices, interfaces, goroutines, channels, ownership,
   traits, etc.) to their C equivalents and note the differences.
4. **Comprehension checkpoints.** After teaching a concept, pose a short question
   and confirm understanding before moving on.
5. **Forbid shortcuts.** Never hold the mutex across an RPC or channel send; never
   skip persistence-before-reply; never drop the `log[N].Term == currentTerm`
   commit clause.
6. **Keep the book updated.** After each stage, append a chapter to `book/` and
   extend the Go-for-C appendix with any new language concepts.

## Repo layout

```
cmd/                  # binaries: kvnode (server), kvctl (CLI client) — later
internal/raft         # consensus core: types.go, interfaces.go (Transport/Persister/StateMachine)
internal/kv           # replicated state machine + client-session dedup — later
internal/transport    # Transport impls: sim-net (tests), gRPC (demo) — sim-net next
internal/storage      # Persister impls: file-backed, optional bbolt/Pebble later
internal/clock        # Clock: real + simulated (virtual time)  [DONE]
web/                  # React dashboard — later stage
test/                 # cluster harness + fault-injection / linearizability tests
book/                 # "Building a Distributed KV Store on Raft" — the project journey (living)
go-for-c-programmers/ # companion Go book (28 chapters)
rust-for-c-programmers/ # companion Rust book (30 chapters)
```

## Build / test / verify

```sh
go build ./...
go vet ./...
gofmt -l .            # must print nothing (canonical formatting)
go test -race ./...   # the gold standard; a race is a bug even if tests pass
```

Go toolchain: 1.26.x. Node 24 / npm 11 for the dashboard later.

## Current status

**Stage 0 (scaffolding) — in progress.** Done: core types & RPC payloads
(Figure 2), the `Transport`/`Persister`/`StateMachine` interfaces, and the
injectable `Clock` (real + deterministic `SimClock` with manually-advanced
virtual time). All builds and passes `-race`.

**Next:** the in-memory simulated network (`internal/transport`) that can delay /
drop / duplicate / reorder / partition messages — the single fault-injection
mechanism used by both tests and (later) the dashboard's partition button. Then
the cluster test harness, then Stage 1 (leader election).

There is an open comprehension question with the owner about `SimClock.Advance`
(why it releases the lock before sending on timer channels).

## Raft conventions (non-negotiable)

- **Figure 2 of the Raft paper is the contract.** When code and Figure 2 disagree,
  Figure 2 wins.
- **Locking model:** one `sync.Mutex` per Raft node guarding all mutable state,
  plus long-lived `ticker()` and `applier()` goroutines. Never hold the lock
  across a `transport.*` call or an `applyCh` send: snapshot under lock, release,
  call, re-acquire, re-validate term/role.
- **Persist** currentTerm/votedFor/log before replying to any RPC that mutated
  them. Save Raft state + snapshot together.
- **Commit rule:** a leader advances commitIndex only to an `N` whose entry is
  from the **current** term (Figure 8). Never commit a prior-term entry by replica
  count alone.
- **Election restriction:** compare last-log **term first, then index**. It's a
  safety mechanism, not an optimization.
- **No Raft library.** Std lib only (`net`, `sync`, `context`, `time`,
  `encoding/gob`), plus gRPC + Porcupine later. Do NOT import hashicorp/raft or
  etcd/raft.
- The canonical bug checklist lives in `book/appendix-bug-checklist.md` — consult
  it when a test fails.

## The books

- `book/` is a living textbook that grows with the project (theory + every code
  file explained + trade-offs). **Markdown is the source of truth.** Render with
  `cd book && make epub|html|pdf` (pandoc installed; PDF needs `brew install
  tectonic`). Append a chapter each stage.
- `go-for-c-programmers/` and `rust-for-c-programmers/` are companion language
  books. Their `node_modules/`, `target/`, and `mermaid-filter.err` are gitignored;
  their rendered PDFs are tracked explicitly.

## Git / identity

- Remote: `git@github.com:AndyWang-v/RaftKV.git` (the owner's **personal** GitHub).
- Commit identity is set **locally** to `AndyWang-v <erminwang1@gmail.com>` (this
  overrides the global work identity). Verify with `git config user.email` before
  committing.
- **Do not add any AI/Devin attribution** to commits or files (no `Co-Authored-By`
  trailers, no "with Devin as guide" author lines). The owner wants clean personal
  authorship.
- `.devin/config.local.json` is gitignored; `.devin/skills/` (this skill) IS tracked.
