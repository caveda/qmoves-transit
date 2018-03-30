package main

import (
	"log"
	"net/http"

	transit "github.com/caveda/qmoves-transit/lib"
)

func handler(w http.ResponseWriter, r *http.Request) {
	transit.Fetch("www.exampleWeb.com")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
