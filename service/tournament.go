package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
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

// edits a tournament provided in context using dates in the format 2024-01-02 - 2024-01-02. Returns [ErrBadFormat]
// if this is improperly formatted or start is after end
func (t *Tournament) UpdateTournamentDates(ctx context.Context, dates string) (*models.Tournament, error) {
	dateSlice := strings.Split(dates, " - ")
	if len(dateSlice) != 2 {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid number of dates in dates string: %d", len(dateSlice)))
		return nil, errors.Join(ErrBadFormat, fmt.Errorf("invalid number of dates in dates string: %d", len(dateSlice)))
	}
	start, startErr := time.Parse(t.DateFormat(), strings.TrimSpace(dateSlice[0]))
	if startErr != nil {
		err := errors.Join(ErrBadFormat, fmt.Errorf("error parsing start date: %w", startErr))
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	end, endErr := time.Parse(t.DateFormat(), strings.TrimSpace(dateSlice[1]))
	if endErr != nil {
		err := errors.Join(ErrBadFormat, fmt.Errorf("error parsing end date: %w", endErr))
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	if start.After(end) {
		return nil, errors.Join(ErrBadFormat, errors.New("start is after end"))
	}
	tournament, err := t.q.UpdateTournamentDates(ctx, models.UpdateTournamentDatesParams{
		StartDate: sql.NullTime{Time: start, Valid: true},
		EndDate: sql.NullTime{Time: end, Valid: true},
		TournamentId: ctxvar.GetTournament(ctx).ID,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error updating tournament dates", err)
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

func (t *Tournament) DateFormat() string {
	return "Jan 2, 2006"
}