package main

import (
	"context"
	"database/sql"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	config "github.com/caldeirag/go-api/src/db"
)

// TestLineHealth is an opt-in live database health check. It walks every line
// in prodFAssy_config.json, runs the same queries the /graph endpoints run and
// reports how long each takes.
//
// Run it with:
//
//	HEALTHCHECK=1 go test -run TestLineHealth -v .
//
// Optional environment variables:
//
//	HEALTHCHECK_DATE     date to query (YYYY-MM-DD, default today)
//	HEALTHCHECK_SLOW_MS  log lines slower than this (default 2000)
func TestLineHealth(t *testing.T) {
	if os.Getenv("HEALTHCHECK") != "1" {
		t.Skip("live DB health check skipped; set HEALTHCHECK=1 (with DB_* in .env) to run")
	}

	conn, err := connectDB()
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}

	date := time.Now()
	if v := os.Getenv("HEALTHCHECK_DATE"); v != "" {
		date, err = time.ParseInLocation("2006-01-02", v, time.Local)
		if err != nil {
			t.Fatalf("invalid HEALTHCHECK_DATE %q: %v", v, err)
		}
	}

	slowThreshold := 2 * time.Second
	if v := os.Getenv("HEALTHCHECK_SLOW_MS"); v != "" {
		if ms, convErr := strconv.Atoi(v); convErr == nil && ms > 0 {
			slowThreshold = time.Duration(ms) * time.Millisecond
		}
	}

	qb := config.NewProdQueryBuilder()
	lines, err := qb.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	lineIDs := make([]string, 0, len(lines))
	for id := range lines {
		lineIDs = append(lineIDs, id)
	}
	sort.Strings(lineIDs)

	var results []lineHealthResult
	started := time.Now()

	for _, lineID := range lineIDs {
		cfg := lines[lineID]
		if cfg.Error || cfg.DatabaseInUse == "" {
			t.Logf("line %-5s %-28s skipped (not implemented)", lineID, cfg.Name)
			continue
		}
		results = append(results,
			measureLine(ctx, conn, qb, lineID, cfg.Name, "hourly", date),
			measureLine(ctx, conn, qb, lineID, cfg.Name, "hourlynok", date),
			measureLine(ctx, conn, qb, lineID, cfg.Name, "daily", date),
			measureLine(ctx, conn, qb, lineID, cfg.Name, "dailynok", date),
			measureLine(ctx, conn, qb, lineID, cfg.Name, "monthly", date),
			measureLine(ctx, conn, qb, lineID, cfg.Name, "monthlynok", date),
		)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].elapsed > results[j].elapsed })

	t.Logf("line health check for %s (sorted slowest first)", date.Format("2006-01-02"))
	t.Logf("%-6s %-28s %-9s %4s %6s %10s  %s", "line", "name", "endpoint", "q", "rows", "elapsed", "status")

	failures := 0
	var slow []lineHealthResult
	for _, r := range results {
		status := "ok"
		switch {
		case r.err != nil:
			status = "ERR"
			failures++
			t.Errorf("line %s (%s) %s failed: %v", r.lineID, r.name, r.kind, r.err)
		case r.elapsed > slowThreshold:
			status = "SLOW"
			slow = append(slow, r)
		}
		t.Logf("%-6s %-28s %-9s %4d %6d %10s  %s",
			r.lineID, r.name, r.kind, r.queries, r.rows, r.elapsed.Round(time.Millisecond), status)
	}

	t.Logf("%d line/endpoint checks in %s; %d failed, %d slow (> %s)",
		len(results), time.Since(started).Round(time.Millisecond), failures, len(slow), slowThreshold)

	if failures > 0 {
		t.Fatalf("%d of %d line checks failed", failures, len(results))
	}
}

type lineHealthResult struct {
	lineID  string
	name    string
	kind    string
	queries int
	rows    int
	elapsed time.Duration
	err     error
}

// measureLine runs every shift query for one line and endpoint, timing the lot.
func measureLine(
	ctx context.Context,
	conn *sql.DB,
	qb *config.ProdQueryBuilder,
	lineID, name, kind string,
	date time.Time,
) (res lineHealthResult) {
	res = lineHealthResult{lineID: lineID, name: name, kind: kind}

	started := time.Now()
	defer func() { res.elapsed = time.Since(started) }()

	var queries []string

	switch kind {
	case "daily", "dailynok", "monthly", "monthlynok":
		var query string
		var err error
		switch kind {
		case "daily":
			query, err = qb.BuildDayQuery(lineID, date, date)
		case "dailynok":
			query, err = qb.BuildDayQueryNOK(lineID, date, date)
		case "monthly":
			query, err = qb.BuildMonthQuery(lineID, date, date)
		default:
			query, err = qb.BuildMonthQueryNOK(lineID, date, date)
		}
		if err != nil {
			res.err = err
			return res
		}
		queries = append(queries, query)
	default:
		for _, shift := range qb.GetAllShiftsForDate(date, lineID) {
			var query string
			var err error
			if kind == "hourlynok" {
				query, err = qb.GetShiftProductionNOK(lineID, date, shift.ShiftType)
			} else {
				query, err = qb.GetShiftProduction(lineID, date, shift.ShiftType)
			}
			if err != nil {
				res.err = err
				return res
			}
			queries = append(queries, query)
		}
	}
	res.queries = len(queries)

	for _, query := range queries {
		n, err := countRows(ctx, conn, query)
		res.rows += n
		if err != nil {
			res.err = err
			return res
		}
	}

	return res
}

func countRows(ctx context.Context, conn *sql.DB, query string) (int, error) {
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	count := 0
	for rows.Next() {
		// Scan into dynamic destinations so any column shape is supported.
		values := make([]interface{}, len(columns))
		dest := make([]interface{}, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return count, err
		}
		count++
	}
	return count, rows.Err()
}
