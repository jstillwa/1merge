//go:build integration
// +build integration

package op

import (
	"testing"
)

// Note: These tests require the op CLI to be installed and the user to be
// signed into 1Password. Skip these tests in CI environments if needed.

// TestCheckOpInstalled verifies that CheckOpInstalled returns nil when op is in PATH.
func TestCheckOpInstalled(t *testing.T) {
	if err := CheckOpInstalled(); err != nil {
		t.Skipf("Skipping: 1Password CLI (op) not installed or not in PATH: %v", err)
	}
}

// TestRunOpCmdInvalidCommand verifies error handling for invalid commands.
func TestRunOpCmdInvalidCommand(t *testing.T) {
	_, err := DefaultClient.RunOpCmd("invalid-command-xyz")
	if err == nil {
		t.Fatal("DefaultClient.RunOpCmd should have returned an error for invalid command")
	}

	if len(err.Error()) == 0 {
		t.Fatal("Error message should not be empty")
	}
}
