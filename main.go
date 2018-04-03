package main

import (
	"log"
	"net/http"

	transit "github.com/caveda/qmoves-transit/lib"
)

func handler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	setupEnvironment(transit.EnvNameBilbao, "bilbao_transit.dat")
	s := transit.GetSources()
	log.Printf("Sources read: %v", s)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
