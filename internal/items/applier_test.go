package items

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"1merge/internal/models"
	"1merge/internal/op"
)

// Helper function to create test items with common properties
func createTestItem(id, title string, updatedAt time.Time) models.Item {
	return models.Item{
		ID:        id,
		Title:     title,
		UpdatedAt: updatedAt,
		Vault: models.Vault{
			ID:   "test_vault",
			Name: "Test Vault",
		},
		Category: "LOGIN",
		Fields: []models.Field{
			{
				ID:    "username",
				Type:  "username",
				Label: "username",
				Value: "testuser",
			},
			{
				ID:    "password",
				Type:  "password",
				Label: "password",
				Value: "testpass",
			},
		},
		URLs: []models.URL{
			{
				HRef:    "https://example.com",
				Primary: true,
			},
		},
	}
}

type opCall struct {
	name  string
	args  []string
	stdin []byte
}

type fakeOpClient struct {
	calls                []opCall
	runOpCmdErr          error
	runOpCmdWithStdinErr error
}

func (f *fakeOpClient) RunOpCmd(args ...string) ([]byte, error) {
	f.calls = append(f.calls, opCall{name: "RunOpCmd", args: args})
	if f.runOpCmdErr != nil {
		return nil, f.runOpCmdErr
	}
	return nil, nil
}

func (f *fakeOpClient) RunOpCmdWithStdin(stdin []byte, args ...string) error {
	f.calls = append(f.calls, opCall{name: "RunOpCmdWithStdin", args: args, stdin: stdin})
	if f.runOpCmdWithStdinErr != nil {
		return f.runOpCmdWithStdinErr
	}
	return nil
}

func TestApplyMerge_DryRun(t *testing.T) {
	tests := []struct {
		name      string
		winner    models.Item
		losers    []models.Item
		dryRun    bool
		shouldErr bool
	}{
		{
			name: "dry-run with single loser",
			winner: createTestItem(
				"winner1",
				"Winner Item",
				time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			),
			losers: []models.Item{
				createTestItem(
					"loser1",
					"Loser Item",
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				),
			},
			dryRun:    true,
			shouldErr: false,
		},
		{
			name: "dry-run with multiple losers",
			winner: createTestItem(
				"winner1",
				"Winner Item",
				time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			),
			losers: []models.Item{
				createTestItem(
					"loser1",
					"Loser Item 1",
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				),
				createTestItem(
					"loser2",
					"Loser Item 2",
					time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				),
			},
			dryRun:    true,
			shouldErr: false,
		},
		{
			name: "dry-run with no losers",
			winner: createTestItem(
				"winner1",
				"Winner Item",
				time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			),
			losers:    []models.Item{},
			dryRun:    true,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			err := ApplyMerge(tt.winner, tt.losers, tt.dryRun)

			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.shouldErr && err == nil {
				t.Errorf("ApplyMerge() expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ApplyMerge() unexpected error: %v", err)
			}

			if tt.dryRun && len(output) == 0 {
				t.Errorf("ApplyMerge() dry-run should produce output, got none")
			}

			// Verify dry-run messages are present
			if tt.dryRun {
				expectedMsg := fmt.Sprintf("[DRY RUN] Would edit item: %s", tt.winner.ID)
				if !bytes.Contains(buf.Bytes(), []byte(expectedMsg)) {
					t.Errorf("ApplyMerge() output missing dry-run edit message: %s", expectedMsg)
				}

				for _, loser := range tt.losers {
					expectedArchiveMsg := fmt.Sprintf("[DRY RUN] Would archive item: %s", loser.ID)
					if !bytes.Contains(buf.Bytes(), []byte(expectedArchiveMsg)) {
						t.Errorf("ApplyMerge() output missing dry-run archive message: %s", expectedArchiveMsg)
					}
				}
			}
		})
	}
}

