package task

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shresthashim/rest-api-golang/internal/storage"
	"github.com/shresthashim/rest-api-golang/internal/types"
	"github.com/shresthashim/rest-api-golang/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
}
