package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"1merge/internal/items"
	"1merge/internal/models"
	"1merge/internal/op"
	"github.com/spf13/cobra"
)

var (
	vault  string
	dryRun bool
	auto   bool
)

var rootCmd = &cobra.Command{
	Use:   "1merge",
	Short: "Merge duplicate 1Password login entries",
	Long: `1Merge is a CLI tool that helps you identify and merge duplicate login entries
in your 1Password vaults. It can scan your vault, find duplicates, and merge them
automatically or with your confirmation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dryRun {
			fmt.Println("Dry Run Mode Enabled")
		}

		// Verify op CLI is installed and user is signed in
		if err := op.VerifyOpReady(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// Get whoami information to confirm authentication
		_, err := op.GetWhoAmI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// Fetch login items from 1Password
		fetchedItems, err := items.FetchItems(vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching items: %v\n", err)
			return
		}

		fmt.Printf("Found %d login items in vault\n", len(fetchedItems))

		// Group duplicates
		duplicateGroups := items.GroupDuplicates(fetchedItems)

		if len(duplicateGroups) == 0 {
			fmt.Println("No duplicate items found.")
			return
		}

		fmt.Printf("Found %d duplicate groups\n", len(duplicateGroups))

		// Initialize statistics tracking
		processedGroups := 0
		skippedGroups := 0
		failedGroups := 0
		totalMerged := 0

		// Create reader for interactive input (only if not --auto)
		var reader *bufio.Reader
		if !auto {
			reader = bufio.NewReader(os.Stdin)
		}

		keys := make([]string, 0, len(duplicateGroups))
		for groupKey := range duplicateGroups {
			keys = append(keys, groupKey)
		}
		sort.Strings(keys)

		// Loop through duplicate groups in deterministic order
		for _, groupKey := range keys {
			groupItems := duplicateGroups[groupKey]
			displayDuplicateGroup(groupKey, groupItems)

			shouldMerge := false

			// Handle auto mode vs interactive mode
			if auto {
				fmt.Println("[AUTO MODE] Merging group automatically...")
				shouldMerge = true
			} else {
				response, err := promptUser(reader)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
					continue
				}

				if response == "q" {
					fmt.Println("Exiting...")
					break
				}

				if response == "n" {
					skippedGroups++
					fmt.Println("Skipped.")
					continue
				}

				if response == "y" {
					shouldMerge = true
				}
			}

			// Process merge if confirmed
			if shouldMerge {
				// Select winner (most recent item)
				winner := items.SelectWinner(groupItems)

				// Build losers slice (all items except winner)
				losers := []models.Item{}
				for _, item := range groupItems {
					if item.ID != winner.ID {
						losers = append(losers, item)
					}
				}

				// Iteratively merge all losers into winner
				merged := winner
				mergeSuccess := true
				for _, loser := range losers {
					var mergeErr error
					merged, mergeErr = items.CalculateMerge(merged, loser)
					if mergeErr != nil {
						fmt.Fprintf(os.Stderr, "Error merging items: %v\n", mergeErr)
						mergeSuccess = false
						failedGroups++
						break
					}
				}

				if !mergeSuccess {
					continue
				}

				// Apply merge using existing helper
				if err := applyMergeAndReport(os.Stdout, merged, losers, dryRun); err != nil {
					fmt.Fprintf(os.Stderr, "Error applying merge: %v\n", err)
					failedGroups++
					continue
				}

				processedGroups++
				totalMerged += len(losers)
			}
		}

		// Print summary
		fmt.Println("\n=== Summary ===")
		fmt.Printf("Processed groups: %d\n", processedGroups)
		fmt.Printf("Skipped groups: %d\n", skippedGroups)
		fmt.Printf("Failed groups: %d\n", failedGroups)
		fmt.Printf("Total items merged: %d\n", totalMerged)
		if dryRun {
			fmt.Println("(Dry run - no changes were made)")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vault, "vault", "", "Specifies which 1Password vault to scan (uses default vault if not specified)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Prevents any write operations and only prints what would happen")
	rootCmd.PersistentFlags().BoolVar(&auto, "auto", false, "Automatically merges duplicates without prompting")
}
