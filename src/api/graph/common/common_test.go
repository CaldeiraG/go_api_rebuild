package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
)

const testConfig = `{
  "1": {
    "name": "Test Line",
    "databaseInUse": "[DB].[dbo].[Table]",
    "ID": "MINDEX",
    "dateTime": "TIME_STAMP",
    "param": "REJECTED = 0",
    "paramRej": "REJECTED = 1"
  },
  "1110": {
    "name": "GEN5 Test",
    "databaseInUse": "[GEN5ProdStats].[dbo].[Production]",
    "dateTime": "Timestamp",
    "param": "AND Line = 'C2'",
    "queryType": "gen5"
  }
}`

func withConfig(t *testing.T, content string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "prodFAssy_config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("PROD_CONFIG_PATH", path)
}

// graphRequest builds a request with the chi line_id route parameter set.
func graphRequest(target, lineID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("line_id", lineID)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func errorBuilder(*db.ProdQueryBuilder, string, time.Time, db.ShiftType) (string, error) {
	return "", errors.New("boom")
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()

	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error body: %v (body=%q)", err, rec.Body.String())
	}
	return body
}

func TestDailyValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		lineID   string
		target   string
		wantCode int
		wantErr  string
	}{
		{
			name:     "missing parameters",
			lineID:   "1",
			target:   "/graph/api/daily/1",
			wantCode: http.StatusBadRequest,
			wantErr:  "MISSING_PARAMETERS",
		},
		{
			name:     "invalid start date",
			lineID:   "1",
			target:   "/graph/api/daily/1?startDate=2026-13-01&endDate=2026-09-21",
			wantCode: http.StatusBadRequest,
			wantErr:  "INVALID_DATE",
		},
		{
			name:     "invalid end date",
			lineID:   "1",
			target:   "/graph/api/daily/1?startDate=2026-09-21&endDate=2026-13-01",
			wantCode: http.StatusBadRequest,
			wantErr:  "INVALID_DATE",
		},
		{
			name:     "end before start",
			lineID:   "1",
			target:   "/graph/api/daily/1?startDate=2026-09-22&endDate=2026-09-21",
			wantCode: http.StatusBadRequest,
			wantErr:  "INVALID_DATE_RANGE",
		},
		{
			name:     "unknown line",
			lineID:   "999",
			target:   "/graph/api/daily/999?startDate=2026-09-21&endDate=2026-09-21",
			wantCode: http.StatusNotFound,
			wantErr:  "LINE_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withConfig(t, testConfig)

			rec := httptest.NewRecorder()
			Daily(rec, graphRequest(tt.target, tt.lineID), (*db.ProdQueryBuilder).GetShiftProduction)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tt.wantCode, rec.Body.String())
			}
			if got := decodeError(t, rec).Code; got != tt.wantErr {
				t.Errorf("error code = %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func TestDailyReturnsPartialErrorsWithoutDatabase(t *testing.T) {
	withConfig(t, testConfig)

	rec := httptest.NewRecorder()
	Daily(rec, graphRequest(
		"/graph/api/daily/1?startDate=2026-09-21&endDate=2026-09-21", "1"),
		errorBuilder)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}

	var results []DailyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatalf("decode results: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want one per shift (3)", len(results))
	}
	for _, r := range results {
		if r.Error == "" {
			t.Errorf("expected a per-row error, got %+v", r)
		}
	}
}

func TestDailyGen5UsesSingleShiftPerDay(t *testing.T) {
	withConfig(t, testConfig)

	rec := httptest.NewRecorder()
	Daily(rec, graphRequest(
		"/graph/api/dailynok/1110?startDate=2026-09-21&endDate=2026-09-22", "1110"),
		errorBuilder)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}

	var results []DailyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatalf("decode results: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want one per day (2)", len(results))
	}
}
