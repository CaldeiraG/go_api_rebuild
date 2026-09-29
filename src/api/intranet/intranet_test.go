package intranet

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHourlyObjectivesInvalidLineID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/intranet/hourly-objectives/abc", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("line_id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	HourlyObjectives(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body=%q)", rec.Code, rec.Body.String())
	}
}

func TestLinesInvalidAreaID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/lines/abc", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("area_id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	Lines(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body=%q)", rec.Code, rec.Body.String())
	}
}

func jsonKeys(t *testing.T, v interface{}) []string {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	keys := make([]string, 0, len(raw))
	for key := range raw {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func TestAreaJSONKeysMatchColumns(t *testing.T) {
	got := jsonKeys(t, Area{})
	want := []string{"IsActive", "description", "id", "name", "production", "type"}

	if len(got) != len(want) {
		t.Fatalf("keys = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %v, want %v", got, want)
		}
	}
}

func TestHourlyObjectiveJSONKeys(t *testing.T) {
	got := jsonKeys(t, HourlyObjective{})
	want := []string{"hours", "line_id", "line_name"}

	if len(got) != len(want) {
		t.Fatalf("keys = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %v, want %v", got, want)
		}
	}
}

func TestHourlyObjectiveHoursArray(t *testing.T) {
	lineID := 5
	lineName := "Line A"
	first := 12.5

	hours := make([]*float64, 24)
	hours[0] = &first

	data, err := json.Marshal(HourlyObjective{LineID: &lineID, LineName: &lineName, Hours: hours})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	var out []*float64
	if err := json.Unmarshal(raw["hours"], &out); err != nil {
		t.Fatalf("unmarshal hours: %v", err)
	}
	if len(out) != 24 {
		t.Fatalf("hours length = %d, want 24", len(out))
	}
	if out[0] == nil || *out[0] != 12.5 {
		t.Errorf("hours[0] = %v, want 12.5", out[0])
	}
	if out[1] != nil {
		t.Errorf("hours[1] = %v, want null", out[1])
	}
}

func TestLineJSONKeysMatchColumns(t *testing.T) {
	got := jsonKeys(t, Line{})
	want := []string{
		"IsActive", "Obj_Downtime", "Obj_FTT", "Obj_OEE", "Obj_Scrap", "Obj_TCiclo",
		"area_id", "description", "id", "name",
	}

	if len(got) != len(want) {
		t.Fatalf("keys = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %v, want %v", got, want)
		}
	}
}
