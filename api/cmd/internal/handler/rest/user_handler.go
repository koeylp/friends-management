package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	userCtrl "github.com/koeylp/friends-management/cmd/internal/controller/user"
	"github.com/koeylp/friends-management/cmd/internal/handler/rest/response"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	userController userCtrl.UserController
}

// NewUserHandler initializes a new UserHandler with the provided UserController
func NewUserHandler(userController userCtrl.UserController) *UserHandler {
	return &UserHandler{userController: userController}
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
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}
	err = h.userController.CreateUser(context.Background(), &createUserReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	createdResponse := response.NewCREATED(nil)
	createdResponse.Send(w)
}
