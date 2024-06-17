package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

var (
	ErrNotFound        = errors.New("did not find a matching record")
	ErrUnexpected      = errors.New("unexpected error")
	ErrBadFormat       = errors.New("poorly formatted request")
	ErrDateOutOfBounds = errors.New("date is out of the proper bounds")
	ErrGameNotLive     = errors.New("game is not live")
	ErrLineNotReady    = errors.New("line is not ready for point start")
)

func convertAndLogSqlError(ctx context.Context, message string, sqlErr error) error {
	if sqlErr == nil {
		return nil
	}

	slog.ErrorContext(ctx, message, "SQL error", sqlErr)
	if errors.Is(sqlErr, sql.ErrNoRows) {
		return ErrNotFound
	} else {
		return errors.Join(ErrUnexpected, sqlErr)
	}
}
