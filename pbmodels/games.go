package pbmodels

import (
	"cmp"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Games struct {
	models.BaseModel

	Opponent          string         `db:"opponent" form:"opponent" json:"opponent"`
	TeamScore         int            `db:"team_score" form:"team_score" json:"team_score"`
	OpponentScore     int            `db:"opponent_score" form:"opponent_score" json:"opponent_score"`
	HalfCap           int            `db:"half_cap" form:"half_cap" json:"half_cap"`
	WindMph           int            `db:"wind_mph" form:"wind_mph" json:"wind_mph"`
	StartTime         types.DateTime `db:"start_time" json:"start_time"`
	StartTimeTimezone string         `db:"-" form:"start_time_timezone" json:"start_time_timezone"`
	StartTimeDatetime string         `db:"-" form:"start_time_datetime" json:"start_time_datetime"`
	Status            GamesStatus    `db:"status" form:"status" json:"status"`
	Tournament        string         `db:"tournament" form:"tournament" json:"tournament"`
	SoftCap           int            `db:"soft_cap" form:"soft_cap" json:"soft_cap"`
	HardCap           int            `db:"hard_cap" form:"hard_cap" json:"hard_cap"`
	TempF             int            `db:"temp_f" form:"temp_f" json:"temp_f"`
}

type GamesStatus string

const (
	GamesStatusScheduled GamesStatus = "scheduled"
	GamesStatusLive      GamesStatus = "live"
	GamesStatusCompleted GamesStatus = "completed"
)

func (m *Games) TableName() string {
    return "games"
}

func (m *Games) GetStartTimeStr(format string, locName string) string {
	if dt, err := m.GetStartTimeDt(); err != nil || dt.IsZero() {
		return ""
	} else if loc, err := time.LoadLocation(locName); err != nil{
		return ""
	} else {
		return dt.Time().In(loc).Format(format)
	}
}

func (m *Games) GetStartTimeDt() (types.DateTime, error) {
	if m.StartTimeDatetime != "" && m.StartTimeTimezone != "" {
		if loc, err := time.LoadLocation(m.StartTimeTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation("2006-01-02T15:04", m.StartTimeDatetime, cmp.Or(loc, time.Local)); err != nil {
			return types.DateTime{}, err
		} else {
			return types.ParseDateTime(time)
		}
	} else if m.StartTimeDatetime != "" {
		return types.ParseDateTime(m.StartTimeDatetime + ":00")
	} else {
		return m.StartTime, nil
	}
}
