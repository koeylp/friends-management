package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/koeylp/friends-management/dto/user"
	"github.com/koeylp/friends-management/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserReq user.CreateUser

	err := json.NewDecoder(r.Body).Decode(&createUserReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.userService.CreateUser(context.Background(), &createUserReq)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
