# Foundations: State Machines, Consensus, and Terms

> **Learning objectives.** By the end of this chapter you can explain: why
> replication reduces to *agreeing on an ordered log*; why naive approaches to
> that agreement fail; Raft's three subproblems; the three server roles; and why
> a *term* is a logical clock rather than wall-clock time. These ideas are
> load-bearing for everything that follows.

## 1.1 The problem: keeping N copies identical

A single server is easy to reason about but fragile: when it dies, your data and
your service die with it. So we run **N copies** (typically 3 or 5) on different
machines. Now one can crash and the others keep serving.

That introduces the hard problem: **how do you keep all N copies identical** when
clients send writes concurrently, machines crash at arbitrary moments, and the
network delays, drops, duplicates, and reorders messages?

## 1.2 The replicated state machine (RSM)

The classic answer reframes the entire problem, and it is worth internalizing:

> If every copy starts in the same state and applies the **exact same sequence of
> commands in the exact same order**, then every copy ends in the same state.
> Determinism does the rest.

A "state machine" here just means: current state + a command → new state, as a
deterministic function. Our state machine is a key-value map; the commands are
`SET`, `GET`, `DELETE`. Because the function is deterministic, *the inputs fully
determine the outputs* — and the inputs are "the starting state" and "the ordered
list of commands."

So the hard problem "keep N state machines identical" collapses into a crisper one:

> **Get all N servers to agree on a single ordered log of commands.**

That agreement problem is called **consensus**. Raft is a consensus algorithm.
Picture the layering:

```
        Clients
          | "SET x=5", "SET y=3", "DEL x" ...
          v
   +--------------+
   |  Raft layer  |  <- agrees on the ORDER:  [SET x=5][SET y=3][DEL x]
   +------+-------+
          | applies in that order, identically, on every node
          v
   +--------------+  +--------------+  +--------------+
   | KV map node1 |  | KV map node2 |  | KV map node3 |   <- all end identical
   +--------------+  +--------------+  +--------------+
```

The Raft layer owns *what order*. The KV map just executes commands
deterministically. That clean separation is why we keep the KV map trivial: all of
the intellectual weight lives in the consensus layer.

### Why order is non-negotiable

Suppose two nodes apply the same two commands in different orders:

```
Node 1:  SET x=5  then  DEL x   ->  final state: x does not exist
Node 2:  DEL x    then  SET x=5 ->  final state: x = 5
```

Same commands, different order, and now **the two nodes disagree about reality**.
A client asking node 1 hears "x is gone"; a client asking node 2 hears "x = 5". The
replicated state machine is supposed to behave like a *single* machine; now it lies
depending on whom you ask. That is the whole game lost.

Note the subtlety: the failure is **divergent final state**, not "a command had an
unmet prerequisite." `SET x=5` succeeds whether or not `x` already exists. Order is
simply one of the inputs to the deterministic function, and different inputs give
different outputs.

## 1.3 Why the naive approaches fail

It is worth seeing why consensus is genuinely hard, so Raft's machinery feels
earned rather than arbitrary.

- **"Pick a leader and have it tell everyone."** Fine until the leader crashes
  mid-write. Who takes over? How does the successor know what the dead leader had
  already promised to clients? This *is* the core difficulty, not a detail.
- **"Order events by timestamp."** Physical clocks on different machines drift and
  are never perfectly synchronized, so timestamp ordering is simply wrong across
  machines. (This is also why we will refuse to ship "lease reads" as a default
  later — they depend on clock-skew assumptions.)
- **"Two-phase commit."** Blocks forever if the coordinator dies at the wrong
  moment. Not fault-tolerant in the way we need.

## 1.4 Raft's decomposition

Raft's claim to fame, versus its famously brain-bending predecessor Paxos, is that
it is *understandable*: it splits consensus into three subproblems a human can hold
in their head. Our build stages map directly onto them — not by accident.

| Raft subproblem | What it does | Build stage |
|---|---|---|
| **Leader election** | Elect exactly one leader; elect a new one when it dies | Stage 1 |
| **Log replication** | Leader accepts commands and copies them to followers in order | Stage 2 |
| **Safety** | Guarantee a new leader can never erase committed data | woven into 1 and 2 |

The simplifying idea underneath all three: **everything flows through a strong
leader.** Clients talk only to the leader. The leader's log is the source of truth.
Followers are passive — they accept what the leader sends them (after a consistency
check). There is exactly one source of authority at any moment, which is what makes
the algorithm tractable.

## 1.5 The two foundational concepts

### Roles

At any instant, every server is in exactly one of three roles:

```
        +-----------+  times out, starts election   +-----------+
        | Follower  | ----------------------------->  | Candidate |
        +-----------+                                 +-----+-----+
              ^   ^                                          | wins majority
              |   | discovers current leader / higher term  | of votes
   higher term|   +--------------------------------------+  v
              |                                        +-----------+
              +--------------------------------------- |  Leader   |
                              discovers higher term    +-----------+
```

