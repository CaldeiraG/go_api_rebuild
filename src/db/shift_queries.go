package db

import (
	"fmt"
	"strings"
	"time"
)

// ShiftType represents a shift type
type ShiftType string

const (
	Shift1 ShiftType = "1T" // 08:00 - 16:30 same day
	Shift2 ShiftType = "2T" // 16:30 - 01:00 next day
	Shift3 ShiftType = "3T" // 01:00 - 08:00 same day
)

// ShiftConfig holds shift-specific configuration
type ShiftConfig struct {
	ShiftType ShiftType
	StartTime time.Time
	EndTime   time.Time
}

// GetShiftConfig returns the ShiftConfig for a given shift type and date
func GetShiftConfig(shiftType ShiftType, date time.Time) ShiftConfig {
	switch shiftType {
	case Shift1:
		// Shift 1: 08:00 to 16:30 same day
		return ShiftConfig{
			ShiftType: Shift1,
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 16, 30, 0, 0, date.Location()),
		}
	case Shift2:
		// Shift 2: 16:30 to 01:00 next day
		return ShiftConfig{
			ShiftType: Shift2,
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 16, 30, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day()+1, 1, 0, 0, 0, date.Location()),
		}
	case Shift3:
		// Shift 3: 01:00 to 08:00 same day
		return ShiftConfig{
			ShiftType: Shift3,
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 1, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, date.Location()),
		}
	default:
		return ShiftConfig{}
	}
}

// ShiftForHour returns the production shift an hour of the day belongs to:
// Shift 1 08:00-16:30, Shift 2 16:30-01:00, Shift 3 01:00-08:00.
//
// Hourly data cannot represent the 16:30 boundary exactly, so hour 16 is
// attributed to shift 2 (the shift that starts within that hour).
func ShiftForHour(hour int) ShiftType {
	switch {
	case hour >= 1 && hour < 8:
		return Shift3
	case hour >= 8 && hour < 16:
		return Shift1
	default: // 16..23 and 0 (00:00-01:00 belongs to shift 2)
		return Shift2
	}
}

// addCondition appends a WHERE condition, ensuring a single " AND " separator.
// A leading "AND" is stripped so both "REJECTED = 0" and "AND REJECTED = 0"
// work, which matters now that the model prefix that used to supply the
// connector is disabled.
func addCondition(whereClause, condition string) string {
	condition = strings.TrimSpace(condition)
	if len(condition) >= 4 && strings.EqualFold(condition[:4], "and ") {
		condition = strings.TrimSpace(condition[4:])
	}
	if condition == "" {
		return whereClause
	}
	return whereClause + " AND " + condition
}

// standardHourlyQuery builds the hourly breakdown query used by standard
// (non-GEN5) lines. Lines without a model_id fall back to a literal 'n/a'
// model and group by the hour alone, since an empty model expression would
// otherwise produce invalid SQL.
func standardHourlyQuery(lineConfig *ProdQueryConfig, whereClause string) string {
	modelExpr := lineConfig.ModelID
	groupBy := fmt.Sprintf("DATEPART(hh,%s)", lineConfig.DateTime)
	if modelExpr == "" {
		modelExpr = "'n/a'"
	} else {
		groupBy += ", " + modelExpr
	}

	return fmt.Sprintf(`
		SELECT 
			DATEPART(hh,%s) AS hora,
			COUNT(%s) AS prod,
			%s as model
		FROM %s
		WHERE %s
		GROUP BY %s
	`, lineConfig.DateTime, lineConfig.ID, modelExpr, lineConfig.DatabaseInUse, whereClause, groupBy)
}

// BuildShiftQuery builds a query for a specific 24-hour shift period
func (qb *ProdQueryBuilder) BuildShiftQuery(
	lineID string,
	date time.Time,
	shiftType ShiftType,
) (string, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		return "", fmt.Errorf("line ID %s not found in configuration", lineID)
	}

	// Calculate shift time range based on shift type
	var startHour, endHour int

	switch shiftType {
	case Shift1:
		// 08:00 to 16:30 same day
		startHour = 8
		endHour = 16
	case Shift2:
		// 16:30 to 01:00 next day
		startHour = 16
		endHour = 1
	case Shift3:
		// 01:00 to 08:00 same day
		startHour = 1
		endHour = 8
	}

	// Format date for SQL
	dateStr := date.Format("2006-01-02")
	dateStrEnd := date.AddDate(0, 0, 1).Format("2006-01-02")

	// Build WHERE clause with exact shift times
	var whereClause string

	if lineConfig.QueryType == "gen5" {
		// GEN5 uses Timestamp
		whereClause = fmt.Sprintf(
			"Timestamp >= '%s 00:00' AND Timestamp < '%s 00:00'",
			dateStr, dateStrEnd,
		)
	} else {
		// Standard lines use DateTime with exact shift times
		if startHour == 1 && endHour == 8 {
			// Shift 3: 01:00 to 08:00
			whereClause = fmt.Sprintf(
				"%s >= '%s 01:00' AND %s < '%s 08:00'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStr,
			)
		} else if startHour == 16 && endHour == 1 {
			// Shift 2: 16:30 to 01:00
			whereClause = fmt.Sprintf(
				"%s >= '%s 16:30' AND %s < '%s 01:00'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStrEnd,
			)
		} else {
			// Shift 1: 08:00 to 16:30
			whereClause = fmt.Sprintf(
				"%s >= '%s 08:00' AND %s < '%s 16:30'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStr,
			)
		}
	}

	// paramModel is intentionally not applied yet: its value still contains the
	// legacy $modelFAssy placeholder, which needs the model-check endpoint to be
	// substituted. Re-enable by appending lineConfig.ParamModel here.

	// Add param condition
	whereClause = addCondition(whereClause, lineConfig.Param)

	// Add station filter if paramRejSta exists
	whereClause = addCondition(whereClause, lineConfig.ParamRejSta)

	query := ""

	if lineConfig.QueryType == "gen5" {
		// For GEN5, we use the specific GEN5 query builder
		query = fmt.Sprintf(`
		SELECT 
			Hour as hora,
			max(OK) as prod,
			model
		FROM %s
		WHERE %s
		GROUP BY Hour, Model
	`, lineConfig.DatabaseInUse, whereClause)
	} else {

		// Build query with GROUP BY hour for hourly breakdown
		query = standardHourlyQuery(lineConfig, whereClause)
	}

	return query, nil
}

