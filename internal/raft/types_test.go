package raft

import "testing"

// TestRoleString is a table-driven test: each case is one {input, want} row, and
// the loop runs them all. This is the idiomatic Go testing pattern we'll use
// throughout the project. t.Run names each subtest so a failure points at the
// exact row.
func TestRoleString(t *testing.T) {
	cases := []struct {
		name string
		role Role
		want string
	}{
		{"follower", Follower, "Follower"},
		{"candidate", Candidate, "Candidate"},
		{"leader", Leader, "Leader"},
		{"unknown", Role(99), "Unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.role.String(); got != tc.want {
				t.Errorf("Role(%d).String() = %q, want %q", tc.role, got, tc.want)
			}
		})
	}
}
