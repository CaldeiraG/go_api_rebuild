package daily

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// DailyProduction handles the daily production totals endpoint.
// GET /graph/api/daily/{line_id}?startDate=2026-09-25&endDate=2026-09-25
func DailyProduction(w http.ResponseWriter, r *http.Request) {
	common.Daily(w, r, (*db.ProdQueryBuilder).BuildDayQuery)
}
