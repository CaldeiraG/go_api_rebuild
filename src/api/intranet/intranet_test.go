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
