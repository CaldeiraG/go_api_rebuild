package scrap

import (
	"database/sql"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"net/http"
	"time"
)

func UpdateTicketPerson(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticket")
	person := r.URL.Query().Get("person")

	if ticketID == "" || person == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required parameters: ticket, person"))
		return
	}

	query, err := config.DB.Prepare(config.SqlUpdatePerson)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer query.Close()

	_, err = query.Exec(sql.Named("ticket", ticketID), sql.Named("person", person), sql.Named("lastupdated", time.Now()))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Person updated!"))
}
