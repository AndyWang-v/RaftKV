# Appendix B: The Canonical Raft Bug Checklist

Consult this list whenever a test fails. These are the mistakes everyone hits,
roughly ordered by how often they bite. Each one is also a lesson in *why* Raft is
shaped the way it is — debugging them is the curriculum, not a chore.

1. **Lock held across an RPC or `applyCh` send → deadlock.** Snapshot the state you
   need under the lock, release the lock, make the call, re-acquire, then re-validate
   (term and role may have changed).

2. **Election timer reset on the wrong events → livelock or no leader.** Reset only
   when (a) you grant a vote, or (b) you receive a current-term leader's
   `AppendEntries`. Never on a rejected `RequestVote` or a stale `AppendEntries`.

3. **Missing the up-to-date check in `RequestVote` → committed entries overwritten.**
   Compare last-log *term* first, then index.

4. **Committing a previous-term entry by replica count → Figure 8 safety violation.**
   A leader may only advance `commitIndex` to an entry from its **current** term;
   older entries commit transitively underneath it. Dropping the
   `log[N].Term == currentTerm` clause is rare, catastrophic, and the bug most often
   introduced by accident.

5. **Not persisting before replying → data loss on crash.** Persist `currentTerm`,
   `votedFor`, and `log` before granting a vote or acknowledging an append.

6. **Applier not monotonic / applies under lock / out of order → divergent or
   duplicated state.** One applier goroutine; apply in index order, exactly once;
   release the lock while delivering on `applyCh`.

7. **Stale RPC replies applied → corruption.** After unlocking for an RPC, re-check
   that term and role are unchanged before acting on the reply; step down on any
   higher term observed.

8. **Blind log truncation on `AppendEntries` → drops valid later entries.** Truncate
   only on a real term conflict (same index, different term); do not truncate entries
   that already match.

9. **`nextIndex` / `matchIndex` mis-initialized on becoming leader → backtracking or
   false commits.** Set `nextIndex = lastLogIndex + 1`, `matchIndex = 0`.

10. **Index/offset confusion after snapshotting → panics or wrong entries.** Route
    every absolute-index ↔ slice-offset translation through the `logStore`.

11. **Sending `AppendEntries` for already-snapshotted indices.** Send
    `InstallSnapshot` instead when `nextIndex[f]` is at or below the snapshot index.

12. **Forgetting the leader's no-op on election → inherited entries never commit;
    early ReadIndex reads unsafe.** A new leader appends a no-op entry in its term to
    commit inherited prior-term entries transitively.
