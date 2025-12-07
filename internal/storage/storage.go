package storage

type Storage interface {
	CreateTask(title, description string) (int, error)
}
