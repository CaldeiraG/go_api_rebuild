package production

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type productionYesterdayStruct struct {
	Prod  int
	Model *string `json:"-"`
}

/*type inTime struct {
	start string
	end   string
	check string
}*/

func ProductionYesterday(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	production := productionYesterdayStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))

	var query *sql.Stmt
	var err error

	switch lineID {
	case 44:
		query, err = config.DB.Prepare(config.SqlVS14)
	case 47:
		query, err = config.DB.Prepare(config.SqlGEN3)
	case 53:
		query, err = config.DB.Prepare(config.SqlInv3)
	case 52:
		query, err = config.DB.Prepare(config.SqlYF)
	case 90:
		query, err = config.DB.Prepare(config.SqlR744)
	case 83:
		query, err = config.DB.Prepare(config.SqlInv4)
	case 85:
		query, err = config.DB.Prepare(config.SqlInv42)
	case 91:
		query, err = config.DB.Prepare(config.SqlInv43)
	}

	//query, err := config.DB.Prepare(sqlSchedule)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer query.Close()

	//var shift1 = "08:00"
	//var shift2 = "16:30"
	//var shift3 = "01:00"

	//currentTime := time.Now()

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

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
