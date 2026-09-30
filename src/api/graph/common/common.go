// Package common holds the shared implementation of the graph hourly
// production/NOK endpoints. Both endpoints have identical request parsing,
// validation and row handling; only the query builder used differs.
package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
)

// HourlyResponse is one hourly production/NOK data point.
type HourlyResponse struct {
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

// QueryFunc builds the hourly query for a line, date and shift, returning the
// SQL and its bound arguments.
type QueryFunc func(*db.ProdQueryBuilder, string, time.Time, db.ShiftType) (string, []interface{}, error)

// Hourly handles an hourly graph request:
// GET /graph/api/{hourly|hourlynok}/{line_id}?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD
func Hourly(w http.ResponseWriter, r *http.Request, build QueryFunc) {
	lineID, startDate, endDate, models, ok := parseLineRange(w, r)
	if !ok {
		return
	}

	qb := db.NewProdQueryBuilder().WithModels(models)

	line, exists, err := qb.LineConfig(lineID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to load line configuration", "CONFIG_ERROR")
		return
	}
	if !exists {
		SendErrorResponse(w, http.StatusNotFound, "Line not found", "LINE_NOT_FOUND")
		return
	}
	isGen5 := line.QueryType == "gen5"

	var results []HourlyResponse

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		for _, shift := range qb.GetAllShiftsForDate(date, lineID) {
			query, args, err := build(qb, lineID, date, shift.ShiftType)
			if err != nil {
				results = append(results, errorResult(lineID, date, shift.ShiftType, err))
				continue
			}

			rows, err := db.DB.QueryContext(r.Context(), query, args...)
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

				value := int64(0)
				if prod.Valid {
					value = prod.Int64
				}

				rowShifts := []db.ShiftType{shift.ShiftType}
				if isGen5 {
					rowShifts = gen5RowShifts(hora)
				}

				for i, rowShift := range rowShifts {
					row := HourlyResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Shift:  string(rowShift),
						Prod:   splitValue(value, len(rowShifts), i),
					}
					if model.Valid && model.String != "" {
						row.Model = model.String
					} else if !prod.Valid {
						row.Model = "n/a"
					}
					results = append(results, row)
				}
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

// gen5RowShifts returns the shift(s) an hourly GEN5 row belongs to. GEN5
// returns the whole day from a single query, so the shift is derived from the
// hour instead of the queried window. Hour 16 (16:00-17:00) straddles the
// 16:30 shift boundary, so it is split between shift 1 and shift 2.
func gen5RowShifts(hora int) []db.ShiftType {
	if hora == 16 {
		return []db.ShiftType{db.Shift1, db.Shift2}
	}
	return []db.ShiftType{db.ShiftForHour(hora)}
}

// splitValue distributes an hourly value across the shifts it belongs to.
// The remainder is given to the last shift so the total is preserved.
func splitValue(value int64, parts, index int) int64 {
	if parts <= 1 {
		return value
	}
	base := value / int64(parts)
	if index == parts-1 {
		return value - base*int64(parts-1)
	}
	return base
}

func errorResult(lineID string, date time.Time, shift db.ShiftType, err error) HourlyResponse {
	return HourlyResponse{
		LineID: lineID,
		Date:   date.Format("2006-01-02"),
		Shift:  string(shift),
		Error:  err.Error(),
	}
}

// DailyResponse is one per-day, per-model total.
type DailyResponse struct {
	LineID string `json:"line_id"`
	Date   string `json:"date"`
	Prod   int64  `json:"prod"`
	Model  string `json:"model,omitempty"`
	Error  string `json:"error,omitempty"`
}

// MonthlyResponse is one per-month, per-model total.
type MonthlyResponse struct {
	LineID string `json:"line_id"`
	Month  string `json:"month"`
	Prod   int64  `json:"prod"`
	Model  string `json:"model,omitempty"`
	Error  string `json:"error,omitempty"`
}

// PeriodQueryFunc builds a whole-period (day or month) query for a line over
// an inclusive range.
type PeriodQueryFunc func(*db.ProdQueryBuilder, string, time.Time, time.Time) (string, []interface{}, error)

type periodRow struct {
	period string
	model  string
	prod   int64
}

// Daily handles a daily totals graph request:
// GET /graph/api/{daily|dailynok}/{line_id}?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD
func Daily(w http.ResponseWriter, r *http.Request, build PeriodQueryFunc) {
	lineID, rows, ok := fetchPeriodTotals(w, r, build)
	if !ok {
		return
	}

	results := make([]DailyResponse, 0, len(rows))
	for _, row := range rows {
		results = append(results, DailyResponse{
			LineID: lineID,
			Date:   row.period,
			Prod:   row.prod,
			Model:  row.model,
		})
	}
	writeJSON(w, results)
}

// Monthly handles a monthly totals graph request:
// GET /graph/api/{monthly|monthlynok}/{line_id}?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD
func Monthly(w http.ResponseWriter, r *http.Request, build PeriodQueryFunc) {
	lineID, rows, ok := fetchPeriodTotals(w, r, build)
	if !ok {
		return
	}

	results := make([]MonthlyResponse, 0, len(rows))
	for _, row := range rows {
		results = append(results, MonthlyResponse{
			LineID: lineID,
			Month:  row.period,
			Prod:   row.prod,
			Model:  row.model,
		})
	}
	writeJSON(w, results)
}

// parseLineRange reads and validates the {line_id}, startDate and endDate
// shared by the graph endpoints, plus the optional model filter, writing an
// error response on failure.
func parseLineRange(w http.ResponseWriter, r *http.Request) (string, time.Time, time.Time, []string, bool) {
	lineID := chi.URLParam(r, "line_id")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if lineID == "" || startDateStr == "" || endDateStr == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Missing required parameters", "MISSING_PARAMETERS")
		return "", time.Time{}, time.Time{}, nil, false
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid startDate format", "INVALID_DATE")
		return "", time.Time{}, time.Time{}, nil, false
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid endDate format", "INVALID_DATE")
		return "", time.Time{}, time.Time{}, nil, false
	}

	if endDate.Before(startDate) {
		SendErrorResponse(w, http.StatusBadRequest, "endDate must be after or equal to startDate", "INVALID_DATE_RANGE")
		return "", time.Time{}, time.Time{}, nil, false
	}

	return lineID, startDate, endDate, parseModels(r), true
}

// parseModels reads the optional model filter. It accepts repeated ?model=
// parameters and/or comma-separated values, trimming and dropping empties.
func parseModels(r *http.Request) []string {
	var models []string
	for _, value := range r.URL.Query()["model"] {
		for _, part := range strings.Split(value, ",") {
			if model := strings.TrimSpace(part); model != "" {
				models = append(models, model)
			}
		}
	}
	return models
}

// ModelsQueryFunc builds a models query for a line over an inclusive range.
type ModelsQueryFunc func(*db.ProdQueryBuilder, string, time.Time, time.Time) (string, []interface{}, error)

// ModelsResponse is one observed model.
type ModelsResponse struct {
	Model string `json:"model"`
}

// Models handles a distinct-models graph request:
// GET /graph/api/models/{line_id}?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD
//
// It returns a JSON array of model names observed for the line in the range.
func Models(w http.ResponseWriter, r *http.Request, build ModelsQueryFunc) {
	lineID, startDate, endDate, models, ok := parseLineRange(w, r)
	if !ok {
		return
	}

	qb := db.NewProdQueryBuilder().WithModels(models)
	if _, exists, err := qb.LineConfig(lineID); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to load line configuration", "CONFIG_ERROR")
		return
	} else if !exists {
		SendErrorResponse(w, http.StatusNotFound, "Line not found", "LINE_NOT_FOUND")
		return
	}

	query, args, err := build(qb, lineID, startDate, endDate)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to build query", "QUERY_ERROR")
		return
	}

	rows, err := db.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Database query failed", "DATABASE_ERROR")
		return
	}
	defer rows.Close()

	results := []ModelsResponse{}
	for rows.Next() {
		var model sql.NullString
		if err := rows.Scan(&model); err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, "Row scan failed", "SCAN_ERROR")
			return
		}
		if model.Valid && model.String != "" {
			results = append(results, ModelsResponse{Model: model.String})
		}
	}

	if err := rows.Err(); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Row iteration failed", "ROWS_ERROR")
		return
	}

	writeJSON(w, results)
}

