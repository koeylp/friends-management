package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/koeylp/friends-management/internal/handlers"
	"github.com/koeylp/friends-management/internal/repository"
	"github.com/koeylp/friends-management/internal/services"
)

func InitUserHandler(db *sql.DB) *handlers.UserHandler {
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	return handlers.NewUserHandler(userService)
}

func InitRelationshipHandler(db *sql.DB) *handlers.RelationshipHandler {
	relationshipRepo := repository.NewRelationshipRepository(db)
	userRepo := repository.NewUserRepository(db)
	relationshipService := services.NewRelationshipService(relationshipRepo, userRepo)
	return handlers.NewRelationshipHandler(relationshipService)
}

func SetupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	userHandler := InitUserHandler(db)
	relationshipHandler := InitRelationshipHandler(db)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUserHandler)
		})
		r.Route("/friends", func(r chi.Router) {
			r.Post("/", relationshipHandler.CreateFriendHandler)
		})
	})

	return r
}
