package main

import (

	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func(app *Config) routes() http.Handler {
    mux := chi.NewRouter()

    mux.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/autheticate", app.Autheticate)

    // Add GET handler for /
    mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Auth service is running!"))
    })

    // Add POST handler for /
    mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("POST request received!"))
    })


    return mux
}
