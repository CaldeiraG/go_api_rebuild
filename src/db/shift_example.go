package db

import (
	"fmt"
	"time"
)

// ExampleShiftQueries demonstrates shift-based production queries
func ExampleShiftQueries() {
	qb := NewProdQueryBuilder()

	// Example date: 2026-09-21
	date := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)

	// Get all shifts for the date
	shifts := qb.GetAllShiftsForDate(date)

	fmt.Println("=== Shift Periods for 2026-09-21 ===")
	for _, shift := range shifts {
		fmt.Printf("\nShift %s:\n", shift.ShiftType)
		fmt.Printf("  Start: %s\n", shift.StartTime.Format("15:04"))
		fmt.Printf("  End:   %s\n", shift.EndTime.Format("15:04"))
	}

	// Build query for Line 45 (Gen3.8 Line A) - Shift 2
	lineID := "45"
	query, err := qb.GetShiftProduction(lineID, date, Shift2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n=== Query for Line 45, Shift 2 ===")
	fmt.Println(query)

	// Build queries for all shifts
	fmt.Println("\n=== All Shift Queries for Line 45 ===")
	_, err = qb.GetShiftProductionForDate(lineID, date)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, shiftType := range []ShiftType{"shift1", "shift2", "shift3"} {
		query, err := qb.GetShiftProduction(lineID, date, shiftType)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", shiftType, err)
			continue
		}
		fmt.Printf("\n--- %s ---\n", shiftType)
		fmt.Println(query)
	}

	// Example with different line
	fmt.Println("\n=== Query for Line 53 (Gen3.8 Inverter), Shift 1 ===")
	query53, err := qb.GetShiftProduction("53", date, Shift1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(query53)

	// Example with GEN5 line (should fail)
	fmt.Println("\n=== Query for Line 1112 (GEN5 D-Total) - Should fail ===")
	queryGen5, err := qb.GetShiftProduction("1112", date, Shift2)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		fmt.Println(queryGen5)
	}
}

// ExampleGetAllShifts demonstrates getting all shifts for multiple dates
func ExampleGetAllShifts() {
	qb := NewProdQueryBuilder()

	startDate := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)
	lineID := "45"

	fmt.Println("=== Shift 2 Queries for Multiple Dates ===")
	for i := 0; i < 5; i++ {
		date := startDate.AddDate(0, 0, i)
		query, err := qb.GetShiftProduction(lineID, date, Shift2)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", date.Format("2006-01-02"), err)
			continue
		}

		fmt.Printf("\n%s:\n", date.Format("2006-01-02"))
		fmt.Printf("%s\n", query)
	}
}

// ExampleGetShiftInfo demonstrates getting shift information from config
func ExampleGetShiftInfo() {
	qb := NewProdQueryBuilder()

	lineID := "45"
	configMap, err := qb.LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	config, exists := configMap[lineID]
	if !exists {
		fmt.Printf("Line ID %s not found\n", lineID)
		return
	}

	fmt.Printf("Line: %s\n", config.Name)
	fmt.Printf("Shift Start: %s\n", config.ShiftStart)
	fmt.Printf("Shift End: %s\n", config.ShiftEnd)

	// For a standard line, these are typically 08:00 - 16:30
	// But shift queries use 01:00 - 01:00 periods
}

// ExampleShift24HourPeriod demonstrates the 24-hour shift concept
func ExampleShift24HourPeriod() {
	qb := NewProdQueryBuilder()

	// Example: Get Shift 1, 2, 3 for 2026-09-21
	date := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)

	fmt.Println("=== 24-Hour Shift Periods ===")
	fmt.Println("Shift 1: 2026-09-20 01:00 to 2026-09-21 01:00")
	fmt.Println("Shift 2: 2026-09-21 01:00 to 2026-09-22 01:00")
	fmt.Println("Shift 3: 2026-09-22 01:00 to 2026-09-23 01:00")

	// Build queries for all shifts
	fmt.Println("\n=== Generated Queries ===")
	for _, shiftType := range []ShiftType{"shift1", "shift2", "shift3"} {
		query, err := qb.GetShiftProduction("45", date, shiftType)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", shiftType, err)
			continue
		}
		fmt.Printf("\n%s:\n", shiftType)
		fmt.Printf("%s\n", query)
	}
}

// ExampleProductionByShift demonstrates real-world usage
func ExampleProductionByShift() {
	qb := NewProdQueryBuilder()

	// User wants to see production for 2026-09-21
	date := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)
	lineID := "45"

	fmt.Printf("Production report for %s - Line %s\n", date.Format("2006-01-02"), lineID)
	fmt.Println("==========================================================")

	// Get all shifts
	for _, shiftType := range []ShiftType{"shift1", "shift2", "shift3"} {
		query, err := qb.GetShiftProduction(lineID, date, shiftType)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", shiftType, err)
			continue
		}
		fmt.Printf("%s:\n", shiftType)
		fmt.Printf("%s\n", query)
		fmt.Println("  [Query would be executed here]")
		fmt.Println("  Expected result: hora=X, prod=Y")
	}
}
