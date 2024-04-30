package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Games struct {
	models.BaseModel

	HalfCap           int            `db:"half_cap" json:"half_cap" form:"half_cap"`
	SoftCap           int            `db:"soft_cap" json:"soft_cap" form:"soft_cap"`
	StartTime         types.DateTime `db:"start_time" json:"start_time"`
	StartTimeTimezone string         `form:"start_time_timezone" json:"start_time_timezone"`
	StartTimeDatetime string         `form:"start_time_datetime" json:"start_time_datetime"`
	Status            GamesStatus    `db:"status" json:"status" form:"status"`
	Tournament        string         `db:"tournament" json:"tournament" form:"tournament"`
	Opponent          string         `form:"opponent" db:"opponent" json:"opponent"`
	TeamScore         int            `form:"team_score" db:"team_score" json:"team_score"`
	OpponentScore     int            `db:"opponent_score" json:"opponent_score" form:"opponent_score"`
	HardCap           int            `db:"hard_cap" json:"hard_cap" form:"hard_cap"`
	WindMph           float32        `db:"wind_mph" json:"wind_mph" form:"wind_mph"`
	TempF             float32        `db:"temp_f" json:"temp_f" form:"temp_f"`
}

type GamesStatus string

const (
	GamesStatusScheduled GamesStatus = "scheduled"
	GamesStatusLive     GamesStatus = "live"
	GamesStatusCompleted GamesStatus = "completed"
)

func (m *Games) TableName() string {
    return "games"
}

