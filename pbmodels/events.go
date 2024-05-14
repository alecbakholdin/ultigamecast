package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Events struct {
	models.BaseModel

	Game       string          `db:"game" form:"game" json:"game"`
	PointType  EventsPointType `db:"point_type" form:"point_type" json:"point_type"`
	Type       EventsType      `db:"type" form:"type" json:"type"`
	Player     string          `db:"player" form:"player" json:"player"`
	IsOpponent bool            `db:"is_opponent" form:"is_opponent" json:"is_opponent"`
	Message    string          `db:"message" form:"message" json:"message"`
}

type EventsPointType string

const (
	EventsPointTypeO EventsPointType = "O"
	EventsPointTypeD EventsPointType = "D"
)

type EventsType string

const (
	EventsTypeStartingLine  EventsType = "Starting Line"
	EventsTypeSubbedIn      EventsType = "Subbed In"
	EventsTypeGoal          EventsType = "Goal"
	EventsTypeAssist        EventsType = "Assist"
	EventsTypeBlock         EventsType = "Block"
	EventsTypeTurn          EventsType = "Turn"
	EventsTypeDrop          EventsType = "Drop"
	EventsTypePointStart    EventsType = "Point Start"
	EventsTypeSubIn         EventsType = "Sub In"
	EventsTypeSubOut        EventsType = "Sub Out"
	EventsTypeHalftime      EventsType = "Halftime"
	EventsTypeHalfCap       EventsType = "Half Cap"
	EventsTypeSoftCap       EventsType = "Soft Cap"
	EventsTypeHardCap       EventsType = "Hard Cap"
	EventsTypeGameEnd       EventsType = "Game End"
)


func (d *Events) CopyFrom(s *Events) *Events {
	d.Game = s.Game
	d.PointType = s.PointType
	d.Type = s.Type
	d.Player = s.Player
	d.IsOpponent = s.IsOpponent
	d.Message = s.Message
	return d
}

func (m *Events) Copy() *Events {
	return (&Events{}).CopyFrom(m)
}

func (m *Events) TableName() string {
    return "events"
}
