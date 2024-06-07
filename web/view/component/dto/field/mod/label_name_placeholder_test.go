package mod

import (
	"testing"
	"ultigamecast/web/view/component/dto/field"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

func TestFieldNameGuesser(t *testing.T) {
	t.Run("fields are properly calculated", func(t *testing.T) {
		fc := &field.FieldConfig{DtoField: "Test3okName123", InputAttributes: templ.Attributes{}}
		FieldNameGuesser().Apply(fc)
		assert.Equal(t, "Test3ok Name123", fc.Label)
		assert.Equal(t, "Test3ok Name123", fc.InputAttributes["placeholder"])
		assert.Equal(t, "test3ok_name123", fc.InputAttributes["name"])
	})

	t.Run("doesnt overwrite existing labels, names, placeholders", func(t *testing.T) {
		fc := &field.FieldConfig{
			DtoField: "Field",
			Label: "Label",
			InputAttributes: templ.Attributes{
				"placeholder": "Placeholder",
				"name": "name",
			},
		}
		FieldNameGuesser().Apply(fc)
		assert.Equal(t, "Label", fc.Label)
		assert.Equal(t, "Placeholder", fc.InputAttributes["placeholder"])
		assert.Equal(t, "name", fc.InputAttributes["name"])
	})

	t.Run("panics if any nonalphanumeric characters are present", func(t *testing.T) {
		fc := &field.FieldConfig{DtoField: "_", InputAttributes: templ.Attributes{}}
		assert.Panics(t, func() {FieldNameGuesser().Apply(fc)})
	})
}