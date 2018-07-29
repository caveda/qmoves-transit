package transit

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
)

// LinesParser implements the signature of type Decorator.
// It's responsible for creating the basic list of lines with stops.
func LinesParser(l *[]Line, s TransitSource) error {
	f, err := os.Open(s.Path)
	if err != nil {
		log.Printf("Error reading %v. Error: %v ", s.Path, err)
		return err
	}
	defer f.Close()

	// Source file comes in encoded in ISO-8859/Windows-1252. Need to be transformed to UTF-8
	r, err := charset.NewReader(f, "windows-1252")
	if err != nil {
		log.Printf("Error converting file content to utf-8. Error: %v ", err)
		return err
	}

	log.Printf("File %v opened", s.Path)
	// Process csv file
	csvr := csv.NewReader(r)
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

	decorateConnections(l, lines)

	return nil
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
		stops[stopId] = Stop{stopId, row[5], row[7], Timetable{"", "", ""}, Coordinates{"", ""}}
	}

	stopOrder := row[3]
	direction, prefix := GetLineDirection(row[1], row[2])
	lineId := row[0] + prefix
	if stopOrder == "1" { // Every stop order equals to 1, we need to create a new line
		lines[lineId] = Line{lineId, row[0], row[1], direction, nil, nil}
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

func decorateConnections(lineList *[]Line, linesMap map[string]Line) error {
	for _, l := range *lineList {
		for i, s := range l.Stops {
			if len(s.Connections) > 0 {
				l.Stops[i].Connections = addDirectionToConnections(s, linesMap, l.Id)
			}
		}
	}
	return nil
}

func addDirectionToConnections(s Stop, lines map[string]Line, stopLineId string) string {
	var connections []string
	stopConnections := strings.Split(s.Connections, ",")
	for _, c := range stopConnections {
		for _, d := range DirectionsPrefixes {
			l, exists := lines[buildLineIdWithDirection(c, d)]
			if exists && belongsToLine(s.Id, l) {
				connections = append(connections, buildLineIdWithDirection(c, d))
				break
			} else {
				log.Printf("Error: Line %v not mapped", buildLineIdWithDirection(c, d))
			}
		}
	}
	result := strings.Join(connections, ",")
	if len(connections) != len(stopConnections) {
		log.Printf("Error: Could not determined direction for all connections. %v - %v. Expected: %v. Got: %v", stopLineId, s.Id, s.Connections, result)
	}
	return result
}

func belongsToLine(stopId string, l Line) bool {
	belongs := false
	for _, s := range l.Stops {
		if s.Id == stopId {
			belongs = true
			break
		}
	}
	return belongs
}

func buildLineIdWithDirection(id, direction string) string {
	return id + direction
}
