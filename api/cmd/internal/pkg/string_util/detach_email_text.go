package utils

import (
	"regexp"
)

// GetEmailFromText extracts email addresses from the provided text.
// It uses a regular expression to identify and return all email addresses found in the input string.
//
// The regular expression pattern used is designed to match:
// - A sequence of characters that can include letters, numbers, dots, underscores, percent signs, plus signs, and hyphens
// - Followed by the '@' symbol
// - Followed by a domain name that can include letters, numbers, dots, and hyphens
// - Ending with a dot followed by a top-level domain (TLD) that consists of at least two alphabetic characters
//
// Example usage:
// emails := GetEmailFromText("Contact us at support@example.com or sales@example.org")
// This would return a slice containing the email addresses found: ["support@example.com", "sales@example.org"]
func GetEmailFromText(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := re.FindAllString(text, -1)
	return emails
}
