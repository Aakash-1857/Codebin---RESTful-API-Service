// cmd/web/main.go
package main

import (

	"log/slog"

	"os"


	"github.com/patrickmn/go-cache"

	// Make sure you have these imports
	"github.com/aakash-1857/codebin/internal/repository"

)

// A single, consolidated config struct.
type config struct {
	port int
	env  string
	db   struct {
		dsn string // dsn stands for "data source name"
	}
	jwt struct {
		secret string //our signing key
	}
}

// The application struct now holds the config, along with other dependencies.
type application struct {
	config   config
	logger   *slog.Logger
	users    *repository.UserRepository
	snippets *repository.SnippetRepository
	cache    *cache.Cache
}

// In cmd/web/main.go

func main() {
	// Initialize the logger.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Read the environment variable from the Render environment.
	dbURL := os.Getenv("DATABASE_URL")

	// Print exactly what we found.
	logger.Info("--- FINAL DEBUG TEST ---", "DATABASE_URL_VALUE", dbURL)

	// If the variable is empty, print a loud, clear error.
	if dbURL == "" {
		logger.Error("!!! CRITICAL FAILURE: The DATABASE_URL environment variable is EMPTY or NOT SET in the Render environment. Please double-check the Environment tab and ensure the changes are saved. !!!")
	} else {
		logger.Info("+++ SUCCESS: The DATABASE_URL was found. +++")
	}
}
