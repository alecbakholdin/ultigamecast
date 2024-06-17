-- +goose Up
-- +goose StatementBegin
DROP TABLE events;
CREATE TABLE events (
    id VARCHAR(36) PRIMARY KEY,
    batch VARCHAR(36),
    "created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    team_score INTEGER NOT NULL,
    opponent_score INTEGER NOT NULL,
    game INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    player INTEGER,
    previous_game_state TEXT NOT NULL,
    previous_event TEXT,
    FOREIGN KEY(game) REFERENCES game(id) ON DELETE CASCADE,
    FOREIGN KEY(player) REFERENCES player(id) ON DELETE CASCADE
);
ALTER TABLE games DROP COLUMN schedule_status;
ALTER TABLE games ADD COLUMN schedule_status TEXT NOT NULL DEFAULT "Scheduled";
ALTER TABLE games ADD COLUMN live_status TEXT NOT NULL DEFAULT "Pre Point";
ALTER TABLE games ADD COLUMN active_players TEXT;
ALTER TABLE games ADD COLUMN last_event TEXT REFERENCES events(id);
ALTER TABLE games ADD COLUMN team_score INT NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN opponent_score INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    game INTEGER NOT NULL,
    "type" TEXT CHECK (
        "type" IN ('assist', 'goal', 'turn')
    ) NOT NULL,
    FOREIGN KEY (game) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE games DROP column schedule_status;
ALTER TABLE games ADD COLUMN schedule_status TEXT CHECK(schedule_status IN ('scheduled', 'live', 'final'));
ALTER TABLE games DROP COLUMN live_status;
ALTER TABLE games DROP COLUMN active_players;
ALTER TABLE games DROP last_event;
ALTER TABLE games DROP COLUMN team_score;
ALTER TABLE games DROP COLUMN opponent_score;
-- +goose StatementEnd