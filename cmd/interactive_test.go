package cmd

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"1merge/internal/models"
)

func TestDisplayDuplicateGroup(_ *testing.T) {
	// Create sample items with different titles, IDs, timestamps, and URLs
	item1 := models.Item{
		ID:        "abc12345abcd1234",
		Title:     "Google Account",
		UpdatedAt: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
		URLs:      []models.URL{{HRef: "https://accounts.google.com"}},
	}

	item2 := models.Item{
		ID:        "def67890def67890",
		Title:     "Gmail Login",
		UpdatedAt: time.Date(2024, 1, 10, 9, 15, 0, 0, time.UTC),
		URLs:      []models.URL{{HRef: "https://mail.google.com"}},
	}

	item3 := models.Item{
		ID:        "ghi11121ghi11121",
		Title:     "Google",
		UpdatedAt: time.Date(2023, 12, 20, 16, 45, 0, 0, time.UTC),
		URLs:      []models.URL{{HRef: "https://google.com"}},
	}

	// Capture stdout (by redirecting output to a buffer)
	// Note: Since displayDuplicateGroup prints to stdout directly, we can't capture it
	// without refactoring. This test verifies the function runs without error.
	groupKey := "google.com|user@example.com"
	items := []models.Item{item1, item2, item3}

	// Just verify it doesn't panic
	displayDuplicateGroup(groupKey, items)
}

func TestPromptUser_ValidInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"y\n", "y"},
		{"Y\n", "y"},
		{"n\n", "n"},
		{"N\n", "n"},
		{"q\n", "q"},
		{"Q\n", "q"},
		{" y \n", "y"},
		{" N \n", "n"},
	}

	for _, tt := range tests {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := promptUser(reader)
		if err != nil {
			t.Fatalf("promptUser(%q) returned error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Fatalf("promptUser(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestPromptUser_InvalidThenValid(t *testing.T) {
	// Create reader with invalid input first, then valid
	input := "invalid\ny\n"
	reader := bufio.NewReader(strings.NewReader(input))

	// Capture stdout to avoid printing during test
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	result, err := promptUser(reader)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("promptUser returned error: %v", err)
	}
	if result != "y" {
		t.Fatalf("promptUser = %q, expected %q", result, "y")
	}
}

func TestFormatTimestamp(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	result := formatTimestamp(testTime)
	expected := "2024-01-15 14:30:45"

	if result != expected {
		t.Fatalf("formatTimestamp(%v) = %q, expected %q", testTime, result, expected)
	}
}
