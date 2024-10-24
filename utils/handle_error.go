package utils

import (
	"errors"
	"net/http"

	"github.com/koeylp/friends-management/internal/responses"
)

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
