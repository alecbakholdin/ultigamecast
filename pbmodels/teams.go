package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Teams struct {
	models.BaseModel

	Slug         string `json:"slug" param:"teamsSlug" db:"slug"`
	Managers     string `db:"managers" json:"managers" form:"managers"`
	Logo         string `json:"logo" form:"logo" db:"logo"`
	Organization string `db:"organization" json:"organization" form:"organization"`
	Name         string `db:"name" json:"name" form:"name"`
}

func (m *Teams) TableName() string {
    return "teams"
}

func (m *Teams) GetLogoPath() string {
	if m.Logo == "" || m.GetId() == "" {
		return ""
	} else {
		return path.Join("846ykkxqtaqjxst", m.GetId(), m.Logo)
	}
}

