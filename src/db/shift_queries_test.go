package db

import (
	"strings"
	"testing"
	"time"
)

var testDate = time.Date(2026, time.September, 21, 12, 0, 0, 0, time.UTC)

func TestGetShiftConfig(t *testing.T) {
	tests := []struct {
		name      string
		shift     ShiftType
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "shift 1 same day",
			shift:     Shift1,
			wantStart: time.Date(2026, 9, 21, 8, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 9, 21, 16, 30, 0, 0, time.UTC),
		},
		{
			name:      "shift 2 crosses midnight",
			shift:     Shift2,
			wantStart: time.Date(2026, 9, 21, 16, 30, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 9, 22, 1, 0, 0, 0, time.UTC),
		},
		{
			name:      "shift 3 same day",
			shift:     Shift3,
			wantStart: time.Date(2026, 9, 21, 1, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, 9, 21, 8, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetShiftConfig(tt.shift, testDate)
			if !got.StartTime.Equal(tt.wantStart) {
				t.Errorf("StartTime = %s, want %s", got.StartTime, tt.wantStart)
			}
			if !got.EndTime.Equal(tt.wantEnd) {
				t.Errorf("EndTime = %s, want %s", got.EndTime, tt.wantEnd)
			}
		})
	}
}

func TestGetAllShiftsForDate(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	standard := qb.GetAllShiftsForDate(testDate, "1")
	if len(standard) != 3 {
		t.Fatalf("standard line returned %d shifts, want 3", len(standard))
	}

	gen5 := qb.GetAllShiftsForDate(testDate, "1110")
	if len(gen5) != 1 {
		t.Fatalf("gen5 line returned %d shifts, want 1", len(gen5))
	}
}

func TestBuildShiftQueryStandard(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	shift1, err := qb.BuildShiftQuery("1", testDate, Shift1)
	if err != nil {
		t.Fatalf("BuildShiftQuery: %v", err)
	}
	for _, want := range []string{
		"[DB].[dbo].[Table]",
		"TIME_STAMP >= '2026-09-21 08:00'",
		"TIME_STAMP < '2026-09-21 16:30'",
		"REJECTED = 0",
		"AND PN like '%X%' AND ",
		"AND STATION = 1",
		"GROUP BY DATEPART(hh,TIME_STAMP)",
	} {
		if !strings.Contains(shift1, want) {
			t.Errorf("shift 1 query missing %q:\n%s", want, shift1)
		}
	}
	if strings.Contains(shift1, "REJECTED = 1") {
		t.Errorf("shift 1 query should not contain the NOK condition:\n%s", shift1)
	}

	shift2, err := qb.BuildShiftQuery("1", testDate, Shift2)
	if err != nil {
		t.Fatalf("BuildShiftQuery: %v", err)
	}
	if !strings.Contains(shift2, "TIME_STAMP >= '2026-09-21 16:30'") ||
		!strings.Contains(shift2, "TIME_STAMP < '2026-09-22 01:00'") {
		t.Errorf("shift 2 query has wrong window:\n%s", shift2)
	}

	shift3, err := qb.BuildShiftQuery("1", testDate, Shift3)
	if err != nil {
		t.Fatalf("BuildShiftQuery: %v", err)
	}
	if !strings.Contains(shift3, "TIME_STAMP >= '2026-09-21 01:00'") ||
		!strings.Contains(shift3, "TIME_STAMP < '2026-09-21 08:00'") {
		t.Errorf("shift 3 query has wrong window:\n%s", shift3)
	}
}

func TestBuildShiftQueryGen5UsesFullDay(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	query, err := qb.BuildShiftQuery("1110", testDate, Shift1)
	if err != nil {
		t.Fatalf("BuildShiftQuery: %v", err)
	}
	if !strings.Contains(query, "Timestamp >= '2026-09-21 00:00'") ||
		!strings.Contains(query, "Timestamp < '2026-09-22 00:00'") {
		t.Errorf("gen5 query should cover the full day:\n%s", query)
	}
	if !strings.Contains(query, "GROUP BY Hour, Model") {
		t.Errorf("gen5 query missing group by:\n%s", query)
	}
}

func TestBuildShiftQueryNOKGen5MatchesOKWindow(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	ok, err := qb.BuildShiftQuery("1110", testDate, Shift1)
	if err != nil {
		t.Fatalf("BuildShiftQuery: %v", err)
	}
	nok, err := qb.BuildShiftQueryNOK("1110", testDate, Shift1)
	if err != nil {
		t.Fatalf("BuildShiftQueryNOK: %v", err)
	}

	for _, query := range []string{ok, nok} {
		if !strings.Contains(query, "Timestamp >= '2026-09-21 00:00'") ||
			!strings.Contains(query, "Timestamp < '2026-09-22 00:00'") {
			t.Errorf("gen5 query should use the full-day window:\n%s", query)
		}
	}
	if !strings.Contains(nok, "max(NOK)") {
		t.Errorf("NOK query should aggregate NOK, got:\n%s", nok)
	}
}

func TestBuildShiftQueryUnknownLine(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	if _, err := qb.BuildShiftQuery("nope", testDate, Shift1); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
	if _, err := qb.BuildShiftQueryNOK("nope", testDate, Shift1); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
}
