package hourly

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// HourlyProduction handles the hourly production endpoint.
// GET /graph/api/hourly/{line_id}?startDate=2026-09-21&endDate=2026-09-21
func HourlyProduction(w http.ResponseWriter, r *http.Request) {
	common.Hourly(w, r, (*db.ProdQueryBuilder).GetShiftProduction)
}
