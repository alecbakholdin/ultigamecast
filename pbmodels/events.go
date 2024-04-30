package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Events struct {
	models.BaseModel

	IsOpponent bool            `db:"is_opponent" json:"is_opponent" form:"is_opponent"`
	Type       EventsType      `db:"type" json:"type" form:"type"`
	Player     string          `db:"player" json:"player" form:"player"`
	Message    string          `db:"message" json:"message" form:"message"`
	Game       string          `db:"game" json:"game" form:"game"`
	PointType  EventsPointType `json:"point_type" form:"point_type" db:"point_type"`
}

type EventsType string

const (
	EventsTypeStartingLine EventsType = "Starting Line"
	EventsTypeSubbedIn     EventsType = "Subbed In"
	EventsTypeGoal         EventsType = "Goal"
	EventsTypeAssist       EventsType = "Assist"
	EventsTypeBlock        EventsType = "Block"
	EventsTypeTurn         EventsType = "Turn"
	EventsTypeDrop         EventsType = "Drop"
	EventsTypePointStart   EventsType = "Point Start"
	EventsTypeSubIn        EventsType = "Sub In"
	EventsTypeSubOut       EventsType = "Sub Out"
	EventsTypeHalftime     EventsType = "Halftime"
	EventsTypeHalfCap      EventsType = "Half Cap"
	EventsTypeSoftCap      EventsType = "Soft Cap"
	EventsTypeHardCap      EventsType = "Hard Cap"
	EventsTypeGameEnd      EventsType = "Game End"
)

type EventsPointType string

const (
	EventsPointTypeO EventsPointType = "O"
	EventsPointTypeD EventsPointType = "D"
)

func (m *Events) TableName() string {
	return "events"
}
