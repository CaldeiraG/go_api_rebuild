package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/caldeirag/go-api/src/api"
	config "github.com/caldeirag/go-api/src/db"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mvrilo/go-redoc"
	"log"
	"net/http"
	"os"
)

func main() {
	// Redoc Environment
	// ================
	doc := redoc.Redoc{
		DocsPath: "/docs",
		// Change SpecPath && SpecFile to ./static/swagger.json when developing locally!
		SpecPath:    "./static/swagger.json",
		SpecFile:    "./static/swagger.json",
		Title:       "OnlyTunes API Template",
		Description: "API Documentation for OnlyTunes API Template",
	}
	// =================
	// Declaring Environment variables
	// =================
	var DBHost, DBUser, DBPass, DBName, DBPort string
	var DBHost2, DBUser2, DBPass2, DBName2, DBPort2 string
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, Using container variables: ERR: %v", err)
	}
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")

	DBHost2 = os.Getenv("DB_HOST2")
	DBUser2 = os.Getenv("DB_USER2")
	DBPass2 = os.Getenv("DB_PASS2")
	DBName2 = os.Getenv("DB_NAME2")
	DBPort2 = os.Getenv("DB_PORT2")
	// =================
	// Declaring Database Connection
	// =================
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable", DBHost, DBUser, DBPass, DBPort, DBName)

	config.DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = config.DB.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	connString2 := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable", DBHost2, DBUser2, DBPass2, DBPort2, DBName2)

	config.DB2, err = sql.Open("sqlserver", connString2)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx2 := context.Background()
	err = config.DB2.PingContext(ctx2)
	if err != nil {
		log.Fatal(err.Error())
	}

	// =================
	// Initialize Router and WebServer
	// =================
	r := chi.NewRouter()
	// Uncomment these during development / debugging
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

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
	// =================
	// Initialize API Documentation
	// =================
	// Change the http.Dir to ./static for local development!
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("/srv/static"))))
	r.Handle("/docs", doc.Handler())
	// =================
	// Start WebServer
	// =================
	fmt.Printf("Starting Server on port 4000\n")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		log.Fatal(err)
	}
}
