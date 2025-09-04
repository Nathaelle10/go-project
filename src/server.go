package main

import (
	"fmt"
	"log"
	"net/http"
)

func StartWebServer() {
	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)

	port := 8080
	fmt.Printf("Dashboard dispo sur http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
