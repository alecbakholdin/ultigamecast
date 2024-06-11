package service

import (
	"context"
	"database/sql"
	"errors"
	"ultigamecast/internal/app/service/slug"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
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

func (t *Team) IsTeamAdmin(ctx context.Context) bool {
	team := ctxvar.GetTeam(ctx)
	user := ctxvar.GetUser(ctx)
	return user != nil && team != nil && team.Owner == user.ID
}

func (t *Team) GetTeam(ctx context.Context, slug string) (*models.Team, error) {
	if team, err := t.q.GetTeam(ctx, slug); err != nil {
		return nil, convertAndLogSqlError(ctx, "error getting team by slug", err)
	} else {
		return &team, nil
	}
}

func (t *Team) GetTeams(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team = make([]models.Team, 0)
	user := ctxvar.GetUser(ctx)

	if ownedTeams, err := t.q.ListOwnedTeams(ctx, user.ID); err == nil {
		teams = append(teams, ownedTeams...)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return []models.Team{}, convertAndLogSqlError(ctx, "error getting owned teams", err)
	}

	if followedTeams, err := t.q.ListFollowedTeams(ctx, user.ID); err == nil {
		teams = append(teams, followedTeams...)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return []models.Team{}, convertAndLogSqlError(ctx, "error getting followed teams", err)
	}

	return teams, nil
}

func (t *Team) CreateTeam(ctx context.Context, name, organization string) (*models.Team, error) {
	slug := slug.From(name)
	if _, err := t.q.GetTeam(ctx, slug); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error fetching team while creating another team", err)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTeamExists
	}
	team, err := t.q.CreateTeam(ctx, models.CreateTeamParams{
		Owner:        ctxvar.GetUser(ctx).ID,
		Name:         name,
		Slug:         slug,
		Organization: sql.NullString{String: organization, Valid: organization != ""},
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating team", err)
	}
	return &team, nil
}

func (t *Team) UpdateName(ctx context.Context, name string) (*models.Team, error) {
	oldSlug := ctxvar.GetTeam(ctx).Slug
	slug := slug.From(name)
	if slug != oldSlug {
		if _, err := t.q.GetTeam(ctx, slug); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, convertAndLogSqlError(ctx, "error fetching team while updating another team", err)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamExists
		}
	}
	newTeam, err := t.q.UpdateTeamName(ctx, models.UpdateTeamNameParams{
		Name:    name,
		OldSlug: oldSlug,
		NewSlug: slug,
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error updating team", err)
	}
	return &newTeam, nil
}

func (t *Team) UpdateOrganization(ctx context.Context, organization string) (*models.Team, error) {
	team := ctxvar.GetTeam(ctx)
	if team.Organization.String == organization {
		return team, nil
	}
	updatedTeam, err := t.q.UpdateTeamOrganization(ctx, models.UpdateTeamOrganizationParams{
		Slug:         team.Slug,
		Organization: sql.NullString{String: organization, Valid: organization != ""},
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error updating organization", err)
	}
	return &updatedTeam, nil
}
