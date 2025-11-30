package items

import (
	"testing"

	"1merge/internal/models"
)

func TestGroupDuplicates(t *testing.T) {
	tests := []struct {
		name           string
		items          []models.Item
		expectedGroups int
		expectedKeys   map[string]int // key -> expected count of items in that group
	}{
		{
			name: "exact duplicates",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "3",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 1,
			expectedKeys: map[string]int{
				"google.com|user@example.com": 3,
			},
		},
		{
			name: "no duplicates",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user1@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://amazon.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user2@example.com"},
					},
				},
				{
					ID: "3",
					URLs: []models.URL{
						{HRef: "https://facebook.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user3@example.com"},
					},
				},
			},
			expectedGroups: 0,
			expectedKeys:   map[string]int{},
		},
		{
			name: "mixed scenario",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "3",
					URLs: []models.URL{
						{HRef: "https://amazon.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "admin@company.com"},
					},
				},
				{
					ID: "4",
					URLs: []models.URL{
						{HRef: "https://amazon.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "admin@company.com"},
					},
				},
				{
					ID: "5",
					URLs: []models.URL{
						{HRef: "https://facebook.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "unique@example.com"},
					},
				},
			},
			expectedGroups: 2,
			expectedKeys: map[string]int{
				"google.com|user@example.com":      2,
				"amazon.com|admin@company.com":     2,
			},
		},
		{
			name: "case insensitivity",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "User@Example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 1,
			expectedKeys: map[string]int{
				"google.com|user@example.com": 2,
			},
		},
		{
			name: "subdomain handling",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://mail.google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://accounts.google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 1,
			expectedKeys: map[string]int{
				"google.com|user@example.com": 2,
			},
		},
		{
			name: "missing username field",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "password", Value: "secret"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "password", Value: "secret"},
					},
				},
			},
			expectedGroups: 0,
			expectedKeys:   map[string]int{},
		},
		{
			name: "missing URLs",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 0,
			expectedKeys:   map[string]int{},
		},
		{
			name:           "empty input",
			items:          []models.Item{},
			expectedGroups: 0,
			expectedKeys:   map[string]int{},
		},
		{
			name: "primary URL selection",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://secondary.google.com", Primary: false},
						{HRef: "https://primary.google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://primary.google.com", Primary: true},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 1,
			expectedKeys: map[string]int{
				"google.com|user@example.com": 2,
			},
		},
		{
			name: "first URL when no primary",
			items: []models.Item{
				{
					ID: "1",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: false},
						{HRef: "https://other.google.com", Primary: false},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
				{
					ID: "2",
					URLs: []models.URL{
						{HRef: "https://google.com", Primary: false},
					},
					Fields: []models.Field{
						{Type: "username", Value: "user@example.com"},
					},
				},
			},
			expectedGroups: 1,
			expectedKeys: map[string]int{
				"google.com|user@example.com": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GroupDuplicates(tt.items)

			if len(result) != tt.expectedGroups {
				t.Errorf("GroupDuplicates() returned %d groups, expected %d", len(result), tt.expectedGroups)
			}

			for expectedKey, expectedCount := range tt.expectedKeys {
				group, exists := result[expectedKey]
				if !exists {
					t.Errorf("expected key %q not found in result", expectedKey)
					continue
				}
				if len(group) != expectedCount {
					t.Errorf("key %q has %d items, expected %d", expectedKey, len(group), expectedCount)
				}
			}

			// Verify no unexpected keys
			for resultKey := range result {
				if _, expected := tt.expectedKeys[resultKey]; !expected {
					t.Errorf("unexpected key %q in result", resultKey)
				}
			}
		})
	}
}

func TestExtractUsername(t *testing.T) {
	tests := []struct {
		name     string
		item     models.Item
		expected string
	}{
		{
			name: "item with username field",
			item: models.Item{
				Fields: []models.Field{
					{Type: "username", Value: "user@example.com"},
				},
			},
			expected: "user@example.com",
		},
		{
			name: "item with multiple fields returns username",
			item: models.Item{
				Fields: []models.Field{
					{Type: "password", Value: "secret"},
					{Type: "username", Value: "myuser"},
					{Type: "email", Value: "user@example.com"},
				},
			},
			expected: "myuser",
		},
		{
			name: "item without username field",
			item: models.Item{
				Fields: []models.Field{
					{Type: "password", Value: "secret"},
					{Type: "email", Value: "user@example.com"},
				},
			},
			expected: "",
		},
		{
			name: "username with mixed case is lowercased",
			item: models.Item{
				Fields: []models.Field{
					{Type: "username", Value: "UserName@Example.COM"},
				},
			},
			expected: "username@example.com",
		},
		{
			name: "username with whitespace is trimmed",
			item: models.Item{
				Fields: []models.Field{
					{Type: "username", Value: "  username  "},
				},
			},
			expected: "username",
		},
		{
			name: "empty fields slice",
			item: models.Item{
				Fields: []models.Field{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractUsername(tt.item)
			if result != tt.expected {
				t.Errorf("extractUsername() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestGetPrimaryURL(t *testing.T) {
	tests := []struct {
		name     string
		item     models.Item
		expected string
	}{
		{
			name: "item with primary URL marked",
			item: models.Item{
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: false},
					{HRef: "https://primary.example.com", Primary: true},
				},
			},
			expected: "https://primary.example.com",
		},
		{
			name: "item with multiple URLs but no primary returns first",
			item: models.Item{
				URLs: []models.URL{
					{HRef: "https://first.example.com", Primary: false},
					{HRef: "https://second.example.com", Primary: false},
				},
			},
			expected: "https://first.example.com",
		},
		{
			name: "item with no URLs",
			item: models.Item{
				URLs: []models.URL{},
			},
			expected: "",
		},
		{
			name: "item with single URL",
			item: models.Item{
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: false},
				},
			},
			expected: "https://example.com",
		},
		{
			name: "item with single URL marked as primary",
			item: models.Item{
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			expected: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPrimaryURL(tt.item)
			if result != tt.expected {
				t.Errorf("getPrimaryURL() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
