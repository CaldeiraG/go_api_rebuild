package db

import (
	"strings"
	"testing"
)

func TestBuildModelsQueryStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildModelsQuery("52", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildModelsQuery: %v", err)
	}

	for _, want := range []string{
		"SELECT DISTINCT RTRIM(LTRIM(MINDEX)) AS model",
		"[DB].[dbo].[Table]",
		"TIME_STAMP >= '2026-09-21 00:00'",
		"TIME_STAMP < '2026-09-22 00:00'",
		"RTRIM(LTRIM(MINDEX)) IS NOT NULL AND RTRIM(LTRIM(MINDEX)) <> ''",
		"ORDER BY model",
	} {
		if !strings.Contains(query, want) {
			t.Errorf("models query missing %q:\n%s", want, query)
		}
	}
	// "Observed" means no OK/NOK condition is applied.
	if strings.Contains(query, "REJECTED = 0") || strings.Contains(query, "REJECTED = 1") {
		t.Errorf("models query should not filter OK/NOK:\n%s", query)
	}
}

func TestBuildModelsQueryEmptyModelID(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildModelsQuery("42", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildModelsQuery: %v", err)
	}
	if !strings.Contains(query, "'n/a' AS model") {
		t.Errorf("models query should fall back to 'n/a':\n%s", query)
	}
}

func TestBuildModelsQueryGen5(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildModelsQuery("1107", testDate, testDate)
	if err != nil {
		t.Fatalf("BuildModelsQuery: %v", err)
	}
	for _, want := range []string{
		"SELECT DISTINCT RTRIM(LTRIM(Model)) AS model",
		"[GEN5ProdStats].[dbo].[Production]",
		"Timestamp >= '2026-09-21 00:00'",
		"AND Line = 'A'",
	} {
		if !strings.Contains(query, want) {
			t.Errorf("gen5 models query missing %q:\n%s", want, query)
		}
	}
}

func TestBuildModelsQueryUnknownLine(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	if _, err := qb.BuildModelsQuery("nope", testDate, testDate); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
}
