package utils

import (
	"regexp"
)

func GetEmailFromText(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := re.FindAllString(text, -1)
	return emails
}
