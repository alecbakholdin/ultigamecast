-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT
);
CREATE TABLE teams (
    id INTEGER PRIMARY KEY,
    "owner" INTEGER NOT NULL,
    "name" TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    organization TEXT
);
CREATE TABLE team_managers (
    team INTEGER NOT NULL,
    user INTEGER NOT NULL,
    FOREIGN KEY (team) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE
 );
CREATE TABLE team_follow (
    team INTEGER NOT NULL,
    user INTEGER NOT NULL,
    FOREIGN KEY (team) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    team INTEGER NOT NULL,
    slug TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "order" INTEGER NOT NULL,
    FOREIGN KEY (team) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(team, slug)
);
CREATE TABLE tournaments (
    id INTEGER PRIMARY KEY,
    team INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    slug TEXT NOT NULL GENERATED ALWAYS AS (LOWER(REPLACE("name", ' ', '-'))) STORED,
    "start_date" DATETIME,
    "end_date" DATETIME,
    "location" TEXT,
    FOREIGN KEY (team) REFERENCES teams(id) ON DELETE CASCADE
);
CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    tournament INTEGER NOT NULL,
    opponent TEXT NOT NULL,
    schedule_status TEXT CHECK(
        schedule_status IN ('scheduled', 'live', 'final')
    ),
    "start" DATETIME,
    start_timezone TEXT,
    wind INTEGER,
    temp INTEGER,
    half_cap INTEGER,
    soft_cap INTEGER,
    hard_cap INTEGER,
    FOREIGN KEY (tournament) REFERENCES tournaments(id) ON DELETE CASCADE
);
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    game INTEGER NOT NULL,
    "type" TEXT CHECK (
        "type" IN ('assist', 'goal', 'turn')
    ) NOT NULL,
    FOREIGN KEY (game) REFERENCES game(id) ON DELETE CASCADE
);