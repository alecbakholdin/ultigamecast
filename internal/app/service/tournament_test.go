package service

import (
	"errors"
	"testing"
	"ultigamecast/test/testdb"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreateTournamentWithDuplicateNames(t *testing.T) {
	to := NewTournament(testdb.DB())
	ctx := testdb.LoadTeam(to.q)
	t1, err := to.CreateTournament(ctx, "Tournament  Name")
	if err != nil {
		t.Fatalf("error creating tourment: %s", err)
	}
	assert.Equal(t, "tournament--name", t1.Slug)

	t2, err := to.CreateTournament(ctx, "Tournament  Name")
	if err != nil {
		t.Fatalf("error creating tourment 2: %s", err)
	}
	assert.Equal(t, "tournament--name-2", t2.Slug)
}

func TestUpdateDataOrder(t *testing.T) {
	to := NewTournament(testdb.DB())
	ctx := testdb.LoadTeam(to.q)
	tournament, err := to.CreateTournament(ctx, "Tournament Name")
	if err != nil {
		t.Fatalf("error creating tournament: %s", err)
	}
	ctx = testdb.LoadCtxValue(ctx, tournament)
	var ids []int64
	for range 5 {
		datum, err := to.NewDatum(ctx)
		if err != nil {
			t.Fatalf("error creating datum: %s", err)
		}
		ids = append(ids, datum.ID)
	}

	dataPreOrder, err := to.Data(ctx)
	if err != nil {
		t.Fatalf("error fetching data pre-order: %s", err)
	}
	idsPreOrder := make([]int64, len(dataPreOrder))
	for i, d := range dataPreOrder {
		idsPreOrder[i] = d.ID
	}
	assert.Equal(t, ids, idsPreOrder)

	temp := ids[1]
	ids[1] = ids[2]
	ids[2] = temp
	if err := to.UpdateDataOrder(ctx, ids); err != nil {
		t.Fatalf("error updating order: %s", err)
	}

	dataPostOrder, err := to.Data(ctx)
	if err != nil {
		t.Fatalf("error fetching data post-order: %s", err)
	}
	idsPostOrder := make([]int64, len(dataPostOrder))
	for i, d := range dataPostOrder {
		idsPostOrder[i] = d.ID
	}
	assert.Equal(t, ids, idsPostOrder)
}

func TestUpdateDataOrderFailsWhenNotProperTournament(t *testing.T) {
	to := NewTournament(testdb.DB())
	withTeam := testdb.LoadTeam(to.q)
	t1, err := to.CreateTournament(withTeam, "random name")
	if err != nil {
		t.Fatalf("error creating tournament: %s", err)
	}
	withT1 := testdb.LoadCtxValue(withTeam, t1)
	datum1, err := to.NewDatum(withT1)
	if err != nil {
		t.Fatalf("error creating datumw for t1: %s", err)
	}

	t2, err := to.CreateTournament(withTeam, "random name 2")
	if err != nil {
		t.Fatalf("error creating tournament 2: %s", err)
	}
	withT2 := testdb.LoadCtxValue(withTeam, t2)
	_, err = to.NewDatum(withT2)
	if err != nil {
		t.Fatalf("error creating datumw for t1: %s", err)
	}

	if err = to.UpdateDataOrder(withT2, []int64{datum1.ID}); !errors.Is(ErrNotFound, err) {
		t.Fatalf("wrong or no error updating order: %s", err)
	}
}
