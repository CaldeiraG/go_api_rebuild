package scrap

import (
	"database/sql"
	"fmt"
	config "github.com/caldeirag/go-api/src/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type TicketDetails struct {
	Ticket string
	Date   time.Time
	Person string
	Price  float64
}

func InsertTicketList(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	//ticketStruct := TicketDetails{}
	ticketID := r.URL.Query().Get("ticket")
	dateTicket := chi.URLParam(r, "date")
	price := chi.URLParam(r, "price")
	costCenter := chi.URLParam(r, "costcenter")

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

func StringToFloat(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return i
}
