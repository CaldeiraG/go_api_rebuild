package scrap

import (
	"database/sql"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"net/http"
	"time"
)

func updateTicketPerson(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	//ticketStruct := TicketDetails{}
	ticketID := r.URL.Query().Get("ticket")
	person := r.URL.Query().Get("person")

	query, err := config.DB.Prepare(config.SqlUpdatePerson)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	_, err = query.Exec(sql.Named("ticket", ticketID), sql.Named("person", person), sql.Named("lastupdated", time.Now()))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Ticket inserted!"))
}
