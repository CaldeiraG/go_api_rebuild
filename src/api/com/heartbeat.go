package com

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	config "github.com/caldeirag/go-api/src/db"
)

func HeartbeatInsert(w http.ResponseWriter, r *http.Request) {

	machineName := r.URL.Query().Get("machine")
	ip := r.RemoteAddr
	app := r.URL.Query().Get("app")
	timestamp := time.Now()

	query, err := config.DB.Prepare(config.SqlHeartbeat)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer query.Close()

	_, err = query.Exec(sql.Named("machine", machineName), sql.Named("ip", ip), sql.Named("app", app), sql.Named("timestamp", timestamp))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Something Went Wrong!"))
		fmt.Printf("Error: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Heartbeat inserted!"))

}
