// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package models

import (
	"context"
	"database/sql"
)

const createPlayer = `-- name: CreatePlayer :one
INSERT INTO players (team, slug, "name", "order")
VALUES (
        ?1,
        ?2,
        ?3,
        (
            SELECT 1 + IFNULL(MAX("order"), -1)
            FROM players p
            WHERE p.team = ?1
        )
    )
RETURNING id, team, slug, name, "order"
`

type CreatePlayerParams struct {
	Team int64  `db:"team" json:"team"`
	Slug string `db:"slug" json:"slug"`
	Name string `db:"name" json:"name" validate:"required,max=64"`
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (Player, error) {
	row := q.db.QueryRowContext(ctx, createPlayer, arg.Team, arg.Slug, arg.Name)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Team,
		&i.Slug,
		&i.Name,
		&i.Order,
	)
	return i, err
}

const createTeam = `-- name: CreateTeam :one
INSERT INTO teams ("owner", "name", "slug", "organization")
VALUES (?, ?, ?, ?)
RETURNING id, owner, name, slug, organization
`

type CreateTeamParams struct {
	Owner        int64          `db:"owner" json:"owner"`
	Name         string         `db:"name" json:"name" validate:"required,max=64"`
	Slug         string         `db:"slug" json:"slug"`
	Organization sql.NullString `db:"organization" json:"organization"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, createTeam,
		arg.Owner,
		arg.Name,
		arg.Slug,
		arg.Organization,
	)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Slug,
		&i.Organization,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (LOWER(?), ?)
RETURNING id, email, password_hash
`

type CreateUserParams struct {
	Email        string         `db:"email" json:"email"`
	PasswordHash sql.NullString `db:"password_hash" json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash)
	return i, err
}

const getPlayer = `-- name: GetPlayer :one
SELECT id, team, slug, name, "order"
FROM players
WHERE team = ?1
    AND slug = ?2
`

type GetPlayerParams struct {
	TeamId int64  `db:"teamId" json:"teamId"`
	Slug   string `db:"slug" json:"slug"`
}

func (q *Queries) GetPlayer(ctx context.Context, arg GetPlayerParams) (Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayer, arg.TeamId, arg.Slug)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Team,
		&i.Slug,
		&i.Name,
		&i.Order,
	)
	return i, err
}

const getTeam = `-- name: GetTeam :one
SELECT id, owner, name, slug, organization
FROM teams
WHERE slug = LOWER(?1)
`

func (q *Queries) GetTeam(ctx context.Context, slug string) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeam, slug)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Slug,
		&i.Organization,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, password_hash
FROM users
WHERE email = LOWER(?1)
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash)
	return i, err
}

const listFollowedTeams = `-- name: ListFollowedTeams :many
SELECT t.id, t.owner, t.name, t.slug, t.organization
FROM team_follow tf
    INNER JOIN teams t ON t.id = tf.team
WHERE tf.user = ?1
ORDER BY t.name ASC
`

func (q *Queries) ListFollowedTeams(ctx context.Context, userid int64) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, listFollowedTeams, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Team
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Name,
			&i.Slug,
			&i.Organization,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOwnedTeams = `-- name: ListOwnedTeams :many
SELECT id, owner, name, slug, organization
FROM teams
WHERE teams.owner = ?1
ORDER BY id DESC
`

func (q *Queries) ListOwnedTeams(ctx context.Context, userid int64) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, listOwnedTeams, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Team
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Name,
			&i.Slug,
			&i.Organization,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTeamPlayers = `-- name: ListTeamPlayers :many
SELECT id, team, slug, name, "order"
FROM players
WHERE team = ?1
ORDER BY p.order ASC
`

func (q *Queries) ListTeamPlayers(ctx context.Context, teamid int64) ([]Player, error) {
	rows, err := q.db.QueryContext(ctx, listTeamPlayers, teamid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.Team,
			&i.Slug,
			&i.Name,
			&i.Order,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTournaments = `-- name: ListTournaments :many
SELECT tournaments.id, tournaments.team, tournaments.name, tournaments.slug, tournaments.start_date, tournaments.end_date, tournaments.location
FROM tournaments
    INNER JOIN teams ON tournaments.team = teams.id
WHERE teams.slug = LOWER(?1)
`

func (q *Queries) ListTournaments(ctx context.Context, slug string) ([]Tournament, error) {
	rows, err := q.db.QueryContext(ctx, listTournaments, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tournament
	for rows.Next() {
		var i Tournament
		if err := rows.Scan(
			&i.ID,
			&i.Team,
			&i.Name,
			&i.Slug,
			&i.StartDate,
			&i.EndDate,
			&i.Location,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePlayer = `-- name: UpdatePlayer :one
UPDATE players
SET "name" = ?, slug = ?
WHERE id = ?
RETURNING id, team, slug, name, "order"
`

type UpdatePlayerParams struct {
	Name string `db:"name" json:"name" validate:"required,max=64"`
	Slug string `db:"slug" json:"slug"`
	ID   int64  `db:"id" json:"id"`
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) (Player, error) {
	row := q.db.QueryRowContext(ctx, updatePlayer, arg.Name, arg.Slug, arg.ID)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Team,
		&i.Slug,
		&i.Name,
		&i.Order,
	)
	return i, err
}

const updateTeam = `-- name: UpdateTeam :one
UPDATE teams
SET "name" = ?,
    "slug" = ?,
    organization = ?
WHERE teams.slug = ?
RETURNING id, owner, name, slug, organization
`

type UpdateTeamParams struct {
	Name         string         `db:"name" json:"name" validate:"required,max=64"`
	Slug         string         `db:"slug" json:"slug"`
	Organization sql.NullString `db:"organization" json:"organization"`
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeam,
		arg.Name,
		arg.Slug,
		arg.Organization,
		arg.Slug,
	)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Slug,
		&i.Organization,
	)
	return i, err
}

const updateTeamOwner = `-- name: UpdateTeamOwner :one
UPDATE teams
SET "owner" = ?
WHERE teams.slug = ?
RETURNING id, owner, name, slug, organization
`

type UpdateTeamOwnerParams struct {
	Owner int64  `db:"owner" json:"owner"`
	Slug  string `db:"slug" json:"slug"`
}

func (q *Queries) UpdateTeamOwner(ctx context.Context, arg UpdateTeamOwnerParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeamOwner, arg.Owner, arg.Slug)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.Slug,
		&i.Organization,
	)
	return i, err
}
