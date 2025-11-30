package models

import "time"

// Item represents a 1Password login item
type Item struct {
	ID                    string    `json:"id"`
	Title                 string    `json:"title"`
	URLs                  []URL     `json:"urls"`
	Vault                 Vault     `json:"vault"`
	Category              string    `json:"category"`
	Fields                []Field   `json:"fields"`
	UpdatedAt             time.Time `json:"updated_at"`
	AdditionalInformation string    `json:"additional_information"`
}

// URL represents a URL associated with an item
type URL struct {
	Label   string `json:"label,omitempty"`
	HRef    string `json:"href"`
	Primary bool   `json:"primary"`
}

// Field represents a field within an item
type Field struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Label   string   `json:"label"`
	Value   string   `json:"value"`
	Section *Section `json:"section,omitempty"`
}

// Section represents a section grouping for fields
type Section struct {
	ID string `json:"id"`
}

// Vault represents vault information
type Vault struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
