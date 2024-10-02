package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	StoragePath string
	Addr        string
	Timeout     time.Duration
	IdleTimeout time.Duration
	LogLevel    string
}

func ConfigLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	addres := os.Getenv("address")
	timeOut, _ := time.ParseDuration(os.Getenv("timeout"))
	idleTimeout, _ := time.ParseDuration(os.Getenv("idle_timeout"))
	level := os.Getenv("log_level")

	storagePath := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("PG_CONTAINER"), os.Getenv("DB"))

	cfg := Config{
		StoragePath: storagePath,
		Addr:        addres,
		Timeout:     timeOut,
		IdleTimeout: idleTimeout,
		LogLevel:    level,
	}

	return &cfg
}
