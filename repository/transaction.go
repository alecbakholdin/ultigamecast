package repository

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
)

type Transaction struct {
	dao *daos.Dao
}

func NewTransaction(app core.App) *Transaction {
	return &Transaction{
		dao: app.Dao(),
	}
}

func (t *Transaction) Run(fn func(txDao *daos.Dao) error) error {
	return t.dao.RunInTransaction(fn)
}
