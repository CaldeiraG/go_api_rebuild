package nok

import (
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
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
	var results []dailyNOKResponse

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		// For each date, we query Shift 1 (08:00-16:30) which is within the date
		// Shift 2 and 3 extend into the next day, so we only query them for multi-day ranges
		
		// Query Shift 1 for this date (08:00 to 16:30 same day)
		query1, err := qb.GetShiftProduction(lineID, date, Shift1)
		if err != nil {
			results = append(results, dailyNOKResponse{
				LineID: lineID,
				Date:   date.Format("2006-01-02"),
				Hora:   0,
				Prod:   0,
				Shift:  string(Shift1),
			})
			continue
		}

		// Execute query for Shift 1
		rows1, err := config.DB.Query(query1)
		if err != nil {
			results = append(results, dailyNOKResponse{
				LineID: lineID,
				Date:   date.Format("2006-01-02"),
				Hora:   0,
				Prod:   0,
				Shift:  string(Shift1),
			})
			continue
		}

		// Iterate through all rows for Shift 1
		for rows1.Next() {
			var hora int
			var prod sql.NullInt64
			err = rows1.Scan(&hora, &prod)
			if err != nil {
				continue
			}

			if prod.Valid {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   hora,
					Prod:   prod.Int64,
					Shift:  string(Shift1),
				})
			} else {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   hora,
					Prod:   0,
					Shift:  string(Shift1),
				})
			}
		}
		rows1.Close()

		// For Shift 2 (16:30 to 01:00 next day) - only query if endDate allows
		// Shift 2 ends on the NEXT day, so check if we have space for it
		if !date.Equal(endDate) {
			query2, err := qb.GetShiftProduction(lineID, date, Shift2)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(Shift2),
				})
				continue
			}

			rows2, err := config.DB.Query(query2)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(Shift2),
				})
				continue
			}

			for rows2.Next() {
				var hora int
				var prod sql.NullInt64
				err = rows2.Scan(&hora, &prod)
				if err != nil {
					continue
				}

				if prod.Valid {
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   prod.Int64,
						Shift:  string(Shift2),
					})
				} else {
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   0,
						Shift:  string(Shift2),
					})
				}
			}
			rows2.Close()
		}

		// For Shift 3 (01:00 to 08:00 next day) - only query if there's room for next day
		if !date.Equal(endDate) && date.AddDate(0, 0, 1).Before(endDate.AddDate(0, 0, 1)) {
			query3, err := qb.GetShiftProduction(lineID, date, Shift3)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(Shift3),
				})
				continue
			}

			rows3, err := config.DB.Query(query3)
			if err != nil {
				results = append(results, dailyNOKResponse{
					LineID: lineID,
					Date:   date.Format("2006-01-02"),
					Hora:   0,
					Prod:   0,
					Shift:  string(Shift3),
				})
				continue
			}

			for rows3.Next() {
				var hora int
				var prod sql.NullInt64
				err = rows3.Scan(&hora, &prod)
				if err != nil {
					continue
				}

				if prod.Valid {
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   prod.Int64,
						Shift:  string(Shift3),
					})
				} else {
					results = append(results, dailyNOKResponse{
						LineID: lineID,
						Date:   date.Format("2006-01-02"),
						Hora:   hora,
						Prod:   0,
						Shift:  string(Shift3),
					})
				}
			}
			rows3.Close()
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
