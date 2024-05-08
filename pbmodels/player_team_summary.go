package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type PlayerTeamSummary struct {
	models.BaseModel

	Drops      float32 `db:"drops" form:"drops" json:"drops"`
	PlayerId   string  `db:"player_id" form:"player_id" json:"player_id"`
	TeamSlug   string  `db:"team_slug" form:"team_slug" json:"team_slug"`
	PlayerName string  `db:"player_name" form:"player_name" json:"player_name"`
	Assists    float32 `db:"assists" form:"assists" json:"assists"`
	Turns      float32 `db:"turns" form:"turns" json:"turns"`
	TeamId     string  `db:"team_id" form:"team_id" json:"team_id"`
	Points     float32 `db:"points" form:"points" json:"points"`
	Goals      float32 `db:"goals" form:"goals" json:"goals"`
}

func (m *PlayerTeamSummary) TableName() string {
    return "player_team_summary"
}
