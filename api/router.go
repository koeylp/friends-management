package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/koeylp/friends-management/internal/handlers"
	"github.com/koeylp/friends-management/internal/repositories"
	"github.com/koeylp/friends-management/internal/services"
	"go.uber.org/fx"
)

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}

func RegisterRoutes(r *chi.Mux, userHandler *handlers.UserHandler, relationshipHandler *handlers.RelationshipHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUserHandler)
		})
		r.Route("/friends", func(r chi.Router) {
			r.Post("/", relationshipHandler.CreateFriendHandler)
		})
	})
}

var Module = fx.Options(
	fx.Provide(
		NewRouter,
		repositories.NewUserRepository,
		repositories.NewRelationshipRepository,
		services.NewUserService,
		services.NewRelationshipService,
		handlers.NewUserHandler,
		handlers.NewRelationshipHandler,
	),
	fx.Invoke(RegisterRoutes),
)

func StartServer(db *sql.DB) {
	app := fx.New(
		Module,
		fx.Supply(db),
		fx.Invoke(func(r *chi.Mux) {
			log.Fatal(http.ListenAndServe(":8080", r))
		}),
	)

	app.Run()
}
