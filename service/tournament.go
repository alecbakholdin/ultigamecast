package service

import (
	"context"
	"database/sql"
	"ultigamecast/models"
)

type Tournament struct {
	q *models.Queries
	db *sql.DB
}

func NewTournament(q *models.Queries, db *sql.DB) *Tournament {
	return &Tournament{
		q: q,
		db: db,
	}
}

func CreateTournament(ctx context.Context) {

}