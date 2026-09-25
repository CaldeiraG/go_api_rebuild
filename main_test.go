package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mvrilo/go-redoc"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()

	return buildRouter(redoc.Redoc{
		DocsPath: "/docs",
		SpecPath: "/static/swagger.json",
		SpecFile: "./static/swagger.json",
		Title:    "test",
	})
}

func TestRouterRouting(t *testing.T) {
	r := testRouter(t)

	tests := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		{"heartbeat", http.MethodGet, "/ping", http.StatusOK},
		{"docs", http.MethodGet, "/docs", http.StatusOK},
		{"graph daily validates before db", http.MethodGet, "/graph/api/daily/1", http.StatusBadRequest},
		{"graph nok validates before db", http.MethodGet, "/graph/api/dailynok/1", http.StatusBadRequest},
		{"legacy graph path removed", http.MethodGet, "/graph/daily/1", http.StatusNotFound},
		{"unknown route", http.MethodGet, "/nope", http.StatusNotFound},
		{"wrong method", http.MethodPost, "/graph/api/daily/1", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.want {
				t.Errorf("%s %s = %d, want %d", tt.method, tt.path, rec.Code, tt.want)
			}
		})
	}
}

func TestSwaggerSpecContainsRoutes(t *testing.T) {
	data, err := os.ReadFile("./static/swagger.json")
	if err != nil {
		t.Fatalf("read swagger.json: %v", err)
	}

	var spec struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("parse swagger.json: %v", err)
	}

	for _, path := range []string{
		"/graph/api/daily/{line_id}",
		"/graph/api/dailynok/{line_id}",
		"/scrap/insertTicket",
		"/scrap/updatePerson",
		"/scrap/updateRequester",
		"/com/heartbeat",
	} {
		if _, ok := spec.Paths[path]; !ok {
			t.Errorf("swagger.json is missing path %s", path)
		}
	}
}
