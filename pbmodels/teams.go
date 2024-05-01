package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Teams struct {
	models.BaseModel

	Logo         string `db:"logo" json:"logo" form:"logo"`
	Organization string `form:"organization" db:"organization" json:"organization"`
	Name         string `db:"name" json:"name" form:"name"`
	Slug         string `db:"slug" json:"slug" param:"teamsSlug"`
	Managers     string `db:"managers" json:"managers" form:"managers"`
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

