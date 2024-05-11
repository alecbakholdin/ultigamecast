package pbmodels

import (
	"cmp"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"strings"
	"time"
)

type Tournaments struct {
	models.BaseModel

	Team          string         `db:"team" form:"team" json:"team"`
	Start         types.DateTime `db:"start" json:"start"`
	StartTimezone string         `db:"-" form:"start_timezone" json:"start_timezone"`
	StartDatetime string         `db:"-" form:"start_datetime" json:"start_datetime"`
	End           types.DateTime `db:"end" json:"end"`
	EndTimezone   string         `db:"-" form:"end_timezone" json:"end_timezone"`
	EndDatetime   string         `db:"-" form:"end_datetime" json:"end_datetime"`
	Location      string         `db:"location" form:"location" json:"location"`
	Name          string         `db:"name" form:"name" json:"name"`
	Slug          string         `db:"slug" json:"slug" param:"tournamentsSlug"`
}


func (d *Tournaments) CopyFrom(s *Tournaments) *Tournaments {
	d.Team = s.Team
	d.Start = s.Start
	d.StartTimezone = s.StartTimezone
	d.StartDatetime = s.StartDatetime
	d.End = s.End
	d.EndTimezone = s.EndTimezone
	d.EndDatetime = s.EndDatetime
	d.Location = s.Location
	d.Name = s.Name
	d.Slug = s.Slug
	return d
}
func (m *Tournaments) TableName() string {
    return "tournaments"
}

func (m *Tournaments) GetStartStr(format string, locName string) string {
	if dt, err := m.GetStartDt(); err != nil || dt.IsZero() {
		return ""
	} else if loc, err := time.LoadLocation(locName); err != nil{
		return ""
	} else {
		return dt.Time().In(loc).Format(format)
	}
}

func (m *Tournaments) GetStartDt() (types.DateTime, error) {
	m.StartDatetime = strings.TrimSpace(m.StartDatetime)
	m.StartTimezone = strings.TrimSpace(m.StartTimezone)
	if m.StartDatetime != "" && m.StartTimezone != "" {
		var datetimeFormat string
		if len(m.StartDatetime) == len("2006-01-02") {
			datetimeFormat = "2006-01-02"
		} else {
			datetimeFormat = "2006-01-02T15:04"
		}
		if loc, err := time.LoadLocation(m.StartTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation(datetimeFormat, m.StartDatetime, cmp.Or(loc, time.Local)); err != nil {
			return types.DateTime{}, err
		} else {
			return types.ParseDateTime(time)
		}
	} else if m.StartDatetime != "" {
		return types.ParseDateTime(m.StartDatetime + ":00")
	} else {
		return m.Start, nil
	}
}

func (m *Tournaments) GetEndStr(format string, locName string) string {
	if dt, err := m.GetEndDt(); err != nil || dt.IsZero() {
		return ""
	} else if loc, err := time.LoadLocation(locName); err != nil{
		return ""
	} else {
		return dt.Time().In(loc).Format(format)
	}
}

func (m *Tournaments) GetEndDt() (types.DateTime, error) {
	m.EndDatetime = strings.TrimSpace(m.EndDatetime)
	m.EndTimezone = strings.TrimSpace(m.EndTimezone)
	if m.EndDatetime != "" && m.EndTimezone != "" {
		var datetimeFormat string
		if len(m.EndDatetime) == len("2006-01-02") {
			datetimeFormat = "2006-01-02"
		} else {
			datetimeFormat = "2006-01-02T15:04"
		}
		if loc, err := time.LoadLocation(m.EndTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation(datetimeFormat, m.EndDatetime, cmp.Or(loc, time.Local)); err != nil {
			return types.DateTime{}, err
		} else {
			return types.ParseDateTime(time)
		}
	} else if m.EndDatetime != "" {
		return types.ParseDateTime(m.EndDatetime + ":00")
	} else {
		return m.End, nil
	}
}
