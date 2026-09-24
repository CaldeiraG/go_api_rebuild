package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ProdQueryConfig represents the configuration for a production line
type ProdQueryConfig struct {
	Name         string `json:"name"`
	DatabaseInUse string `json:"databaseInUse"`
	ID           string `json:"ID"`
	DateTime     string `json:"dateTime"`
	Param        string `json:"param"`
	ParamRej     string `json:"paramRej"`
	ParamModel   string `json:"paramModel"`
	ModelID      string `json:"model_id"`
	ParamRejSta  string `json:"paramRejSta"`
	SpecialModel map[string]string `json:"specialModel,omitempty"`
	Error        bool   `json:"error,omitempty"`
	QueryType    string `json:"queryType,omitempty"` // "standard" or "gen5"
	ShiftStart   string `json:"shiftStart"`
	ShiftEnd     string `json:"shiftEnd"`
}

// ProdQueryBuilder builds queries for production statistics
type ProdQueryBuilder struct {
	configPath string
}

// NewProdQueryBuilder creates a new query builder instance
func NewProdQueryBuilder() *ProdQueryBuilder {
	// Get the directory of this file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller information")
	}

	// Get directory path and resolve to project root
	dir := filepath.Dir(filename)
	configPath := filepath.Join(dir, "..", "..", "prodFAssy_config.json")

	return &ProdQueryBuilder{
		configPath: configPath,
	}
}

// LoadConfig loads and parses the JSON configuration
func (qb *ProdQueryBuilder) LoadConfig() (map[string]*ProdQueryConfig, error) {
	data, err := os.ReadFile(qb.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	result := make(map[string]*ProdQueryConfig)
	for lineID, lineData := range config {
		lineMap, ok := lineData.(map[string]interface{})
		if !ok {
			continue
		}

		cfg := &ProdQueryConfig{
			Name:         "",
			DatabaseInUse: "",
			ID:           "",
			DateTime:     "",
			Param:        "",
			ParamRej:     "",
			ParamModel:   "",
			ModelID:      "",
			ParamRejSta:  "",
			Error:        false,
			QueryType:    "",
			ShiftStart:   "",
			ShiftEnd:     "",
		}

		// Safely extract fields - check both existence and type
		if v, exists := lineMap["name"]; exists {
			if s, ok := v.(string); ok {
				cfg.Name = s
			}
		}
		if v, exists := lineMap["databaseInUse"]; exists {
			if s, ok := v.(string); ok {
				cfg.DatabaseInUse = s
			}
		}
		if v, exists := lineMap["ID"]; exists {
			if s, ok := v.(string); ok {
				cfg.ID = s
			}
		}
		if v, exists := lineMap["dateTime"]; exists {
			if s, ok := v.(string); ok {
				cfg.DateTime = s
			}
		}
		if v, exists := lineMap["param"]; exists {
			if s, ok := v.(string); ok {
				cfg.Param = s
			}
		}
		if v, exists := lineMap["paramRej"]; exists {
			if s, ok := v.(string); ok {
				cfg.ParamRej = s
			}
		}
		if v, exists := lineMap["paramModel"]; exists {
			if s, ok := v.(string); ok {
				cfg.ParamModel = s
			}
		}
		if v, exists := lineMap["model_id"]; exists {
			if s, ok := v.(string); ok {
				cfg.ModelID = s
			}
		}
		if v, exists := lineMap["paramRejSta"]; exists {
			if s, ok := v.(string); ok {
				cfg.ParamRejSta = s
			}
		}
		if v, exists := lineMap["specialModel"]; exists {
			if sm, ok := v.(map[string]interface{}); ok {
				cfg.SpecialModel = make(map[string]string)
				for key, val := range sm {
					if s, ok := val.(string); ok {
						cfg.SpecialModel[key] = s
					}
				}
			}
		}
		if v, exists := lineMap["queryType"]; exists {
			if s, ok := v.(string); ok {
				cfg.QueryType = s
			}
		}
		if v, exists := lineMap["shiftStart"]; exists {
			if s, ok := v.(string); ok {
				cfg.ShiftStart = s
			}
		}
		if v, exists := lineMap["shiftEnd"]; exists {
			if s, ok := v.(string); ok {
				cfg.ShiftEnd = s
			}
		}

		result[lineID] = cfg
	}

	return result, nil
}

// BuildProdQuery constructs the production statistics query
func (qb *ProdQueryBuilder) BuildProdQuery(
	lineID string,
	dataInit time.Time,
	dataFinal time.Time,
	station string,
	models []string,
) (*ProdQueryConfig, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	lineConfig, exists := configMap[lineID]
	if !exists {
		return nil, fmt.Errorf("line ID %s not found in configuration", lineID)
	}

	// Handle special model logic (case 53 - Gen3.8 Inverter)
	if lineConfig.QueryType == "" && len(lineConfig.SpecialModel) > 0 {
		// Check for BMW/VW specific models
		for modelName, modelPattern := range lineConfig.SpecialModel {
			if contains(models, modelName) {
				lineConfig.ParamModel = modelPattern
				break
			}
		}
	}

	// Build WHERE clause for standard queries
	dataInitStr := dataInit.Format("2006-01-02 15:04")
	dataFinalStr := dataFinal.Format("2006-01-02 15:04")

	whereParts := []string{
		fmt.Sprintf("%s > '%s 00:00'", lineConfig.DateTime, dataInitStr),
		fmt.Sprintf("%s < '%s 16:30'", lineConfig.DateTime, dataFinalStr),
	}

	// Add paramModel if exists and not empty
	if lineConfig.ParamModel != "" {
		whereParts = append(whereParts, lineConfig.ParamModel)
	}

	// Add param condition
	if lineConfig.Param != "" {
		whereParts = append(whereParts, lineConfig.Param)
	}

	// Build the final WHERE clause
	whereClause := fmt.Sprintf("WHERE %s", joinStrings(whereParts, " AND "))

	// Add station filter if paramRejSta exists
	if station != "" && lineConfig.ParamRejSta != "" {
		stationClause := fmt.Sprintf(lineConfig.ParamRejSta, station)
		if stationClause != "" {
			whereClause += " AND " + stationClause
		}
	}

	// Construct the full SQL query
	query := fmt.Sprintf(`
		SET DATEFIRST 1; 
		SELECT MAX(DATEPART(hh,%s)) AS hora, COUNT(%s) AS prod
		FROM [%s]
		%s
		GROUP BY DATEPART(hh,%s);
	`, lineConfig.DateTime, lineConfig.ID, lineConfig.DatabaseInUse, whereClause, lineConfig.DateTime)

	// Log the query
	fmt.Printf("[DEBUG] Building query for %s (%s):\n  Query: %s\n  Where: %s\n",
		lineID, lineConfig.Name, query, whereClause)

	return lineConfig, nil
}

// GetAllLines returns all available production lines
func (qb *ProdQueryBuilder) GetAllLines() (map[string]*ProdQueryConfig, error) {
	return qb.LoadConfig()
}

// contains checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// joinStrings joins multiple strings with a separator
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
