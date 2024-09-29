package config

import (
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
}

func ConfigLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	storagePath := os.Getenv("storage_path")
	addres := os.Getenv("address")
	timeOut, _ := time.ParseDuration(os.Getenv("timeout"))
	idleTimeout, _ := time.ParseDuration(os.Getenv("idle_timeout"))

	cfg := Config{
		StoragePath: storagePath,
		Addr:        addres,
		Timeout:     timeOut,
		IdleTimeout: idleTimeout,
	}

	return &cfg
}
