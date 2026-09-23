package daily

import (
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"net/http"
	"time"
)

type dailyProductionResponse struct {
	LineID   string  `json:"line_id"`
	Date     string  `json:"date"`
	Hora     int     `json:"hora"`
	Prod     int64   `json:"prod"`
	Error    string  `json:"error,omitempty"`
}

// DailyProduction handles the daily production endpoint
// GET /api/graph/daily/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func DailyProduction(w http.ResponseWriter, r *http.Request) {
	// Get URL parameters
	lineID := getUrlParam(r, "line_id")
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
					Error:  fmt.Sprintf("Failed to build query for shift %s: %v", shift.ShiftType, err),
				})
				continue
			}

			// Execute query
			var hora int
			var prod sql.NullInt64
			err = config.DB.QueryRow(query).Scan(&hora, &prod)
			if err != nil {
				results = append(results, dailyProductionResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Error:  fmt.Sprintf("Failed to execute query: %v", err),
				})
				continue
			}

			if prod.Valid {
				results = append(results, dailyProductionResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   hora,
					Prod:   prod.Int64,
				})
			} else {
				results = append(results, dailyProductionResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   hora,
					Prod:   0,
				})
			}
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

// getUrlParam is a helper to get URL parameters
func getUrlParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
