package service

import (
	"context"
	"database/sql"
	"slices"
	"ultigamecast/app/ctxvar"
	"ultigamecast/models"
	"ultigamecast/service/slug"
)

type Player struct {
	q  *models.Queries
	db *sql.DB
}

func NewPlayer(q *models.Queries, db *sql.DB) *Player {
	return &Player{
		q:  q,
		db: db,
	}
}

func (p *Player) GetPlayer(ctx context.Context, slug string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	
}

func (p *Player) CreatePlayer(ctx context.Context, name string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	player, err := p.q.CreatePlayer(ctx, models.CreatePlayerParams{
		Team: team.ID,
		Name: name,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating player", err)
	}
	return &player, nil
}

func (p *Player) UpdatePlayer(ctx context.Context, name string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	existingPlayer := ctxvar.GetPlayer(ctx)

	teamPlayers, err := p.q.ListTeamPlayers(ctx, team.ID)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error fetching team players when renaming", err)
	}

	player, err := p.q.UpdatePlayer(ctx, models.CreatePlayerParams{
		Team: team.ID,
		Name: name,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating player", err)
	}
	return &player, nil
}

func (p * Player) nextAvailableSlug(ctx context.Context, teamId int64, name string) (string, error) {
	teamPlayers, err := p.q.ListTeamPlayers(ctx, teamId)
	if err != nil {
		return "", nil
	}
	newSlug := slug.From(name)
	num := 0
}