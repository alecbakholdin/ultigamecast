package ctxvar

import (
	"context"
	"fmt"
	"ultigamecast/models"
)

type ContextVar string

const (
	HttpMethod ContextVar = "method"
	Path       ContextVar = "path"
	ReqId      ContextVar = "RequestId"
	User       ContextVar = "User"
	Team       ContextVar = "Team"
	Tournament ContextVar = "Tournament"
	Game       ContextVar = "Game"
)

var LogMessageVars = []ContextVar{HttpMethod, Path}
var LogAttrVars = []ContextVar{ReqId, User}

func IsAuthenticated(ctx context.Context) bool {
	return GetUser(ctx) != nil
}

func GetUser(ctx context.Context) *models.User {
	return getModel[models.User](ctx, User)
}

func GetTeam(ctx context.Context) *models.Team {
	return getModel[models.Team](ctx, Team)
}

func GetTournament(ctx context.Context) *models.Tournament {
	return getModel[models.Tournament](ctx, Tournament)
}

func GetGame(ctx context.Context) *models.Game {
	return getModel[models.Game](ctx, Game)
}

func getModel[T interface{}](ctx context.Context, key ContextVar) *T {
	if m, ok := ctx.Value(key).(*T); ok {
		return m
	}
	return nil
}

func GetValue(ctx context.Context, key ContextVar) string {
	switch key {
	case User:
		if user := GetUser(ctx); user != nil {
			return fmt.Sprintf("[%d] %s", user.ID, user.Email)
		}
	case Team:
		if team := GetTeam(ctx); team != nil {
			return fmt.Sprintf("[%d] %s", team.ID, team.Name)
		}
	case Tournament:
		if tournament := GetTournament(ctx); tournament != nil {
			return fmt.Sprintf("[%d] %s", tournament.ID, tournament.Name)
		}
	case Game:
		if game := GetGame(ctx); game != nil {
			return fmt.Sprintf("[%d] %s", game.ID, game.Opponent)
		}

	default:
		if val, ok := ctx.Value(key).(string); ok {
			return val
		}
	}
	return ""
}
