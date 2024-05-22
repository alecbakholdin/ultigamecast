package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

var (
	ErrNotFound = errors.New("did not find a matching record")
	ErrUnexpected = errors.New("unexpected error")
)

func convertAndLogSqlError(ctx context.Context, message string, sqlErr error) error {
	if sqlErr == nil {
		return nil
	}

	slog.ErrorContext(ctx, message, "error", sqlErr)
	if errors.Is(sql.ErrNoRows, sqlErr) {
		return ErrNotFound
	} else {
		return ErrUnexpected
	}
}