package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// ProdQueryConfig represents the configuration for a production line
type ProdQueryConfig struct {
	Name          string            `json:"name"`
	DatabaseInUse string            `json:"databaseInUse"`
	ID            string            `json:"ID"`
	DateTime      string            `json:"dateTime"`
	Param         string            `json:"param"`
	ParamRej      string            `json:"paramRej"`
	ParamModel    string            `json:"paramModel"`
	ModelID       string            `json:"model_id"`
	ParamRejSta   string            `json:"paramRejSta"`
	SpecialModel  map[string]string `json:"specialModel,omitempty"`
	Error         bool              `json:"error,omitempty"`
	QueryType     string            `json:"queryType,omitempty"` // "standard" or "gen5"
	ShiftStart    string            `json:"shiftStart"`
	ShiftEnd      string            `json:"shiftEnd"`
}

// ProdQueryBuilder builds queries for production statistics
type ProdQueryBuilder struct {
	configPath string
}

// Config cache. The config file is read on every query build; caching by
// modification time avoids re-reading/parsing it while still picking up edits
// made to the file at runtime.
var (
	configCacheMu   sync.RWMutex
	configCache     map[string]*ProdQueryConfig
	configCacheMod  time.Time
	configCacheSize int64
	configCachePath string
)

// NewProdQueryBuilder creates a new query builder instance
func NewProdQueryBuilder() *ProdQueryBuilder {
	return &ProdQueryBuilder{
		configPath: resolveConfigPath(),
	}
}

// resolveConfigPath resolves prodFAssy_config.json in a way that works both
// from a source checkout and from a compiled/deployed binary:
//  1. PROD_CONFIG_PATH environment variable (explicit override)
//  2. next to the running executable
//  3. project root relative to this source file (local development fallback)
func resolveConfigPath() string {
	if p := os.Getenv("PROD_CONFIG_PATH"); p != "" {
		return p
	}

	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), "prodFAssy_config.json")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	if _, filename, _, ok := runtime.Caller(0); ok {
		p := filepath.Join(filepath.Dir(filename), "..", "..", "prodFAssy_config.json")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "prodFAssy_config.json"
}

// LineConfig returns the configuration for a line and whether it exists.
func (qb *ProdQueryBuilder) LineConfig(lineID string) (*ProdQueryConfig, bool, error) {
	configMap, err := qb.LoadConfig()
	if err != nil {
		return nil, false, err
	}
	cfg, exists := configMap[lineID]
	return cfg, exists, nil
}

// LineExists reports whether a line is present in the production config.
func (qb *ProdQueryBuilder) LineExists(lineID string) (bool, error) {
	_, exists, err := qb.LineConfig(lineID)
	return exists, err
}

// LoadConfig loads and parses the JSON configuration, caching the result until
// the underlying file changes.
func (qb *ProdQueryBuilder) LoadConfig() (map[string]*ProdQueryConfig, error) {
	var modTime time.Time
	var size int64
	if info, err := os.Stat(qb.configPath); err == nil {
		modTime = info.ModTime()
		size = info.Size()

		configCacheMu.RLock()
		if configCache != nil && configCachePath == qb.configPath &&
			configCacheSize == size && configCacheMod.Equal(modTime) {
			cached := configCache
			configCacheMu.RUnlock()
			return cached, nil
		}
		configCacheMu.RUnlock()
	}

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
			Name:          "",
			DatabaseInUse: "",
			ID:            "",
			DateTime:      "",
			Param:         "",
			ParamRej:      "",
			ParamModel:    "",
			ModelID:       "",
			ParamRejSta:   "",
			Error:         false,
			QueryType:     "",
			ShiftStart:    "",
			ShiftEnd:      "",
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

	configCacheMu.Lock()
	configCache = result
	configCacheMod = modTime
	configCacheSize = size
	configCachePath = qb.configPath
	configCacheMu.Unlock()

	return result, nil
}
