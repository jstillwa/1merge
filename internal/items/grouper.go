package items

import (
	"strings"

	"1merge/internal/domain"
	"1merge/internal/models"
)

// GroupDuplicates groups items by matching base domain and username combinations.
// It returns a map where keys are "baseDomain|username" (lowercased) and values are
// slices of items that share the same domain and username. Only groups with 2 or more
// items (actual duplicates) are returned in the map.
func GroupDuplicates(items []models.Item) map[string][]models.Item {
	groups := make(map[string][]models.Item)

	for _, item := range items {
		username := extractUsername(item)
		url := getPrimaryURL(item)

		// Skip items with missing username or URL
		if username == "" || url == "" {
			continue
		}

		// Extract base domain from URL
		baseDomain, err := domain.GetBaseDomain(url)
		if err != nil {
			// Skip items with invalid URLs
			continue
		}

		// Generate grouping key: baseDomain|username (case-insensitive)
		key := strings.ToLower(baseDomain) + "|" + username

		// Append item to the group
		groups[key] = append(groups[key], item)
	}

	// Filter out single-item groups (not duplicates)
	for key, group := range groups {
		if len(group) < 2 {
			delete(groups, key)
		}
	}

	return groups
}

// extractUsername extracts the username from an item.
// When using "op item list", the username is in AdditionalInformation.
// When using "op item get", it's in the Fields array with Type="username".
// Returns empty string if no username is found.
func extractUsername(item models.Item) string {
	// First check AdditionalInformation (from "op item list")
	if item.AdditionalInformation != "" {
		return strings.ToLower(strings.TrimSpace(item.AdditionalInformation))
	}

	// Fall back to Fields array (from "op item get")
	for _, field := range item.Fields {
		if field.Type == "username" {
			return strings.ToLower(strings.TrimSpace(field.Value))
		}
	}
	return ""
}

// getPrimaryURL returns the primary URL from an item, or the first URL if no primary is marked.
// Returns empty string if the item has no URLs.
func getPrimaryURL(item models.Item) string {
	if len(item.URLs) == 0 {
		return ""
	}

	// Look for primary URL
	for _, u := range item.URLs {
		if u.Primary {
			return u.HRef
		}
	}

	// Return first URL if no primary found
	return item.URLs[0].HRef
}
