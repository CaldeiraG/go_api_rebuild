package gen5

import (
	"context"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	mssql "github.com/denisenkom/go-mssqldb"
)

type productionDailyStruct struct {
	Hour      int
	Timestamp mssql.DateTime1
	LineAOK   int
	LineANOK  int
	LineBOK   int
	LineBNOK  int
	LineC1OK  int
	LineC1NOK int
	LineC2OK  int
	LineC2NOK int
	LineDOK   int
	LineDNOK  int
}

func ProductionGEN5_cron() {
	ctx := context.Background()
	production := productionDailyStruct{}

	query, err := config.DB.Prepare(config.ProductionGEN5)
	if err != nil {
		fmt.Printf("Could not prepare query: %v\n", err)
		return
	}

	defer query.Close()

	rows, err := query.QueryContext(ctx)
	if err != nil {
		//w.WriteHeader(500)
		//w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Could not execute query: %v\n", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&production.Hour, &production.Timestamp, &production.LineAOK, &production.LineANOK, &production.LineBOK, &production.LineBNOK, &production.LineC1OK, &production.LineC1NOK, &production.LineC2OK, &production.LineC2NOK, &production.LineDOK, &production.LineDNOK)
		if err != nil {
			fmt.Printf("Could not be scanned to struct: %v\n", err)
			return
		}

		fmt.Println(production)
	}

	//fmt.Println(production)
}
