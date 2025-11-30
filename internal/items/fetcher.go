package items

import (
	"encoding/json"
	"fmt"

	"1merge/internal/models"
	"1merge/internal/op"
)

// FetchItems retrieves login items from 1Password
func FetchItems(vault string) ([]models.Item, error) {
	// Build command arguments
	args := []string{"item", "list", "--categories", "LOGIN", "--format", "json"}

	// Append vault flag if specified
	if vault != "" {
		args = append(args, "--vault", vault)
	}

	// Execute the op command
	output, err := op.RunOpCmd(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items from 1Password: %w", err)
	}

	// Unmarshal JSON response into Item slice
	var items []models.Item
	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal 1Password items: %w", err)
	}

	return items, nil
}
