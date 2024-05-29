-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournament_data (
    id INTEGER NOT NULL PRIMARY KEY,
    tournament INTEGER NOT NULL,
    icon TEXT NOT NULL DEFAULT 'description',
    title TEXT NOT NULL DEFAULT 'No name',
    show_in_preview INTEGER DEFAULT 0,
    text_preview TEXT NOT NULL DEFAULT '',
    data_type INTEGER NOT NULL DEFAULT 0,
    value_text TEXT,
    value_link TEXT,
    "order" INTEGER NOT NULL,
    FOREIGN KEY(tournament) REFERENCES tournament(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tournament_data;
-- +goose StatementEnd
