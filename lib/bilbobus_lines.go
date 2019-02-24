package transit

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
	"strconv"
)

// Constants
const GeneratedBaseNumber int = 90000
const ConnectionsInputSeparator string = ","
const ConnectionsOutputSeparator string = " "

// Globals to this file
var currentLineDirection string
var currentLinePrefixDirection string
var currentGeneratedNumberOrdinal int = 0
var lineNumberIdMap map[string]int


// LinesParser implements the signature of type Parse.
// It's responsible for creating the basic list of lines with stops.
func LinesParser(l *[]Line, s TransitSource) error {
	currentGeneratedNumberOrdinal=0
	lineNumberIdMap = make(map[string]int)
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
		digestLineStopRow(row, lines)
	}

	*l = toLineSlice(lines)

	decorateConnections(l, lines)

	return nil
}

func addStopToLine(lines map[string]Line, lineId string, s Stop) {
	line := lines[lineId]
	line.Stops = append(line.Stops, s)
	lines[lineId] = line
}

func digestLineStopRow(row []string, lines map[string]Line) {
	stopId := row[4]

	stop := &Stop{stopId, row[5], row[7], Timetable{"", "", "", "", ""}, Coordinates{"", ""}}

	lineId := buildLineIdWithDirectionPrefix(row[0], currentLinePrefixDirection)
	if row[3] == "1" { // Every stop order equals to 1, we need to create a new line
		updateCurrentLineDirection(row[1], row[2], row[0], lines)
		lineId = buildLineIdWithDirectionPrefix(row[0],currentLinePrefixDirection)
		isNightly := false
		lines[lineId] = Line{lineId, row[0], toLineNumber(row[0]), strings.ToUpper(row[2]), currentLineDirection, nil, nil,  &isNightly}
	}

	addStopToLine(lines, lineId, *stop)
}

func toLineSlice(m map[string]Line) []Line {
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
	stopConnections := strings.Split(s.Connections, ConnectionsInputSeparator)
	for _, c := range stopConnections {
		for _, d := range Directions {
			l, exists := lines[buildLineIdWithDirection(c, d)]
			if !exists {
				log.Printf("Error: Line %v not mapped", buildLineIdWithDirection(c, d))
			} else if belongsToLine(s.Id, l) {
				connections = append(connections, buildLineIdWithDirection(c, d))
				break
			}
		}
	}
	result := strings.Join(connections, ConnectionsOutputSeparator)
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
	return buildLineIdWithDirectionPrefix(id, ToDirectionPrefix(direction))
}

func buildLineIdWithDirectionPrefix(id, directionPrefix string) string {
	return directionPrefix+id
}

func getDirectionByAppearance(id string, lines map[string]Line) (long string, short string) {
	long = DirectionForward

	if _, exists := lines[buildLineIdWithDirection(id, DirectionForward)]; exists {
		long = DirectionBackward
	}
	return long, ToDirectionPrefix(long)
}

func updateCurrentLineDirection(name, rawDirection, id string, lines map[string]Line) {
	var err error
	currentLineDirection, currentLinePrefixDirection, err = GetLineDirection(name, rawDirection)
	if err != nil {
		currentLineDirection, currentLinePrefixDirection = getDirectionByAppearance(id, lines)
	}
}

func GetLineDirection(name string, rawDirection string) (long string, short string, err error) {
	nameParts := strings.Split(name, separator)
	directionParts := strings.Split(rawDirection, separator)
	if len(nameParts) != 2 || len(directionParts) != 2 {
		err = errors.New("Error: Name or direction could not be splitted in two chunks. Direction: " + rawDirection + " .Name: " + name)
		return "", "", err
	}

	nameBegin := strings.ToUpper(strings.TrimSpace(nameParts[0]))
	directionBegin := strings.ToUpper(strings.TrimSpace(directionParts[0]))
	nameEnd := strings.ToUpper(strings.TrimSpace(nameParts[1]))
	directionEnd := strings.ToUpper(strings.TrimSpace(directionParts[1]))
	long = ""
	err = nil

	// Check error control
	if directionEnd == nameEnd || directionBegin == nameBegin {
		long = DirectionForward
	} else if directionBegin == nameEnd || directionEnd == nameBegin {
		long = DirectionBackward
	} else {
		err = errors.New("Error: Name and direction do not match. Direction: " + directionBegin + "," + directionEnd + " .Name: " + nameBegin + "," + nameEnd)
		return "", "", err
	}

	return long, ToDirectionPrefix(long), err
}

// LinesListParser implements the signature of type Parse.
// It's responsible for creating the list of daily/nightly lines.
func LinesListParser(l *[]Line, s TransitSource) error {
	currentGeneratedNumberOrdinal=0
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
		isNightly := false
		lines[row[0]] = Line{"", row[0],toLineNumber(row[0]), row[1], "", nil, nil, &isNightly}
	}

	*l = toLineSlice(lines)
	return nil
}


func toLineNumber (s string) int {
	var n int
	n, err := strconv.Atoi(s)
	if err!=nil {
		i, exists := lineNumberIdMap[s]
		if !exists {
			currentGeneratedNumberOrdinal++
			n = GeneratedBaseNumber + currentGeneratedNumberOrdinal
			lineNumberIdMap[s]=n
		} else {
			n=i
		}
	}
	return n
}