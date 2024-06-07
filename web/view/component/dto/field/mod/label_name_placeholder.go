package mod

import (
	"fmt"
	"strings"
	"ultigamecast/internal/assert"
	"ultigamecast/web/view/component/dto/field"
)

// modifier guesses the label, name, and placeholder given the DTO FieldName of the FieldConfig object
// expects the FieldName to be Pascal Case. Does not overwrite existing label, name, or placeholder.
// panics if provided name contains non-alphanumeric characters
// label and placeholder: Pascal Case but space-separated
// name: snake_case
type FieldNameModifier struct {
}

// modifier guesses the label, name, and placeholder given the DTO FieldName of the FieldConfig object
// expects the FieldName to be Pascal Case. Does not overwrite existing label, name, or placeholder.
// panics if provided name contains non-alphanumeric characters
// label and placeholder: Pascal Case but space-separated
// name: snake_case
func FieldNameGuesser() *FieldNameModifier {
	return &FieldNameModifier{}
}

func (l *FieldNameModifier) Apply(fc *field.FieldConfig) {
	assert.That(fc.InputAttributes != nil, "FieldConfig InputAttributes is nil")

	tokens := tokenizePascalCase(fc.DtoField)
	
	spaceSeparated := strings.Join(tokens, " ")
	if fc.Label == "" {
		fc.Label = spaceSeparated
	}
	if _, ok := fc.InputAttributes["placeholder"]; !ok {
		fc.InputAttributes["placeholder"] = spaceSeparated
	}

	if _, ok := fc.InputAttributes["name"]; !ok {
		snakeCase := strings.ToLower(strings.Join(tokens, "_"))
		fc.InputAttributes["name"] = snakeCase
	}
}

// converts pascal-case TestString into a slice of strings ['Test', 'String']
func tokenizePascalCase(s string) []string {
	tokens := []string{}
	token := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			if token != "" {
				tokens = append(tokens, token)
			}
			token = string(c)
		} else if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			token += string(c)
		} else {
			panic(fmt.Sprintf("Unexpected character %c", c))
		}
	}
	if token != "" {
		tokens = append(tokens, token)
	}
	return tokens
}
