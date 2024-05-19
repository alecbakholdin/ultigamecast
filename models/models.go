// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"database/sql"
)

type Event struct {
	ID   int64  `db:"id" json:"id"`
	Game int64  `db:"game" json:"game"`
	Type string `db:"type" json:"type"`
}

type Game struct {
	ID             int64          `db:"id" json:"id"`
	Tournament     int64          `db:"tournament" json:"tournament"`
	Opponent       string         `db:"opponent" json:"opponent"`
	ScheduleStatus sql.NullString `db:"schedule_status" json:"schedule_status"`
	Start          sql.NullTime   `db:"start" json:"start"`
	StartTimezone  sql.NullString `db:"start_timezone" json:"start_timezone"`
	Wind           sql.NullInt64  `db:"wind" json:"wind"`
	Temp           sql.NullInt64  `db:"temp" json:"temp"`
	HalfCap        sql.NullInt64  `db:"half_cap" json:"half_cap"`
	SoftCap        sql.NullInt64  `db:"soft_cap" json:"soft_cap"`
	HardCap        sql.NullInt64  `db:"hard_cap" json:"hard_cap"`
}

type Player struct {
	ID    int64  `db:"id" json:"id"`
	Team  int64  `db:"team" json:"team"`
	Name  string `db:"name" json:"name" validate:"required,max=64"`
	Order int64  `db:"order" json:"order"`
}

type Team struct {
	ID           int64          `db:"id" json:"id"`
	Owner        int64          `db:"owner" json:"owner"`
	Name         string         `db:"name" json:"name" validate:"required,max=64"`
	Slug         string         `db:"slug" json:"slug"`
	Organization sql.NullString `db:"organization" json:"organization"`
}

type TeamManager struct {
	Team int64 `db:"team" json:"team"`
	User int64 `db:"user" json:"user"`
}

type Tournament struct {
	ID        int64          `db:"id" json:"id"`
	Team      int64          `db:"team" json:"team"`
	Name      string         `db:"name" json:"name" validate:"required,max=64"`
	Slug      string         `db:"slug" json:"slug"`
	StartDate sql.NullTime   `db:"start_date" json:"start_date"`
	EndDate   sql.NullTime   `db:"end_date" json:"end_date"`
	Location  sql.NullString `db:"location" json:"location"`
}

type User struct {
	ID           int64          `db:"id" json:"id"`
	Email        string         `db:"email" json:"email"`
	PasswordHash sql.NullString `db:"password_hash" json:"password_hash"`
}
