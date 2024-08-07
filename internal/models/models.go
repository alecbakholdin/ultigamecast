// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"database/sql"

	"ultigamecast/internal/models/tournament_data_types"
)

type Event struct {
	ID                string         `db:"id" json:"id"`
	Batch             sql.NullString `db:"batch" json:"batch"`
	Created           sql.NullTime   `db:"created" json:"created"`
	TeamScore         int64          `db:"team_score" json:"team_score"`
	OpponentScore     int64          `db:"opponent_score" json:"opponent_score"`
	Game              int64          `db:"game" json:"game"`
	Type              EventType      `db:"type" json:"type"`
	Player            sql.NullInt64  `db:"player" json:"player"`
	PreviousGameState GameLiveStatus `db:"previous_game_state" json:"previous_game_state"`
	PreviousEvent     sql.NullString `db:"previous_event" json:"previous_event"`
}

type Game struct {
	ID             int64              `db:"id" json:"id"`
	Tournament     int64              `db:"tournament" json:"tournament"`
	Slug           string             `db:"slug" json:"slug"`
	Opponent       string             `db:"opponent" json:"opponent"`
	Start          sql.NullTime       `db:"start" json:"start"`
	StartTimezone  sql.NullString     `db:"start_timezone" json:"start_timezone"`
	Wind           sql.NullInt64      `db:"wind" json:"wind"`
	Temp           sql.NullInt64      `db:"temp" json:"temp"`
	HalfCap        sql.NullInt64      `db:"half_cap" json:"half_cap"`
	SoftCap        sql.NullInt64      `db:"soft_cap" json:"soft_cap"`
	HardCap        sql.NullInt64      `db:"hard_cap" json:"hard_cap"`
	ScheduleStatus GameScheduleStatus `db:"schedule_status" json:"schedule_status"`
	LiveStatus     GameLiveStatus     `db:"live_status" json:"live_status"`
	ActivePlayers  sql.NullString     `db:"active_players" json:"active_players"`
	LastEvent      sql.NullString     `db:"last_event" json:"last_event"`
	TeamScore      int64              `db:"team_score" json:"team_score"`
	OpponentScore  int64              `db:"opponent_score" json:"opponent_score"`
}

type Player struct {
	ID    int64  `db:"id" json:"id"`
	Team  int64  `db:"team" json:"team"`
	Slug  string `db:"slug" json:"slug"`
	Name  string `db:"name" json:"name"`
	Order int64  `db:"order" json:"order"`
}

type Team struct {
	ID           int64          `db:"id" json:"id"`
	Owner        int64          `db:"owner" json:"owner"`
	Name         string         `db:"name" json:"name"`
	Slug         string         `db:"slug" json:"slug"`
	Organization sql.NullString `db:"organization" json:"organization"`
}

type TeamFollow struct {
	Team int64 `db:"team" json:"team"`
	User int64 `db:"user" json:"user"`
}

type TeamManager struct {
	Team int64 `db:"team" json:"team"`
	User int64 `db:"user" json:"user"`
}

type Tournament struct {
	ID        int64          `db:"id" json:"id"`
	Team      int64          `db:"team" json:"team"`
	Name      string         `db:"name" json:"name"`
	Slug      string         `db:"slug" json:"slug"`
	StartDate sql.NullTime   `db:"start_date" json:"start_date"`
	EndDate   sql.NullTime   `db:"end_date" json:"end_date"`
	Location  sql.NullString `db:"location" json:"location"`
}

type TournamentDatum struct {
	ID            int64                        `db:"id" json:"id"`
	Tournament    int64                        `db:"tournament" json:"tournament"`
	Icon          string                       `db:"icon" json:"icon"`
	Title         string                       `db:"title" json:"title"`
	ShowInPreview sql.NullInt64                `db:"show_in_preview" json:"show_in_preview"`
	TextPreview   string                       `db:"text_preview" json:"text_preview"`
	DataType      tournament_data_types.Option `db:"data_type" json:"data_type"`
	ValueText     sql.NullString               `db:"value_text" json:"value_text"`
	ValueLink     sql.NullString               `db:"value_link" json:"value_link"`
	Order         int64                        `db:"order" json:"order"`
}

type User struct {
	ID           int64          `db:"id" json:"id"`
	Email        string         `db:"email" json:"email"`
	PasswordHash sql.NullString `db:"password_hash" json:"password_hash"`
}
