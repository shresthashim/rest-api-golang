package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shresthashim/rest-api-golang/internal/config"
)

type SQLiteStorage struct {
	Db *sql.DB
}

func NewSQLiteStorage(cfg *config.Config) (*SQLiteStorage, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	)`)

	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{Db: db}, nil

}
