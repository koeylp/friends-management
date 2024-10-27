package utils

import (
	"errors"
	"net/http"

	responses "github.com/koeylp/friends-management/cmd/internal/handler/rest/response"
)

// HandleError is a utility function that processes errors and sends the appropriate HTTP response.
// It checks the type of error and responds accordingly:
// - If the error is of type NotFoundError, it sends a 404 Not Found response.
// - If the error is of type BadRequestError, it sends a 400 Bad Request response.
// - For all other errors, it sends a 500 Internal Server Error response.
//
// Parameters:
// - w: http.ResponseWriter used to write the HTTP response.
// - err: the error that needs to be handled.
func HandleError(w http.ResponseWriter, err error) {
	var notFoundErr *responses.NotFoundError
	var badRequestErr *responses.BadRequestError
	switch {
	case errors.As(err, &notFoundErr):
		notFoundErr.Send(w)
	case errors.As(err, &badRequestErr):
		badRequestErr.Send(w)
	default:
		responses.NewInternalServerError(err.Error()).Send(w)
	}
}
