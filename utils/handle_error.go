package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/koeylp/friends-management/internal/responses"
)

func HandleError(w http.ResponseWriter, err error) {
	var notFoundErr *responses.NotFoundError
	var badRequestErr *responses.BadRequestError
	fmt.Println(&notFoundErr)
	switch {
	case errors.As(err, &notFoundErr):
		responses.NewNotFoundError(notFoundErr.Error()).Send(w)
	case errors.As(err, &badRequestErr):
		responses.NewBadRequestError(badRequestErr.Error()).Send(w)
	default:
		responses.NewInternalServerError(err.Error()).Send(w)
	}
}
