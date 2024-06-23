package models

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameFieldNamesHaventChanged(t *testing.T) {
	t.Run("Game.ScheduleStatus", func(t *testing.T) {
		assert.EqualValues(t, GameJsonFieldScheduleStatus, jsonField(t, reflect.TypeFor[Game](), "ScheduleStatus"))
	})
}

func jsonField(t *testing.T, r reflect.Type, goField string) string {
	f, ok := r.FieldByName(goField)
	assert.True(t, ok)
	tag := f.Tag.Get("json")
	assert.NotEmpty(t, tag)
	return tag
}
