package db

import (
	"fmt"
	"time"
)

// BuildQuery builds a query for either standard or GEN5 lines
func (qb *ProdQueryBuilder) BuildQuery(
	lineID string,
	startTime, endTime time.Time,
	station string,
	models []string,
) (string, error) {
	config, err := qb.BuildProdQuery(lineID, startTime, endTime, station, models)
	if err != nil {
		return "", err
	}

	// Check if this is a GEN5 query
	if config.QueryType == "gen5" {
		return qb.GetGen5Query(lineID, startTime, endTime)
	}

	// Standard query
	return qb.BuildStandardQuery(lineID, startTime, endTime, station, models)
}

// BuildStandardQuery builds a standard production query
func (qb *ProdQueryBuilder) BuildStandardQuery(
	lineID string,
	startTime, endTime time.Time,
	station string,
	models []string,
) (string, error) {
	config, err := qb.BuildProdQuery(lineID, startTime, endTime, station, models)
	if err != nil {
		return "", err
	}

	// Format time for SQL
	startStr := startTime.Format("2006-01-02 15:04")
	endStr := endTime.Format("2006-01-02 15:04")

	// Check if this is a GEN5 query
	if config.QueryType == "gen5" {
		return qb.GetGen5Query(lineID, startTime, endTime)
	}

	query := fmt.Sprintf(`
		SET DATEFIRST 1;
		SELECT 
			MAX(DATEPART(hh,%s)) AS hora,
			COUNT(%s) AS prod
		FROM [%s]
		WHERE 
			%s > '%s 00:00'
			AND %s < '%s 16:30'
			%s
			%s
		GROUP BY DATEPART(hh,%s);
	`,
		config.DateTime,
		config.ID,
		config.DatabaseInUse,
		config.DateTime,
		startStr,
		config.DateTime,
		endStr,
		config.ParamModel,
		config.Param,
		config.DateTime,
	)

	return query, nil
}
