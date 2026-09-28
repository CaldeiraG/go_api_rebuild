package db

import (
	"fmt"
	"time"
)

// BuildDayQuery builds a per-day, per-model OK total over an inclusive date
// range. Days are calendar days (00:00-24:00).
func (qb *ProdQueryBuilder) BuildDayQuery(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildDayQuery(lineID, startDate, endDate, false)
}

// BuildDayQueryNOK builds the NOK equivalent of BuildDayQuery.
func (qb *ProdQueryBuilder) BuildDayQueryNOK(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildDayQuery(lineID, startDate, endDate, true)
}

func (qb *ProdQueryBuilder) buildDayQuery(lineID string, startDate, endDate time.Time, nok bool) (string, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		return "", fmt.Errorf("line ID %s not found in configuration", lineID)
	}

	timeCol := lineConfig.DateTime
	startStr := startDate.Format("2006-01-02")
	endExclusiveStr := endDate.AddDate(0, 0, 1).Format("2006-01-02")

	// Calendar-day window [start 00:00, end+1 00:00).
	whereClause := fmt.Sprintf(
		"%s >= '%s 00:00' AND %s < '%s 00:00'",
		timeCol, startStr, timeCol, endExclusiveStr,
	)

	// Style 23 renders a date as YYYY-MM-DD, avoiding driver type quirks.
	dayExpr := fmt.Sprintf("CONVERT(varchar(10), %s, 23)", timeCol)

	if lineConfig.QueryType == "gen5" {
		// GEN5 expresses its line selector in param and aggregates OK/NOK.
		whereClause = addCondition(whereClause, lineConfig.Param)

		agg := "SUM(OK)"
		if nok {
			agg = "SUM(NOK)"
		}

		return fmt.Sprintf(`
		SELECT 
			%s AS [date],
			Model AS model,
			%s AS prod
		FROM %s
		WHERE %s
		GROUP BY %s, Model
		ORDER BY [date], model
	`, dayExpr, agg, lineConfig.DatabaseInUse, whereClause, dayExpr), nil
	}

	if nok {
		whereClause = addCondition(whereClause, lineConfig.ParamRej)
	} else {
		whereClause = addCondition(whereClause, lineConfig.Param)
	}
	whereClause = addCondition(whereClause, lineConfig.ParamRejSta)

	modelExpr := lineConfig.ModelID
	if modelExpr == "" {
		modelExpr = "'n/a'"
	}

	return fmt.Sprintf(`
		SELECT 
			%s AS [date],
			%s AS model,
			COUNT(%s) AS prod
		FROM %s
		WHERE %s
		GROUP BY %s, %s
		ORDER BY [date], model
	`, dayExpr, modelExpr, lineConfig.ID, lineConfig.DatabaseInUse, whereClause, dayExpr, modelExpr), nil
}
