package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// QueryBuilder handles dynamic SQL query construction from JSON config
type QueryBuilder struct {
	configPath string
}

// QueryResult holds the query result
type QueryResult struct {
	Hora  int    `json:"hora"`
	Prod  int64  `json:"prod"`
	Error string `json:"error,omitempty"`
}

// LoadConfig loads the JSON configuration
func (qb *QueryBuilder) LoadConfig() (map[string]interface{}, error) {
	// Get the directory of this file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get caller information")
	}

	// Get directory path
	dir := filepath.Dir(filename)
	configPath := filepath.Join(dir, "..", "..", "prodFAssy_config.json")

	// Check if config file exists in current directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try relative path from current working directory
		configPath = filepath.Join("..", "..", "prodFAssy_config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %s", configPath)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return config, nil
}

// BuildQuery constructs the SQL query for a specific line
func (qb *QueryBuilder) BuildQuery(lineID string, dataInit time.Time, station string, models []string) (string, error) {
	config, err := qb.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	lineConfig, exists := config[lineID]
	if !exists {
		return "", fmt.Errorf("line ID %s not found in config", lineID)
	}

	lineConfigMap := lineConfig.(map[string]interface{})

	// Extract configuration values
	databaseInUse, _ := lineConfigMap["databaseInUse"].(string)
	idCol, _ := lineConfigMap["ID"].(string)
	dateTimeCol, _ := lineConfigMap["dateTime"].(string)
	paramModel, _ := lineConfigMap["paramModel"].(string)

	// Build WHERE conditions
	// Format: 'YYYY-MM-DD HH:MM'
	dataInitStr := dataInit.Format("2006-01-02 15:04")
	whereClause := fmt.Sprintf(`WHERE %s > '%s 00:00' AND %s < '%s 16:30'`, dateTimeCol, dataInitStr, dateTimeCol, dataInitStr)

	// Add paramModel if not empty
	if paramModel != "" {
		whereClause += " " + paramModel
	}

	// Add param if exists (from config)
	param, exists := lineConfigMap["param"].(string)
	if exists && param != "" {
		whereClause += " " + param
	}

	// Add station filter if paramRejSta exists
	paramRejSta, exists := lineConfigMap["paramRejSta"].(string)
	if exists && paramRejSta != "" && station != "" {
		// Replace $station placeholder with actual station value
		paramRejSta = fmt.Sprintf(paramRejSta, station)
		whereClause += " " + paramRejSta
	}

	// Build the main query
	query := fmt.Sprintf(`
		SET DATEFIRST 1; 
		SELECT MAX(DATEPART(hh,%s)) AS hora, COUNT(%s) AS prod
		FROM [%s]
		%s
		GROUP BY DATEPART(hh,%s);
	`, dateTimeCol, idCol, databaseInUse, whereClause, dateTimeCol)

	log.Printf("Built query for line %s: %s", lineID, query)

	return query, nil
}

// Query executes the built query and returns results
func (qb *QueryBuilder) Query(lineID string, dataInit time.Time, station string, models []string) ([]QueryResult, error) {
	query, err := qb.BuildQuery(lineID, dataInit, station, models)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Note: In actual implementation, you would execute this query using your database connection
	// This is a placeholder to show the structure
	log.Printf("Executing query: %s", query)

	// Example of how you would execute it:
	// rows, err := DB.Query(query)
	// if err != nil {
	//     return nil, err
	// }
	// defer rows.Close()

	// var results []QueryResult
	// for rows.Next() {
	//     var result QueryResult
	//     if err := rows.Scan(&result.Hora, &result.Prod); err != nil {
	//         return nil, err
	//     }
	//     results = append(results, result)
	// }

	return []QueryResult{}, nil
}

// GetConfigValue gets a specific value from the config
func (qb *QueryBuilder) GetConfigValue(lineID string, key string) (interface{}, error) {
	config, err := qb.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	lineConfig, exists := config[lineID]
	if !exists {
		return nil, fmt.Errorf("line ID %s not found", lineID)
	}

	lineConfigMap := lineConfig.(map[string]interface{})
	value, exists := lineConfigMap[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found for line %s", key, lineID)
	}

	return value, nil
}

// PrintConfig prints the entire configuration
func (qb *QueryBuilder) PrintConfig() {
	config, err := qb.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(config)
}

// ValidateLineID validates if a line ID exists in the config
func (qb *QueryBuilder) ValidateLineID(lineID string) bool {
	config, err := qb.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return false
	}

	_, exists := config[lineID]
	return exists
}
