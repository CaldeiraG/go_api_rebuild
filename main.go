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
	// Declaring Environment variables
	// =================
	var DBHost, DBUser, DBPass, DBName, DBPort string
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, Using container variables: ERR: %v", err)
	}
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")

	// =================
	// Declaring Database Connection
	// =================
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;app name=HanonSystemsAPI", DBHost, DBUser, DBPass, DBPort, DBName)

	config.DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	config.DB.SetMaxOpenConns(25)
	config.DB.SetMaxIdleConns(25)
	config.DB.SetConnMaxLifetime(5 * time.Minute)

	ctx := context.Background()
	err = config.DB.PingContext(ctx)
	if err != nil {
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
