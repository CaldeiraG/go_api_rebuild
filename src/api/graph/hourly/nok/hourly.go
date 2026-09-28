package nok

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// HourlyNOK handles the hourly NOK (non-conforming) production endpoint.
// GET /graph/api/hourlynok/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func HourlyNOK(w http.ResponseWriter, r *http.Request) {
	common.Hourly(w, r, (*db.ProdQueryBuilder).GetShiftProductionNOK)
}
