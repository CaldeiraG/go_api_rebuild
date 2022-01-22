package production

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type productionStruct struct {
	Prod  int
	Model string
}

type modelStruct struct {
	Model string
}

/*type inTime struct {
	start string
	end   string
	check string
}*/

func Production(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	production := productionStruct{}
	model := modelStruct{}
	lineID := StringToInt(chi.URLParam(r, "line_id"))

	var query *sql.Stmt
	var err error

	switch lineID {
	case 44:
		query, err = config.DB.Prepare(config.SqlVS14)
	case 47:
		query, err = config.DB.Prepare(config.SqlGEN3)
	case 53:
		query, err = config.DB2.Prepare(config.SqlInv3)
	case 52:
		query, err = config.DB2.Prepare(config.SqlYF)
	case 90:
		query, err = config.DB2.Prepare(config.SqlR744)
	case 83:
		query, err = config.DB2.Prepare(config.SqlInv4)
	case 85:
		query, err = config.DB2.Prepare(config.SqlInv42)
	case 91:
		query, err = config.DB2.Prepare(config.SqlInv43)
	}

	//query, err := config.DB.Prepare(sqlSchedule)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer query.Close()

	var shift1 = "08:00"
	var shift2 = "16:30"
	var shift3 = "01:00"

	currentTime := time.Now()

	var newRecord *sql.Row

	if currentTime.Format("15:04") >= "08:00" && currentTime.Format("15:04") < "16:30" {
		dataInicial := ConvertTimeCurrentDate(shift1)
		dataFinal := ConvertTimeCurrentDate(shift2)
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	} else if currentTime.Format("15:04") >= "16:30" && currentTime.Format("15:04") < "00:00" {
		dataInicial := ConvertTimeCurrentDate(shift2)
		dataFinal := ConvertTimeCurrentDate("00:00").AddDate(0, 0, 1)
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	} else if currentTime.Format("15:04") >= "00:00" && currentTime.Format("15:04") < "01:00" {
		dataInicial := ConvertTimeCurrentDate("00:00").AddDate(0, 0, -1)
		dataFinal := ConvertTimeCurrentDate(shift3)
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	} else if currentTime.Format("15:04") >= "01:00" && currentTime.Format("15:04") < "08:00" {
		dataInicial := ConvertTimeCurrentDate(shift3)
		dataFinal := ConvertTimeCurrentDate(shift1)
		newRecord = query.QueryRowContext(ctx, sql.Named("dataInicial", dataInicial), sql.Named("dataFinal", dataFinal))
	}

	err = newRecord.Scan(&production.Prod, &production.Model)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	queryModel, err := config.DB.Prepare(config.SqlmodelCheck)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer queryModel.Close()

	var newRecord2 *sql.Row

	switch lineID {
	case 47: //Gen3
		newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
	case 53: //Inv3
		newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
	case 52: //YF
		newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
	case 83: //Inv4
		newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
	case 85: //Inv42
		newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
		//case 91: //Inv43
		//newRecord2 = queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))
	}

	//newRecord2 := queryModel.QueryRowContext(ctx, sql.Named("model", production.Model), sql.Named("line_id", lineID))

	if newRecord2 != nil {
		err = newRecord2.Scan(&model.Model)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Something Went Wrong!"))
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	if model.Model != "" {
		production.Model = model.Model
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

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ConvertTimeCurrentDate(s string) time.Time {
	oldTime, _ := time.Parse("15:04", s)
	hour, min, sec := oldTime.Clock()
	now := time.Now().UTC()
	year, month, day := now.Date()
	oldTime = time.Date(year, month, day, hour, min, sec, 0, time.UTC)
	return oldTime
}

// This is an example of a basic Get Request with Go-chi
// More info can be found here: https://go-chi.io/ && https://github.com/go-chi/chi/tree/master/_examples
