package api

import (
	"github.com/caldeirag/go-api/src/api/production"
	"github.com/caldeirag/go-api/src/api/schedule"
	"github.com/go-chi/chi/v5"
)

func ScheduleRouter() chi.Router {

	// New Chi SubRouter
	scheduleRoute := chi.NewRouter()

	// Set up sub-routes
	scheduleRoute.Get("/shift/{line_id}", schedule.Shift)
	scheduleRoute.Get("/now/{line_id}", schedule.Now)
	scheduleRoute.Get("/day/{line_id}", schedule.ShiftDay)
	scheduleRoute.Get("/yesterday/{line_id}", schedule.ShiftYesterday)

	// Return the Sub-Route back to the main API Router in main.go
	return scheduleRoute
}

func ProductionRouter() chi.Router {

	// New Chi SubRouter
	productionRoute := chi.NewRouter()

	// Set up sub-routes
	productionRoute.Get("/now/{line_id}", production.Production)
	productionRoute.Get("/day/{line_id}", production.DayProduction)
	productionRoute.Get("/yesterday/{line_id}", production.ProductionYesterday)

	// Return the Sub-Route back to the main API Router in main.go
	return productionRoute
}
