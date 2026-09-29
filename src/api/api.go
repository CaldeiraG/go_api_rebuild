package api

import (
	"github.com/caldeirag/go-api/src/api/com"
	"github.com/caldeirag/go-api/src/api/graph/daily"
	dailyNok "github.com/caldeirag/go-api/src/api/graph/daily/nok"
	"github.com/caldeirag/go-api/src/api/graph/hourly"
	hourlyNok "github.com/caldeirag/go-api/src/api/graph/hourly/nok"
	"github.com/caldeirag/go-api/src/api/graph/models"
	"github.com/caldeirag/go-api/src/api/graph/monthly"
	monthlyNok "github.com/caldeirag/go-api/src/api/graph/monthly/nok"
	"github.com/caldeirag/go-api/src/api/intranet"
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

func IntranetRouter() chi.Router {

	// New Chi SubRouter
	intranetRoute := chi.NewRouter()

	// Set up sub-routes
	intranetRoute.Get("/areas", intranet.Areas)
	intranetRoute.Get("/lines", intranet.Lines)
	intranetRoute.Get("/lines/{area_id}", intranet.Lines)
	intranetRoute.Get("/hourly-objectives", intranet.HourlyObjectives)
	intranetRoute.Get("/hourly-objectives/{line_id}", intranet.HourlyObjectives)

	// Return the Sub-Route back to the main API Router in main.go
	return intranetRoute
}

func GraphRouter() chi.Router {

	// New Chi SubRouter
	graphRoute := chi.NewRouter()

	// Set up sub-routes
	graphRoute.Get("/hourly/{line_id}", hourly.HourlyProduction)
	graphRoute.Get("/hourlynok/{line_id}", hourlyNok.HourlyNOK)
	graphRoute.Get("/daily/{line_id}", daily.DailyProduction)
	graphRoute.Get("/dailynok/{line_id}", dailyNok.DailyNOK)
	graphRoute.Get("/monthly/{line_id}", monthly.MonthlyProduction)
	graphRoute.Get("/monthlynok/{line_id}", monthlyNok.MonthlyNOK)
	graphRoute.Get("/models/{line_id}", models.Models)

	// Return the Sub-Route back to the main API Router in main.go
	return graphRoute
}
