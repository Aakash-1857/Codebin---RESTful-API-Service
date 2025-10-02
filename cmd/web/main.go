// cmd/web/main.go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrickmn/go-cache"

	// Make sure you have these imports
	"github.com/aakash-1857/codebin/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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

func main() {
	var cfg config

	// --- CONFIGURATION ---
	// All settings are now read from command-line flags.
	// This makes the application flexible and easy to deploy.
    // Define the flags for local development.
    flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN") // We'll get this from a flag OR an env var.
    flag.StringVar(&cfg.jwt.secret, "jwt-secret", "your-super-secret-key", "JWT secret")
    flag.Parse()

    // --- NEW AND CORRECT LOGIC ---
    // Read the DATABASE_URL from the Render environment.
    envDSN := os.Getenv("DATABASE_URL")
    if envDSN != "" {
        // If it exists, use it and override any -db-dsn flag.
        // This makes the code work perfectly on Render.
        cfg.db.dsn = envDSN
    }


	// --- LOGGER ---
	// A structured logger for machine-readable logs.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SERVER STARTUP - VERSION 2")

	// --- DATABASE CONNECTION ---
	// We now connect using the DSN provided in the command-line flag.
	dbpool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	// Cache
	appCache := cache.New(5*time.Minute, 10*time.Minute)
	// --- APPLICATION INSTANCE ---
	// Initialize our application struct with all its dependencies.
	app := &application{
		config:   cfg,
		logger:   logger,
		users:    repository.NewUserRepository(dbpool),
		snippets: repository.NewSnippetRepository(dbpool),
		cache:    appCache,
	}

	// Create the HTTP server.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	// Create a channel to receive shutdown errors.
	shutdownError := make(chan error)

	// Start a background goroutine to listen for shutdown signals.
	go func() {
		// Create a channel to receive OS signals.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Block until a signal is received.
		s := <-quit
		app.logger.Info("shutting down server", "signal", s.String())

		// Create a context with a 5-second timeout for graceful shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call Shutdown() on our server, which will gracefully handle active connections.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("server is starting", "addr", srv.Addr, "env", app.config.env)

	// Start the server. This will block until a fatal error occurs.
	err = srv.ListenAndServe()
	// If the error is not http.ErrServerClosed, it's an unexpected fatal error.
	// We check against the error sent from the shutdown goroutine.
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Wait for the graceful shutdown to complete.
	err = <-shutdownError
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app.logger.Info("server stopped")

}
