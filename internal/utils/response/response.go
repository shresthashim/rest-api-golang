package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type ValidationResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {

	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}

}

func ValidationError(errs validator.ValidationErrors) ValidationResponse {

	var errMessages []string

	for _, err := range errs {

		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, err.Field()+" is required")
		default:
			errMessages = append(errMessages, err.Field()+" is not valid")
		}
	}

	return ValidationResponse{
		Status:  StatusError,
		Message: "Validation failed",
		Errors:  errMessages,
	}
}
