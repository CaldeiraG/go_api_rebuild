package intranet

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	config "github.com/caldeirag/go-api/src/db"
)

// Line mirrors a row of [Intranet].[dbo].[lines].
type Line struct {
	ID          int      `json:"id"`
	AreaID      *int     `json:"area_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	IsActive    *int     `json:"IsActive"`
	ObjTCiclo   *float64 `json:"Obj_TCiclo"`
	ObjOEE      *float64 `json:"Obj_OEE"`
	ObjFTT      *float64 `json:"Obj_FTT"`
	ObjDowntime *float64 `json:"Obj_Downtime"`
	ObjScrap    *float64 `json:"Obj_Scrap"`
}

// Lines handles GET /lines and GET /lines/{area_id}. The area filter is
// optional and may also be passed as ?area_id=.
func Lines(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "area_id")
	if areaID == "" {
		areaID = r.URL.Query().Get("area_id")
	}

	query := config.SqlLines
	var args []interface{}

	if areaID != "" {
		id, err := strconv.Atoi(areaID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid area_id"))
			return
		}
		query = config.SqlLinesByArea
		args = append(args, sql.Named("area_id", id))
	}

	rows, err := config.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeError(w, err)
		return
	}
	defer rows.Close()

	lines := []Line{}
	for rows.Next() {
		var line Line
		if err := rows.Scan(
			&line.ID,
			&line.AreaID,
			&line.Name,
			&line.Description,
			&line.IsActive,
			&line.ObjTCiclo,
			&line.ObjOEE,
			&line.ObjFTT,
			&line.ObjDowntime,
			&line.ObjScrap,
		); err != nil {
			writeError(w, err)
			return
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lines)
}
