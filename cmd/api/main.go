package main

import (
	"encoding/json" // For formatting JSON response
	"fmt"
	"net/http"
	"os"
	"time"

	// Import our logger package
	"github.com/joshuaking1/aegis_finance_core/internal/platform/logger" // CHANGE 'yourusername'

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log" // Now using zerolog's global logger
)

func main() {
	// Initialize our structured logger
	// In a real app, "debug" might come from a config file or env var
	logger.Init("debug") // Or "info" for less verbosity

	port := os.Getenv("AEGIS_PORT")
	if port == "" {
		port = "8080"
		log.Info().Msg("AEGIS_PORT not set, defaulting to 8080")
	}

	r := chi.NewRouter()

	// === Middlewares ===
	// Replace chi's default logger with zerolog-based one for consistency
	// Chi's built-in middleware.Logger prints to standard log, which we are replacing.
	// So, we can create our own or use a community one like:
	// go get github.com/go-chi/httplog
	// For now, let's stick to simple logging within handlers to see zerolog in action
	// and remove middleware.Logger. We can add a dedicated zerolog HTTP logger middleware later.
	// r.Use(middleware.Logger) // Remove this for now to avoid double logging

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// === Define Routes ===
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		log.Info().Str("request_id", reqID).Str("path", r.URL.Path).Msg("Root handler invoked")
		fmt.Fprintf(w, "Welcome to Aegis Finance Core API! The future is autonomous. RequestID: %s", reqID)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		log.Info().Str("request_id", reqID).Str("path", r.URL.Path).Msg("Health check invoked")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"status": "UP", "service": "Aegis Finance Core API"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error().Err(err).Str("request_id", reqID).Msg("Failed to write health check response")
		}
	})

	r.Route("/users", func(subRouter chi.Router) {
		subRouter.Get("/{userID}", func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			userID := chi.URLParam(r, "userID")
			log.Info().Str("request_id", reqID).Str("user_id", userID).Msg("User details handler invoked")
			fmt.Fprintf(w, "Details for user %s. RequestID: %s", userID, reqID)
		})
	})

	log.Info().Str("port", port).Msg("Aegis Finance API server starting")

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		// Use zerolog for fatal errors too
		log.Fatal().Err(err).Msg("Could not start server")
	}
}