func assertOpCalls(t *testing.T, got []opCall, expected []opCall) {
	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("expected %d op calls, got %d", len(expected), len(got))
	}

	for i := range expected {
		if got[i].name != expected[i].name {
			t.Fatalf("call %d: expected name %s, got %s", i, expected[i].name, got[i].name)
		}
		if !slices.Equal(got[i].args, expected[i].args) {
			t.Fatalf("call %d: expected args %v, got %v", i, expected[i].args, got[i].args)
		}
	}
}

func TestApplyMerge_NonDryRun_UsesOpClient(t *testing.T) {
	winner := createTestItem(
		"winner1",
		"Winner Item",
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)
	losers := []models.Item{
		createTestItem(
			"loser1",
			"Loser Item 1",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		),
		createTestItem(
			"loser2",
			"Loser Item 2",
			time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		),
	}

	tests := []struct {
		name                 string
		client               *fakeOpClient
		expectedCalls        []opCall
		expectErr            bool
		expectedErrSubstring string
	}{
		{
			name:   "successful merge executes edit and archive commands",
			client: &fakeOpClient{},
			expectedCalls: []opCall{
				{name: "RunOpCmdWithStdin", args: []string{"item", "edit", winner.ID}},
				{name: "RunOpCmd", args: []string{"item", "delete", losers[0].ID, "--archive"}},
				{name: "RunOpCmd", args: []string{"item", "delete", losers[1].ID, "--archive"}},
			},
		},
		{
			name:   "edit failure propagates error and stops deletes",
			client: &fakeOpClient{runOpCmdWithStdinErr: errors.New("edit failed")},
			expectedCalls: []opCall{
				{name: "RunOpCmdWithStdin", args: []string{"item", "edit", winner.ID}},
			},
			expectErr:            true,
			expectedErrSubstring: "failed to edit item",
		},
		{
			name:   "archive failure propagates error and halts further deletes",
			client: &fakeOpClient{runOpCmdErr: errors.New("archive failed")},
			expectedCalls: []opCall{
				{name: "RunOpCmdWithStdin", args: []string{"item", "edit", winner.ID}},
				{name: "RunOpCmd", args: []string{"item", "delete", losers[0].ID, "--archive"}},
			},
			expectErr:            true,
			expectedErrSubstring: "failed to archive item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetOpClient(tt.client)
			t.Cleanup(func() { SetOpClient(op.DefaultClient) })

			err := ApplyMerge(winner, losers, false)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("ApplyMerge() expected error, got nil")
				}
				if tt.expectedErrSubstring != "" && !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Fatalf("ApplyMerge() error %q does not contain expected substring %q", err, tt.expectedErrSubstring)
				}
			} else if err != nil {
				t.Fatalf("ApplyMerge() unexpected error: %v", err)
			}

			assertOpCalls(t, tt.client.calls, tt.expectedCalls)
			for _, call := range tt.client.calls {
				if call.name == "RunOpCmdWithStdin" && len(call.stdin) == 0 {
					t.Fatalf("RunOpCmdWithStdin should be provided JSON data")
				}
			}
		})
	}
}

func TestApplyMerge_EmptyLosers(t *testing.T) {
	tests := []struct {
		name      string
		winner    models.Item
		losers    []models.Item
		shouldErr bool
	}{
		{
			name: "dry-run with empty losers",
			winner: createTestItem(
				"winner1",
				"Winner Item",
				time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			),
			losers:    []models.Item{},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			err := ApplyMerge(tt.winner, tt.losers, true)

			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)

			if tt.shouldErr && err == nil {
				t.Errorf("ApplyMerge() expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ApplyMerge() unexpected error: %v", err)
			}

			// Verify dry-run edit message is present
			expectedMsg := fmt.Sprintf("[DRY RUN] Would edit item: %s", tt.winner.ID)
			if !bytes.Contains(buf.Bytes(), []byte(expectedMsg)) {
				t.Errorf("ApplyMerge() output missing expected message: %s", expectedMsg)
			}
		})
	}
}

func TestApplyMerge_MarshalError(t *testing.T) {
	// Create a test item with all nil/empty fields to ensure marshaling works
	winner := models.Item{
		ID:    "test1",
		Title: "Test",
	}

	// This test verifies that the marshaling doesn't fail with valid items
	// Actual marshaling errors are unlikely with valid models.Item structs
	err := ApplyMerge(winner, []models.Item{}, true)
	if err != nil {
		t.Errorf("ApplyMerge() should handle empty items without error: %v", err)
	}
}

