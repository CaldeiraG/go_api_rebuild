package monthly

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// MonthlyProduction handles the monthly production totals endpoint.
// GET /graph/api/monthly/{line_id}?startDate=2026-01-01&endDate=2026-12-31
func MonthlyProduction(w http.ResponseWriter, r *http.Request) {
	common.Monthly(w, r, (*db.ProdQueryBuilder).BuildMonthQuery)
}
