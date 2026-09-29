package db

import (
	"strings"
	"testing"
	"time"
)

func TestBuildDayQueryStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildDayQuery("52", testDate, testDate.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("BuildDayQuery: %v", err)
	}

	for _, want := range []string{
		"CONVERT(varchar(10), TIME_STAMP, 23) AS [date]",
		"RTRIM(LTRIM(MINDEX)) AS model",
		"COUNT(MINDEX) AS prod",
		"[DB].[dbo].[Table]",
		"TIME_STAMP >= '2026-09-21 00:00'",
		"TIME_STAMP < '2026-09-23 00:00'", // endDate + 1 day, exclusive
		"REJECTED = 0",
		"AND STATION = 1",
		"GROUP BY CONVERT(varchar(10), TIME_STAMP, 23), RTRIM(LTRIM(MINDEX))",
	} {
		if !strings.Contains(query, want) {
			t.Errorf("day query missing %q:\n%s", want, query)
		}
	}
	if strings.Contains(query, "REJECTED = 1") {
		t.Errorf("OK day query should not contain the NOK condition:\n%s", query)
	}
}

func TestBuildDayQueryNOKStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildDayQueryNOK("52", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildDayQueryNOK: %v", err)
	}
	if !strings.Contains(query, "REJECTED = 1") {
		t.Errorf("NOK day query should use paramRej:\n%s", query)
	}
	if strings.Contains(query, "REJECTED = 0") {
		t.Errorf("NOK day query should not use the OK param:\n%s", query)
	}
}

func TestBuildDayQueryEmptyModelID(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	for _, build := range []struct {
		name  string
		query func(string, time.Time, time.Time) (string, error)
	}{
		{"daily", qb.BuildDayQuery},
		{"dailynok", qb.BuildDayQueryNOK},
	} {
		query, err := build.query("42", testDate, testDate)
		if err != nil {
			t.Fatalf("%s: %v", build.name, err)
		}
		if !strings.Contains(query, "'n/a' AS model") {
			t.Errorf("%s should fall back to an 'n/a' model:\n%s", build.name, query)
		}
		if !strings.Contains(query, "GROUP BY CONVERT(varchar(10), DateTime, 23), 'n/a'") {
			t.Errorf("%s should group by day and the model literal:\n%s", build.name, query)
		}
	}
}

func TestBuildDayQueryGen5(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	ok, err := qb.BuildDayQuery("1107", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildDayQuery: %v", err)
	}
	for _, want := range []string{
		"CONVERT(varchar(10), Timestamp, 23) AS [date]",
		"RTRIM(LTRIM(Model)) AS model",
		"SUM(OK) AS prod",
		"[GEN5ProdStats].[dbo].[Production]",
		"Timestamp >= '2026-09-21 00:00'",
		"Timestamp < '2026-09-22 00:00'",
		"AND Line = 'A'",
		"GROUP BY CONVERT(varchar(10), Timestamp, 23), RTRIM(LTRIM(Model))",
	} {
		if !strings.Contains(ok, want) {
			t.Errorf("gen5 day query missing %q:\n%s", want, ok)
		}
	}

	nok, err := qb.BuildDayQueryNOK("1107", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildDayQueryNOK: %v", err)
	}
	if !strings.Contains(nok, "SUM(NOK) AS prod") {
		t.Errorf("gen5 NOK day query should sum NOK:\n%s", nok)
	}
	if !strings.Contains(nok, "AND Line = 'A'") {
		t.Errorf("gen5 NOK day query should be scoped to the line:\n%s", nok)
	}
}

func TestBuildDayQueryUnknownLine(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	if _, err := qb.BuildDayQuery("nope", testDate, testDate); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
	if _, err := qb.BuildDayQueryNOK("nope", testDate, testDate); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
}
