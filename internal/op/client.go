package op

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// Client defines the operations needed to execute op CLI commands.
type Client interface {
	RunOpCmd(args ...string) ([]byte, error)
}

// DefaultClient executes real op CLI commands.
var DefaultClient Client = commandClient{}

type commandClient struct{}

func (commandClient) RunOpCmd(args ...string) ([]byte, error) {
	return runOpCmdInternal(args...)
}

// runOpCmdInternal executes an op CLI command and returns stdout bytes.
// Used by RunOpCmd to handle command execution.
func runOpCmdInternal(args ...string) ([]byte, error) {
	cmd := exec.Command("op", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("op command failed: %w\nstderr: %s", err, stderr.String())
		}
		return nil, fmt.Errorf("op command failed: %w", err)
	}

	return stdout.Bytes(), nil
}

// CheckOpInstalled verifies that the op binary exists in PATH.
func CheckOpInstalled() error {
	_, err := exec.LookPath("op")
	if err != nil {
		return fmt.Errorf("1Password CLI (op) not found in PATH. Please install it from: https://developer.1password.com/docs/cli/get-started/")
	}
	return nil
}

// CheckOpSignedIn verifies that the user is authenticated with the op CLI.
func CheckOpSignedIn() error {
	_, err := DefaultClient.RunOpCmd("whoami")
	if err != nil {
		return fmt.Errorf("failed to verify 1Password CLI sign-in; please run 'op signin': %w", err)
	}
	return nil
}

// VerifyOpReady performs both installation and authentication checks.
func VerifyOpReady() error {
	if err := CheckOpInstalled(); err != nil {
		return err
	}
	if err := CheckOpSignedIn(); err != nil {
		return err
	}
	return nil
}
