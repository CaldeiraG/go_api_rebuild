// Package common holds the shared implementation of the graph daily
// production/NOK endpoints. Both endpoints have identical request parsing,
// validation and row handling; only the query builder used differs.
package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
)

// DailyResponse is one hourly production/NOK data point.
type DailyResponse struct {
	LineID string `json:"line_id"`
	Date   string `json:"date"`
	Hora   int    `json:"hora"`
	Prod   int64  `json:"prod"`
	Shift  string `json:"shift"`
	Model  string `json:"model,omitempty"`
	Error  string `json:"error,omitempty"`
}

// errorResponse is the JSON body returned for non-200 responses.
type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// SendErrorResponse writes a JSON error body with the given HTTP status.
func SendErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message, Code: code})
}

// QueryFunc builds the hourly query for a line, date and shift.
type QueryFunc func(*db.ProdQueryBuilder, string, time.Time, db.ShiftType) (string, error)

// Daily handles a daily graph request:
// GET /graph/api/{daily|dailynok}/{line_id}?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD
func Daily(w http.ResponseWriter, r *http.Request, build QueryFunc) {
	lineID := chi.URLParam(r, "line_id")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if lineID == "" || startDateStr == "" || endDateStr == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Missing required parameters", "MISSING_PARAMETERS")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid startDate format", "INVALID_DATE")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid endDate format", "INVALID_DATE")
		return
	}

	if endDate.Before(startDate) {
		SendErrorResponse(w, http.StatusBadRequest, "endDate must be after or equal to startDate", "INVALID_DATE_RANGE")
		return
	}

	qb := db.NewProdQueryBuilder()

	exists, err := qb.LineExists(lineID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to load line configuration", "CONFIG_ERROR")
		return
	}
	if !exists {
		SendErrorResponse(w, http.StatusNotFound, "Line not found", "LINE_NOT_FOUND")
		return
	}

	var results []DailyResponse

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		for _, shift := range qb.GetAllShiftsForDate(date, lineID) {
			query, err := build(qb, lineID, date, shift.ShiftType)
			if err != nil {
				results = append(results, errorResult(lineID, date, shift.ShiftType, err))
				continue
			}

			rows, err := db.DB.QueryContext(r.Context(), query)
			if err != nil {
				results = append(results, errorResult(lineID, date, shift.ShiftType,
					fmt.Errorf("database query failed: %w", err)))
				continue
			}

			for rows.Next() {
				var hora int
				var prod sql.NullInt64
				var model sql.NullString
				if err := rows.Scan(&hora, &prod, &model); err != nil {
					results = append(results, errorResult(lineID, date, shift.ShiftType,
						fmt.Errorf("row scan failed: %w", err)))
					continue
				}

				row := DailyResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   hora,
					Shift:  string(shift.ShiftType),
				}
				if prod.Valid {
					row.Prod = prod.Int64
				}
				if model.Valid && model.String != "" {
					row.Model = model.String
				} else if !prod.Valid {
					row.Model = "n/a"
				}
				results = append(results, row)
			}

			if err := rows.Err(); err != nil {
				results = append(results, errorResult(lineID, date, shift.ShiftType,
					fmt.Errorf("row iteration failed: %w", err)))
			}
			rows.Close()
		}
	}

	if len(results) == 0 {
		SendErrorResponse(w, http.StatusNotFound, "No production data found", "NO_DATA_FOUND")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		fmt.Printf("Failed to encode response: %v\n", err)
	}
}

func errorResult(lineID string, date time.Time, shift db.ShiftType, err error) DailyResponse {
	return DailyResponse{
		LineID: lineID,
		Date:   date.Format("2006-01-02"),
		Shift:  string(shift),
		Error:  err.Error(),
	}
}
