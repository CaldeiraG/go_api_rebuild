package nok

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// DailyNOK handles the daily NOK (non-conforming) production endpoint.
// GET /graph/api/dailynok/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func DailyNOK(w http.ResponseWriter, r *http.Request) {
	common.Daily(w, r, (*db.ProdQueryBuilder).GetShiftProductionNOK)
}
