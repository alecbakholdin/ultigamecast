package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Games struct {
	models.BaseModel

	WindMph           float32        `json:"wind_mph" form:"wind_mph" db:"wind_mph"`
	StartTime         types.DateTime `db:"start_time" json:"start_time"`
	StartTimeTimezone string         `form:"start_time_timezone" json:"start_time_timezone"`
	StartTimeDatetime string         `form:"start_time_datetime" json:"start_time_datetime"`
	Tournament        string         `db:"tournament" json:"tournament" form:"tournament"`
	TeamScore         int            `form:"team_score" db:"team_score" json:"team_score"`
	OpponentScore     int            `db:"opponent_score" json:"opponent_score" form:"opponent_score"`
	SoftCap           int            `db:"soft_cap" json:"soft_cap" form:"soft_cap"`
	HardCap           int            `db:"hard_cap" json:"hard_cap" form:"hard_cap"`
	Opponent          string         `db:"opponent" json:"opponent" form:"opponent"`
	HalfCap           int            `form:"half_cap" db:"half_cap" json:"half_cap"`
	TempF             float32        `db:"temp_f" json:"temp_f" form:"temp_f"`
	Status            GamesStatus    `db:"status" json:"status" form:"status"`
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

func (m *Games) GetStartTimeStr(format string) string {
	if dt, err := m.GetStartTimeDt(); err != nil || dt.IsZero() {
		return ""
	} else {
		return dt.Time().Format(format)
	}
}

func (m *Games) GetStartTimeDt() (types.DateTime, error) {
	if m.StartTimeDatetime != "" && m.StartTimeTimezone != "" {
		return types.ParseDateTime(m.StartTimeDatetime + ":00" + m.StartTimeTimezone)
	} else if m.StartTimeDatetime != "" {
		return types.ParseDateTime(m.StartTimeDatetime + ":00")
	} else {
		return m.StartTime, nil
	}
}

