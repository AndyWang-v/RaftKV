# Leader Election

> **Learning objectives.** Explain how a follower detects a dead leader; trace the
> `RequestVote` RPC and its receiver logic; explain why randomized timeouts prevent
> livelock; and articulate the two subtle rules where correctness lives — the
> election-timer reset rule and the up-to-date (election-restriction) safety check.

This is the first Raft subproblem and the first thing we will implement (Stage 1).
The goal: start with everyone a Follower and no leader, and reliably converge to
exactly one leader — and re-converge whenever the leader dies.

## 2.1 Detecting a dead leader: the election timer

A follower cannot reliably *ask* "are you alive?" — the network might just be slow.
So Raft uses a timeout, the one place wall-clock time enters the algorithm:

- Each follower runs an **election timer** (a countdown, e.g. 300 ms).
- Every time it hears a valid heartbeat (an `AppendEntries` from the current leader)
  it **resets** the timer.
- If the timer ever reaches zero — "I have not heard from a leader in 300 ms" — the
  follower assumes the leader is dead and starts an election.

```
Leader alive, heartbeats every 100 ms:
  --H------H------H------H------H--    follower keeps resetting, stays Follower
    reset  reset  reset  reset  reset

Leader dies:
  --H------H-----------X  (silence)
    reset  reset       |
                       +-- timer counts down... 300 ms... -> becomes CANDIDATE
```

This timeout governs **liveness, not safety**. If it fires while the leader is
actually alive (timer too short, or a slow network), nothing is corrupted — we just
trigger an unnecessary election that costs a little time. Keep that reassurance in
mind; it is why exact timing is "tune it," not "prove it."

## 2.2 Becoming a candidate

When a Follower (or a Candidate whose election stalled) times out, it does the
following, in order:

1. **Increment `currentTerm`** (new election, new term).
2. **Transition to Candidate.**
3. **Vote for itself** (`votedFor = self`).
4. **Persist** `currentTerm` and `votedFor` (they must survive a crash — see the
   persistence chapter for why).
5. **Reset its own election timer** (so a failed election will retry).
6. **Send `RequestVote` RPCs to all other servers, in parallel.**

## 2.3 The `RequestVote` RPC

The payload (this is Figure 2 of the paper):

```go
type RequestVoteArgs struct {
    Term         int     // candidate's term
    CandidateID  NodeID  // who is asking
    LastLogIndex int     // index of candidate's last log entry   (safety)
    LastLogTerm  int     // term  of candidate's last log entry   (safety)
}
type RequestVoteReply struct {
    Term        int     // voter's currentTerm, so a stale candidate steps down
    VoteGranted bool
}
```

The receiver logic, annotated line by line:

```
reply.Term = currentTerm

// (a) Stale candidate? Reject.
if args.Term < currentTerm:
    reply.VoteGranted = false; return

// (b) Higher term seen? Step down first (the universal term rule).
if args.Term > currentTerm:
    currentTerm = args.Term; votedFor = None; become Follower; persist

// (c) The up-to-date check (SAFETY — see 2.6).
upToDate = (args.LastLogTerm > myLastLogTerm) ||
           (args.LastLogTerm == myLastLogTerm && args.LastLogIndex >= myLastLogIndex)

// (d) Grant only if we have not already voted elsewhere this term AND the
//     candidate's log is at least as up-to-date as ours.
if (votedFor == None || votedFor == args.CandidateID) && upToDate:
    votedFor = args.CandidateID
    persist
    reset election timer
    reply.VoteGranted = true
else:
    reply.VoteGranted = false
```

The `votedFor == args.CandidateID` clause makes a re-delivered (duplicate)
`RequestVote` idempotent: voting "again" for the same candidate is harmless.

## 2.4 The three ways an election ends

After sending vote requests, exactly one of three things happens:

1. **It wins** — a **majority** grant their votes (counting its own). It becomes
   **Leader** and immediately heartbeats so everyone else resets and stops trying.
2. **Someone else wins** — it receives an `AppendEntries` from a server with term ≥
   its own. That is a legitimate leader, so it steps down to **Follower**.
3. **Nobody wins** (split vote) — it gets neither a majority nor a valid leader, and
   its own timer fires again, so it starts a fresh election in a higher term.

## 2.5 Split votes and randomized timeouts

Outcome 3 is dangerous if it repeats forever. Imagine all servers boot with
*identical* 300 ms timers: they all time out together, all become candidates in the
same term, all vote for themselves, and split the vote. No majority. They retry —
and time out together *again*. That is **livelock**: endless elections, no leader.

Raft's fix is delightfully cheap: **randomize the election timeout**, re-randomized
each election (e.g. uniform in [300 ms, 600 ms]).

