package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DatabaseUrl string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"

	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/chat?sslmode=disable"
	}

	return &Config{
		Addr:        addr,
		DatabaseUrl: dbURL,
	}
}
