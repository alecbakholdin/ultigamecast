package service

import (
	"context"
	"database/sql"
	"ultigamecast/app/ctxvar"
	"ultigamecast/models"
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

func CreateTournament(ctx context.Context) {

}
