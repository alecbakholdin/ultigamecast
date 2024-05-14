package pbmodels

import (
	"cmp"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"strings"
	"time"
)

type Games struct {
	models.BaseModel

	Opponent          string          `db:"opponent" form:"opponent" json:"opponent"`
	Status            GamesStatus     `db:"status" form:"status" json:"status"`
	OpponentScore     int             `db:"opponent_score" form:"opponent_score" json:"opponent_score"`
	HardCap           int             `db:"hard_cap" form:"hard_cap" json:"hard_cap"`
	WindMph           int             `db:"wind_mph" form:"wind_mph" json:"wind_mph"`
	Tournament        string          `db:"tournament" form:"tournament" json:"tournament"`
	StartTime         types.DateTime  `db:"start_time" json:"start_time"`
	StartTimeTimezone string          `db:"-" form:"start_time_timezone" json:"start_time_timezone"`
	StartTimeDatetime string          `db:"-" form:"start_time_datetime" json:"start_time_datetime"`
	SoftCap           int             `db:"soft_cap" form:"soft_cap" json:"soft_cap"`
	TempF             int             `db:"temp_f" form:"temp_f" json:"temp_f"`
	LiveStatus        GamesLiveStatus `db:"live_status" form:"live_status" json:"live_status"`
	TeamScore         int             `db:"team_score" form:"team_score" json:"team_score"`
	HalfCap           int             `db:"half_cap" form:"half_cap" json:"half_cap"`
}

type GamesStatus string

const (
	GamesStatusScheduled GamesStatus = "scheduled"
	GamesStatusLive      GamesStatus = "live"
	GamesStatusCompleted GamesStatus = "completed"
)

type GamesLiveStatus string

const (
	GamesLiveStatusBeforePoint         GamesLiveStatus = "Before Point"
	GamesLiveStatusHalftime            GamesLiveStatus = "Halftime"
	GamesLiveStatusFinal               GamesLiveStatus = "Final"
	GamesLiveStatusTeamPossession      GamesLiveStatus = "Team Possession"
	GamesLiveStatusOpponentPossession  GamesLiveStatus = "Opponent Possession"
)


func (d *Games) CopyFrom(s *Games) *Games {
	d.Opponent = s.Opponent
	d.Status = s.Status
	d.OpponentScore = s.OpponentScore
	d.HardCap = s.HardCap
	d.WindMph = s.WindMph
	d.Tournament = s.Tournament
	d.StartTime = s.StartTime
	d.StartTimeTimezone = s.StartTimeTimezone
	d.StartTimeDatetime = s.StartTimeDatetime
	d.SoftCap = s.SoftCap
	d.TempF = s.TempF
	d.LiveStatus = s.LiveStatus
	d.TeamScore = s.TeamScore
	d.HalfCap = s.HalfCap
	return d
}

func (m *Games) Copy() *Games {
	return (&Games{}).CopyFrom(m)
}

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
	m.StartTimeDatetime = strings.TrimSpace(m.StartTimeDatetime)
	m.StartTimeTimezone = strings.TrimSpace(m.StartTimeTimezone)
	if m.StartTimeDatetime != "" && m.StartTimeTimezone != "" {
		var datetimeFormat string
		if len(m.StartTimeDatetime) == len("2006-01-02") {
			datetimeFormat = "2006-01-02"
		} else {
			datetimeFormat = "2006-01-02T15:04"
		}
		if loc, err := time.LoadLocation(m.StartTimeTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation(datetimeFormat, m.StartTimeDatetime, cmp.Or(loc, time.Local)); err != nil {
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
