package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournaments struct {
	models.BaseModel

	End           types.DateTime `db:"end" json:"end"`
	EndTimezone   string         `db:"-" form:"end_timezone" json:"end_timezone"`
	EndDatetime   string         `db:"-" form:"end_datetime" json:"end_datetime"`
	Location      string         `json:"location" form:"location" db:"location"`
	Team          string         `db:"team" json:"team" form:"team"`
	Name          string         `db:"name" json:"name" form:"name"`
	Slug          string         `db:"slug" json:"slug" param:"tournamentsSlug"`
	Start         types.DateTime `db:"start" json:"start"`
	StartTimezone string         `db:"-" form:"start_timezone" json:"start_timezone"`
	StartDatetime string         `form:"start_datetime" json:"start_datetime" db:"-"`
}

func (m *Tournaments) TableName() string {
    return "tournaments"
}

func (m *Tournaments) GetEndStr(format string) string {
	if dt, err := m.GetEndDt(); err != nil || dt.IsZero() {
		return ""
	} else {
		return dt.Time().Format(format)
	}
}

func (m *Tournaments) GetEndDt() (types.DateTime, error) {
	if m.EndDatetime != "" && m.EndTimezone != "" {
		return types.ParseDateTime(m.EndDatetime + ":00" + m.EndTimezone)
	} else if m.EndDatetime != "" {
		return types.ParseDateTime(m.EndDatetime + ":00")
	} else {
		return m.End, nil
	}
}

func (m *Tournaments) GetStartStr(format string) string {
	if dt, err := m.GetStartDt(); err != nil || dt.IsZero() {
		return ""
	} else {
		return dt.Time().Format(format)
	}
}

func (m *Tournaments) GetStartDt() (types.DateTime, error) {
	if m.StartDatetime != "" && m.StartTimezone != "" {
		return types.ParseDateTime(m.StartDatetime + ":00" + m.StartTimezone)
	} else if m.StartDatetime != "" {
		return types.ParseDateTime(m.StartDatetime + ":00")
	} else {
		return m.Start, nil
	}
}

