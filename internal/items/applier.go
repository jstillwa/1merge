package items

import (
	"encoding/json"
	"fmt"
	"os"

	"1merge/internal/models"
	"1merge/internal/op"
)

// opClient is the shared injectable client for op CLI interactions in both applier and fetcher, overridden in tests.
var opClient op.Client = op.DefaultClient

// SetOpClient allows callers to override the shared op client for both applier and fetcher operations (useful for testing).
func SetOpClient(client op.Client) {
	if client == nil {
		opClient = op.DefaultClient
		return
	}
	opClient = client
}

// ApplyMerge orchestrates the actual 1Password vault modifications.
// It updates the winner item with merged data and archives all loser items.
// If dryRun is true, it prints what would be changed without executing any op commands.
func ApplyMerge(winner models.Item, losers []models.Item, dryRun bool) error {
	// Marshal winner to JSON
	jsonBytes, err := json.MarshalIndent(winner, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal winner item to JSON: %w", err)
	}

	// Handle dry-run mode
	if dryRun {
		fmt.Printf("[DRY RUN] Would edit item: %s (%s)\n", winner.ID, winner.Title)
		fmt.Println(string(jsonBytes))
		for _, loser := range losers {
			fmt.Printf("[DRY RUN] Would archive item: %s (%s)\n", loser.ID, loser.Title)
		}
		return nil
	}

	// Create temp file for item JSON template
	tempFile, err := os.CreateTemp("", "1merge-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write JSON to temp file
	if _, err := tempFile.Write(jsonBytes); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Execute edit command using template file
	if _, err := opClient.RunOpCmd("item", "edit", winner.ID, "--template", tempFile.Name()); err != nil {
		return fmt.Errorf("failed to edit item %s: %w", winner.ID, err)
	}

	// Execute archive commands
	for _, loser := range losers {
		if _, err := opClient.RunOpCmd("item", "delete", loser.ID, "--archive"); err != nil {
			return fmt.Errorf("failed to archive item %s: %w", loser.ID, err)
		}
	}

	return nil
}
