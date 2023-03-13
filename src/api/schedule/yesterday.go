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

func ShiftYesterday(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	schedule := DayStruct{}
	scheduleTest := DayTestStruct{}
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

	currentTime := time.Now()

	if currentTime.Format("15:04") >= "00:00" && currentTime.Format("15:04") < "01:00" {
		oldTime := ConvertTimeCurrentDate("00:00")
		diff := currentTime.Sub(oldTime)
		shiftNow = int(diff.Minutes() * DailySchedule / 60)
	} else {
		oldTime := ConvertTimeCurrentDate("01:00")
		diff := currentTime.Sub(oldTime)
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
