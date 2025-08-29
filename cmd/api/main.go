package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/angelospanag/roe/internal/app"
	"github.com/angelospanag/roe/internal/db"
	"github.com/angelospanag/roe/internal/features/feeds"
	appMiddleware "github.com/angelospanag/roe/pkg/middleware"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Read environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	// Database
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databaseName := os.Getenv("DATABASE_NAME")

	ctx := context.Background()
	dbConn, err := pgx.Connect(
		ctx,
		fmt.Sprintf("user=%s password=%s host=%s dbname=%s", databaseUser,
			databasePassword,
			databaseHost,
			databaseName),
	)
	if err != nil {
		log.Fatalf("Error opening database connection %v", err.Error())
	}
	defer func(dbConn *pgx.Conn, ctx context.Context) {
		err := dbConn.Close(ctx)
		if err != nil {
			log.Fatalf("Error closing database connection %v", err.Error())
		}
	}(dbConn, ctx)

	queries := db.New(dbConn)

	// Create a new router & API
	addr := "0.0.0.0:8000"
	router := chi.NewMux()

	config := huma.DefaultConfig("Roe API", "1.0.0")

	// Get CORS allowed origins from environment variable
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	// If no origins are specified, default to allowing all
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"*"}
	}

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Healthcheck endpoint
	router.Use(middleware.Heartbeat("/health"))

	api := humachi.New(router, config)

	// Logging middleware
	api.UseMiddleware(appMiddleware.LoggingMiddleware)

	application := &app.App{
		DBQueries: queries,
	}

	// Register feature routes
	feeds.RegisterRoutes(api, application)

	// Start the server
	slog.Info("Starting server",
		slog.String("address", addr),
		slog.String("allowedOrigins", strings.Join(allowedOrigins, ",")),
	)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
