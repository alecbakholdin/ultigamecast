package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type PlayerTournamentSummary struct {
	models.BaseModel

	TournamentId string  `db:"tournament_id" form:"tournament_id" json:"tournament_id"`
	PlayerName   string  `db:"player_name" form:"player_name" json:"player_name"`
	Points       float32 `db:"points" form:"points" json:"points"`
	Goals        float32 `db:"goals" form:"goals" json:"goals"`
	Assists      float32 `db:"assists" form:"assists" json:"assists"`
	Turns        float32 `db:"turns" form:"turns" json:"turns"`
	Drops        float32 `db:"drops" form:"drops" json:"drops"`
	PlayerId     string  `db:"player_id" form:"player_id" json:"player_id"`
	TeamId       string  `db:"team_id" form:"team_id" json:"team_id"`
}

func (m *PlayerTournamentSummary) TableName() string {
    return "player_tournament_summary"
}
