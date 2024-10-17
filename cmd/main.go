package main

import (
	"context"
	"log"
	"net/http"

	"github.com/koeylp/friends-management/api"
	"github.com/koeylp/friends-management/database"
)

func main() {
	var dbConn, err = database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.CloseDB(context.Background())

	r := api.SetupRouter(dbConn)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
