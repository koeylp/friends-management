package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/koeylp/friends-management/config"
	"github.com/volatiletech/sqlboiler/boil"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	dbConfig := config.GetDBConfig()
	connStr := dbConfig.GetConnectionString()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}
	DB = db
	boil.SetDB(DB)

	return DB, nil
}

func CloseDB(ctx context.Context) {
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close the database: %v", err)
	}
}
