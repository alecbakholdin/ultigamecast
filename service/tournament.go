package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"ultigamecast/app/ctxvar"
	"ultigamecast/models"
	"ultigamecast/service/slug"
)

type Tournament struct {
	q  *models.Queries
	db *sql.DB
}

func NewTournament(q *models.Queries, db *sql.DB) *Tournament {
	return &Tournament{
		q:  q,
		db: db,
	}
}

func (t *Tournament) GetTournament(ctx context.Context, slug string) (*models.Tournament, error) {
	team := ctxvar.GetTeam(ctx)
	tournament, err := t.q.GetTournament(ctx, models.GetTournamentParams{
		TeamId: team.ID,
		Slug:   slug,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error fetching tournament", err)
	}
	return &tournament, nil
}

func (t *Tournament) GetTeamTournaments(ctx context.Context) ([]models.Tournament, error) {
	team := ctxvar.GetTeam(ctx)
	tournaments, err := t.q.ListTournaments(ctx, team.ID)
	if err != nil && !errors.Is(sql.ErrNoRows, err){
		return nil, convertAndLogSqlError(ctx, "error fetching team tournaments", err)
	}
	return tournaments, nil
}

func (t *Tournament) CreateTournament(ctx context.Context, name string) (*models.Tournament, error) {
	slug, err := t.getSafeSlug(ctx, -1, name)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating safe slug", err)
	}
	team := ctxvar.GetTeam(ctx)
	tournament, err := t.q.CreateTournament(ctx, models.CreateTournamentParams{
		TeamId: team.ID,
		Slug: slug,
		Name: name,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating tournament", err)
	}
	return &tournament, nil
}

func (t *Tournament) getSafeSlug(ctx context.Context, tournamentId int64, name string) (string, error) {
	team := ctxvar.GetTeam(ctx)
	tournaments, err := t.q.ListTournaments(ctx, team.ID)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", err
	}
	s := slug.From(name)
	num := 2
	for slices.ContainsFunc(tournaments, func(to models.Tournament) bool {return to.ID != tournamentId && to.Slug == s}) {
		s = fmt.Sprintf("%s-%d", slug.From(name), num)
		num++
	}
	return s, nil
}
