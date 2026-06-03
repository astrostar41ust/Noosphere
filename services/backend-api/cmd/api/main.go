package main

import (
	"fmt"
	"net/http"
	"noosphere/backend-api/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "noosphere-backend"}`))
	})

	serverAddress := ":" + cfg.Port
	fmt.Printf("Noosphere Engine online. Listening on port %s...\n", cfg.Port)
	
	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Failed to start HTTP server: %v", err))
	}
}