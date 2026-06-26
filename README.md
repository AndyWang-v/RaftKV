# RaftKV

A linearizable, fault-tolerant **key-value store backed by a from-scratch
implementation of the Raft consensus algorithm**, written in Go — built stage by
stage as a deep learning project, not assembled from an existing library.

The key-value map is deliberately trivial; it exists to make the consensus layer
*observable*. The Raft layer — leader election, log replication, persistence,
snapshotting, and linearizable client semantics — is the heart of the project. A
React dashboard (a later stage) will visualize the live cluster: the current
leader, log replication in real time, and a button to inject a partition and watch
re-election.

## Status

**Stage 0 — scaffolding (in progress).** Done so far:

- Core data types and RPC payloads (Raft's Figure 2).
- The seam interfaces: `Transport`, `Persister`, `StateMachine`.
- An injectable `Clock` — a real implementation and a deterministic simulated one
  (virtual time), so tests run thousands of fault-injection scenarios in seconds.

All packages build and pass `go test -race`. Next: the in-memory simulated network
(delay / drop / duplicate / reorder / partition) and the cluster test harness.

## Repository layout

```
cmd/                 # binaries: kvnode (server), kvctl (CLI client)
internal/raft        # the consensus core
internal/kv          # the replicated state machine + client-session dedup
internal/transport   # Transport implementations (sim-net, later gRPC)
internal/storage     # Persister implementations (file-backed, later embedded engine)
internal/clock       # Clock: real + simulated (virtual time)
web/                 # React dashboard (later stage)
test/                # cluster harness + fault-injection / linearizability tests
book/                # "Building a Distributed KV Store on Raft" — the project's journey
go-for-c-programmers/# "Go for C Programmers" — a companion book learning the language
```

## The two books

This repo doubles as a learning archive:

- **`book/`** — a living textbook that grows with the project: the Raft theory, the
  design decisions, and every piece of code, explained. Start at
  [`book/00-preface.md`](book/00-preface.md).
- **`go-for-c-programmers/`** — a 28-chapter book teaching Go from a C programmer's
  perspective, with runnable examples. Start at
  [`go-for-c-programmers/README.md`](go-for-c-programmers/README.md).

## Building and testing

```sh
go build ./...
go vet ./...
go test -race ./...
```

## Scope (non-goals, by design)

- One Raft group — no sharding / multi-Raft.
- Crash-stop failures only — no Byzantine fault tolerance.
- File-backed persistence is sufficient (with an optional embedded-engine upgrade).

## References

The canonical Raft paper (Ongaro & Ousterhout) and other sources are listed in
[`book/references.md`](book/references.md).
