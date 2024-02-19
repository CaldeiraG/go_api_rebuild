package scrap

import (
	"database/sql"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"net/http"
	"time"
)

func InsertTicketList(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticket")
	dateTicket := r.URL.Query().Get("date")
	price := r.URL.Query().Get("price")
	costCenter := r.URL.Query().Get("costcenter")

	query, err := config.DB.Prepare(config.SqlInsertTicket)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	_, err = query.Exec(sql.Named("ticket", ticketID), sql.Named("date", dateTicket), sql.Named("lastupdated", time.Now()), sql.Named("price", price), sql.Named("costcenter", costCenter))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Ticket inserted!"))
}
