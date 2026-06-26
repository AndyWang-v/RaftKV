// Package raft implements the Raft consensus algorithm from scratch.
//
// This file defines the core data types (design doc §5). It is pure data:
// no behavior, no goroutines, no locking. The RPC payload structs below ARE
// Figure 2 of the Raft paper ("In Search of an Understandable Consensus
// Algorithm"). When the implementation and Figure 2 disagree, Figure 2 wins.
package raft

// NodeID identifies one server in the cluster. It is a distinct type (not a
// bare int) so the compiler stops us from accidentally mixing a node id with,
// say, a log index or a count.
type NodeID int

// None is the sentinel meaning "no node" — used for votedFor when we have not
// voted in the current term, and for leaderId when we don't know the leader.
const None NodeID = -1

// Role is the server's current role in the Raft state machine. Every server is
// in exactly one role at a time (paper §5.1). The zero value is Follower, which
// is intentional: a freshly constructed server boots as a Follower.
type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// String makes Role satisfy the fmt.Stringer interface, so %v / %s prints
// "Follower" instead of "0". Purely for readable logs and test failures.
func (r Role) String() string {
	switch r {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// LogEntry is one command in the replicated log.
//
// Term is the leader's term when the entry was created; it is the key to all of
// Raft's consistency checks (the Log Matching property: same index + same term
// => identical entry and identical history before it).
//
// Command is opaque to Raft — Raft only orders and replicates these bytes; it
// never interprets them. Only the StateMachine (the KV layer) decodes them.
// This opacity is the seam that keeps consensus and application logic separate.
type LogEntry struct {
	Term    int
	Command []byte
}

// --- RPC payloads: this section IS Figure 2 ---

// RequestVoteArgs is sent by a Candidate soliciting votes (paper Figure 2).
type RequestVoteArgs struct {
	Term         int    // candidate's term
	CandidateID  NodeID // candidate requesting the vote
	LastLogIndex int    // index of candidate's last log entry  (election restriction)
	LastLogTerm  int    // term of candidate's last log entry   (election restriction)
}

// RequestVoteReply is the voter's response.
type RequestVoteReply struct {
	Term        int  // voter's currentTerm, so a stale candidate can step down
	VoteGranted bool // true means the candidate received this vote
}

// AppendEntriesArgs is sent by the Leader to replicate log entries and as a
// heartbeat when Entries is empty (paper Figure 2).
type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderID     NodeID     // so a follower can redirect clients to the leader
	PrevLogIndex int        // index of the entry immediately preceding Entries
	PrevLogTerm  int        // term of the PrevLogIndex entry (consistency check)
	Entries      []LogEntry // new entries to store (empty for heartbeat)
	LeaderCommit int        // leader's commitIndex
}

// AppendEntriesReply is the follower's response.
type AppendEntriesReply struct {
	Term    int  // follower's currentTerm, for the leader to update itself
	Success bool // true if follower contained an entry matching PrevLogIndex/Term

	// Fast-backup hints. Not in Figure 2, but a standard optimization so the
	// leader can skip many indices at once instead of decrementing nextIndex
	// one entry per round trip. Populated only when Success is false.
	ConflictTerm  int // term of the conflicting entry, or -1 if none at PrevLogIndex
	ConflictIndex int // first index the leader should try next
}

// InstallSnapshotArgs is sent when a follower is so far behind that the leader
// has already discarded (snapshotted away) the entries the follower needs
// (paper Figure 13). v1 ships the whole snapshot in one message (no chunking).
type InstallSnapshotArgs struct {
	Term              int
	LeaderID          NodeID
	LastIncludedIndex int    // snapshot replaces all entries up through this index
	LastIncludedTerm  int    // term of LastIncludedIndex
	Data              []byte // raw serialized snapshot bytes (opaque to Raft)
}

// InstallSnapshotReply lets the leader detect it is stale (reply.Term > its own).
type InstallSnapshotReply struct {
	Term int
}

// ApplyMsg is how Raft hands committed work up to the application. Raft's
// applier goroutine sends these on applyCh, strictly in log order, exactly once
// per committed entry. Either CommandValid OR SnapshotValid is true, never both.
type ApplyMsg struct {
	// A committed command to apply to the state machine.
	CommandValid bool
	Command      []byte
	CommandIndex int

	// A snapshot to install (used in Stage 4). When the leader catches a far-
	// behind follower up via InstallSnapshot, the follower replays it here.
	SnapshotValid bool
	Snapshot      []byte
	SnapshotIndex int
}

// RaftState is the durable state that MUST survive a crash (design doc §4, §7.6).
// It is persisted together with the snapshot so the two can never disagree after
// a crash. currentTerm/votedFor/log are the three persistent fields from
// Figure 2; LastIncluded{Index,Term} are snapshot metadata added in Stage 4.
type RaftState struct {
	CurrentTerm       int
	VotedFor          NodeID // None (-1) if we have not voted this term
	Log               []LogEntry
	LastIncludedIndex int
	LastIncludedTerm  int
}
