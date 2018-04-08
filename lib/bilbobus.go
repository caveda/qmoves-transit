package transit

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
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
			log.Printf("Error while parsing input: %v ", envData)
			return nil
		}
	}
	return sources
}

// Process the data files in folder dataPath and build the data model.
func (p *Bilbobus) Parse(dataPath string) error {
	p.lines = digestLines(dataPath)
	return nil
}

// Process the data files in folder dataPath and build the data model.
func (p Bilbobus) Lines() []Line {
	return p.lines
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
		stops[stopId] = Stop{stopId, row[5], row[7], ""}
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

func digestLines(dataFilePath string) []Line {
	f, err := os.Open(dataFilePath)
	if err != nil {
		log.Printf("Error reading %v. Error: %v ", dataFilePath, err)
		return nil
	}
	defer f.Close()

	log.Printf("File %v opened", dataFilePath)
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
			return nil
		}

		if !firstIgnored {
			firstIgnored = true
			continue
		}

		// Process each row
		log.Printf("Processing row %v", row)
		digestLineStopRow(row, lines, stops)
	}

	return toSlice(lines)
}
