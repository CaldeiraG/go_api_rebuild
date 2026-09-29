package intranet

import (
	"fmt"
	"net/http"
)

// writeError reports a server-side failure using the same plain-text shape as
// the other non-graph handlers.
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something Went Wrong!"))
	fmt.Printf("Error: %v\n", err)
}
