package db

import (
	"fmt"
	"time"
)

// ShiftType represents a shift type
type ShiftType string

const (
	Shift1 ShiftType = "shift1" // 01:00 - 01:00 (24h period 1)
	Shift2 ShiftType = "shift2" // 01:00 - 01:00 (24h period 2)
	Shift3 ShiftType = "shift3" // 01:00 - 01:00 (24h period 3)
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
		// Shift 1: 01:00 previous day to 01:00 today
		return ShiftConfig{
			ShiftType: Shift1,
			StartTime: time.Date(date.Year(), date.Month(), date.Day()-1, 1, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 1, 0, 0, 0, date.Location()),
		}
	case Shift2:
		// Shift 2: 01:00 today to 01:00 next day
		return ShiftConfig{
			ShiftType: Shift2,
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 1, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day()+1, 1, 0, 0, 0, date.Location()),
		}
	case Shift3:
		// Shift 3: 01:00 next day to 01:00 2 days later
		return ShiftConfig{
			ShiftType: Shift3,
			StartTime: time.Date(date.Year(), date.Month(), date.Day()+1, 1, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day()+2, 1, 0, 0, 0, date.Location()),
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

	// Calculate shift time range (01:00 to 01:00 next day)
	// Each shift is 24 hours from 01:00 to 01:00
	// Shift 1: -1 day, Shift 2: 0 day, Shift 3: +1 day
	dateOffset := 0
	if shiftType == "shift1" {
		dateOffset = -1
	} else if shiftType == "shift3" {
		dateOffset = 1
	}

	startDate := date.AddDate(0, 0, dateOffset)
	endDate := date.AddDate(0, 0, dateOffset+1)

	// Format date for SQL
	dateStr := startDate.Format("2006-01-02")
	dateStrEnd := endDate.Format("2006-01-02")

	// Build WHERE clause
	var whereClause string

	if lineConfig.QueryType == "gen5" {
		// GEN5 uses Timestamp
		whereClause = fmt.Sprintf(
			"WHERE Timestamp >= '%s 01:00' AND Timestamp < '%s 01:00'",
			dateStr, dateStrEnd,
		)
	} else {
		// Standard lines use DateTime
		whereClause = fmt.Sprintf(
			"WHERE %s >= '%s 01:00' AND %s < '%s 01:00'",
			lineConfig.DateTime, dateStr, lineConfig.DateTime, dateStrEnd,
		)
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

	// Build query
	query := fmt.Sprintf(`
		SELECT 
			MAX(DATEPART(hh,%s)) AS hora,
			COUNT(%s) AS prod
		FROM [%s]
		%s
	`, lineConfig.DateTime, lineConfig.ID, lineConfig.DatabaseInUse, whereClause)

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
