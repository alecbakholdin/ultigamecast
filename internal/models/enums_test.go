package models

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameFieldNamesHaventChanged(t *testing.T) {
	fields := []struct {
		Enum  GameField
		Field string
	}{
		{GameFieldScheduleStatus, "ScheduleStatus"},
		{GameFieldHalfCap, "HalfCap"},
		{GameFieldSoftCap, "SoftCap"},
		{GameFieldHardCap, "HardCap"},
		{GameFieldStart, "Start"},
	}
	gameType := reflect.TypeFor[Game]()
	for _, f := range fields {
		t.Run(fmt.Sprintf("Game.%s", f.Field), func(t *testing.T) {
			assert.EqualValues(t, f.Field, f.Enum)
			_, ok := gameType.FieldByName(f.Field)
			assert.True(t, ok)
		})
	}
}
