package service

import (
	"context"
	"database/sql"
	"errors"
	"ultigamecast/app/ctx_var"
	"ultigamecast/models"
	"ultigamecast/service/slug"
)

type Team struct {
	q  *models.Queries
	db *sql.DB
}

func NewTeam(q *models.Queries, db *sql.DB) *Team {
	return &Team{
		q:  q,
		db: db,
	}
}

var (
	ErrTeamExists = errors.New("a team with that name already exists")
)

func (t *Team) GetBySlug(ctx context.Context, slug string) (*models.Team, error) {
	if team, err := t.q.GetTeam(ctx, slug); err != nil {
		return nil, convertAndLogSqlError(ctx, "error getting team by slug", err)
	} else {
		return &team, nil
	}
}

func (t *Team) GetTeams(ctx context.Context) ([]*models.Team, error) {
	return []*models.Team{}, nil
}

func (t *Team) CreateTeam(ctx context.Context, name, organization string) (*models.Team, error) {
	slug := slug.From(name)
	if _, err := t.q.GetTeam(ctx, slug); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, convertAndLogSqlError(ctx, "error fetching team while creating another team", err)
	} else if !errors.Is(sql.ErrNoRows, err) {
		return nil, ErrTeamExists
	}
	team, err := t.q.CreateTeam(ctx, models.CreateTeamParams{
		Owner:        ctx_var.GetUser(ctx).ID,
		Name:         name,
		Slug:         slug,
		Organization: sql.NullString{String: organization, Valid: organization != ""},
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating team", err)
	}
	return &team, nil
}
