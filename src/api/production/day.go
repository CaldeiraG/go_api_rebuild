package production

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type productionDayStruct struct {
	Prod  int
	Model *string `json:"-"`
}

func DayProduction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	production := productionDayStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))

	query, err := prepareProductionStmt(lineID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer query.Close()

	currentTime := time.Now()

	var newRecord *sql.Row

	if currentTime.Format("15:04") >= "00:00" && currentTime.Format("15:04") < "01:00" {
		dataInicial := ConvertTimeCurrentDate("01:00").AddDate(0, 0, -1)
		dataFinal := currentTime
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	} else {
		dataInicial := ConvertTimeCurrentDate("01:00")
		dataFinal := currentTime.Add(time.Hour * 1)
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	}
	err = newRecord.Scan(&production.Prod, &production.Model)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	errJSON := json.NewEncoder(w).Encode(production)
	if errJSON != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		return
	}
}
