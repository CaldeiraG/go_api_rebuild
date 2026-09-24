package api

import (
	"github.com/caldeirag/go-api/src/api/com"
	"github.com/caldeirag/go-api/src/api/graph/daily"
	"github.com/caldeirag/go-api/src/api/graph/daily/nok"
	"github.com/caldeirag/go-api/src/api/production"
	"github.com/caldeirag/go-api/src/api/schedule"
	"github.com/caldeirag/go-api/src/api/scrap"
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

func ScrapRouter() chi.Router {

	// New Chi SubRouter
	scrapRoute := chi.NewRouter()

	// Set up sub-routes
	//scrapRoute.Get("/insertticket/?ticket={ticket}&date={date}&costcenter={costcenter}&price={price}", scrap.InsertTicketList)
	scrapRoute.Get("/insertTicket", scrap.InsertTicketList)
	scrapRoute.Get("/updatePerson", scrap.UpdateTicketPerson)
	scrapRoute.Get("/updateRequester", scrap.UpdateTicketRequester)
	//scrapRoute.Get("/yesterday/{line_id}", production.ProductionYesterday)

	// Return the Sub-Route back to the main API Router in main.go
	return scrapRoute
}

func HeartbeatRouter() chi.Router {

	// New Chi SubRouter
	heartbeatRoute := chi.NewRouter()

	// Set up sub-routes
	//scrapRoute.Get("/insertticket/?ticket={ticket}&date={date}&costcenter={costcenter}&price={price}", scrap.InsertTicketList)
	heartbeatRoute.Get("/heartbeat", com.HeartbeatInsert)
	//scrapRoute.Get("/yesterday/{line_id}", production.ProductionYesterday)

	// Return the Sub-Route back to the main API Router in main.go
	return heartbeatRoute
}

func GraphRouter() chi.Router {

	// New Chi SubRouter
	graphRoute := chi.NewRouter()

	// Set up sub-routes
	graphRoute.Get("/daily/{line_id}", daily.DailyProduction)
	graphRoute.Get("/dailynok/{line_id}", nok.DailyNOK)

	// Return the Sub-Route back to the main API Router in main.go
	return graphRoute
}
