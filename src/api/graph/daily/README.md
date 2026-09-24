# Daily Production API

## Overview

The Daily Production endpoint retrieves production data for a specific line over a date range. It queries production data for all 3 shifts (01:00 to 01:00 next day) for each date in the requested range.

## Endpoint

```
GET /graph/api/daily/{line_id}?startDate={startDate}&endDate={endDate}
```

## Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `line_id` | string | Yes | Line ID from config (e.g., "45", "1112") |
| `startDate` | string | Yes | Start date in YYYY-MM-DD format |
| `endDate` | string | Yes | End date in YYYY-MM-DD format |

## Request Example

```bash
curl "http://localhost:4000/graph/api/daily/45?startDate=2026-09-21&endDate=2026-09-21"
```

## Response Format

### Success (200 OK)

```json
[
  {
    "line_id": "45",
    "date": "2026-09-21",
    "hora": 9,
    "prod": 125,
    "shift": "shift2"
  },
  {
    "line_id": "45",
    "date": "2026-09-21",
    "hora": 10,
    "prod": 156,
    "shift": "shift2"
  },
  {
    "line_id": "45",
    "date": "2026-09-21",
    "hora": 11,
    "prod": 143,
    "shift": "shift2"
  }
]
```

### Error Responses

#### Missing Parameters (400)
```json
{"error": "Missing required parameters: line_id, startDate, endDate"}
```

#### Invalid Date Format (400)
```json
{"error": "Invalid startDate format. Use YYYY-MM-DD"}
```

#### Invalid Date Range (400)
```json
{"error": "endDate must be after or equal to startDate"}
```

#### Query Error (500)
```json
{"error": "Failed to execute query: [error details]"}
```

## Usage Examples

### Single Date
```bash
# Get daily production for Line 45 on 2026-09-21
curl "http://localhost:4000/graph/api/daily/45?startDate=2026-09-21&endDate=2026-09-21"
```

### Date Range
```bash
# Get daily production for Line 45 from 2026-09-21 to 2026-09-23
curl "http://localhost:4000/graph/api/daily/45?startDate=2026-09-21&endDate=2026-09-23"
```

### Multiple Lines
```bash
# Get daily production for multiple lines (use separate requests)
curl "http://localhost:4000/graph/api/daily/45?startDate=2026-09-21&endDate=2026-09-21"
curl "http://localhost:4000/graph/api/daily/53?startDate=2026-09-21&endDate=2026-09-21"
curl "http://localhost:4000/graph/api/daily/1112?startDate=2026-09-21&endDate=2026-09-21"
```

## Implementation Details

### Query Logic

For each date in the range, the endpoint:
1. Builds 3 shift queries (Shift 1, 2, 3)
2. Each shift covers a 24-hour period from 01:00 to 01:00
3. Executes the query using `config.DB`
4. Returns results in the format: `hora, prod` per hour

### Shift Periods

- **Shift 1 (1T):** 08:00 to 16:30 (same day)
- **Shift 2 (2T):** 16:30 to 01:00 (next day)
- **Shift 3 (3T):** 01:00 to 08:00 (next day)

### Example Query (Shift 2 - 16:30 to 01:00)
```sql
SELECT 
    MAX(DATEPART(hh,Data)) AS hora,
    COUNT(ID) AS prod
FROM [GEN3_Clone].[dbo].[GEN3_LINE_A_PROD]
WHERE Data >= '2026-09-21 16:30' AND Data < '2026-09-22 01:00'
    AND PALETE_DB_SN_SNCH like '$modelFAssy%' AND 
    PALETE_DB_STATUS_OK_NOK like '1'
```

### Response Format

Each record represents **one hour** of production data:
- `hora`: Hour of day (0-23)
- `prod`: Production count for that hour
- `shift`: Shift identifier (shift1, shift2, shift3)

## Notes

- Returns results for all 3 shifts for each date in the range
- If a query returns no results, `hora` and `prod` will be 0
- Supports both standard lines and GEN5 lines (1107-1114)
- Uses the dynamic query builder from `src/db/query_builder.go`

## Files

- `day.go` - Main endpoint handler
- `router.go` - Route registration
- `README.md` - This file
