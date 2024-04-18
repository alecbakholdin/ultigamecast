package repository

import "github.com/pocketbase/pocketbase/models"

func toArr[T any](records []*models.Record, fn func(r *models.Record) *T) []*T {
	arr := make([]*T, len(records))
	for i, r := range records {
		arr[i] = fn(r)
	}
	return arr
}