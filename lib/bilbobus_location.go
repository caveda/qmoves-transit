package transit

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path"
)

// LocationParser implements the signature of type Decorator.
// It's responsible for decorating lines with the location of the stops.
func LocationParser(l *[]Line, ts TransitSource) error {
	baseDir := path.Dir(ts.Path)
	p := path.Join(baseDir, gtsfStopsFileName)
	err := UnzipFromArchive(ts.Path, gtsfStopsFileName, baseDir)
	if err != nil {
		log.Printf("Error unzipping %v while parsing input: %v ", p, err)
		return err
	}

	// parse stops txt
	stopsLocation, errParse := parseGTFSStops(p)
	if errParse != nil {
		log.Printf("Error parsing GTFS stops file: %v", err)
		return errParse
	}

	for i, line := range *l {
		for j, s := range line.Stops {
			(*l)[i].Stops[j].Location = stopsLocation[s.Id]
		}
		addRoute(&(*l)[i])
	}
	return nil
}

// parseGTFSStops parses the stops GTFS files producing a map with
// stopId as key and Coordinates as value.
func parseGTFSStops(filePath string) (map[string]Coordinates, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error reading file %v. Error: %v ", filePath, err)
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.FieldsPerRecord = -1 // No checks
	csvr.LazyQuotes = true

	// Prepare containers
	stops := make(map[string]Coordinates)
	firstIgnored := false
	for {
		// Start reading csv
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			log.Printf("Read csv %v", err)
			return nil, err
		}

		if !firstIgnored {
			firstIgnored = true
			continue
		}
		stops[row[0]] = Coordinates{row[4], row[5]}
	}

	return stops, nil
}

func addRoute(l *Line) {
	r := make([]Coordinates, len(l.Stops))
	for j, s := range l.Stops {
		r[j] = s.Location
	}
	(*l).MapRoute = r
}
