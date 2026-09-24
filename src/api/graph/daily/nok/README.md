# Daily NOK Production API

## Overview

The Daily NOK Production endpoint retrieves **NOK (Not OK)** production data for a specific line over a date range. It queries production data for all 3 shifts using the `paramRej` condition instead of `param`.

## Endpoint

```
GET /graph/api/dailynok/{line_id}?startDate={startDate}&endDate={endDate}
```

## Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `line_id` | string | Yes | Line ID from config (e.g., "45", "1112") |
| `startDate` | string | Yes | Start date in YYYY-MM-DD format |
| `endDate` | string | Yes | End date in YYYY-MM-DD format |

## Request Example

```bash
curl "http://localhost:4000/graph/api/dailynok/45?startDate=2026-09-21&endDate=2026-09-21"
```

## Response Format

### Success (200 OK)

```json
[
  {
    "line_id": "45",
    "date": "2026-09-21",
    "hora": 8,
    "prod": 119,
    "shift": "1T"
  },
  {
    "line_id": "45",
    "date": "2026-09-21",
    "hora": 17,
    "prod": 80,
    "shift": "2T"
  },
  {
    "line_id": "45",
    "date": "2026-09-22",
    "hora": 2,
    "prod": 45,
    "shift": "3T"
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
# Get daily NOK production for Line 45 on 2026-09-21
curl "http://localhost:4000/graph/api/dailynok/45?startDate=2026-09-21&endDate=2026-09-21"
```

### Date Range
```bash
# Get daily NOK production for Line 45 from 2026-09-21 to 2026-09-23
curl "http://localhost:4000/graph/api/dailynok/45?startDate=2026-09-21&endDate=2026-09-23"
```

### Multiple Lines
```bash
# Get daily NOK production for multiple lines (use separate requests)
curl "http://localhost:4000/graph/api/dailynok/45?startDate=2026-09-21&endDate=2026-09-21"
curl "http://localhost:4000/graph/api/dailynok/53?startDate=2026-09-21&endDate=2026-09-21"
curl "http://localhost:4000/graph/api/dailynok/1112?startDate=2026-09-21&endDate=2026-09-21"
```

## Implementation Details

### Query Logic

For each date in the range, the endpoint:
1. Builds 3 shift queries (Shift 1, 2, 3)
2. Each shift covers a specific time period
3. Uses `paramRej` condition instead of `param` for NOK data
4. Executes the query using `config.DB`
5. Returns results in the format: `hora, prod` per hour

### Shift Periods

- **Shift 1 (1T):** 08:00 to 16:30 (same day)
- **Shift 2 (2T):** 16:30 to 01:00 (next day)
- **Shift 3 (3T):** 01:00 to 08:00 (next day)

### Example Query (Shift 1 - NOK)

```sql
SELECT 
    DATEPART(hh,DateTime) AS hora,
    COUNT(ID) AS prod
FROM [Database].[dbo].[Table]
WHERE DateTime >= '2026-09-21 08:00' AND DateTime < '2026-09-21 16:30'
  AND GlobalResult=0  -- paramRej condition
  AND ...other conditions...
GROUP BY DATEPART(hh,DateTime)
ORDER BY hora;
```

### Response Format

Each record represents **one hour** of NOK production data:
- `hora`: Hour of day (0-23)
- `prod`: NOK production count for that hour
- `shift`: Shift identifier (1T, 2T, or 3T)

## Notes

- Returns results for all 3 shifts for each date in the range
- If a query returns no results, `hora` and `prod` will be 0
- Supports both standard lines and GEN5 lines (1107-1114)
- Uses the `paramRej` condition for NOK/defect tracking
- Separate endpoint from `/graph/api/daily` (which uses `param` for OK production)

## Files

- `day.go` - Main endpoint handler
- `router.go` - Route registration
- `README.md` - This file
