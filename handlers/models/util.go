package models

import (
	"regexp"
	"strings"
)

var whitespaceRegex = regexp.MustCompile(`[\s+]`)

func convertToSlug(s string) string {
	noWhitespace := whitespaceRegex.ReplaceAllString(s, "-")
	return strings.ToLower(noWhitespace)
}
