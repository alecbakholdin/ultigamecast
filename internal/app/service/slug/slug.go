package slug

import "strings"

const (
	validChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func From(name string) string {
	name = strings.ToLower(name)
	builder := strings.Builder{}
	for _, c := range name {
		if strings.ContainsRune(validChars, c) {
			builder.WriteRune(c)
		} else if c == ' ' {
			builder.WriteRune('-')
		}
	}
	return builder.String()
}