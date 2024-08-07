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
	"ultigamecast/internal/app/service/slug"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
)

type Tournament struct {
	q  *models.Queries
	db *sql.DB
}

func NewTournament(db *sql.DB) *Tournament {
	return &Tournament{
		q:  models.New(db),
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

func (t *Tournament) GetSchedule(ctx context.Context) ([]models.TournamentSummary, error) {
	team := ctxvar.GetTeam(ctx)
	tournaments, err := t.q.ListTournaments(ctx, team.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error getting team tournaments", err)
	}
	games, err := t.q.ListTeamGames(ctx, team.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error getting team games", err)
	}
	data, err := t.q.ListTeamTournamentData(ctx, team.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error getting team tournament data", err)
	}

	tg := make([]models.TournamentSummary, len(tournaments))
	for i, t := range tournaments {
		tg[i] = models.TournamentSummary{
			Tournament: &t,
			Games:      []models.Game{},
			Data:       []models.TournamentDatum{},
		}
		for _, g := range games {
			if g.Tournament == t.ID {
				tg[i].Games = append(tg[i].Games, g)
			}
		}
		for _, td := range data {
			if td.Tournament == t.ID {
				tg[i].Data = append(tg[i].Data, td)
			}
		}
	}

	return tg, nil
}

func (t *Tournament) CreateTournament(ctx context.Context, name, dates string) (*models.TournamentSummary, error) {
	start, end, err := t.parseDates(ctx, dates)
	if err != nil {
		return nil, err
	}
	slug, err := t.getSafeSlug(ctx, -1, name)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating safe slug", err)
	}
	team := ctxvar.GetTeam(ctx)
	tournament, err := t.q.CreateTournament(ctx, models.CreateTournamentParams{
		TeamId: team.ID,
		Slug:   slug,
		Name:   name,
		StartDate: sql.NullTime{Time: start, Valid: !start.IsZero()},
		EndDate: sql.NullTime{Time: end, Valid: !end.IsZero()},
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating tournament", err)
	}
	return &models.TournamentSummary{
		Tournament: &tournament,
		Games: []models.Game{},
		Data: []models.TournamentDatum{},
	}, nil
}

// edits a tournament provided in context using dates in the format DateFormat(). Returns [ErrBadFormat]
// if this is improperly formatted or start is after end
func (t *Tournament) UpdateTournamentDates(ctx context.Context, dates string) (*models.Tournament, error) {
	start, end, err := t.parseDates(ctx, dates)
	if err != nil {
		return nil, err
	}
	tournament, err := t.q.UpdateTournamentDates(ctx, models.UpdateTournamentDatesParams{
		StartDate:    sql.NullTime{Time: start, Valid: !start.IsZero()},
		EndDate:      sql.NullTime{Time: end, Valid: !end.IsZero()},
		TournamentId: ctxvar.GetTournament(ctx).ID,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error updating tournament dates", err)
	}
	return &tournament, nil
}

func (t *Tournament) parseDates(ctx context.Context, dates string) (time.Time, time.Time, error) {
	if dates == "" {
		return time.Time{}, time.Time{}, nil
	}
	dateSlice := strings.Split(dates, " - ")
	if len(dateSlice) != 2 {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid number of dates in dates string: %d", len(dateSlice)))
		return time.Time{}, time.Time{}, errors.Join(ErrBadFormat, fmt.Errorf("invalid number of dates in dates string: %d", len(dateSlice)))
	}
	start, startErr := time.Parse(t.DateFormat(), strings.TrimSpace(dateSlice[0]))
	if startErr != nil {
		err := errors.Join(ErrBadFormat, fmt.Errorf("error parsing start date: %w", startErr))
		slog.ErrorContext(ctx, err.Error())
		return time.Time{}, time.Time{}, err
	}
	end, endErr := time.Parse(t.DateFormat(), strings.TrimSpace(dateSlice[1]))
	if endErr != nil {
		err := errors.Join(ErrBadFormat, fmt.Errorf("error parsing end date: %w", endErr))
		slog.ErrorContext(ctx, err.Error())
		return time.Time{}, time.Time{}, err
	}
	if start.After(end) {
		return time.Time{}, time.Time{}, errors.Join(ErrBadFormat, errors.New("start is after end"))
	}
	return start, end, nil
}

// get all the tournament data for the tournament found in ctx
func (t *Tournament) Data(ctx context.Context) ([]models.TournamentDatum, error) {
	tournament := ctxvar.GetTournament(ctx)
	data, err := t.q.ListTournamentData(ctx, tournament.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error fetching tournament data", err)
	}
	return data, nil
}

// create a new datum with default values for the tournament in ctx
func (t *Tournament) NewDatum(ctx context.Context) (*models.TournamentDatum, error) {
	if datum, err := t.q.CreateTournamentDatum(ctx, ctxvar.GetTournament(ctx).ID); err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating datum", err)
	} else {
		return &datum, err
	}
}

func (t *Tournament) UpdateDataOrder(ctx context.Context, ids []int64) (error) {
	tid := ctxvar.GetTournament(ctx).ID
	tx, err := t.db.Begin()
	if err != nil {
		return convertAndLogSqlError(ctx, "error opening tx", err)
	}
	db := t.q.WithTx(tx)
	defer tx.Rollback()
	for i, id := range ids {
		if _, err = db.GetTournamentDatum(ctx, models.GetTournamentDatumParams{DataId: id, TournamentId: tid}); err != nil {
			return convertAndLogSqlError(ctx, fmt.Sprintf("error updating order of id %d to %d", id, i), err)
		}
		if err = db.UpdateTournamentDatumOrder(ctx, models.UpdateTournamentDatumOrderParams{
			Order: int64(i),
			DataId: id,
			TournamentId: tid,
		}); err != nil {
			return convertAndLogSqlError(ctx, fmt.Sprintf("error updating order of id %d to %d", id, i), err)
		}
	}
	if err := tx.Commit(); err != nil {
		return convertAndLogSqlError(ctx, "error committing tx", err)
	}

	return nil
}

func (t *Tournament) getSafeSlug(ctx context.Context, tournamentId int64, name string) (string, error) {
	team := ctxvar.GetTeam(ctx)
	tournaments, err := t.q.ListTournaments(ctx, team.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	s := slug.From(name)
	num := 2
	for slices.ContainsFunc(tournaments, func(to models.Tournament) bool { return to.ID != tournamentId && to.Slug == s }) {
		s = fmt.Sprintf("%s-%d", slug.From(name), num)
		num++
	}
	return s, nil
}

func (t *Tournament) DateFormat() string {
	return "Jan 2, 2006"
}
