package items

import (
	"encoding/json"
	"fmt"

	"1merge/internal/models"
	"1merge/internal/op"
)

// opClient is the injectable client for op CLI interactions, overridden in tests.
var opClient op.Client = op.DefaultClient

// SetOpClient allows callers to override the op client (useful for testing).
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

	// Execute edit command
	if err := opClient.RunOpCmdWithStdin(jsonBytes, "item", "edit", winner.ID); err != nil {
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
