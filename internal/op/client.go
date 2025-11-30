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
	RunOpCmdWithStdin(stdin []byte, args ...string) error
}

// DefaultClient executes real op CLI commands.
var DefaultClient Client = commandClient{}

type commandClient struct{}

func (commandClient) RunOpCmd(args ...string) ([]byte, error) {
	return runOpCmdInternal(nil, args...)
}

func (commandClient) RunOpCmdWithStdin(stdin []byte, args ...string) error {
	_, err := runOpCmdInternal(stdin, args...)
	return err
}

// runOpCmdInternal executes an op CLI command with optional stdin and returns stdout bytes.
// Shared by RunOpCmd and RunOpCmdWithStdin to keep behavior consistent.
func runOpCmdInternal(stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command("op", args...)

	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("op command failed: %s\nstderr: %s", err, stderr.String())
		}
		return nil, fmt.Errorf("op command failed: %w", err)
	}

	return stdout.Bytes(), nil
}

// RunOpCmd executes an op CLI command with the given arguments
// and returns the stdout output as bytes.
func RunOpCmd(args ...string) ([]byte, error) {
	return runOpCmdInternal(nil, args...)
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
	_, err := RunOpCmd("whoami")
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

// GetWhoAmI returns the JSON output from op whoami command.
func GetWhoAmI() ([]byte, error) {
	return RunOpCmd("whoami")
}

// RunOpCmdWithStdin executes an op CLI command with JSON data piped via stdin.
// It returns an error if the command fails.
func RunOpCmdWithStdin(stdin []byte, args ...string) error {
	_, err := runOpCmdInternal(stdin, args...)
	return err
}
