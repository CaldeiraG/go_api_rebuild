package db

import (
	"fmt"
	"time"
)

// ShiftType represents a shift type
type ShiftType string

const (
	Shift1 ShiftType = "1T" // 08:00 - 16:30 same day
	Shift2 ShiftType = "2T" // 16:30 - 01:00 next day
	Shift3 ShiftType = "3T" // 01:00 - 08:00 next day
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
		// Shift 3: 01:00 to 08:00 next day
		return ShiftConfig{
			ShiftType: Shift3,
			StartTime: time.Date(date.Year(), date.Month(), date.Day()+1, 1, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day()+1, 8, 0, 0, 0, date.Location()),
		}
	default:
		return ShiftConfig{}
	}
}

// GetShiftFromConfig returns the shift type from config shift times
func GetShiftFromConfig(startTime, endTime string) (ShiftType, error) {
	// Default is Shift2 (01:00 - 01:00)
	return Shift2, nil
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
		// 01:00 to 08:00 next day
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
			"WHERE Timestamp >= '%s 08:00' AND Timestamp < '%s 08:00'",
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

	// Add paramModel if exists
	if lineConfig.ParamModel != "" {
		whereClause += " " + lineConfig.ParamModel
	}

	// Add param condition
	if lineConfig.Param != "" {
		whereClause += " " + lineConfig.Param
	}

	// Add station filter if paramRejSta exists
	if lineConfig.ParamRejSta != "" {
		stationParam := lineConfig.ParamRejSta
		if stationParam != "" {
			whereClause += " " + stationParam
		}
	}

	// Build query with GROUP BY hour for hourly breakdown
	query := fmt.Sprintf(`
		SELECT 
			DATEPART(hh,%s) AS hora,
			COUNT(%s) AS prod
		FROM %s
		WHERE %s
		GROUP BY DATEPART(hh,%s)
		ORDER BY hora;
	`, lineConfig.DateTime, lineConfig.ID, lineConfig.DatabaseInUse, whereClause, lineConfig.DateTime)

	// Log the full query for debugging
	fmt.Printf("[DEBUG] BuildShiftQuery:\n  LineID: %s\n  Shift: %s\n  Date: %s\n  Query: %s\n",
		lineID, shiftType, dateStr, query)

	return query, nil
}

// GetShiftProduction builds production query for a 24-hour shift
func (qb *ProdQueryBuilder) GetShiftProduction(
	lineID string,
	date time.Time,
	shiftType ShiftType,
) (string, error) {
	return qb.BuildShiftQuery(lineID, date, shiftType)
}

// GetAllShiftsForDate returns all 3 shift periods for a date
func (qb *ProdQueryBuilder) GetAllShiftsForDate(date time.Time) []ShiftConfig {
	return []ShiftConfig{
		GetShiftConfig(Shift1, date),
		GetShiftConfig(Shift2, date),
		GetShiftConfig(Shift3, date),
	}
}

// GetShiftProductionForDate builds queries for all shifts on a date
func (qb *ProdQueryBuilder) GetShiftProductionForDate(
	lineID string,
	date time.Time,
) (map[ShiftType]string, error) {
	shifts := qb.GetAllShiftsForDate(date)
	result := make(map[ShiftType]string)

	for _, shift := range shifts {
		query, err := qb.BuildShiftQuery(lineID, date, shift.ShiftType)
		if err != nil {
			return nil, err
		}
		result[shift.ShiftType] = query
	}

	return result, nil
}
