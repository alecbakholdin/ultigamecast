package testdb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"ultigamecast/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	testDbPath = filepath.Join(basepath, "..", "test.db")
)
var testdb *sql.DB

func DB() (*models.Queries, *sql.DB) {
	if testdb == nil {
		Init()
	}
	return models.New(testdb), testdb
}

func Init() {
	ogTestDb, err := os.Open(testDbPath)
	if err != nil {
		panic(fmt.Errorf("error opening file: %w", err))
	}
	defer ogTestDb.Close()
	newTestDb, err := os.CreateTemp(os.TempDir(), "*.db")
	if err != nil {
		panic(fmt.Errorf("error creating temp file: %w", err))
	}

	_, err = io.Copy(newTestDb, ogTestDb)
	newTestDb.Close()
	if err != nil {
		panic(fmt.Errorf("error copying to new file: %w", err))
	}

	testdb, err = sql.Open("sqlite3", newTestDb.Name())
	if err != nil {
		panic(err)
	}
}