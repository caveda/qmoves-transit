package transit

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	lines []Line
}

// Constants
const EnvNameBilbao string = "BILBAO_TRANSIT"
const DirectionForward string = "FORWARD"
const DirectionBackward string = "BACKWARD"
const separator string = "-"
const gtsfStopsFileName string = "stops.txt"

// Read the data sources from the env var.
// Returns the list of sources.
func (p Bilbobus) GetSources() []TransitSource {
	envData := os.Getenv(EnvNameBilbao)
	// Base case
	if len(envData) == 0 {
		log.Printf("Warning: Env variable %v is empty!", EnvNameBilbao)
		return make([]TransitSource, 0)
	}

	var sources []TransitSource
	dec := json.NewDecoder(strings.NewReader(envData))
	dec.DisallowUnknownFields()
	for {
		if err := dec.Decode(&sources); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while parsing input: %v ", err)
			return nil
		}
	}
	return sources
}

// Process the data files in folder dataPath and build the data model.
func (p *Bilbobus) Digest(sources []TransitSource) error {
	var err error
	for _, s := range sources {
		parser, e := getParser(s)
		if e != nil {
			log.Printf("Error getting parser for source %v: %v", s.Id, e)
			continue
		}
		err = parser(&p.lines, s)
		if err != nil {
			log.Printf("Error while processing source %v from path %v. Error: %v ", s.Id, s.Path, err)
			continue
		}
	}
	return nil
}

// Observer that returns the list of lines for this transit.
func (p Bilbobus) Lines() []Line {
	return p.lines
}

// getParser returns the proper parser for the given transitSource.
func getParser(s TransitSource) (Parse, error) {
	switch s.Id {
	case SourceLines:
		return linesParser, nil
	case SourceLocation:
		return locationParser, nil
	case SourceSchedule:
		return scheduleParser, nil
	default:
		return nil, errors.New("Unknown source id " + s.Id)
	}
}

// LocationParser implements the signature of type Decorator.
// It's responsible for decorating lines with the location of the stops.
func locationParser(l *[]Line, ts TransitSource) error {
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

	for _, line := range *l {
		for _, s := range line.Stops {
			s.Location = stopsLocation[s.Id]
		}
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

// scheduleParser implements the signature of type Decorator.
// It's responsible for decorating lines with the location of the stops.
func scheduleParser(l *[]Line, ts TransitSource) error {
	return errors.New("Not implemented")
}

func getLineDirection(name string, rawDirection string) string {
	nameBegin := strings.TrimSpace(strings.Split(name, separator)[0])
	directionBegin := strings.TrimSpace(strings.Split(rawDirection, separator)[0])
	returnedDirection := DirectionBackward
	if strings.Contains(directionBegin, nameBegin) {
		returnedDirection = DirectionForward
	}
	return returnedDirection
}

func addStopToLine(lines map[string]Line, lineId string, s Stop) {
	line := lines[lineId]
	line.Stops = append(line.Stops, s)
	lines[lineId] = line
}

func digestLineStopRow(row []string, lines map[string]Line, stops map[string]Stop) {

	stopId := row[4]
	_, stopPresent := stops[stopId]
	if !stopPresent {
		stops[stopId] = Stop{stopId, row[5], row[7], "", Coordinates{"", ""}}
	}

	lineId := row[0]
	_, linePresent := lines[lineId]
	if !linePresent {
		lines[lineId] = Line{lineId, row[1], getLineDirection(row[1], row[2]), nil}
	}

	addStopToLine(lines, lineId, stops[stopId])
}

func toSlice(m map[string]Line) []Line {
	v := make([]Line, len(m))
	index := 0
	for _, value := range m {
		v[index] = value
		index++
	}
	return v
}

func linesParser(l *[]Line, s TransitSource) error {
	f, err := os.Open(s.Path)
	if err != nil {
		log.Printf("Error reading %v. Error: %v ", s.Path, err)
		return err
	}
	defer f.Close()

	log.Printf("File %v opened", s.Path)
	// Process csv file
	csvr := csv.NewReader(f)
	csvr.FieldsPerRecord = -1 // No checks
	csvr.LazyQuotes = true
	csvr.Comma = ';'

	// Prepare containers
	lines := make(map[string]Line)
	stops := make(map[string]Stop)

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
			return err
		}

		if !firstIgnored {
			firstIgnored = true
			continue
		}

		// Process each row
		log.Printf("Processing row %v", row)
		digestLineStopRow(row, lines, stops)
	}

	*l = toSlice(lines)
	return nil
}