func TestApplyMerge_ItemsWithComplexData(t *testing.T) {
	tests := []struct {
		name      string
		winner    models.Item
		losers    []models.Item
		dryRun    bool
		shouldErr bool
	}{
		{
			name: "complex item with multiple fields and URLs",
			winner: models.Item{
				ID:        "complex_winner",
				Title:     "Complex Winner",
				UpdatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Vault: models.Vault{
					ID:   "vault_id",
					Name: "Test Vault",
				},
				Category: "LOGIN",
				Fields: []models.Field{
					{
						ID:    "f1",
						Type:  "username",
						Label: "username",
						Value: "user1",
					},
					{
						ID:    "f2",
						Type:  "password",
						Label: "password",
						Value: "pass1",
					},
					{
						ID:    "f3",
						Type:  "email",
						Label: "email",
						Value: "user@example.com",
					},
				},
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
					{HRef: "https://backup.example.com", Primary: false},
				},
			},
			losers: []models.Item{
				createTestItem("loser1", "Loser 1", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				createTestItem("loser2", "Loser 2", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
			dryRun:    true,
			shouldErr: false,
		},
		{
			name: "item with custom sections in fields",
			winner: models.Item{
				ID:        "winner_sections",
				Title:     "Winner With Sections",
				UpdatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Vault: models.Vault{
					ID:   "vault_id",
					Name: "Test Vault",
				},
				Category: "LOGIN",
				Fields: []models.Field{
					{
						ID:    "f1",
						Type:  "username",
						Label: "username",
						Value: "user1",
					},
					{
						ID:    "f2",
						Type:  "custom",
						Label: "security_question",
						Value: "What is your pet's name?",
						Section: &models.Section{
							ID: "security_questions",
						},
					},
				},
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			losers:    []models.Item{},
			dryRun:    true,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			err := ApplyMerge(tt.winner, tt.losers, tt.dryRun)

			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)

			if tt.shouldErr && err == nil {
				t.Errorf("ApplyMerge() expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ApplyMerge() unexpected error: %v", err)
			}

			if tt.dryRun {
				expectedMsg := fmt.Sprintf("[DRY RUN] Would edit item: %s", tt.winner.ID)
				if !bytes.Contains(buf.Bytes(), []byte(expectedMsg)) {
					t.Errorf("ApplyMerge() output missing dry-run edit message")
				}
			}
		})
	}
}

func TestApplyMerge_OutputFormat(t *testing.T) {
	winner := createTestItem(
		"test_winner",
		"Test Winner",
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)
	loser := createTestItem(
		"test_loser",
		"Test Loser",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	// Capture stdout
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	err := ApplyMerge(winner, []models.Item{loser}, true)

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("ApplyMerge() unexpected error: %v", err)
	}

	// Verify output contains expected components
	if !bytes.Contains(buf.Bytes(), []byte("[DRY RUN]")) {
		t.Errorf("ApplyMerge() output missing [DRY RUN] prefix")
	}

	if !bytes.Contains(buf.Bytes(), []byte("test_winner")) {
		t.Errorf("ApplyMerge() output missing winner ID")
	}

	if !bytes.Contains(buf.Bytes(), []byte("test_loser")) {
		t.Errorf("ApplyMerge() output missing loser ID")
	}

	// Verify JSON is present in output
	if !bytes.Contains(buf.Bytes(), []byte("{")) && !bytes.Contains(buf.Bytes(), []byte("}")) {
		t.Errorf("ApplyMerge() output missing JSON data")
	}

	// Verify the output contains the item ID
	if !bytes.Contains(buf.Bytes(), []byte("\"id\"")) {
		t.Errorf("ApplyMerge() JSON output missing 'id' field")
	}

	if len(output) == 0 {
		t.Errorf("ApplyMerge() produced no output in dry-run mode")
	}
}
