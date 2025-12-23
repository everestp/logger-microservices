package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)


func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// CORS
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all for testing
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	mux.Use(middleware.Heartbeat("/ping"))

	// Routes
	mux.Post("/", app.SendMail)
	
	return mux
}