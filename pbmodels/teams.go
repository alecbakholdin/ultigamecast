package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Teams struct {
	models.BaseModel

	Name         string `json:"name" form:"name" db:"name"`
	Slug         string `db:"slug" json:"slug" param:"teamsSlug"`
	Managers     string `json:"managers" form:"managers" db:"managers"`
	Logo         string `db:"logo" json:"logo" form:"logo"`
	Organization string `json:"organization" form:"organization" db:"organization"`
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

