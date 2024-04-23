package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

func toArr[T any](records []*models.Record, fn func(r *models.Record) T) []T {
	arr := make([]T, len(records))
	for i, r := range records {
		arr[i] = fn(r)
	}
	return arr
}

func toValArr[T any](records []*models.Record, fn func(r *models.Record) *T) []T {
	arr := make([]T, len(records))
	for i, r := range records {
		arr[i] = *fn(r)
	}
	return arr
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, sql.ErrNoRows)
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "+"
	SortDirectionDesc SortDirection= "-"
)
func GetSortString(d SortDirection, field string) string {
	return fmt.Sprintf("%s%s", d, field)
}

