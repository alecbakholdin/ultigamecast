dev:
	air -c config/.air.toml

create:
	goose -dir ./config/migrations create new_migration sql

up:
	goose -dir ./config/migrations sqlite3 data.db up

down:
	goose -dir ./config/migrations sqlite3 data.db down