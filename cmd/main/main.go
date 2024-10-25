package main

import (
	"context"
	"log"

	"github.com/koeylp/friends-management/internal/database"
)

func main() {
	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.CloseDB(context.Background())

	StartServer(dbConn)
}
