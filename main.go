package main

import (
	"fmt"
	"log"
	"net/http"

	transit "github.com/caveda/qmoves-transit/lib"
)

const downloadFolder string = "./download"

var bilboBus transit.Bilbobus

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received")
	j, err := transit.JsonPresenter{}.FormatList(bilboBus.Lines())
	if err != nil {
		http.Error(w, "Can't get lines", 500)
		return
	}

	fmt.Fprintf(w, j)
	return
}

func main() {
	prepare()
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func prepare() {
	setupEnvironment("./setupEnv.sh")
	sources := bilboBus.GetSources()
	for _, s := range sources {
		log.Printf("Sources read: %v", s)
		transit.Download(s.Uri, s.Path)
	}
	bilboBus.Digest(sources)

	transit.Publish(bilboBus.Lines(), "./gen", transit.JsonPresenter{})
}
