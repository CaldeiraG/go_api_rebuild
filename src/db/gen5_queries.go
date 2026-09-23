package db

import (
	"fmt"
	"time"
)

// BuildGen5Query builds the GEN5 production query
func (qb *ProdQueryBuilder) BuildGen5Query(
	config *ProdQueryConfig,
	dataInit time.Time,
	dataFinal time.Time,
) error {
	// Format time for SQL
	dataInitStr := dataInit.Format("2006-01-02 00:00")
	dataFinalStr := dataFinal.Format("2006-01-02 00:00")

	// Build the WHERE clause
	whereClause := fmt.Sprintf(
		"WHERE Hour < 17 AND Timestamp >= '%s 00:00' AND Timestamp < '%s 00:00'",
		dataInitStr, dataFinalStr,
	)

	// Construct the query
	query := fmt.Sprintf(`
		SELECT 
			Hour as hora,
			%s as prod
		FROM [%s]
		%s
	`, config.ID, config.DatabaseInUse, whereClause)

	// Log the query
	fmt.Printf("[DEBUG] Building GEN5 query for %s (%s):\n  Query: %s\n",
		config.Name, config.Name, query)

	// Return the config with the query built
	// In actual implementation, you would execute this query using DB.Query(query)
	return nil
}

// GetGen5Query builds a complete GEN5 query string
func (qb *ProdQueryBuilder) GetGen5Query(
	lineID string,
	dataInit time.Time,
	dataFinal time.Time,
) (string, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		return "", fmt.Errorf("line ID %s not found in configuration", lineID)
	}

	if lineConfig.QueryType != "gen5" {
		return "", fmt.Errorf("line ID %s is not a GEN5 query type", lineID)
	}

	// Format time for SQL
	dataInitStr := dataInit.Format("2006-01-02 00:00")
	dataFinalStr := dataFinal.Format("2006-01-02 00:00")

	// Build the WHERE clause
	whereClause := fmt.Sprintf(
		"WHERE Hour < 17 AND Timestamp >= '%s 00:00' AND Timestamp < '%s 00:00'",
		dataInitStr, dataFinalStr,
	)

	// Construct the query
	query := fmt.Sprintf(`
		SELECT 
			Hour as hora,
			%s as prod
		FROM [%s]
		%s
	`, lineConfig.ID, lineConfig.DatabaseInUse, whereClause)

	// Log the query
	fmt.Printf("[DEBUG] GEN5 Query for %s (%s):\n  Query: %s\n",
		lineConfig.Name, lineConfig.Name, query)

	return query, nil
}

// GetGen5Metrics returns all GEN5 metric IDs
func (qb *ProdQueryBuilder) GetGen5Metrics() []string {
	return []string{
		"1107", // A_Total_OK
		"1108", // B_Total_OK
		"1109", // C1_Total_OK
		"1110", // C2_Total_OK
		"1112", // D_Total_OK
		"1113", // BracketVW_Total_OK
		"1114", // E_Total_OK
	}
}

// IsGen5Line checks if a line ID is a GEN5 metric
func (qb *ProdQueryBuilder) IsGen5Line(lineID string) bool {
	gen5Metrics := qb.GetGen5Metrics()
	for _, metric := range gen5Metrics {
		if metric == lineID {
			return true
		}
	}
	return false
}
