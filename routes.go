package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// This basically Adds routes to the existing routes.
func (a *application) routes() *chi.Mux {
	// Add routes here
	a.get("/", a.Handlers.Home)

	// Static Routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
