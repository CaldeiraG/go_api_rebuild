package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caldeirag/go-api/src/api"
	"github.com/caldeirag/go-api/src/api/gen5"
	config "github.com/caldeirag/go-api/src/db"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mvrilo/go-redoc"
)

func main() {
	// Redoc Environment
	// ================
	doc := redoc.Redoc{
		DocsPath: "/docs",
		// SpecPath is the URL the browser requests (served by the /static file server).
		// SpecFile is the local file Redoc reads at startup.
		SpecPath:    "/static/swagger.json",
		SpecFile:    "./static/swagger.json",
		Title:       "Hanon Systems API",
		Description: "API Documentation for Hanon Systems",
	}
	// =================
	// Declaring Database Connection
	// =================
	db, err := connectDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	config.DB = db

	ctx := context.Background()
	if err := config.DB.PingContext(ctx); err != nil {
		log.Fatal(err.Error())
	}

	gen5.ProductionGEN5_cron()

	// =================
	// Start WebServer
	// =================
	r := buildRouter(doc)

	currentTime := time.Now()
	log.Printf("Starting Server on port 4000")
	log.Printf("Time: %s", currentTime.Format("2006-01-02 15:04:05"))

	srv := &http.Server{
		Addr:              ":4000",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// connectDB loads the environment and opens the SQL Server connection pool.
// It is shared by main and the live health check test.
func connectDB() (*sql.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file, Using container variables: ERR: %v", err)
	}

	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;app name=HanonSystemsAPI",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// buildRouter wires up middleware, API routes and documentation. It is kept
// separate from main so it can be exercised by tests without a database.
func buildRouter(doc redoc.Redoc) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	// Request logging is opt-in; set LOG_REQUESTS=true to enable it.
	if os.Getenv("LOG_REQUESTS") == "true" {
		r.Use(middleware.Logger)
	}

	// =================
	// Set Custom status messages for 404 && 405
	// =================
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("404 - Page not found - Check the API Docs for more info"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("405 - Method not allowed - Check the API Docs for more info"))
	})
	// =================
	// Initialize API Middleware
	// =================
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.NoCache)
	// =================
	// Initialize API Routes
	// =================
	r.Mount("/schedule", api.ScheduleRouter())
	r.Mount("/production", api.ProductionRouter())
	r.Mount("/scrap", api.ScrapRouter())
	r.Mount("/com", api.HeartbeatRouter())
	r.Mount("/graph/api", api.GraphRouter())
	// =================
	// Initialize API Documentation
	// =================
	// Change the http.Dir to ./static for local development!
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.Handle("/docs", doc.Handler())

	return r
}
