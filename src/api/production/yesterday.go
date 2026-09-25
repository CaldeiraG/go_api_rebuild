package production

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type productionYesterdayStruct struct {
	Prod  int
	Model *string `json:"-"`
}

func ProductionYesterday(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	production := productionYesterdayStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))

	query, err := prepareProductionStmt(lineID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer query.Close()

	var newRecord *sql.Row

	dataInicial := ConvertTimeCurrentDate("01:00").AddDate(0, 0, -1)
	dataFinal := ConvertTimeCurrentDate("01:00")
	newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))

	err = newRecord.Scan(&production.Prod, &production.Model)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	errJSON := json.NewEncoder(w).Encode(production)
	if errJSON != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}
}
