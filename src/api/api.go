package api

import (
	"github.com/caldeirag/go-api/src/api/schedule"
	"github.com/go-chi/chi/v5"
)

func ScheduleRouter() chi.Router {

	// New Chi SubRouter
	scheduleRoute := chi.NewRouter()

	// Set up sub-routes
	scheduleRoute.Get("/shift/{line_id}", schedule.Shift)
	scheduleRoute.Get("/now/{line_id}", schedule.Now)
	scheduleRoute.Post("/post", schedule.Post)
	scheduleRoute.Delete("/delete", schedule.Delete)

	// Return the Sub-Route back to the main API Router in main.go
	return scheduleRoute
}
