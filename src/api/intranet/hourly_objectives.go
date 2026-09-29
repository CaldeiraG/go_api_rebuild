package intranet

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	config "github.com/caldeirag/go-api/src/db"
)

// hourlyObjectiveHours is the number of hour columns in Hourly_Objective.
const hourlyObjectiveHours = 24

// HourlyObjective mirrors a row of [Intranet].[dbo].[Hourly_Objective]. The
// hourly values are returned as a nested "hours" array indexed 0-23.
type HourlyObjective struct {
	LineID   *int       `json:"line_id"`
	LineName *string    `json:"line_name"`
	Hours    []*float64 `json:"hours"`
}

// HourlyObjectives handles GET /intranet/hourly-objectives and
// GET /intranet/hourly-objectives/{line_id}. The line filter is optional and
// may also be passed as ?line_id=.
func HourlyObjectives(w http.ResponseWriter, r *http.Request) {
	lineID := chi.URLParam(r, "line_id")
	if lineID == "" {
		lineID = r.URL.Query().Get("line_id")
	}

	query := config.SqlHourlyObjectives
	var args []interface{}

	if lineID != "" {
		id, err := strconv.Atoi(lineID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid line_id"))
			return
		}
		query = config.SqlHourlyObjectivesByLine
		args = append(args, sql.Named("line_id", id))
	}

	rows, err := config.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeError(w, err)
		return
	}
	defer rows.Close()

	objectives := []HourlyObjective{}
	for rows.Next() {
		var lineID *int
		var lineName *string
		hours := make([]*float64, hourlyObjectiveHours)

		dest := make([]interface{}, 0, hourlyObjectiveHours+2)
		dest = append(dest, &lineID, &lineName)
		for i := range hours {
			dest = append(dest, &hours[i])
		}

		if err := rows.Scan(dest...); err != nil {
			writeError(w, err)
			return
		}

		objectives = append(objectives, HourlyObjective{
			LineID:   lineID,
			LineName: lineName,
			Hours:    hours,
		})
	}
	if err := rows.Err(); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objectives)
}
