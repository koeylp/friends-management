package database

import (
	"context"
	"log"

	"github.com/koeylp/friends-management/config"

	"github.com/jackc/pgx/v5"
)

func InitDB() (*pgx.Conn, error) {
	dbURL := config.GetDBConfig()
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database!")
	return conn, nil
}
