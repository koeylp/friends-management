package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	services "github.com/koeylp/friends-management/cmd/internal/controller"
	"github.com/koeylp/friends-management/cmd/internal/dto/user"
	"github.com/koeylp/friends-management/cmd/internal/responses"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler initializes a new UserHandler with the provided UserService
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserHandler handles the creation of a new user
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserReq user.CreateUser

	err := json.NewDecoder(r.Body).Decode(&createUserReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := user.ValidateCreateUserRequest(&createUserReq); err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}
	err = h.userService.CreateUser(context.Background(), &createUserReq)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	createdResponse := responses.NewCREATED(nil)
	createdResponse.Send(w)
}
