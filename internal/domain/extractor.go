package domain

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// GetBaseDomain extracts the base domain (eTLD+1) from a URL string.
// It is used to identify duplicate entries by comparing login URLs across subdomains.
// For example:
//   - "https://mail.google.com/foo" returns "google.com"
//   - "https://aws.amazon.com" returns "amazon.com"
//   - "https://www.bbc.co.uk" returns "bbc.co.uk"
//
// Edge cases:
//   - IP addresses are returned as-is (e.g., "192.168.1.1")
//   - localhost and similar hostnames are returned as-is
//   - Results are always lowercase for consistent comparison
//
// If the URL is invalid or cannot be parsed, an error is returned.
func GetBaseDomain(urlStr string) (string, error) {
	if urlStr == "" {
		return "", errors.New("empty URL string")
	}

	// Add scheme if missing to allow url.Parse to correctly identify the hostname
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract hostname (strips port numbers automatically)
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return "", errors.New("invalid URL: no hostname found")
	}

	// Check if hostname is an IP address
	if ip := net.ParseIP(hostname); ip != nil {
		return hostname, nil
	}

	// Check for localhost hostname
	if hostname == "localhost" {
		return hostname, nil
	}

	// Extract base domain using public suffix list
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to extract base domain: %w", err)
	}

	// Return in lowercase for consistent comparison
	return strings.ToLower(baseDomain), nil
}
