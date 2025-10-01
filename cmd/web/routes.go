package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler{
	/*mux:=http.NewServeMux()
	mux.HandleFunc("/",app.home)
	mux.HandleFunc("/snippet/view",app.snippetView)
	mux.HandleFunc("/snippet/create",app.snippetCreate)*/

	mux:=chi.NewRouter()

	mux.Use(app.recoverPanic)

	mux.Get("/healthcheck",app.healthcheck)

	// Add the new snippet routes
	mux.Post("/snippets",app.snippetCreate)
	mux.Get("/snippets/{id}",app.snippetView)
	mux.Get("/snippets",app.snippetLatest)
	mux.Post("/users",app.registerUserHandler)
	mux.Post("/tokens/authentication",app.createAuthenticationTokenHandler)

	// Protected routes
	mux.Group(func(r chi.Router){
		r.Use(app.requireAuthentication)
		r.Post("/snippets",app.snippetCreate)
	})
	return mux
}
