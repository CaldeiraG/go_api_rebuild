package intranet

import (
	"encoding/json"
	"net/http"

	config "github.com/caldeirag/go-api/src/db"
)

// Area mirrors a row of [Intranet].[dbo].[areas].
type Area struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Production  *int    `json:"production"`
	Type        *string `json:"type"`
	IsActive    *int    `json:"IsActive"`
}

// Areas handles GET /areas.
func Areas(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.QueryContext(r.Context(), config.SqlAreas)
	if err != nil {
		writeError(w, err)
		return
	}
	defer rows.Close()

	areas := []Area{}
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Name, &area.Description, &area.Production, &area.Type, &area.IsActive); err != nil {
			writeError(w, err)
			return
		}
		areas = append(areas, area)
	}
	if err := rows.Err(); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}
