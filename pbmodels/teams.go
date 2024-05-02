package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Teams struct {
	models.BaseModel

	Organization string `db:"organization" form:"organization" json:"organization"`
	Name         string `db:"name" form:"name" json:"name"`
	Slug         string `db:"slug" json:"slug" param:"teamsSlug"`
	Managers     string `db:"managers" form:"managers" json:"managers"`
	Logo         string `db:"logo" form:"logo" json:"logo"`
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
