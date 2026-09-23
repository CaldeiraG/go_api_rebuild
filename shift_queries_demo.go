//go:build ignore

package main

import (
	"fmt"
	"time"
)

// This file demonstrates the shift-based production queries
// Run with: go run shift_queries_demo.go

func main() {
	// Import the db package (this would be in your main project)
	// import "github.com/caldeirag/go-api/src/db"
	
	qb := db.NewProdQueryBuilder()

	// Example date: 2026-09-21
	date := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)

	fmt.Println("==========================================================")
	fmt.Println("SHIFT-BASED PRODUCTION QUERIES DEMO")
	fmt.Println("Date: 2026-09-21")
	fmt.Println("==========================================================\n")

	// Get all shifts for the date
	shifts := qb.GetAllShiftsForDate(date)

	fmt.Println("=== Shift Periods ===")
	for i, shift := range shifts {
		fmt.Printf("%d. Shift %s\n", i+1, shift.ShiftType)
		fmt.Printf("   Start: %s\n", shift.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   End:   %s\n", shift.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Build queries for Line 45 (Gen3.8 Line A)
	lineID := "45"
	fmt.Printf("\n=== Line %s Queries ===\n", lineID)
	
	// Shift 1
	query1, _ := qb.GetShiftProduction(lineID, date, Shift1)
	fmt.Printf("\nShift 1 (01:00 to 01:00):\n%s\n", query1)

	// Shift 2
	query2, _ := qb.GetShiftProduction(lineID, date, Shift2)
	fmt.Printf("\nShift 2 (01:00 to 01:00):\n%s\n", query2)

	// Shift 3
	query3, _ := qb.GetShiftProduction(lineID, date, Shift3)
	fmt.Printf("\nShift 3 (01:00 to 01:00):\n%s\n", query3)

	// Example with different line
	fmt.Println("\n=== Line 53 (Gen3.8 Inverter) - Shift 2 ===")
	query53, _ := qb.GetShiftProduction("53", date, Shift2)
	fmt.Println(query53)

	// Example with GEN5 line (should fail)
	fmt.Println("\n=== Line 1112 (GEN5 D-Total) - Should fail ===")
	_, err := qb.GetShiftProduction("1112", date, Shift2)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	fmt.Println("\n==========================================================")
	fmt.Println("DEMO COMPLETE")
	fmt.Println("==========================================================")
}

// ShiftType is defined in the db package
type ShiftType string

const (
	Shift1 ShiftType = "shift1"
	Shift2 ShiftType = "shift2"
	Shift3 ShiftType = "shift3"
)

// ProdQueryBuilder is defined in the db package
type ProdQueryBuilder struct{}

// GetShiftProduction returns a query string
func (qb *ProdQueryBuilder) GetShiftProduction(lineID string, date time.Time, shiftType ShiftType) (string, error) {
	// Implementation in db package
	return "", nil
}

// GetAllShiftsForDate returns all shift configs for a date
func (qb *ProdQueryBuilder) GetAllShiftsForDate(date time.Time) []ShiftConfig {
	return []ShiftConfig{}
}

// ShiftConfig holds shift-specific configuration
type ShiftConfig struct {
	ShiftType ShiftType
	StartTime time.Time
	EndTime   time.Time
}
