package api

import (
	"github.com/caldeirag/go-api/src/api/test"
	"github.com/go-chi/chi/v5"
)

func TestRouter() chi.Router {

	// New Chi SubRouter
	testRoute := chi.NewRouter()

	// Set up sub-routes
	testRoute.Get("/get", test.Get)
	testRoute.Post("/post", test.Post)
	testRoute.Delete("/delete", test.Delete)

	// Return the Sub-Route back to the main API Router in main.go
	return testRoute
}
