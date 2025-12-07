package storage

import "github.com/shresthashim/rest-api-golang/internal/types"

type Storage interface {
	CreateTask(title, description string) (int, error)
	GetTasks() ([]types.Task, error)
	GetTask(id int) (types.Task, error)
}
