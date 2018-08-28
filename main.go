package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caveda/qmoves-transit/lib"
	)

const downloadFolder string = "./download"

var bilboBus transit.Bilbobus

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received")
	// j, err := transit.JsonPresenter{}.FormatList(bilboBus.Data().lines)
	// if err != nil {
	// 	http.Error(w, "Can't get lines", 500)
	// 	return
	// }

	// fmt.Fprintf(w, j)
	return
}

func main() {
	prepare()
	//http.HandleFunc("/", handler)
	//
	//log.Fatal(http.ListenAndServe(":8081", nil))
}

func prepare() {
	setupEnvironment("./setupEnv.sh")
	sources := bilboBus.GetSources()
	for _, s := range sources {
		log.Printf("Sources read: %v", s)
		if !transit.UseCachedData() || !transit.Exists(s.Path) {
			transit.Download(s.Uri, s.Path)
		}

	}
	bilboBus.Digest(sources)

	transit.Publish(bilboBus.Data(), "./gen", transit.JsonPresenter{})
	log.Printf("Ready to serve the transit information")
}
