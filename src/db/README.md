# Database Query Builder

This package provides a dynamic SQL query builder for production statistics based on JSON configuration.

## Files

- `db.go` - Database connection and SQL constants
- `db2.go` - Additional SQL constants (GEN2 specific)
- `query_builder.go` - Generic query builder with helper methods
- `prod_query_builder.go` - Production query builder with hourly statistics
- `gen5_queries.go` - GEN5-specific query builders
- `query_helper.go` - Helper methods for multiple query types
- `query_builder_example.go` - Example usage demonstrations

## Configuration

The query builder uses `prodFAssy_config.json` to define query parameters for each production line.

### JSON Structure

```json
{
  "42": {
    "name": "VS14 Line A",
    "databaseInUse": "[VS14LineA_Production].[dbo].[VS14AResults]",
    "ID": "ID",
    "dateTime": "DateTime",
    "param": "GlobalResult=1",
    "paramRej": "GlobalResult=0",
    "paramModel": "AND model_id in ($modelFAssy) AND ",
    "model_id": "",
    "paramRejSta": "GlobalResult=0 AND rejected_station='$station'"
  },
  "1112": {
    "name": "GEN5 D-Total",
    "databaseInUse": "[GEN5ProdStats].[dbo].[Production]",
    "ID": "D_Total_OK",
    "dateTime": "Timestamp",
    "param": "",
    "paramRej": "",
    "paramModel": "",
    "model_id": "",
    "queryType": "gen5"
  }
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Display name for the production line |
| `databaseInUse` | string | SQL table reference |
| `ID` | string | Column name for product count |
| `dateTime` | string | Column name for timestamp |
| `param` | string | WHERE clause for good records |
| `paramRej` | string | WHERE clause for rejected records |
| `paramModel` | string | Model filtering condition |
| `model_id` | string | Column used for model filtering |
| `paramRejSta` | string | Station-based rejection filter |
| `specialModel` | object | Special model logic (only for line 53) |
| `queryType` | string | Query type: "standard" (default) or "gen5" |

## Usage

### Basic Query Building

```go
package main

import (
    "fmt"
    "time"
    "github.com/caldeirag/go-api/src/db"
)

func main() {
    // Create query builder
    qb := db.NewProdQueryBuilder()
    
    // Define date range
    dataInit := time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local)
    dataFinal := time.Date(2026, 9, 21, 23, 59, 59, 0, time.Local)
    
    // Build query for line 45 (Gen3.8 Line A)
    config, err := qb.BuildProdQuery("45", dataInit, dataFinal, "290", nil)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Build the SQL query
    query, err := qb.BuildProdQueryHourly("45", dataInit, dataFinal, "290", nil)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("%s\n", query)
}
```

### Output Example

```sql
SET DATEFIRST 1;
SELECT 
    MAX(DATEPART(hh,DateDataSave)) AS hora,
    COUNT(ID) AS prod
FROM [GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]
WHERE 
    DateDataSave > '2026-09-21 00:00'
    AND DateDataSave < '2026-09-21 16:30'
    AND CodeCeHousing like '$modelFAssy%' AND 
    LineNumber=1 AND LineB1Good=1
