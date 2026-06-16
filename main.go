package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/example/sstpinger/pkg/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/add", api.AddHandler)
	mux.HandleFunc("/api/list", api.ListHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
