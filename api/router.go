package api

import (
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// r.Post("/friends", CreateFriendConnection)
	// r.Get("/friends/{email}", GetFriendsList)
	// r.Post("/subscribe", SubscribeToUpdates)
	// r.Post("/block", BlockUpdates)
	// r.Get("/updates/{email}", GetUpdateRecipients)

	return r
}
