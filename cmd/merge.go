package cmd

import (
	"fmt"
	"io"

	"1merge/internal/items"
	"1merge/internal/models"
)

// applyMergeAndReport delegates merging to items.ApplyMerge and handles user-facing success logging.
func applyMergeAndReport(out io.Writer, winner models.Item, losers []models.Item, dryRun bool) error {
	if err := items.ApplyMerge(winner, losers, dryRun); err != nil {
		return err
	}

	if !dryRun {
		fmt.Fprintf(out, "Successfully merged %d items into %s\n", len(losers), winner.ID)
	}

	return nil
}
