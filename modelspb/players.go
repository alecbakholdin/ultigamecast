// Code generated by pb-codegen v0.1 DO NOT EDIT.
package modelspb

import (
    "github.com/pocketbase/pocketbase/models"
)

// Players is a wrapper around models.Record for type safe operations on the collection players.
type Players struct {
    Record *models.Record
}

// GetTeam returns the value of the "team" field
// Relation collection related : 846ykkxqtaqjxst
func (m *Players) GetTeam() string {
    return m.Record.GetString("team")
}

// SetTeam sets the value of the "team" field
// Relation collection related : 846ykkxqtaqjxst
func (m *Players) SetTeam(val string)  {
    m.Record.Set("team", val)
}

// GetName returns the value of the "name" field
func (m *Players) GetName() string {
    return m.Record.GetString("name")
}

// SetName sets the value of the "name" field
func (m *Players) SetName(val string)  {
    m.Record.Set("name", val)
}

// GetOrder returns the value of the "order" field
func (m *Players) GetOrder() int {
    return m.Record.GetInt("order")
}

// SetOrder sets the value of the "order" field
func (m *Players) SetOrder(val int)  {
    m.Record.Set("order", val)
}
