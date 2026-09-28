package db

import (
	"fmt"
	"time"
)

// BuildDayQuery builds a per-day, per-model OK total over an inclusive date
// range. Days are calendar days (00:00-24:00).
func (qb *ProdQueryBuilder) BuildDayQuery(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildAggregateQuery(lineID, startDate, endDate, false, "day")
}

// BuildDayQueryNOK builds the NOK equivalent of BuildDayQuery.
func (qb *ProdQueryBuilder) BuildDayQueryNOK(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildAggregateQuery(lineID, startDate, endDate, true, "day")
}

// BuildMonthQuery builds a per-month, per-model OK total over an inclusive
// date range. Months are calendar months (YYYY-MM).
func (qb *ProdQueryBuilder) BuildMonthQuery(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildAggregateQuery(lineID, startDate, endDate, false, "month")
}

// BuildMonthQueryNOK builds the NOK equivalent of BuildMonthQuery.
func (qb *ProdQueryBuilder) BuildMonthQueryNOK(lineID string, startDate, endDate time.Time) (string, error) {
	return qb.buildAggregateQuery(lineID, startDate, endDate, true, "month")
}

// buildAggregateQuery aggregates a line by calendar period ("day" or "month")
// and model over an inclusive date range. The period is rendered with
// CONVERT(..., 23) so days/months never collide across years.
func (qb *ProdQueryBuilder) buildAggregateQuery(lineID string, startDate, endDate time.Time, nok bool, period string) (string, error) {
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

	// Window [start 00:00, end+1 00:00).
	whereClause := fmt.Sprintf(
		"%s >= '%s 00:00' AND %s < '%s 00:00'",
		timeCol, startStr, timeCol, endExclusiveStr,
	)

	periodExpr := fmt.Sprintf("CONVERT(varchar(10), %s, 23)", timeCol)
	periodAlias := "[date]"
	orderBy := "[date], model"
	if period == "month" {
		periodExpr = fmt.Sprintf("CONVERT(varchar(7), %s, 23)", timeCol)
		periodAlias = "[month]"
		orderBy = "[month], model"
	}

	if lineConfig.QueryType == "gen5" {
		// GEN5 expresses its line selector in param and aggregates OK/NOK.
		whereClause = addCondition(whereClause, lineConfig.Param)

		agg := "SUM(OK)"
		if nok {
			agg = "SUM(NOK)"
		}

		return fmt.Sprintf(`
		SELECT 
			%s AS %s,
			Model AS model,
			%s AS prod
		FROM %s
		WHERE %s
		GROUP BY %s, Model
		ORDER BY %s
	`, periodExpr, periodAlias, agg, lineConfig.DatabaseInUse, whereClause, periodExpr, orderBy), nil
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
			%s AS %s,
			%s AS model,
			COUNT(%s) AS prod
		FROM %s
		WHERE %s
		GROUP BY %s, %s
		ORDER BY %s
	`, periodExpr, periodAlias, modelExpr, lineConfig.ID, lineConfig.DatabaseInUse, whereClause, periodExpr, modelExpr, orderBy), nil
}
