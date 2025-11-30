package items

import (
	"1merge/internal/models"
)

// SelectWinner identifies the most recent item from a slice of items by comparing their UpdatedAt timestamps.
// If the slice is empty, it returns a zero-value Item.
// If multiple items have the same most recent timestamp, the first one is returned.
func SelectWinner(items []models.Item) models.Item {
	if len(items) == 0 {
		return models.Item{}
	}

	winner := items[0]
	for i := 1; i < len(items); i++ {
		if items[i].UpdatedAt.After(winner.UpdatedAt) {
			winner = items[i]
		}
	}

	return winner
}

// CalculateMerge implements the Superset merge strategy, combining the winner and loser items.
// It deep-copies the winner item and adds unique fields and URLs from the loser.
// Conflicting fields (same label) are placed in an "Archived Conflicts" section.
// Duplicate URLs are skipped.
func CalculateMerge(winner models.Item, loser models.Item) (models.Item, error) {
	// Deep copy the winner
	merged := models.Item{
		ID:                    winner.ID,
		Title:                 winner.Title,
		Vault:                 winner.Vault,
		Category:              winner.Category,
		UpdatedAt:             winner.UpdatedAt,
		AdditionalInformation: winner.AdditionalInformation,
	}

	// Deep copy fields from winner, including deep copy of Section pointers
	merged.Fields = make([]models.Field, len(winner.Fields))
	for i, field := range winner.Fields {
		merged.Fields[i] = field
		if field.Section != nil {
			sectionCopy := *field.Section
			merged.Fields[i].Section = &sectionCopy
		}
	}

	// Deep copy URLs from winner and track if winner has a primary URL
	winnerHasPrimary := false
	merged.URLs = make([]models.URL, len(winner.URLs))
	for i, url := range winner.URLs {
		merged.URLs[i] = url
		if url.Primary {
			winnerHasPrimary = true
		}
	}

	// Process loser's fields
	for _, loserField := range loser.Fields {
		exists, existingField := fieldExistsByLabel(merged.Fields, loserField.Label)
		if !exists {
			// Unique field, add it
			merged.Fields = append(merged.Fields, loserField)
		} else {
			sameValue := existingField.Value == loserField.Value
			sameType := existingField.Type == loserField.Type
			sameSection := (existingField.Section == nil && loserField.Section == nil) ||
				(existingField.Section != nil && loserField.Section != nil && existingField.Section.ID == loserField.Section.ID)

			if sameValue && sameType && sameSection {
				// Identical field, skip to avoid duplicate/conflict
				continue
			}
			// Conflicting field, add to "Archived Conflicts" section
			section := getOrCreateArchivedConflictsSection(merged.Fields)
			loserFieldCopy := loserField
			loserFieldCopy.Section = section
			merged.Fields = append(merged.Fields, loserFieldCopy)
		}
	}

	// Process loser's URLs
	for _, loserURL := range loser.URLs {
		if !urlExists(merged.URLs, loserURL.HRef) {
			// Unique URL, add it
			// If loser's URL is primary but winner already has a primary, demote loser's
			if loserURL.Primary && winnerHasPrimary {
				demotedURL := loserURL
				demotedURL.Primary = false
				merged.URLs = append(merged.URLs, demotedURL)
			} else {
				merged.URLs = append(merged.URLs, loserURL)
			}
		}
	}

	return merged, nil
}

// fieldExistsByLabel searches the fields slice for a field matching the given label (case-sensitive).
// Returns true and the matching field if found, false and zero-value Field otherwise.
func fieldExistsByLabel(fields []models.Field, label string) (bool, models.Field) {
	for _, field := range fields {
		if field.Label == label {
			return true, field
		}
	}
	return false, models.Field{}
}

// urlExists checks if a URL with the given href already exists in the URLs slice (exact string match).
// Returns true if found, false otherwise.
func urlExists(urls []models.URL, href string) bool {
	for _, url := range urls {
		if url.HRef == href {
			return true
		}
	}
	return false
}

// getOrCreateArchivedConflictsSection searches for or creates an "Archived Conflicts" section.
// All conflicting fields are grouped under this section.
// Note: This creates a section reference without explicitly defining it in a sections array.
// The 1Password CLI automatically creates sections when fields reference them during item edit.
func getOrCreateArchivedConflictsSection(fields []models.Field) *models.Section {
	// Search for existing "archived_conflicts" section
	for i := range fields {
		if fields[i].Section != nil && fields[i].Section.ID == "archived_conflicts" {
			return fields[i].Section
		}
	}

	// Create new section
	section := &models.Section{ID: "archived_conflicts"}
	return section
}
