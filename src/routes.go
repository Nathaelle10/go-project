package main

import (
	"fmt"
	"log"
	"net/http"
)

func StartAPIServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", StatusHandler)
	mux.HandleFunc("/api/scan", ScanHandler)

	fmt.Println("API démarrée sur http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", WithCORS(mux)))
}