- **Follower** (everyone boots here): passive; responds to leaders and candidates.
  If it hears nothing for a while, it suspects the leader is dead and becomes a
  candidate.
- **Candidate**: trying to get elected; solicits votes.
- **Leader**: handles all client requests, replicates the log, and sends periodic
  **heartbeats** so followers know it is alive.

### Terms: a logical clock

Time in Raft is divided into **terms**, numbered with consecutive integers. Each
term begins with an election.

```
 term 1        term 2     term 3            term 4
+--------+   +--------+  +--------+        +--------------
|election|   |election|  |election|       |election| ...
| + work |   | + work |  |(no win)|       | + work |
+--------+   +--------+  +--------+        +--------------
                          ^ split vote, no leader -> next term
```

A term has **at most one leader** (sometimes zero, if its election fails). The term
number is the project's most important idea after the RSM model: it is a **logical
clock**, a plain counter that only ever increases and that ticks on exactly one
event — *starting a new election*. It has nothing to do with real time; term 5
might last three milliseconds or three hours.

Why a logical clock instead of timestamps? Because the only question Raft ever needs
to answer is **"is this information stale?"**, and "stale" means nothing more than
"from a smaller term." A message stamped term 4 arriving at a node already in term 7
is obviously old news. No clock comparison, no skew, no NTP — just integer
comparison. This yields the single most universal rule in the algorithm:

> **The term rule (applies to every message, both directions).** On seeing a term
> `T > currentTerm`, adopt `T`, clear your vote, revert to Follower. On receiving a
> request with a term `< currentTerm`, reject it as stale and reply with your own
> term (so the stale sender learns it is behind and steps down).

Wall-clock time enters the algorithm in exactly *one* fuzzy place: "I have not heard
from a leader in a while, maybe I should start an election." Even there, exact
timing never affects **correctness**, only **liveness** (how quickly we recover).
Hold onto that distinction; it recurs constantly.

## 1.6 Why at most one leader per term — quorum intersection

This guarantee is the foundation of Raft's safety, and the mechanism is just
counting:

- To become leader for a term, a candidate must collect votes from a **majority**
  of the cluster (e.g., 3 of 5).
- **Each server votes for at most one candidate per term** (it records `votedFor`).
- **Any two majorities of the same set must overlap in at least one member.** In a
  5-node cluster, two groups of 3 must share a node. That shared node would have had
  to vote for *both* candidates in the same term — but it votes only once.
  Contradiction. **Two leaders in one term is impossible.**

That "any two majorities intersect" fact — a **quorum intersection** — is the trick
behind *all* of Raft's safety arguments. We will meet it again in the commit rules.

**Zero leaders** happens on a **split vote**: candidates divide the votes and nobody
reaches a majority. The term then ends leaderless and everyone tries again in the
next term. This is also why odd cluster sizes (3, 5) are preferred — they are harder
to split evenly, and a cluster of 2N+1 tolerates the same N failures as one of
2N+2 while needing fewer machines.

## 1.7 A real limitation: the leader is a bottleneck

Because every write funnels through one leader, a single Raft group's write
throughput is capped by what one machine can push, and the leader does the most
network I/O. This is a genuine cost of Raft's "strong leader" simplicity, not a bug.
Production systems mitigate it in layers:

| Technique | Idea | In this project? |
|---|---|---|
| **Batching** | Bundle many commands into one replication round | Emerges naturally at Stage 2 |
| **Pipelining** | Send the next batch before the previous ack returns | Stage 2+ |
| **Offload reads** | Reads need not go through the log; serve via ReadIndex/lease | Stage 5 |
| **Follower reads** | Let followers serve checked/stale reads | Advanced, optional |
| **Sharding / multi-Raft** | Run many Raft groups, spreading leadership | **Out of scope (by design)** |

That last row is how real systems (CockroachDB, TiKV, Spanner) scale: thousands of
Raft groups, each owning a slice of the keyspace with its own leader, so no single
machine is *the* leader for everything. We deliberately scope it out — one Raft group
is the right amount of hard for *understanding* consensus. But knowing exactly where
the ceiling is, and how you would raise it, is the kind of judgment that separates
"I used a consensus library" from "I understand consensus."

> **Checkpoint.**
>
> *Q1. Why must commands be applied in the same order on every node?* Because the
> state machine is a deterministic function of (starting state, ordered commands).
> Different orders are different inputs, so they produce different final states, and
> the replicas diverge — the cluster stops behaving like one machine.
>
> *Q2. Why at most one leader per term, but possibly zero?* A leader needs a majority
> of votes; each node votes once per term; two majorities must overlap; so two
> leaders in one term is impossible. Zero happens on a split vote where nobody
> reaches a majority.
>
> *Q3. Term vs. wall-clock time?* A term is a logical clock — a monotonic integer
> that ticks only when an election starts, used purely to detect staleness via
> integer comparison. Wall-clock time is unreliable across machines (drift, no sync)
> and is used only as a fuzzy "haven't heard from a leader" hint that affects
> liveness, never correctness.
