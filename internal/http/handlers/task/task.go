package task

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shresthashim/rest-api-golang/internal/storage"
	"github.com/shresthashim/rest-api-golang/internal/types"
	"github.com/shresthashim/rest-api-golang/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateTask(w, r, storage)
		case http.MethodGet:
			handleGetTasks(w, r, storage)
		default:
			response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(errors.New("method not allowed")))
		}
	}
}

func handleCreateTask(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	var task types.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if errors.Is(err, io.EOF) {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("request body cannot be empty")))
		return
	}

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	err = validator.New().Struct(task)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
		return
	}

	id, err := storage.CreateTask(task.Title, task.Description)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	task.ID = id

	response.WriteJSON(w, http.StatusCreated, task)
}

func NewWithID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(errors.New("method not allowed")))
			return
		}

		// Extract ID from URL path: /tasks/{id}
		path := r.URL.Path
		if !strings.HasPrefix(path, "/tasks/") {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid path")))
			return
		}

		idStr := strings.TrimPrefix(path, "/tasks/")
		if idStr == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("task ID is required")))
			return
		}

		// Convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid task ID")))
			return
		}

		handleGetTask(w, r, storage, id)
	}
}

func handleGetTasks(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	tasks, err := storage.GetTasks()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, tasks)
}

func handleGetTask(w http.ResponseWriter, r *http.Request, storage storage.Storage, id int) {
	task, err := storage.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			response.WriteJSON(w, http.StatusNotFound, response.GeneralError(errors.New("task not found")))
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, task)
}
