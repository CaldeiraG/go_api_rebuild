package models

import (
	"net/http"

	"github.com/caldeirag/go-api/src/api/graph/common"
	"github.com/caldeirag/go-api/src/db"
)

// Models handles the distinct-models endpoint.
// GET /graph/api/models/{line_id}?startDate=2026-09-01&endDate=2026-09-30
func Models(w http.ResponseWriter, r *http.Request) {
	common.Models(w, r, (*db.ProdQueryBuilder).BuildModelsQuery)
}