// GetShiftProduction builds production query for a 24-hour shift (using param for OK)
func (qb *ProdQueryBuilder) GetShiftProduction(
	lineID string,
	date time.Time,
	shiftType ShiftType,
) (string, error) {
	return qb.BuildShiftQuery(lineID, date, shiftType)
}

// GetShiftProductionNOK builds NOK production query for a 24-hour shift (using paramRej)
func (qb *ProdQueryBuilder) GetShiftProductionNOK(
	lineID string,
	date time.Time,
	shiftType ShiftType,
) (string, error) {
	return qb.BuildShiftQueryNOK(lineID, date, shiftType)
}

// BuildShiftQueryNOK builds a NOK shift query using paramRej instead of param
func (qb *ProdQueryBuilder) BuildShiftQueryNOK(
	lineID string,
	date time.Time,
	shiftType ShiftType,
) (string, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		return "", fmt.Errorf("line ID %s not found in configuration", lineID)
	}

	// Calculate shift time range based on shift type
	var startHour, endHour int

	switch shiftType {
	case Shift1:
		// 08:00 to 16:30 same day
		startHour = 8
		endHour = 16
	case Shift2:
		// 16:30 to 01:00 next day
		startHour = 16
		endHour = 1
	case Shift3:
		// 01:00 to 08:00 same day
		startHour = 1
		endHour = 8
	}

	// Format date for SQL
	dateStr := date.Format("2006-01-02")
	dateStrEnd := date.AddDate(0, 0, 1).Format("2006-01-02")

	// Build WHERE clause with exact shift times using paramRej
	var whereClause string

	if lineConfig.QueryType == "gen5" {
		// GEN5 uses Timestamp. Keep the same full-day window as the OK query so
		// /daily and /dailynok report the same 24h period.
		whereClause = fmt.Sprintf(
			"Timestamp >= '%s 00:00' AND Timestamp < '%s 00:00'",
			dateStr, dateStrEnd,
		)
	} else {
		// Standard lines use DateTime with exact shift times
		if startHour == 1 && endHour == 8 {
			// Shift 3: 01:00 to 08:00
			whereClause = fmt.Sprintf(
				"%s >= '%s 01:00' AND %s < '%s 08:00'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStr,
			)
		} else if startHour == 16 && endHour == 1 {
			// Shift 2: 16:30 to 01:00
			whereClause = fmt.Sprintf(
				"%s >= '%s 16:30' AND %s < '%s 01:00'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStrEnd,
			)
		} else {
			// Shift 1: 08:00 to 16:30
			whereClause = fmt.Sprintf(
				"%s >= '%s 08:00' AND %s < '%s 16:30'",
				lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStr,
			)
		}
	}

	// paramModel is intentionally not applied yet: its value still contains the
	// legacy $modelFAssy placeholder, which needs the model-check endpoint to be
	// substituted. Re-enable by appending lineConfig.ParamModel here.

	// Add the NOK/line condition. Standard lines use paramRej for defects,
	// while GEN5 configs express their line selector in param and have no
	// paramRej, so reuse param to keep the query scoped to the requested line.
	if lineConfig.QueryType == "gen5" {
		whereClause = addCondition(whereClause, lineConfig.Param)
	} else {
		whereClause = addCondition(whereClause, lineConfig.ParamRej)
	}

	// Add station filter if paramRejSta exists
	whereClause = addCondition(whereClause, lineConfig.ParamRejSta)

	query := ""

	if lineConfig.QueryType == "gen5" {
		// For GEN5, we use the specific GEN5 query builder
		query = fmt.Sprintf(`
		SELECT 
			Hour as hora,
			max(NOK) as prod,
			model
		FROM %s
		WHERE %s
		GROUP BY Hour, Model
	`, lineConfig.DatabaseInUse, whereClause)
	} else {
		// Build query with GROUP BY hour for hourly breakdown
		query = standardHourlyQuery(lineConfig, whereClause)
	}

	return query, nil
}

// GetAllShiftsForDate returns all 3 shift periods for a date
func (qb *ProdQueryBuilder) GetAllShiftsForDate(date time.Time, lineID string) []ShiftConfig {

	configMap, err := qb.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return nil
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		fmt.Printf("line ID %s not found in configuration\n", lineID)
		return nil
	}

	if lineConfig.QueryType == "gen5" {
		// For GEN5, we can return a single shift covering the whole day
		return []ShiftConfig{GetShiftConfig(Shift1, date)}
	} else {

		return []ShiftConfig{
			GetShiftConfig(Shift1, date),
			GetShiftConfig(Shift2, date),
			GetShiftConfig(Shift3, date),
		}
	}
}
