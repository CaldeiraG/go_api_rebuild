package nok

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

type dailyNOKResponse struct {
	LineID string `json:"line_id"`
	Date   string `json:"date"`
	Hora   int    `json:"hora"`
	Prod   int64  `json:"prod"`
	Shift  string `json:"shift"`
	Model  string `json:"model,omitempty"`
	Error  string `json:"error,omitempty"` // Detailed error message when query fails
}

// errorResponse represents a detailed error message
type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// sendErrorResponse sends a detailed error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{
		Error:   message,
		Message: code,
		Code:    "ERROR",
	})
}

// DailyNOK handles the daily NOK production endpoint (using paramRej)
// GET /graph/api/dailynok/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func DailyNOK(w http.ResponseWriter, r *http.Request) {
	// Get URL parameters - use chi for path params
	lineID := chi.URLParam(r, "line_id")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	// Validate parameters
	if lineID == "" || startDateStr == "" || endDateStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Missing required parameters", "MISSING_PARAMETERS")
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid startDate format", fmt.Sprintf("INVALID_DATE: %v", err))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid endDate format", fmt.Sprintf("INVALID_DATE: %v", err))
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		sendErrorResponse(w, http.StatusBadRequest, "endDate must be after or equal to startDate", "INVALID_DATE_RANGE")
		return
	}

	// Create query builder
	qb := config.NewProdQueryBuilder()

	// Build queries for all shifts (Shift 1, 2, 3) using ParamRej for NOK
	var results []dailyNOKResponse

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		// Build queries for all 3 shifts
		shifts := qb.GetAllShiftsForDate(date, lineID)

		for _, shift := range shifts {
			query, err := qb.GetShiftProductionNOK(lineID, date, shift.ShiftType)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(shift.ShiftType),
					Model:  "",
					Error:  err.Error(), // Add error details to the response
				})
				continue
			}

			// Execute query - get hora and prod for ALL hours
			rows, err := config.DB.Query(query)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(shift.ShiftType),
					Model:  "",
					Error:  fmt.Sprintf("Database query failed: %v", err), // Add detailed error
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
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   0,
						Prod:   0,
						Shift:  string(shift.ShiftType),
						Model:  "",
						Error:  fmt.Sprintf("Row scan failed: %v", err), // Add detailed error
					})
					continue
				}

				if prod.Valid {
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   prod.Int64,
						Shift:  string(shift.ShiftType),
						Model:  model,
					})
				} else {
					results = append(results, dailyNOKResponse{
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
		sendErrorResponse(w, http.StatusNotFound, "No NOK data found", "NO_DATA_FOUND")
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to encode response", fmt.Sprintf("JSON_ERROR: %v", err))
		return
	}
}
