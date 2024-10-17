package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/koeylp/friends-management/api/handlers"
	"github.com/koeylp/friends-management/repository"
	"github.com/koeylp/friends-management/services"
)

func SetupRouter(db *sql.DB) *chi.Mux {
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUserHandler)

	return r
}
