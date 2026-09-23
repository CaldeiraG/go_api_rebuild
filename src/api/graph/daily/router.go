package daily

import (
	"github.com/go-chi/chi/v5"
)

// Router sets up the daily production routes
func Router() chi.Router {
	// New Chi SubRouter
	dailyRoute := chi.NewRouter()

	// Set up sub-routes
	// GET /api/graph/daily/{line_id}?startDate=2026-09-21&endDate=2026-09-21
	dailyRoute.Get("/api/daily/{line_id}", DailyProduction)

	return dailyRoute
}
