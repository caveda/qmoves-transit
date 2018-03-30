package main

import (
	"log"
	"net/http"

	"github.com/caveda/qmoves-transit/pkg/fetcher"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fetcher.Fetch("www.exampleWeb.com")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
