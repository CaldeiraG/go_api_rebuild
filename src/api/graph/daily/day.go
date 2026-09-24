package daily

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
)

// Import shift constants from db package
const (
	Shift1 = config.Shift1
	Shift2 = config.Shift2
	Shift3 = config.Shift3
)

type dailyProductionResponse struct {
	LineID string `json:"line_id"`
	Date   string `json:"date"`
	Hora   int    `json:"hora"`
	Prod   int64  `json:"prod"`
	Shift  string `json:"shift"`
	Model  string `json:"model,omitempty"`
}

// DailyProduction handles the daily production endpoint
// GET /graph/api/daily/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func DailyProduction(w http.ResponseWriter, r *http.Request) {
	// Get URL parameters - use chi for path params
	lineID := chi.URLParam(r, "line_id")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	// Validate parameters
	if lineID == "" || startDateStr == "" || endDateStr == "" {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Missing required parameters: line_id, startDate, endDate"}`))
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"error": "Invalid startDate format. Use YYYY-MM-DD"}`)))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"error": "Invalid endDate format. Use YYYY-MM-DD"}`)))
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "endDate must be after or equal to startDate"}`))
		return
	}

	// Create query builder
	qb := config.NewProdQueryBuilder()

	// Build queries for all shifts (Shift 1, 2, 3)
	var results []dailyProductionResponse

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		// Build queries for all 3 shifts
		shifts := qb.GetAllShiftsForDate(date)

		for _, shift := range shifts {
			query, err := qb.GetShiftProduction(lineID, date, shift.ShiftType)
			if err != nil {
				results = append(results, dailyProductionResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(shift.ShiftType),
					Model:  "",
				})
				continue
			}

			// Execute query - get hora and prod for ALL hours
			rows, err := config.DB.Query(query)
			if err != nil {
				results = append(results, dailyProductionResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(shift.ShiftType),
					Model:  "",
				})
				continue
			}

			// Iterate through all rows
			for rows.Next() {
				var hora int
				var prod sql.NullInt64
				var model string
				err = rows.Scan(&hora, &prod, &model)
				if err != nil {
					continue
				}

				if prod.Valid {
					results = append(results, dailyProductionResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   prod.Int64,
						Shift:  string(shift.ShiftType),
						Model:  model,
					})
				} else {
					results = append(results, dailyProductionResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   0,
						Shift:  string(shift.ShiftType),
						Model:  "n/a",
					})
				}
			}

			rows.Close()
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")

	if len(results) == 0 {
		w.WriteHeader(200)
		w.Write([]byte(`{"error": "No records found"}`))
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"error": "Failed to encode response"}`)))
		return
	}
}
