package raft

import "context"

// This file defines the seams between Raft and the outside world (design doc §4).
// These are interfaces, not concrete types, so the same Raft code runs against a
// deterministic in-memory simulator under test and against real network/disk in
// the demo. Idiomatic Go defines an interface in the package that CONSUMES it
// (here, raft), while the implementations live in their own packages.

// Transport delivers Raft RPCs to peers.
//
// Contract (design doc §6, non-negotiable):
//   - All calls are synchronous request/response.
//   - The CALLER owns timeouts (via ctx) and retries.
//   - The caller MUST NOT hold the Raft mutex while calling any of these. Doing
//     so deadlocks the cluster: a blocked RPC would pin the lock the peer's
//     reply handler needs. Snapshot what you need under the lock, release it,
//     make the call, re-acquire, then re-validate term/role.
//
// Implementations: simnet (deterministic tests) and gRPC/TCP (the demo).
type Transport interface {
	RequestVote(ctx context.Context, peer NodeID, args *RequestVoteArgs) (*RequestVoteReply, error)
	AppendEntries(ctx context.Context, peer NodeID, args *AppendEntriesArgs) (*AppendEntriesReply, error)
	InstallSnapshot(ctx context.Context, peer NodeID, args *InstallSnapshotArgs) (*InstallSnapshotReply, error)
}

// Persister stores the durable Raft state.
//
// Contract (design doc §7.6):
//   - Save MUST be durable (fsync) before the RPC handler that triggered it
//     replies. Persisting after replying loses data on a crash.
//   - Raft state and snapshot are saved TOGETHER so a crash can never leave the
//     two disagreeing (e.g. a log that references a snapshot that wasn't written).
//
// Implementations: in-memory (tests), file-backed (demo), and optionally a
// bbolt/Pebble-backed engine later (roadmap).
type Persister interface {
	// Save atomically persists state and snapshot. A crash leaves either the old
	// pair or the new pair, never a mix.
	Save(state RaftState, snapshot []byte) error
	// Load returns the persisted state and snapshot, or zero values on first boot.
	Load() (RaftState, []byte, error)
	// RaftStateSize reports the serialized size of the persisted Raft state, used
	// to decide when to snapshot (Stage 4).
	RaftStateSize() int
}

// StateMachine is the application that Raft replicates — for us, the KV map.
//
// Contract (design doc §4):
//   - Apply is invoked strictly in log order, exactly once per committed entry.
//   - Apply MUST be deterministic: same command sequence => same state on every
//     node. This determinism is the entire point of the replicated state machine.
//   - Snapshot/Restore serialize and reload the full state for log compaction.
type StateMachine interface {
	Apply(cmd []byte) (result []byte)
	Snapshot() []byte
	Restore(snapshot []byte)
}
