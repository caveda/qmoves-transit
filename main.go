package main

import (
		"log"
		"github.com/caveda/qmoves-transit/lib"
	"os"
)

const downloadFolder string = "./download"

var bilboBus transit.Bilbobus


func main() {
	prepare()
}

func prepare() {
	setupEnvironment("./setupEnv.sh")
	sources := bilboBus.GetSources()
	for _, s := range sources {
		log.Printf("Sources read: %v", s)
		if !transit.UseCachedData() || !transit.Exists(s.Path) {
			transit.Download(s.Uri, s.Path, transit.IsFileSizeGreaterThanZero)
		}

	}
	bilboBus.Digest(sources)

	// Data health analysis
	report, err := transit.CheckConsistency(bilboBus.Data())
	if err!=nil {
		log.Printf("Found consistency errors: %s", err)
		os.Exit(-1)
	}

	log.Printf(report)

	// Publishing
	transit.Publish(bilboBus.Data(), "./gen", transit.JsonPresenter{})
	log.Printf("Ready to serve the transit information")
}
