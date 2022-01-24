package schedule

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type ScheduleStruct struct {
	CurrentWeek int
	Shift1      float64
	Shift2      float64
	Shift3      float64
}

type ScheduleStructTest struct {
	CurrentWeek int
	Shift       float64
}

type ScheduleNowStruct struct {
	CurrentWeek int
	Shift       float64
}

/*type inTime struct {
	start string
	end   string
	check string
}*/

const sqlSchedule = `select piv.cur_week as 'CurrentWeek', ISNULL(piv.[1],0) as 'Shift1', ISNULL(piv.[2],0) as 'Shift2', ISNULL(piv.[3],0) as 'Shift3' from 
					(select schedule as 'schedule',shift as 'shift',datepart(ww,getdate()) as 'cur_week'
					from TESTEProd.dbo.weekly_sched WS

					where line_id = @line_id and timestamp = CAST(GETDATE() as DATE) 
					) as src 
					pivot 
					(
							sum(schedule) for shift in ([1], [2], [3])
					) as piv`

func Shift(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	schedule := ScheduleStruct{}
	scheduleTest := ScheduleStructTest{}
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

	var shiftNow = 0

	currentTime := time.Now()

	if currentTime.Format("15:04") >= "08:00" && currentTime.Format("15:04") < "16:30" {
		shiftNow = int(schedule.Shift1)
	} else if currentTime.Format("15:04") >= "16:30" && currentTime.AddDate(0, 0, -1).Format("15:04") < "01:00" {
		shiftNow = int(schedule.Shift2)
	} else if currentTime.Format("15:04") >= "01:00" && currentTime.Format("15:04") < "08:00" {
		shiftNow = int(schedule.Shift3)
	}

	scheduleTest.Shift = float64(shiftNow)

	errJSON := json.NewEncoder(w).Encode(scheduleTest)
	if errJSON != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}
}

func Now(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	schedule := ScheduleStruct{}
	scheduleNow := ScheduleNowStruct{}
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

	var shift1 = "08:00"
	var shift2 = "16:30"
	var shift3 = "01:00"

	var shiftNow = 0

	currentTime := time.Now()

	if currentTime.Format("15:04") >= "08:00" && currentTime.Format("15:04") < "16:30" {
		oldTime := ConvertTimeCurrentDate(shift1)
		diff := currentTime.Sub(oldTime)
		shiftNow = int(diff.Minutes() * schedule.Shift1 / 510)
	} else if currentTime.Format("15:04") >= "16:30" && currentTime.AddDate(0, 0, 1).Format("15:04") < "01:00" {
		oldTime := ConvertTimeCurrentDate(shift2)
		diff := currentTime.Sub(oldTime)
		shiftNow = int(diff.Minutes() * schedule.Shift2 / 510)
	} else if currentTime.Format("15:04") >= "01:00" && currentTime.Format("15:04") < "08:00" {
		oldTime := ConvertTimeCurrentDate(shift3)
		diff := currentTime.Sub(oldTime)
		shiftNow = int(diff.Minutes() * schedule.Shift3 / 420)
	}

	scheduleNow.Shift = float64(shiftNow)

	errJSON := json.NewEncoder(w).Encode(scheduleNow)
	if errJSON != nil {
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

func ConvertTimeCurrentDate(s string) time.Time {
	oldTime, _ := time.Parse("15:04", s)
	hour, min, sec := oldTime.Clock()
	now := time.Now().UTC()
	year, month, day := now.Date()
	oldTime = time.Date(year, month, day, hour, min, sec, 0, time.UTC)
	return oldTime
}

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
