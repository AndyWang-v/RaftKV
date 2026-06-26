# References

Read these before coding, and re-read them while debugging. Treat Figure 2 of the
Raft paper as the contract: when our code and Figure 2 disagree, Figure 2 wins.

## Primary sources

- **The Raft paper** — Diego Ongaro & John Ousterhout, *"In Search of an
  Understandable Consensus Algorithm (Extended Version)."* Figure 2 is the entire
  RPC specification; Figure 8 is the subtle safety case that breaks naive
  implementations. <https://raft.github.io/raft.pdf>

- **Ongaro's PhD dissertation** — *"Consensus: Bridging Theory and Practice."* The
  authoritative source for snapshotting and membership changes.

## Learning aids

- **The "Students' Guide to Raft"** (MIT PDOS / 6.5840 course material) — a list of
  the exact bugs everyone hits; distilled in Appendix B of this book.

- **Interactive visualizations** — watch leader election and replication animate
  before writing code:
  - <https://thesecretlivesofdata.com/raft/>
  - <https://raft.github.io/>

- **MIT 6.5840 (formerly 6.824), labs 2–4** — the canonical implementation
  progression; our build stages mirror it.

## Tools

- **Porcupine** — a Go linearizability checker, used to *prove* the KV store is
  linearizable. <https://github.com/anishathalye/porcupine>

## On not using a Raft library

The standard library does almost everything we need (`net`, `net/http`, `sync`,
`context`, `time`, `encoding/gob`). For the demo we add gRPC and protobuf, and
Porcupine for linearizability testing. We deliberately do **not** import
`hashicorp/raft` or `etcd/raft` — the entire point is to build consensus from
scratch. Reading those implementations *after* finishing is, however, excellent.
