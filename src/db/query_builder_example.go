package db

import (
	"encoding/json"
	"fmt"
	"time"
)

// Example usage of the ProdQueryBuilder
// This demonstrates how to build and execute queries

func ExampleUsage() {
	// Create a new query builder
	qb := NewProdQueryBuilder()

	// Define the date range (example: 2026-09-21)
	dataInit := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)
	dataFinal := time.Date(2026, 9, 21, 23, 59, 59, 0, time.Local)

	// Example 1: Build query for a specific line (e.g., line 42 - VS14 Line A)
	lineID := "42"
	config, err := qb.BuildProdQuery(lineID, dataInit, dataFinal, "1", nil)
	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Printf("Line: %s\n", config.Name)
	fmt.Printf("Database: %s\n", config.DatabaseInUse)
	fmt.Printf("ID Column: %s\n", config.ID)
	fmt.Printf("DateTime Column: %s\n", config.DateTime)
	fmt.Printf("Model Column: %s\n", config.ModelID)

	// Example 2: Build the actual SQL query
	query, err := qb.BuildProdQueryHourly(lineID, dataInit, dataFinal, "1", nil)
	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Printf("\nGenerated SQL Query:\n%s\n", query)

	// Example 3: Get line information
	info, err := qb.GetLineInfo(lineID)
	if err != nil {
		fmt.Printf("Error getting line info: %v\n", err)
		return
	}

	fmt.Printf("\nLine Information:\n")
	fmt.Printf("  Name: %s\n", info.Name)
	fmt.Printf("  Param: %s\n", info.Param)
	fmt.Printf("  Param Rej: %s\n", info.ParamRej)

	// Example 4: Handle models (if any)
	models := []string{"BMW", "VW"} // Example models
	configWithModels, _ := qb.BuildProdQuery(lineID, dataInit, dataFinal, "1", models)
	fmt.Printf("\nWith Models:\n")
	fmt.Printf("  ParamModel: %s\n", configWithModels.ParamModel)

	// Example 5: Print all available lines
	allLines, _ := qb.GetAllLines()
	fmt.Printf("\nAvailable Lines (%d total):\n", len(allLines))
	for lineID := range allLines {
		fmt.Printf("  - %s\n", lineID)
	}

	// Example 6: JSON output (for debugging)
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("\nConfig as JSON:\n%s\n", string(configJSON))
}
