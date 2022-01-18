package test

import (
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ScheduleStruct struct {
	CurrentWeek int
	Shift1      int
	Shift2      int
	Shift3      int
}

const sqlSchedule = `select piv.cur_week as 'CurrentWeek', piv.[1] as 'Shift1', piv.[2] as 'Shift2', piv.[3] as 'Shift3' from 
					(select schedule as 'schedule',shift as 'shift',datepart(ww,getdate()) as 'cur_week'
					from TESTEProd.dbo.weekly_sched WS

					where line_id = ? and timestamp = CAST(GETDATE() as DATE) 
					) as src 
					pivot 
					(
							sum(schedule) for shift in ([1], [2], [3])
					) as piv`

func Get(w http.ResponseWriter, r *http.Request) {
	schedule := ScheduleStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))
	DBRes := config.DB.Exec(sqlSchedule, lineID).Scan(&schedule)
	fmt.Printf("%v\n", schedule)
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

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
