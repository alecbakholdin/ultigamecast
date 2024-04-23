package modelspb

import (
	"ultigamecast/validation"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GamesDto struct {
	TeamSlug          string `param:"teamSlug"`
	TournamentSlug    string `param:"tournamentSlug"`
	GameID            string `param:"gameId"`
	GameOpponent      string `form:"opponent"`
	GameTeamScore     int    `form:"team_score"`
	GameOpponentScore int    `form:"opponent_score"`
	GameHalfCap       int    `form:"half_cap"`
	GameSoftCap       int    `form:"soft_cap"`
	GameHardCap       int    `form:"hard_cap"`
	GameWindMph       int    `form:"wind_mph"`
	GameTempF         int    `form:"temp_f"`
	GameStartTime     string `form:"start_time"`
	GameStartTimeDt   types.DateTime
	GamesIsCompleted  bool `form:"is_completed"`
}

func BindGameDto(c echo.Context, dto *GamesDto) (err error) {
	if err := c.Bind(dto); err != nil {
		return err
	}

	if dto.GameStartTime != "" {
		if dto.GameStartTimeDt, err = types.ParseDateTime(dto.GameStartTime); err != nil {
			validation.AddFieldErrorString(c, "start_time", "invalid date format")
		}
	}

	return nil
}

func DtoFromGame(game *Games) *GamesDto {
	return &GamesDto{
		GameID:            game.Record.GetId(),
		GameOpponent:      game.GetOpponent(),
		GameTeamScore:     game.GetTeamScore(),
		GameOpponentScore: game.GetOpponentScore(),
		GameHalfCap:       game.GetHalfCap(),
		GameSoftCap:       game.GetSoftCap(),
		GameHardCap:       game.GetHardCap(),
		GameWindMph:       game.GetWindMph(),
		GameTempF:         game.GetTempF(),
	}
}