```
Server A: 312 ms -> times out FIRST, asks for votes
Server B: 487 ms -> still waiting... receives A's RequestVote -> votes for A
Server C: 551 ms -> still waiting... receives A's RequestVote -> votes for A
                    A reaches majority before B or C even wakes up. Clean win.
```

Usually one server's timer fires far enough ahead that it wins before the others
stir. If a split happens anyway, the next round re-randomizes and self-corrects.

## 2.6 The two rules where correctness lives

These are the rules an implementation most often gets subtly wrong. Memorize them.

### Rule A — election-timer reset

You may reset your election timer in **only** two situations:

- you **grant** a vote to a candidate, and
- you receive a **valid `AppendEntries` from the current leader** (term ≥ yours).

You must **not** reset it when you **reject** a `RequestVote`, nor on a **stale**
`AppendEntries` from an old-term leader.

*Why.* The timer is your only mechanism for ejecting a bad leader or out-voting a
bad candidate. If you reset it on stale or rejected traffic, you may never time out
to start a legitimate election → livelock or no leader. The timer must be quieted
only by signs of *legitimate current authority*.

### Rule B — the up-to-date restriction (a safety property, not an optimization)

A candidate may win **only if its log is at least as up-to-date as the voter's**,
defined precisely: compare **last-log term first** (higher wins); only on a tie,
compare **last-log index** (longer wins).

*Why term before index?* Consider two logs:

```
Log A:  [t1][t1][t1][t1][t1]   <- 5 entries, all term 1.  last: index 5, term 1
Log B:  [t1][t2]               <- 2 entries.              last: index 2, term 2
```

Comparing **index first** would let A (length 5) beat B (length 2) — a disaster.
B's last entry is from term 2, which means a term-2 leader was elected *after* term
1 ended, which requires a majority to have moved on to term 2. So A's extra term-1
entries are **uncommitted garbage** from a term-1 leader that crashed before
committing them (had they been committed, B could not have diverged). A higher
last-log term means *your log has seen more recent leadership* — more authoritative
than merely being longer. Length breaks ties only within the same term. Comparing
index-first would let a node full of never-committed entries win and **erase
committed data**. Term first, then index — for safety. (We prove this connects to
the *Leader Completeness* property — "an elected leader holds every committed entry"
— in the log-replication chapter.)

## 2.7 The disruptive server (why Rule A is not hypothetical)

When does a candidate keep *spamming* `RequestVote`? Concretely: a node partitioned
into a minority can never reach a majority, so it never wins; its timer keeps firing,
incrementing its term again and again — climbing to term 50 while the healthy
majority runs happily under a leader at term 7. When the partition heals, the node
rejoins blasting `RequestVote` at term 50. Two things must hold:

1. The healthy nodes adopt term 50 (the universal term rule) — annoying but safe.
2. Critically, the rejoiner's **log is behind** (it missed terms 7–49), so Rule B
   **rejects** its candidacy. Safety holds.

But if the healthy followers reset their timers on that *rejected* `RequestVote`,
they would never time out to elect a real leader while the disruptor keeps spamming
— the cluster wedges. Hence Rule A: rejection must not reset the timer. (Raft's
optional **PreVote** extension prevents the term inflation in the first place; we do
not require it, but now you know what it is for.)

## 2.8 Timing constants

A sound starting point (tune per transport):

- `heartbeatInterval = 100 ms`
- election timeout = random in [300 ms, 600 ms], re-randomized each election
- Invariant: `heartbeatInterval` << election timeout << mean-time-between-failures,
  and broadcast time << election timeout.

## 2.9 Stage 1: definition of done

When we implement this, Stage 1 is complete only when, under `go test -race`:

- exactly one leader is elected in a quiet cluster;
- **at most one leader exists per term, ever** (asserted continuously by the harness
  — this is the quorum-intersection guarantee);
- a new leader is elected within a bounded time after the leader is partitioned away;
- an old leader **steps down** when it rejoins and sees a higher term.

> **Checkpoint.**
>
> *Q1. Why do randomized timeouts prevent livelock?* Identical fixed timeouts make
> everyone time out together, self-vote, and split the vote forever. Randomization
> staggers wakeups so one candidate usually grabs a majority before others stir.
>
> *Q2. A candidate sends term 5; a follower is at term 8. What happens?* The follower
> rejects (5 < 8), replies with term 8, and does **not** reset its timer. Seeing
> term 8, the stale candidate adopts it and steps down to Follower.
>
> *Q3. Why must a follower not reset its timer on a rejected `RequestVote`?* Because a
> disruptive (e.g. partitioned, log-behind) node can spam `RequestVote` forever; if
> rejections reset timers, no healthy node ever times out to elect a real leader.
>
> *Q4. Why compare last-log term before index?* A higher last-log term reflects more
> recent leadership; extra entries from an older term may be uncommitted garbage.
> Comparing index-first could let a stale-but-long log win and erase committed data.
