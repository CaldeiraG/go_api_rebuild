package nok

import (
	"github.com/go-chi/chi/v5"
)

// Router sets up the daily NOK production routes
func Router() chi.Router {
	// New Chi SubRouter
	nokRoute := chi.NewRouter()

	// Set up sub-routes
	// GET /graph/api/dailynok/{line_id}?startDate=2026-09-21&endDate=2026-09-21
	nokRoute.Get("/dailynok/{line_id}", DailyNOK)

	return nokRoute
}
