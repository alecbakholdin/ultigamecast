package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"ultigamecast/internal/app/service/slug"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
)

type Player struct {
	q  *models.Queries
	db *sql.DB
}

func NewPlayer(db *sql.DB) *Player {
	return &Player{
		q:  models.New(db),
		db: db,
	}
}

func (p *Player) GetPlayer(ctx context.Context, slug string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	player, err := p.q.GetPlayer(ctx, models.GetPlayerParams{
		TeamId: team.ID,
		Slug:   slug,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error getting player", err)
	}
	return &player, nil
}

func (p *Player) GetTeamPlayerMap(ctx context.Context) (map[int64]models.Player, error) {
	players, err := p.GetTeamPlayers(ctx)
	if err != nil {
		return nil, err
	}
	playerMap := make(map[int64]models.Player, len(players))
	for _, p := range players {
		playerMap[p.ID] = p
	}
	return playerMap, nil
}

func (p *Player) GetTeamPlayers(ctx context.Context) ([]models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	players, err := p.q.ListTeamPlayers(ctx, team.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error fetching team players", err)
	}
	return players, nil
}

func (p *Player) CreatePlayer(ctx context.Context, name string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	slug, err := p.nextAvailableSlug(ctx, 0, team.ID, name)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error finding next available slug", err)
	}
	player, err := p.q.CreatePlayer(ctx, models.CreatePlayerParams{
		Team: team.ID,
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating player", err)
	}
	return &player, nil
}

func (p *Player) UpdatePlayer(ctx context.Context, name string) (*models.Player, error) {
	team := ctxvar.GetTeam(ctx)
	existingPlayer := ctxvar.GetPlayer(ctx)

	slug, err := p.nextAvailableSlug(ctx, existingPlayer.ID, team.ID, name)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error fetching slug for existing player", err)
	}

	player, err := p.q.UpdatePlayer(ctx, models.UpdatePlayerParams{
		ID:   existingPlayer.ID,
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating player", err)
	}
	return &player, nil
}

func (p *Player) nextAvailableSlug(ctx context.Context, playerId int64, teamId int64, name string) (string, error) {
	teamPlayers, err := p.q.ListTeamPlayers(ctx, teamId)
	if err != nil {
		return "", nil
	}
	newSlug := slug.From(name)
	num := 0
	for slices.ContainsFunc(teamPlayers, func(p models.Player) bool { return p.Slug == newSlug && p.ID != playerId }) {
		num++
		newSlug = fmt.Sprintf("%s-%d", slug.From(name), num+1)
	}
	return newSlug, nil
}

func (p *Player) UpdatePlayerOrder(ctx context.Context, playerIds []int64) error {
	tx, err := p.db.Begin()
	if err != nil {
		return convertAndLogSqlError(ctx, "error starting order transaction", err)
	}
	defer tx.Rollback()
	q := p.q.WithTx(tx)

	team := ctxvar.GetTeam(ctx)
	if teamPlayers, err := q.ListTeamPlayers(ctx, team.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return convertAndLogSqlError(ctx, "error fetching team players", err)
	} else if len(teamPlayers) != len(playerIds) {
		return errors.Join(ErrBadFormat, fmt.Errorf("expected %d players when updating order but found %d", len(teamPlayers), len(playerIds)))
	}

	for i, id := range playerIds {
		if err := q.UpdatePlayerOrder(ctx, models.UpdatePlayerOrderParams{
			Team:  team.ID,
			Order: int64(i),
			ID:    id,
		}); err != nil {
			return convertAndLogSqlError(ctx, "error updating player order", err)
		}
	}

	return convertAndLogSqlError(ctx, "error committing transaction", tx.Commit())
}
