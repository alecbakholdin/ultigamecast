-- name: GetUser :one
SELECT *
FROM users
WHERE email = LOWER(@email);
-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (LOWER(@email), ?)
RETURNING *;
-- name: CreateTeam :one
INSERT INTO teams ("owner", "name", "slug", "organization")
VALUES (?, ?, ?, ?)
RETURNING *;
-- name: GetTeam :one
SELECT *
FROM teams
WHERE slug = LOWER(@slug);
-- name: ListOwnedTeams :many
SELECT *
FROM teams
WHERE teams.owner = @userId
ORDER BY id DESC;
-- name: ListFollowedTeams :many
SELECT t.*
FROM team_follow tf
    INNER JOIN teams t ON t.id = tf.team
WHERE tf.user = @userId
ORDER BY t.name ASC;
-- name: UpdateTeam :one
UPDATE teams
SET "name" = ?,
    "slug" = ?,
    organization = ?
WHERE teams.slug = @slug
RETURNING *;
-- name: UpdateTeamOwner :one
UPDATE teams
SET "owner" = ?
WHERE teams.slug = @slug
RETURNING *;
-- name: CreatePlayer :one
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
RETURNING *;
-- name: GetPlayer :one
SELECT *
FROM players
WHERE team = @teamId
    AND slug = @slug;
-- name: ListTeamPlayers :many
SELECT *
FROM players
WHERE team = @teamId
ORDER BY "order" ASC;
-- name: UpdatePlayer :one
UPDATE players
SET "name" = ?,
    slug = ?
WHERE id = ?
RETURNING *;
-- name: UpdatePlayerOrder :exec
UPDATE players
SET "order" = ?
WHERE id = ?
    AND team = ?;
-- name: GetTournament :one
SELECT *
FROM tournaments
WHERE team = @teamId
    AND slug = @slug;
-- name: ListTournaments :many
SELECT *
FROM tournaments
WHERE team = @teamId;
-- name: CreateTournament :one
INSERT INTO tournaments (team, "name", slug)
VALUES (@teamId, ?, ?)
RETURNING *;
-- name: UpdateTournamentDates :one
UPDATE tournaments
SET "start_date" = ?,
    "end_date" = ?
WHERE id = @tournamentId
RETURNING *;