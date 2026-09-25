package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testConfig = `{
  "52": {
    "name": "Test Line",
    "databaseInUse": "[DB].[dbo].[Table]",
    "ID": "MINDEX",
    "dateTime": "TIME_STAMP",
    "param": "REJECTED = 0",
    "paramRej": "REJECTED = 1",
    "paramModel": "AND PN like '%X%' AND ",
    "paramRejSta": "AND STATION = 1",
    "shiftStart": "08:00",
    "shiftEnd": "16:30"
  },
  "1107": {
    "name": "GEN5 Test",
    "databaseInUse": "[GEN5ProdStats].[dbo].[Production]",
    "dateTime": "Timestamp",
    "param": "AND Line = 'A'",
    "queryType": "gen5",
    "shiftStart": "08:00",
    "shiftEnd": "16:30"
  }
}`

func newTestBuilder(t *testing.T, content string) *ProdQueryBuilder {
	t.Helper()

	path := filepath.Join(t.TempDir(), "prodFAssy_config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return &ProdQueryBuilder{configPath: path}
}

func TestLoadConfigParsesFields(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	cfg, err := qb.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	line, ok := cfg["52"]
	if !ok {
		t.Fatal("expected line 52 in config")
	}
	if line.Name != "Test Line" {
		t.Errorf("Name = %q, want Test Line", line.Name)
	}
	if line.Param != "REJECTED = 0" {
		t.Errorf("Param = %q", line.Param)
	}
	if line.ParamRej != "REJECTED = 1" {
		t.Errorf("ParamRej = %q", line.ParamRej)
	}
	if line.ParamModel != "AND PN like '%X%' AND " {
		t.Errorf("ParamModel = %q", line.ParamModel)
	}
	if line.ParamRejSta != "AND STATION = 1" {
		t.Errorf("ParamRejSta = %q", line.ParamRejSta)
	}
	if line.ShiftStart != "08:00" || line.ShiftEnd != "16:30" {
		t.Errorf("shift = %s/%s", line.ShiftStart, line.ShiftEnd)
	}

	gen5, ok := cfg["1107"]
	if !ok {
		t.Fatal("expected line 1107 in config")
	}
	if gen5.QueryType != "gen5" {
		t.Errorf("QueryType = %q, want gen5", gen5.QueryType)
	}
}

func TestLoadConfigCachesAndReloadsOnChange(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	first, err := qb.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	// Mutate the returned map; if the cache is used, the next call sees it.
	first["__cache_sentinel__"] = &ProdQueryConfig{Name: "sentinel"}

	second, err := qb.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if second["__cache_sentinel__"] == nil {
		t.Fatal("expected second LoadConfig to return the cached map")
	}

	// A different file size must invalidate the cache.
	expanded := strings.Replace(testConfig, `"shiftEnd": "16:30"`,
		`"shiftEnd": "16:30", "extra": "1"`, 1)
	if err := os.WriteFile(qb.configPath, []byte(expanded), 0o600); err != nil {
		t.Fatalf("rewrite config: %v", err)
	}

	third, err := qb.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig after change: %v", err)
	}
	if third["__cache_sentinel__"] != nil {
		t.Fatal("expected cache to be invalidated after the file changed")
	}
	if _, ok := third["52"]; !ok {
		t.Fatal("expected line 52 after reload")
	}
}

func TestLineExists(t *testing.T) {
	qb := newTestBuilder(t, testConfig)

	exists, err := qb.LineExists("52")
	if err != nil {
		t.Fatalf("LineExists: %v", err)
	}
	if !exists {
		t.Error("expected line 52 to exist")
	}

	exists, err = qb.LineExists("does-not-exist")
	if err != nil {
		t.Fatalf("LineExists: %v", err)
	}
	if exists {
		t.Error("expected unknown line to not exist")
	}
}

func TestResolveConfigPathHonorsEnv(t *testing.T) {
	expected := filepath.Join(t.TempDir(), "custom.json")
	t.Setenv("PROD_CONFIG_PATH", expected)

	if got := resolveConfigPath(); got != expected {
		t.Errorf("resolveConfigPath() = %q, want %q", got, expected)
	}
}
