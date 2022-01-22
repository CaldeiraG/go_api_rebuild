package schedule

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type ScheduleDayStruct struct {
	CurrentWeek int
	Shift1      float64
	Shift2      float64
	Shift3      float64
}

type ScheduleDayTestStruct struct {
	CurrentWeek   int
	DailySchedule float64
	Schedule      float64
}

/*type inTime struct {
	start string
	end   string
	check string
}*/

const sqlScheduleDay = `select piv.cur_week as 'CurrentWeek', ISNULL(piv.[1],0) as 'Shift1', ISNULL(piv.[2],0) as 'Shift2', ISNULL(piv.[3],0) as 'Shift3' from 
					(select schedule as 'schedule',shift as 'shift',datepart(ww,getdate()) as 'cur_week'
					from TESTEProd.dbo.weekly_sched WS

					where line_id = @line_id and timestamp = CAST(GETDATE() as DATE) 
					) as src 
					pivot 
					(
							sum(schedule) for shift in ([1], [2], [3])
					) as piv`

func ShiftDay(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	schedule := ScheduleDayStruct{}
	scheduleTest := ScheduleDayTestStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))

	query, err := config.DB.Prepare(sqlSchedule)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer query.Close()

	newRecord := query.QueryRowContext(ctx, sql.Named("line_id", lineID))

	err = newRecord.Scan(&schedule.CurrentWeek, &schedule.Shift1, &schedule.Shift2, &schedule.Shift3)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	var DailySchedule = schedule.Shift1 + schedule.Shift2 + schedule.Shift3

	var shiftNow = 0

	if lineID == 90 {

		fmt.Printf("%v", schedule.Shift1)
		fmt.Printf("%v", schedule.Shift2)
		fmt.Printf("%v", schedule.Shift3)
	}

	/*var ObjScheduleShift1 = schedule.Shift1 / 8.5
	var ObjScheduleShift2 = schedule.Shift2 / 8.5
	var ObjScheduleShift3 = schedule.Shift3 / 7*/

	currentTime := time.Now()

	if currentTime.Format("15:04") >= "00:00" && currentTime.Format("15:04") < "01:00" {
		oldTime := ConvertTimeCurrentDate("00:00")
		diff := currentTime.Sub(oldTime)
		shiftNow = int(diff.Minutes() * DailySchedule / 60)
	} else {
		oldTime := ConvertTimeCurrentDate("01:00")
		diff := currentTime.Sub(oldTime)
		/*if schedule.Shift3 > 0 && schedule.Shift2 > 0 && schedule.Shift1 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 1380)
		} else if schedule.Shift2 > 0 && schedule.Shift1 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 1020)
		} else if schedule.Shift3 > 0 && schedule.Shift1 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 930)
		} else if schedule.Shift1 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 510)
		} else if schedule.Shift2 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 510)
		} else if schedule.Shift3 > 0 {
			shiftNow = int(diff.Minutes() * DailySchedule / 420)
		}*/
		shiftNow = int(diff.Minutes() * DailySchedule / 1380)
	}

	scheduleTest.DailySchedule = float64(DailySchedule)
	scheduleTest.Schedule = float64(shiftNow)

	errJSON := json.NewEncoder(w).Encode(scheduleTest)
	if errJSON != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}
}

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
