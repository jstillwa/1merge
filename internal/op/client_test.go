//go:build integration
// +build integration

package op

import (
	"bytes"
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

// TestGetWhoAmI verifies that GetWhoAmI returns non-empty JSON output.
func TestGetWhoAmI(t *testing.T) {
	if err := CheckOpInstalled(); err != nil {
		t.Skipf("Skipping: 1Password CLI (op) not installed or not in PATH: %v", err)
	}

	if err := CheckOpSignedIn(); err != nil {
		t.Skipf("Skipping: 1Password CLI (op) not signed in: %v", err)
	}

	output, err := GetWhoAmI()
	if err != nil {
		t.Fatalf("GetWhoAmI failed: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("GetWhoAmI returned empty output")
	}

	// Verify output contains expected JSON structure
	outputStr := string(output)
	if !bytes.Contains(output, []byte("account")) && !bytes.Contains(output, []byte("user")) {
		t.Logf("GetWhoAmI output: %s", outputStr)
		t.Fatal("GetWhoAmI output does not contain expected JSON fields")
	}
}

// TestRunOpCmdInvalidCommand verifies error handling for invalid commands.
func TestRunOpCmdInvalidCommand(t *testing.T) {
	_, err := RunOpCmd("invalid-command-xyz")
	if err == nil {
		t.Fatal("RunOpCmd should have returned an error for invalid command")
	}

	if len(err.Error()) == 0 {
		t.Fatal("Error message should not be empty")
	}
}
