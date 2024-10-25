package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// start database config
type DBConfig struct {
	User     string
	Password string
	Host     string
	DBName   string
	SSLMode  string
	Port     string
}

func GetDBConfig() *DBConfig {
	_ = godotenv.Load()

	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Port:     os.Getenv("DB_PORT"),
	}
}

func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// end database config
