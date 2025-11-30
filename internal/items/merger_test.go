package items

import (
	"testing"
	"time"

	"1merge/internal/models"
)

func TestSelectWinner(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected models.Item
	}{
		{
			name: "multiple items with different timestamps",
			items: []models.Item{
				{
					ID:        "item1",
					Title:     "Item 1",
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "item2",
					Title:     "Item 2",
					UpdatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "item3",
					Title:     "Item 3",
					UpdatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: models.Item{
				ID:        "item2",
				Title:     "Item 2",
				UpdatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "two items with identical timestamps returns first",
			items: []models.Item{
				{
					ID:        "item1",
					Title:     "Item 1",
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "item2",
					Title:     "Item 2",
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: models.Item{
				ID:        "item1",
				Title:     "Item 1",
				UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "single item slice returns that item",
			items: []models.Item{
				{
					ID:        "item1",
					Title:     "Item 1",
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: models.Item{
				ID:        "item1",
				Title:     "Item 1",
				UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:     "empty slice returns zero-value Item",
			items:    []models.Item{},
			expected: models.Item{},
		},
		{
			name: "items not pre-sorted selects most recent correctly",
			items: []models.Item{
				{
					ID:        "item3",
					Title:     "Item 3",
					UpdatedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "item1",
					Title:     "Item 1",
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "item2",
					Title:     "Item 2",
					UpdatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: models.Item{
				ID:        "item3",
				Title:     "Item 3",
				UpdatedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectWinner(tt.items)

			if result.ID != tt.expected.ID {
				t.Errorf("SelectWinner() returned ID %q, expected %q", result.ID, tt.expected.ID)
			}
			if result.Title != tt.expected.Title {
				t.Errorf("SelectWinner() returned Title %q, expected %q", result.Title, tt.expected.Title)
			}
			if !result.UpdatedAt.Equal(tt.expected.UpdatedAt) {
				t.Errorf("SelectWinner() returned UpdatedAt %v, expected %v", result.UpdatedAt, tt.expected.UpdatedAt)
			}
		})
	}
}

func TestCalculateMerge(t *testing.T) {
	tests := []struct {
		name               string
		winner             models.Item
		loser              models.Item
		expectedFieldCount int
		expectedURLCount   int
		fieldCheck         func(*testing.T, models.Item) // Custom validation function
	}{
		{
			name: "unique field addition",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{Label: "username", Value: "user1"},
					{Label: "password", Value: "pass1"},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "email", Value: "test@example.com"},
					{Label: "notes", Value: "some notes"},
				},
			},
			expectedFieldCount: 4,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				labels := []string{}
				for _, field := range merged.Fields {
					labels = append(labels, field.Label)
				}
				if len(labels) != 4 {
					t.Errorf("expected 4 fields, got %d", len(labels))
				}
				hasUsername := false
				hasEmail := false
				hasPassword := false
				hasNotes := false
				for _, field := range merged.Fields {
					switch field.Label {
					case "username":
						hasUsername = true
					case "email":
						hasEmail = true
					case "password":
						hasPassword = true
					case "notes":
						hasNotes = true
					}
				}
				if !hasUsername || !hasEmail || !hasPassword || !hasNotes {
					t.Errorf("merged item missing expected fields")
				}
			},
		},
		{
			name: "no conflicts - identical field labels with same values",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{Label: "username", Value: "user@example.com"},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "username", Value: "user@example.com"},
				},
			},
			expectedFieldCount: 1,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 1 {
					t.Fatalf("expected 1 field, got %d", len(merged.Fields))
				}
				field := merged.Fields[0]
				if field.Label != "username" || field.Value != "user@example.com" {
					t.Errorf("merged username field not preserved correctly")
				}
				if field.Section != nil {
					t.Errorf("identical field should not be assigned to a section, got %v", field.Section)
				}
				for _, f := range merged.Fields {
					if f.Section != nil && f.Section.ID == "archived_conflicts" {
						t.Errorf("identical fields should not create archived_conflicts section")
					}
				}
			},
		},
		{
			name: "field conflicts are moved to Archived Conflicts section",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{Label: "notes", Value: "original notes"},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "notes", Value: "conflicting notes"},
				},
			},
			expectedFieldCount: 2,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 2 {
					t.Errorf("expected 2 fields, got %d", len(merged.Fields))
				}
				// First field should be the winner's notes without section
				if merged.Fields[0].Label != "notes" || merged.Fields[0].Value != "original notes" {
					t.Errorf("winner's notes field not preserved correctly")
				}
				if merged.Fields[0].Section != nil {
					t.Errorf("winner's field should not have section")
				}
				// Second field should be the loser's notes in Archived Conflicts section
				if merged.Fields[1].Label != "notes" || merged.Fields[1].Value != "conflicting notes" {
					t.Errorf("loser's notes field not added correctly")
				}
				if merged.Fields[1].Section == nil || merged.Fields[1].Section.ID != "archived_conflicts" {
					t.Errorf("conflicting field should be in archived_conflicts section, got %v", merged.Fields[1].Section)
				}
			},
		},
		{
			name: "URL merging - winner has URL1, loser has URL2",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				URLs: []models.URL{
					{HRef: "https://other.com", Primary: false},
				},
			},
			expectedFieldCount: 0,
			expectedURLCount:   2,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.URLs) != 2 {
					t.Errorf("expected 2 URLs, got %d", len(merged.URLs))
				}
				hasExample := false
				hasOther := false
				for _, url := range merged.URLs {
					if url.HRef == "https://example.com" {
						hasExample = true
					}
					if url.HRef == "https://other.com" {
						hasOther = true
					}
				}
				if !hasExample || !hasOther {
					t.Errorf("merged item missing expected URLs")
				}
			},
		},
		{
			name: "duplicate URL handling - same URL in both",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: false},
				},
			},
			expectedFieldCount: 0,
			expectedURLCount:   1,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.URLs) != 1 {
					t.Errorf("expected 1 URL (no duplicates), got %d", len(merged.URLs))
				}
			},
		},
		{
			name: "mixed scenario - unique fields, conflicts, and URL merging",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{Label: "username", Value: "user1"},
					{Label: "notes", Value: "winner notes"},
				},
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "email", Value: "user@example.com"},
					{Label: "notes", Value: "loser notes"},
				},
				URLs: []models.URL{
					{HRef: "https://other.com", Primary: false},
					{HRef: "https://example.com", Primary: true},
				},
			},
			expectedFieldCount: 4,
			expectedURLCount:   2,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 4 {
					t.Errorf("expected 4 fields, got %d", len(merged.Fields))
				}
				if len(merged.URLs) != 2 {
					t.Errorf("expected 2 URLs, got %d", len(merged.URLs))
				}
				// Verify conflict was placed in Archived Conflicts
				conflictCount := 0
				for _, field := range merged.Fields {
					if field.Section != nil && field.Section.ID == "archived_conflicts" {
						conflictCount++
					}
				}
				if conflictCount != 1 {
					t.Errorf("expected 1 field in Archived Conflicts, got %d", conflictCount)
				}
			},
		},
		{
			name: "empty loser - loser has no fields or URLs",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{Label: "username", Value: "user1"},
				},
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			loser: models.Item{
				ID:     "2",
				Title:  "Loser",
				Fields: []models.Field{},
				URLs:   []models.URL{},
			},
			expectedFieldCount: 1,
			expectedURLCount:   1,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 1 {
					t.Errorf("expected 1 field (winner's), got %d", len(merged.Fields))
				}
				if len(merged.URLs) != 1 {
					t.Errorf("expected 1 URL (winner's), got %d", len(merged.URLs))
				}
			},
		},
		{
			name: "empty winner - winner has no fields or URLs",
			winner: models.Item{
				ID:     "1",
				Title:  "Winner",
				Fields: []models.Field{},
				URLs:   []models.URL{},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "username", Value: "user1"},
				},
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			expectedFieldCount: 1,
			expectedURLCount:   1,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 1 {
					t.Errorf("expected 1 field (loser's), got %d", len(merged.Fields))
				}
				if len(merged.URLs) != 1 {
					t.Errorf("expected 1 URL (loser's), got %d", len(merged.URLs))
				}
			},
		},
		{
			name: "section preservation - winner has custom sections",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{
						Label:   "security_question",
						Value:   "What is your pet's name?",
						Section: &models.Section{ID: "security_questions"},
					},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				Fields: []models.Field{
					{Label: "username", Value: "user1"},
				},
			},
			expectedFieldCount: 2,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if len(merged.Fields) != 2 {
					t.Errorf("expected 2 fields, got %d", len(merged.Fields))
				}
				// Verify first field still has its custom section
				if merged.Fields[0].Section == nil || merged.Fields[0].Section.ID != "security_questions" {
					t.Errorf("custom section not preserved, got %v", merged.Fields[0].Section)
				}
			},
		},
		{
			name: "primary URL flag handling - loser primary demoted when winner has primary",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: true},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				URLs: []models.URL{
					{HRef: "https://other.com", Primary: true},
				},
			},
			expectedFieldCount: 0,
			expectedURLCount:   2,
			fieldCheck: func(t *testing.T, merged models.Item) {
				// Verify winner's primary flag is preserved
				if merged.URLs[0].Primary != true {
					t.Errorf("winner's primary URL flag not preserved")
				}
				// Verify loser's primary URL is demoted to non-primary
				if merged.URLs[1].Primary != false {
					t.Errorf("loser's primary flag should be demoted to false, got true")
				}
			},
		},
		{
			name: "preserves AdditionalInformation from winner",
			winner: models.Item{
				ID:                    "1",
				Title:                 "Winner",
				AdditionalInformation: "Important notes about this account",
				Fields:                []models.Field{{Label: "username", Value: "user1"}},
			},
			loser: models.Item{
				ID:                    "2",
				Title:                 "Loser",
				AdditionalInformation: "Different notes",
				Fields:                []models.Field{{Label: "email", Value: "test@example.com"}},
			},
			expectedFieldCount: 2,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				if merged.AdditionalInformation != "Important notes about this account" {
					t.Errorf("AdditionalInformation not preserved from winner, got %q", merged.AdditionalInformation)
				}
			},
		},
		{
			name: "deep copies sections to prevent shared pointer mutations",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				Fields: []models.Field{
					{
						Label:   "field1",
						Value:   "value1",
						Section: &models.Section{ID: "section1"},
					},
					{
						Label:   "field2",
						Value:   "value2",
						Section: &models.Section{ID: "section1"},
					},
				},
			},
			loser:              models.Item{ID: "2", Title: "Loser"},
			expectedFieldCount: 2,
			expectedURLCount:   0,
			fieldCheck: func(t *testing.T, merged models.Item) {
				// Verify sections are different pointers (deep copy)
				if merged.Fields[0].Section == merged.Fields[1].Section {
					t.Errorf("sections should be deep copied, got same pointer")
				}
				// But have same ID
				if merged.Fields[0].Section.ID != merged.Fields[1].Section.ID {
					t.Errorf("section IDs should match")
				}
			},
		},
		{
			name: "loser primary URL preserved when winner has no primary",
			winner: models.Item{
				ID:    "1",
				Title: "Winner",
				URLs: []models.URL{
					{HRef: "https://example.com", Primary: false},
				},
			},
			loser: models.Item{
				ID:    "2",
				Title: "Loser",
				URLs: []models.URL{
					{HRef: "https://other.com", Primary: true},
				},
			},
			expectedFieldCount: 0,
			expectedURLCount:   2,
			fieldCheck: func(t *testing.T, merged models.Item) {
				// Winner has no primary, so loser's primary should be preserved
				if merged.URLs[1].Primary != true {
					t.Errorf("loser's primary flag should be preserved when winner has no primary")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateMerge(tt.winner, tt.loser)

			if err != nil {
				t.Errorf("CalculateMerge() returned unexpected error: %v", err)
			}

			if len(result.Fields) != tt.expectedFieldCount {
				t.Errorf("CalculateMerge() returned %d fields, expected %d", len(result.Fields), tt.expectedFieldCount)
			}

			if len(result.URLs) != tt.expectedURLCount {
				t.Errorf("CalculateMerge() returned %d URLs, expected %d", len(result.URLs), tt.expectedURLCount)
			}

			// Verify winner's identity is preserved
			if result.ID != tt.winner.ID {
				t.Errorf("merged item ID %q, expected %q (winner's ID)", result.ID, tt.winner.ID)
			}

			// Run custom field checks
			if tt.fieldCheck != nil {
				tt.fieldCheck(t, result)
			}
		})
	}
}

