package main

import (
	"context"
	"log"

	"github.com/koeylp/friends-management/cmd/internal/infra/database/postgres"
)

func main() {
	dbConn, err := postgres.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer postgres.CloseDB(context.Background())

	StartServer(dbConn)
}