// fetchPeriodTotals validates the request, runs the period query and returns
// its rows. On failure it writes the error response and returns ok=false.
func fetchPeriodTotals(w http.ResponseWriter, r *http.Request, build PeriodQueryFunc) (string, []periodRow, bool) {
	lineID, startDate, endDate, models, ok := parseLineRange(w, r)
	if !ok {
		return "", nil, false
	}

	qb := db.NewProdQueryBuilder().WithModels(models)
	if _, exists, err := qb.LineConfig(lineID); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to load line configuration", "CONFIG_ERROR")
		return "", nil, false
	} else if !exists {
		SendErrorResponse(w, http.StatusNotFound, "Line not found", "LINE_NOT_FOUND")
		return "", nil, false
	}

	query, args, err := build(qb, lineID, startDate, endDate)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to build query", "QUERY_ERROR")
		return "", nil, false
	}

	rows, err := db.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Database query failed", "DATABASE_ERROR")
		return "", nil, false
	}
	defer rows.Close()

	var results []periodRow
	for rows.Next() {
		var period string
		var model sql.NullString
		var prod sql.NullInt64
		if err := rows.Scan(&period, &model, &prod); err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, "Row scan failed", "SCAN_ERROR")
			return "", nil, false
		}

		row := periodRow{period: period}
		if prod.Valid {
			row.prod = prod.Int64
		}
		if model.Valid {
			row.model = model.String
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Row iteration failed", "ROWS_ERROR")
		return "", nil, false
	}

	if len(results) == 0 {
		SendErrorResponse(w, http.StatusNotFound, "No production data found", "NO_DATA_FOUND")
		return "", nil, false
	}

	return lineID, results, true
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Printf("Failed to encode response: %v\n", err)
	}
}