GROUP BY DATEPART(hh,DateDataSave);
```

## Methods

### ProdQueryBuilder Methods

- `LoadConfig()` - Load and parse the JSON configuration
- `BuildProdQuery(lineID, dataInit, dataFinal, station, models)` - Build query configuration
- `BuildProdQueryHourly(lineID, startTime, endTime, station, models)` - Build SQL query string
- `GetGen5Query(lineID, dataInit, dataFinal)` - Build GEN5 query
- `GetShiftProduction(lineID, date, shiftType)` - Build query for a specific shift
- `GetAllShiftsForDate(date)` - Get all 3 shift configs for a date
- `GetShiftProductionForDate(lineID, date)` - Build queries for all shifts on a date
- `GetLineInfo(lineID)` - Get configuration for a specific line
- `GetAllLines()` - Get all available production lines
- `ValidateLineID(lineID)` - Check if a line ID exists
- `IsGen5Line(lineID)` - Check if a line is a GEN5 metric

### Gen5Queries Methods

- `BuildGen5Query(config, dataInit, dataFinal)` - Build GEN5 query for a specific config
- `GetGen5Metrics()` - Get all available GEN5 metric IDs

### Generic QueryBuilder Methods

- `BuildQuery(lineID, dataInit, station, models)` - Generic query builder (standard or GEN5)
- `Query(lineID, dataInit, station, models)` - Execute query and return results
- `GetConfigValue(lineID, key)` - Get specific config value
- `PrintConfig()` - Print entire configuration

## Special Cases

### Line 53 (Gen3.8 Inverter)

This line has special model filtering logic for BMW and VW:

```go
// BMW models (T, U, V, Y, Z)
// VW models (A, B, C)
```

Usage:
```go
models := []string{"BMW", "ABC123"}
config, _ := qb.BuildProdQuery("53", startTime, endTime, "", models)
```

### GEN5 Queries (Lines 1107-1114)

GEN5 queries aggregate production metrics from the Production table. These queries:
- Select specific metric columns (A_Total_OK, B_Total_OK, D_Total_OK, etc.)
- Filter by Hour < 17
- Use Timestamp range for date filtering

Available GEN5 metrics:
| Line ID | Metric | Description |
|---------|--------|-------------|
| 1107 | A_Total_OK | Line A total OK |
| 1108 | B_Total_OK | Line B total OK |
| 1109 | C1_Total_OK | Line C1 total OK |
| 1110 | C2_Total_OK | Line C2 total OK |
| 1112 | D_Total_OK | Line D total OK |
| 1113 | BracketVW_Total_OK | VW Bracket total OK |
| 1114 | E_Total_OK | Line E total OK |

Example:
```go
// Get D_Total_OK for 2026-09-21
query, err := qb.BuildProdQueryHourly("1112", startTime, endTime, "", nil)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Output:
// SELECT 
//     Hour as hora,
//     D_Total_OK as prod
// FROM [GEN5ProdStats].[dbo].[Production]
// WHERE Hour < 17 AND Timestamp >= '2026-09-21 00:00' AND Timestamp < '2026-09-22 00:00'
```

### Shift-Based Queries (Per Shift Period)

The query builder supports shift-based queries where each shift covers a 24-hour period from 01:00 to 01:00 the next day.

**Shift Definitions:**
- **Shift 1:** 01:00 (Day N-1) to 01:00 (Day N)
- **Shift 2:** 01:00 (Day N) to 01:00 (Day N+1)
- **Shift 3:** 01:00 (Day N+1) to 01:00 (Day N+2)

Example:
```go
// Get Shift 2 production for 2026-09-21
query, err := qb.GetShiftProduction("45", time.Date(2026, 9, 21, 0, 0, 0, 0, time.Local), Shift2)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Output:
// SELECT 
//     MAX(DATEPART(hh,Data)) AS hora,
//     COUNT(ID) AS prod
// FROM [GEN3_Clone].[dbo].[GEN3_LINE_A_PROD]
// WHERE Data >= '2026-09-21 01:00' AND Data < '2026-09-22 01:00'
//     AND PALETE_DB_SN_SNCH like '$modelFAssy%' AND 
//     PALETE_DB_STATUS_OK_NOK like '1'
```

### Shift Methods

**ProdQueryBuilder Methods:**
- `GetShiftProduction(lineID, date, shiftType)` - Build query for a specific shift
- `GetAllShiftsForDate(date)` - Get all 3 shift configs for a date
- `GetShiftProductionForDate(lineID, date)` - Build queries for all shifts on a date

## Error Handling

All methods return errors for:
- Missing configuration file
- Invalid JSON format
- Unknown line ID
- Missing required fields

Always check errors before using the results.

## Testing

Run the example:
```bash
go run src/db/query_builder_example.go
```

## Dependencies

- Go 1.18+
- JSON configuration file at `prodFAssy_config.json`
