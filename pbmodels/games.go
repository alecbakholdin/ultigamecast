package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Games struct {
	models.BaseModel

	OpponentScore     int            `db:"opponent_score" json:"opponent_score" form:"opponent_score"`
	HalfCap           int            `db:"half_cap" json:"half_cap" form:"half_cap"`
	HardCap           int            `db:"hard_cap" json:"hard_cap" form:"hard_cap"`
	StartTime         types.DateTime `db:"start_time" json:"start_time"`
	StartTimeTimezone string         `form:"start_time_timezone" json:"start_time_timezone"`
	StartTimeDatetime string         `form:"start_time_datetime" json:"start_time_datetime"`
	Status            GamesStatus    `db:"status" json:"status" form:"status"`
	Tournament        string         `form:"tournament" db:"tournament" json:"tournament"`
	Opponent          string         `db:"opponent" json:"opponent" form:"opponent"`
	WindMph           int            `db:"wind_mph" json:"wind_mph" form:"wind_mph"`
	TempF             int            `db:"temp_f" json:"temp_f" form:"temp_f"`
	TeamScore         int            `db:"team_score" json:"team_score" form:"team_score"`
	SoftCap           int            `db:"soft_cap" json:"soft_cap" form:"soft_cap"`
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

