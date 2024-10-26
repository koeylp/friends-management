package responses

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	StatusForbidden    = http.StatusForbidden
	StatusNotFound     = http.StatusNotFound
	StatusBadRequest   = http.StatusBadRequest
	StatusUnauthorized = http.StatusUnauthorized
	StatusInternal     = http.StatusInternalServerError
)

var (
	ReasonBadRequest   = "Bad Request"
	ReasonNotFound     = "Not Found"
	ReasonForbidden    = "Access Denied"
	ReasonUnauthorized = "Unauthorized"
	ReasonInternal     = "Internal Server Error"
)

type ErrorResponse struct {
	Message string
	Status  int `json:"-"`
	Time    time.Time
}

func NewErrorResponse(message string, status int) *ErrorResponse {
	err := &ErrorResponse{
		Message: message,
		Status:  status,
		Time:    time.Now(),
	}
	log.Printf("ERROR: %d -- %s \n", err.Status, err.Message)
	return err
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

func (e *ErrorResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}

type BadRequestError struct {
	*ErrorResponse
}

func NewBadRequestError(message string) *BadRequestError {
	if message == "" {
		message = ReasonBadRequest
	}
	return &BadRequestError{NewErrorResponse(message, StatusBadRequest)}
}

type ForbiddenError struct {
	*ErrorResponse
}

func NewForbiddenError(message string) *ForbiddenError {
	if message == "" {
		message = ReasonForbidden
	}
	return &ForbiddenError{NewErrorResponse(message, StatusForbidden)}
}

type NotFoundError struct {
	*ErrorResponse
}

func NewNotFoundError(message string) *NotFoundError {
	if message == "" {
		message = ReasonNotFound
	}
	return &NotFoundError{NewErrorResponse(message, StatusNotFound)}
}

type UnauthorizedError struct {
	*ErrorResponse
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	if message == "" {
		message = ReasonUnauthorized
	}
	return &UnauthorizedError{NewErrorResponse(message, StatusUnauthorized)}
}

type InternalServerError struct {
	*ErrorResponse
}

func NewInternalServerError(message string) *InternalServerError {
	if message == "" {
		message = ReasonInternal
	}
	return &InternalServerError{NewErrorResponse(message, StatusInternal)}
}
