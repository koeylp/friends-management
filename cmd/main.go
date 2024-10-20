package main

import (
	"context"
	"log"

	"github.com/koeylp/friends-management/api"
	"github.com/koeylp/friends-management/database"
)

func main() {
	var dbConn, err = database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.CloseDB(context.Background())

	api.StartServer(dbConn)
}