func TestFieldExistsByLabel(t *testing.T) {
	tests := []struct {
		name           string
		fields         []models.Field
		label          string
		expectedExists bool
		expectedValue  string
	}{
		{
			name: "field found",
			fields: []models.Field{
				{Label: "username", Value: "user1"},
				{Label: "password", Value: "pass1"},
			},
			label:          "username",
			expectedExists: true,
			expectedValue:  "user1",
		},
		{
			name: "field not found",
			fields: []models.Field{
				{Label: "username", Value: "user1"},
			},
			label:          "email",
			expectedExists: false,
			expectedValue:  "",
		},
		{
			name:           "empty slice",
			fields:         []models.Field{},
			label:          "username",
			expectedExists: false,
			expectedValue:  "",
		},
		{
			name: "case-sensitive matching",
			fields: []models.Field{
				{Label: "Username", Value: "user1"},
			},
			label:          "username",
			expectedExists: false,
			expectedValue:  "",
		},
		{
			name: "case-sensitive matching - correct case",
			fields: []models.Field{
				{Label: "Username", Value: "user1"},
			},
			label:          "Username",
			expectedExists: true,
			expectedValue:  "user1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, field := fieldExistsByLabel(tt.fields, tt.label)

			if exists != tt.expectedExists {
				t.Errorf("fieldExistsByLabel() returned exists=%v, expected %v", exists, tt.expectedExists)
			}

			if exists && field.Value != tt.expectedValue {
				t.Errorf("fieldExistsByLabel() returned value %q, expected %q", field.Value, tt.expectedValue)
			}
		})
	}
}

