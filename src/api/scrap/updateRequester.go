package scrap

import (
	"database/sql"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"net/http"
	"time"
)

func UpdateTicketRequester(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticket")
	requester := r.URL.Query().Get("requester")

	query, err := config.DB.Prepare(config.SqlUpdateRequester)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	_, err = query.Exec(sql.Named("ticket", ticketID), sql.Named("person", requester), sql.Named("lastupdated", time.Now()))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Requester updated!"))
}
