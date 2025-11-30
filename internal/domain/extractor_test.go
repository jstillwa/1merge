package domain

import (
	"testing"
)

func TestGetBaseDomain(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Standard domains
		{
			name:        "simple domain",
			input:       "https://google.com",
			expected:    "google.com",
			expectError: false,
		},
		{
			name:        "domain with www subdomain",
			input:       "https://www.google.com",
			expected:    "google.com",
			expectError: false,
		},
		{
			name:        "domain with multiple subdomains and path",
			input:       "https://mail.google.com/foo",
			expected:    "google.com",
			expectError: false,
		},

		// Subdomains
		{
			name:        "AWS domain with subdomain",
			input:       "https://aws.amazon.com",
			expected:    "amazon.com",
			expectError: false,
		},
		{
			name:        "AWS console with multiple subdomains",
			input:       "https://console.aws.amazon.com",
			expected:    "amazon.com",
			expectError: false,
		},

		// Country-code TLDs
		{
			name:        "BBC UK domain with www",
			input:       "https://www.bbc.co.uk",
			expected:    "bbc.co.uk",
			expectError: false,
		},
		{
			name:        "BBC UK domain with news subdomain",
			input:       "https://news.bbc.co.uk",
			expected:    "bbc.co.uk",
			expectError: false,
		},

		// Edge cases: invalid URLs
		{
			name:        "invalid URL string",
			input:       "not a url",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},

		// Edge cases: IP addresses
		{
			name:        "IPv4 address",
			input:       "http://192.168.1.1",
			expected:    "192.168.1.1",
			expectError: false,
		},
		{
			name:        "IPv4 address with port",
			input:       "http://192.168.1.1:8080",
			expected:    "192.168.1.1",
			expectError: false,
		},

		// Edge cases: localhost
		{
			name:        "localhost",
			input:       "http://localhost:8080",
			expected:    "localhost",
			expectError: false,
		},

		// Edge cases: URLs without scheme
		{
			name:        "URL without scheme",
			input:       "google.com",
			expected:    "google.com",
			expectError: false,
		},

		// Edge cases: URLs with paths and query params
		{
			name:        "URL with complex path and query",
			input:       "https://mail.google.com/mail/u/0/?tab=rm",
			expected:    "google.com",
			expectError: false,
		},

		// Case sensitivity
		{
			name:        "uppercase URL returns lowercase domain",
			input:       "HTTPS://GOOGLE.COM",
			expected:    "google.com",
			expectError: false,
		},
		{
			name:        "mixed case subdomain returns lowercase domain",
			input:       "https://Mail.Google.COM/foo",
			expected:    "google.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetBaseDomain(tt.input)

			if (err != nil) != tt.expectError {
				t.Errorf("GetBaseDomain(%q) error = %v, expectError = %v", tt.input, err, tt.expectError)
				return
			}

			if result != tt.expected {
				t.Errorf("GetBaseDomain(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
