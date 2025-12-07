package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shresthashim/rest-api-golang/internal/config"
	"github.com/shresthashim/rest-api-golang/internal/types"
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

func (s *SQLiteStorage) CreateTask(title, description string) (int, error) {

	stat, err := s.Db.Prepare(`INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)`)

	if err != nil {
		return 0, err
	}

	defer stat.Close()

	result, err := stat.Exec(title, description, false)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteStorage) GetTasks() ([]types.Task, error) {
	rows, err := s.Db.Query(`SELECT id, title, description, completed FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []types.Task
	for rows.Next() {
		var task types.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *SQLiteStorage) GetTask(id int) (types.Task, error) {
	var task types.Task
	err := s.Db.QueryRow(`SELECT id, title, description, completed FROM tasks WHERE id = ?`, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Task{}, errors.New("task not found")
		}
		return types.Task{}, err
	}
	return task, nil
}
