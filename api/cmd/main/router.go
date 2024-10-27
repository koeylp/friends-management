package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	relationshipCtrl "github.com/koeylp/friends-management/cmd/internal/controller/relationship"
	userCtrl "github.com/koeylp/friends-management/cmd/internal/controller/user"
	handler "github.com/koeylp/friends-management/cmd/internal/handler/rest"
	relationshipRepo "github.com/koeylp/friends-management/cmd/internal/repository/relationship"
	userRepo "github.com/koeylp/friends-management/cmd/internal/repository/user"
	"go.uber.org/fx"
)

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}

func RegisterRoutes(r *chi.Mux, userHandler *handler.UserHandler, relationshipHandler *handler.RelationshipHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUserHandler)
		})
		r.Route("/friends", func(r chi.Router) {
			r.Post("/", relationshipHandler.CreateFriendHandler)
			r.Post("/list", relationshipHandler.GetFriendListByEmailHandler)
			r.Post("/common-list", relationshipHandler.GetCommonListHandler)
		})
		r.Route("/subcription", func(r chi.Router) {
			r.Post("/", relationshipHandler.SubscribeHandler)
			r.Post("/recipients", relationshipHandler.GetUpdatableEmailAddressesHandler)
		})
		r.Route("/block", func(r chi.Router) {
			r.Post("/", relationshipHandler.BlockUpdatesHandler)
		})
	})
}

var Module = fx.Options(
	fx.Provide(
		NewRouter,
		userRepo.NewUserRepository,
		relationshipRepo.NewRelationshipRepository,
		userCtrl.NewUserController,
		relationshipCtrl.NewRelationshipController,
		handler.NewUserHandler,
		handler.NewRelationshipHandler,
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
