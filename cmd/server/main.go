package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/asim-sde/rag-search-service/internal/api"
	"github.com/asim-sde/rag-search-service/internal/db"
	"github.com/asim-sde/rag-search-service/internal/embedding"
	"github.com/asim-sde/rag-search-service/internal/rate"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/rag_db?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	database, err := db.NewDB(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer database.Close()

	limiter, err := rate.NewLimiter(redisURL)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer limiter.Close()

	embedService := embedding.NewService()
	handler := api.NewHandler(database, embedService, limiter)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/documents", handler.IngestDocument)
	r.Post("/api/search", handler.Search)

	log.Printf("RAG server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
