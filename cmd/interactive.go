package cmd

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"1merge/internal/models"
)

// displayDuplicateGroup displays information about a duplicate group to help user make merge decisions.
func displayDuplicateGroup(groupKey string, items []models.Item) {
	// Extract domain and username from groupKey (format: domain|username)
	parts := strings.Split(groupKey, "|")
	domain := parts[0]
	username := ""
	if len(parts) > 1 {
		username = parts[1]
	}

	fmt.Printf("\n=== Duplicate Group: %s | %s ===\n", domain, username)
	fmt.Printf("Found %d duplicate items:\n", len(items))

	for i, item := range items {
		maxLen := 8
		if len(item.ID) < maxLen {
			maxLen = len(item.ID)
		}
		fmt.Printf("  %d. %q (ID: %s...) - Updated: %s\n", i+1, item.Title, item.ID[:maxLen], formatTimestamp(item.UpdatedAt))
		if len(item.URLs) > 0 {
			fmt.Printf("     URL: %s\n", item.URLs[0].HRef)
		}
	}
	fmt.Println()
}

// promptUser prompts user for y/n/q input and returns normalized response.
func promptUser(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print("Merge these items? (y/n/q): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		response := strings.ToLower(strings.TrimSpace(line))

		// Validate input
		if response == "y" || response == "n" || response == "q" {
			return response, nil
		}

		// Invalid input, prompt again
		fmt.Println("Invalid input. Please enter 'y', 'n', or 'q'.")
	}
}

// formatTimestamp formats timestamp in human-readable format (YYYY-MM-DD HH:MM:SS).
func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
