package production

import (
	"testing"

	config "github.com/caldeirag/go-api/src/db"
)

func TestProductionSQLLineMapping(t *testing.T) {
	cases := []struct {
		lineID int
		want   string
	}{
		{47, config.SqlGEN3},
		{53, config.SqlInv3},
		{52, config.SqlYF},
		{90, config.SqlR744},
		{83, config.SqlInv4},
		{85, config.SqlInv42},
		{91, config.SqlInv43},
		{1110, config.SqlGEN5},
	}

	for _, c := range cases {
		got, ok := productionSQL(c.lineID)
		if !ok {
			t.Errorf("line %d: expected a mapping", c.lineID)
			continue
		}
		if got != c.want {
			t.Errorf("line %d: got the wrong SQL statement", c.lineID)
		}
	}
}

func TestProductionSQLUnknownLine(t *testing.T) {
	if _, ok := productionSQL(9999); ok {
		t.Fatal("expected unknown line to report ok=false")
	}
}

func TestPrepareProductionStmtUnknownLine(t *testing.T) {
	if _, err := prepareProductionStmt(9999); err == nil {
		t.Fatal("expected an error for an unknown line ID")
	}
}
