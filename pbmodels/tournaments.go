package pbmodels

import (
	"cmp"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Tournaments struct {
	models.BaseModel

	Start         types.DateTime `db:"start" json:"start"`
	StartTimezone string         `db:"-" form:"start_timezone" json:"start_timezone"`
	StartDatetime string         `db:"-" form:"start_datetime" json:"start_datetime"`
	End           types.DateTime `db:"end" json:"end"`
	EndTimezone   string         `db:"-" form:"end_timezone" json:"end_timezone"`
	EndDatetime   string         `db:"-" form:"end_datetime" json:"end_datetime"`
	Location      string         `db:"location" form:"location" json:"location"`
	Team          string         `db:"team" form:"team" json:"team"`
	Name          string         `db:"name" form:"name" json:"name"`
	Slug          string         `db:"slug" json:"slug" param:"tournamentsSlug"`
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
	if m.StartDatetime != "" && m.StartTimezone != "" {
		if loc, err := time.LoadLocation(m.StartTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation("2006-01-02T15:04", m.StartDatetime, cmp.Or(loc, time.Local)); err != nil {
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
	if m.EndDatetime != "" && m.EndTimezone != "" {
		if loc, err := time.LoadLocation(m.EndTimezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation("2006-01-02T15:04", m.EndDatetime, cmp.Or(loc, time.Local)); err != nil {
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
