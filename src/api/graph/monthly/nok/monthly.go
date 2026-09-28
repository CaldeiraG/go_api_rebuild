package nok

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// MonthlyNOK handles the monthly NOK totals endpoint.
// GET /graph/api/monthlynok/{line_id}?startDate=2026-01-01&endDate=2026-12-31
func MonthlyNOK(w http.ResponseWriter, r *http.Request) {
	common.Monthly(w, r, (*db.ProdQueryBuilder).BuildMonthQueryNOK)
}
