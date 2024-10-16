package main

import (
	"context"
	"log"
	"net/http"

	"github.com/koeylp/friends-management/api"
	"github.com/koeylp/friends-management/database"
)

func main() {
	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close(context.Background())

	r := api.SetupRouter()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
