package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type PlayerGameSummary struct {
	models.BaseModel

	PlayerId   string  `db:"player_id" form:"player_id" json:"player_id"`
	GameId     string  `db:"game_id" form:"game_id" json:"game_id"`
	TeamId     string  `db:"team_id" form:"team_id" json:"team_id"`
	Assists    float32 `db:"assists" form:"assists" json:"assists"`
	Turns      float32 `db:"turns" form:"turns" json:"turns"`
	Drops      float32 `db:"drops" form:"drops" json:"drops"`
	PlayerName string  `db:"player_name" form:"player_name" json:"player_name"`
	Points     float32 `db:"points" form:"points" json:"points"`
	Goals      float32 `db:"goals" form:"goals" json:"goals"`
}

func (m *PlayerGameSummary) TableName() string {
    return "player_game_summary"
}
