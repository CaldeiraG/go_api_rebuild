package test

import (
	"encoding/json"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ScheduleStruct struct {
	CurrentWeek int
	Shift1      int
	Shift2      int
	Shift3      int
}

func Get(w http.ResponseWriter, r *http.Request) {
	schedule := ScheduleStruct{}
	lineID := chi.URLParam(r, "line_id")
	DBRes := config.DB.Raw("select *\nfrom \n(\n\tselect schedule,shift,datepart(ww,getdate()) as cur_week\n\tfrom TESTEProd.dbo.weekly_sched WS\n\tinner join TESTEProd.dbo.lines L on L.id = WS.line_id\n\tinner join TESTEProd.dbo.areas A on L.area_id = A.id\n\twhere line_id in (?) and timestamp = CAST(GETDATE() as DATE)\n) src\npivot\n(\n\tsum(schedule)\n\tfor shift in ([1], [2], [3])\n) piv", lineID).Scan(&schedule)

	if DBRes.Error != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(schedule)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}

}

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
