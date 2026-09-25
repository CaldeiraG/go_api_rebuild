package production

import (
	"database/sql"
	"fmt"

	config "github.com/caldeirag/go-api/src/db"
)

// productionSQL maps a line ID to the statement used to read its production.
func productionSQL(lineID int) (string, bool) {
	switch lineID {
	case 47:
		return config.SqlGEN3, true
	case 53:
		return config.SqlInv3, true
	case 52:
		return config.SqlYF, true
	case 90:
		return config.SqlR744, true
	case 83:
		return config.SqlInv4, true
	case 85:
		return config.SqlInv42, true
	case 91:
		return config.SqlInv43, true
	case 1110:
		return config.SqlGEN5, true
	default:
		return "", false
	}
}

// prepareProductionStmt prepares the production statement for a line.
// It returns an explicit error for unknown line IDs instead of leaving the
// caller with a nil statement (which would panic on use).
func prepareProductionStmt(lineID int) (*sql.Stmt, error) {
	query, ok := productionSQL(lineID)
	if !ok {
		return nil, fmt.Errorf("unsupported line ID: %d", lineID)
	}
	return config.DB.Prepare(query)
}