func TestUrlExists(t *testing.T) {
	tests := []struct {
		name           string
		urls           []models.URL
		href           string
		expectedExists bool
	}{
		{
			name: "URL found",
			urls: []models.URL{
				{HRef: "https://example.com", Primary: true},
				{HRef: "https://other.com", Primary: false},
			},
			href:           "https://example.com",
			expectedExists: true,
		},
		{
			name: "URL not found",
			urls: []models.URL{
				{HRef: "https://example.com", Primary: true},
			},
			href:           "https://other.com",
			expectedExists: false,
		},
		{
			name:           "empty slice",
			urls:           []models.URL{},
			href:           "https://example.com",
			expectedExists: false,
		},
		{
			name: "exact string match required",
			urls: []models.URL{
				{HRef: "https://example.com", Primary: true},
			},
			href:           "https://example.com/",
			expectedExists: false,
		},
		{
			name: "case-sensitive matching",
			urls: []models.URL{
				{HRef: "https://Example.com", Primary: true},
			},
			href:           "https://example.com",
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := urlExists(tt.urls, tt.href)

			if exists != tt.expectedExists {
				t.Errorf("urlExists() returned %v, expected %v", exists, tt.expectedExists)
			}
		})
	}
}

func TestGetOrCreateArchivedConflictsSection(t *testing.T) {
	tests := []struct {
		name                  string
		fields                []models.Field
		expectedSectionID     string
		expectedSectionNotNil bool
	}{
		{
			name: "creates new section when none exists",
			fields: []models.Field{
				{Label: "username", Value: "user1"},
			},
			expectedSectionID:     "archived_conflicts",
			expectedSectionNotNil: true,
		},
		{
			name: "returns existing archived_conflicts section",
			fields: []models.Field{
				{
					Label:   "notes",
					Value:   "some notes",
					Section: &models.Section{ID: "archived_conflicts"},
				},
			},
			expectedSectionID:     "archived_conflicts",
			expectedSectionNotNil: true,
		},
		{
			name:                  "empty fields slice creates new section",
			fields:                []models.Field{},
			expectedSectionID:     "archived_conflicts",
			expectedSectionNotNil: true,
		},
		{
			name: "ignores other custom sections",
			fields: []models.Field{
				{
					Label:   "security_question",
					Value:   "What?",
					Section: &models.Section{ID: "security_questions"},
				},
			},
			expectedSectionID:     "archived_conflicts",
			expectedSectionNotNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := getOrCreateArchivedConflictsSection(tt.fields)

			if section == nil && tt.expectedSectionNotNil {
				t.Errorf("getOrCreateArchivedConflictsSection() returned nil section")
			}

			if section != nil && section.ID != tt.expectedSectionID {
				t.Errorf("getOrCreateArchivedConflictsSection() returned section ID %q, expected %q", section.ID, tt.expectedSectionID)
			}
		})
	}
}
