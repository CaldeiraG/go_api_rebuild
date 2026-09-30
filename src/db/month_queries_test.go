package db

import (
	"strings"
	"testing"
	"time"
)

var (
	monthStart = time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	monthEnd   = time.Date(2026, time.September, 30, 0, 0, 0, 0, time.UTC)
)

func TestBuildMonthQueryStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, _, err := qb.BuildMonthQuery("52", monthStart, monthEnd)
	if err != nil {
		t.Fatalf("BuildMonthQuery: %v", err)
	}

	for _, want := range []string{
		"CONVERT(varchar(7), TIME_STAMP, 23) AS [month]",
		"RTRIM(LTRIM(MINDEX)) AS model",
		"COUNT(MINDEX) AS prod",
		"[DB].[dbo].[Table]",
		"TIME_STAMP >= '2026-09-01 00:00'",
		"TIME_STAMP < '2026-10-01 00:00'", // endDate + 1 day, exclusive
		"REJECTED = 0",
		"GROUP BY CONVERT(varchar(7), TIME_STAMP, 23), RTRIM(LTRIM(MINDEX))",
	} {
		if !strings.Contains(query, want) {
			t.Errorf("month query missing %q:\n%s", want, query)
		}
	}
}

func TestBuildMonthQueryNOKStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, _, err := qb.BuildMonthQueryNOK("52", monthStart, monthEnd)
	if err != nil {
		t.Fatalf("BuildMonthQueryNOK: %v", err)
	}
	if !strings.Contains(query, "REJECTED = 1") {
		t.Errorf("NOK month query should use paramRej:\n%s", query)
	}
	if strings.Contains(query, "REJECTED = 0") {
		t.Errorf("NOK month query should not use the OK param:\n%s", query)
	}
}

func TestBuildMonthQueryEmptyModelID(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	for _, build := range []struct {
		name  string
		query func(string, time.Time, time.Time) (string, []interface{}, error)
	}{
		{"monthly", qb.BuildMonthQuery},
		{"monthlynok", qb.BuildMonthQueryNOK},
	} {
		query, _, err := build.query("42", monthStart, monthEnd)
		if err != nil {
			t.Fatalf("%s: %v", build.name, err)
		}
		if !strings.Contains(query, "'n/a' AS model") {
			t.Errorf("%s should fall back to an 'n/a' model:\n%s", build.name, query)
		}
	}
}

func TestBuildMonthQueryGen5(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	ok, _, err := qb.BuildMonthQuery("1107", monthStart, monthEnd)
	if err != nil {
		t.Fatalf("BuildMonthQuery: %v", err)
	}
	for _, want := range []string{
		"CONVERT(varchar(7), Timestamp, 23) AS [month]",
		"RTRIM(LTRIM(Model)) AS model",
		"SUM(OK) AS prod",
		"[GEN5ProdStats].[dbo].[Production]",
		"AND Line = 'A'",
		"GROUP BY CONVERT(varchar(7), Timestamp, 23), RTRIM(LTRIM(Model))",
	} {
		if !strings.Contains(ok, want) {
			t.Errorf("gen5 month query missing %q:\n%s", want, ok)
		}
	}

	nok, _, err := qb.BuildMonthQueryNOK("1107", monthStart, monthEnd)
	if err != nil {
		t.Fatalf("BuildMonthQueryNOK: %v", err)
	}
	if !strings.Contains(nok, "SUM(NOK) AS prod") {
		t.Errorf("gen5 NOK month query should sum NOK:\n%s", nok)
	}
}

func TestBuildMonthQueryUnknownLine(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	if _, _, err := qb.BuildMonthQuery("nope", monthStart, monthEnd); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
	if _, _, err := qb.BuildMonthQueryNOK("nope", monthStart, monthEnd); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
}
