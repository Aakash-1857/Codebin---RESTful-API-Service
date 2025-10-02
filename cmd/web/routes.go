package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors" // Add this import
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Add the CORS middleware. This MUST come before your other middleware and routes.
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://codebin-restful-api-service.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	mux.Use(app.recoverPanic)

	// Public routes
	mux.Get("/healthcheck", app.healthcheck)
	mux.Get("/snippets/{id}", app.snippetView)
	mux.Get("/snippets", app.snippetLatest)
	mux.Post("/users", app.registerUserHandler)
	mux.Post("/tokens/authentication", app.createAuthenticationTokenHandler)

	// Protected routes
	mux.Group(func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Post("/snippets", app.snippetCreate)
	})

	// The file server should be last.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/*", http.StripPrefix("/", fileServer))

	return mux
}