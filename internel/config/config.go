package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DatabaseUrl string
	JwtSecret   string
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return &Config{
		Addr:        addr,
		DatabaseUrl: dbURL,
		JwtSecret:   jwtSecret,
	}
}
