//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Test if the server is running and route is accessible
	fmt.Println("Testing API routes...")
	
	urls := []string{
		"http://localhost:4000/ping",
		"http://localhost:4000/graph/daily/45?startDate=2026-09-21&endDate=2026-09-22",
		"http://localhost:4000/production/day/45",
	}
	
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s: Error - %v\n", url, err)
		} else {
			fmt.Printf("%s: Status %d\n", url, resp.StatusCode)
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
