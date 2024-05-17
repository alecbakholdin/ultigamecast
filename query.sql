-- name: GetTeam :one
SELECT *
FROM teams
WHERE slug = LOWER(@slug);
-- name: UpdateTeam :one
UPDATE teams
SET "name" = ?,
    organization = ?
WHERE teams.slug = @slug
RETURNING *;
-- name: CreatePlayer :one
INSERT INTO players (team, "name", "order")
VALUES (
        ?1,
        ?2,
        (
            SELECT 1 + IFNULL(MAX("order"), -1)
            FROM players p
            WHERE p.team = ?1
        )
    )
RETURNING *;
-- name: ListTeamPlayers :many
SELECT p.*
FROM players p
    INNER JOIN teams t ON p.team = t.id
WHERE t.slug = LOWER(@slug);
-- name: UpdatePlayer :one
UPDATE players
SET "name" = ?
WHERE id = ?
RETURNING *;
-- name: ListTournaments :many
SELECT tournaments.*
FROM tournaments
    INNER JOIN teams ON tournaments.team = teams.id
WHERE teams.slug = LOWER(@slug)