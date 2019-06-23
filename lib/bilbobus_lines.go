package transit

import (
	"fmt"
	"log"
		"errors"
		"path"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Constants
const EnvBilbobusAgencyLines = "BILBOBUS_AGENCY_LINES"
const BilbobusLineListPattern string = `(?m).*<option value="\S{2}">(\S{2})[\s,-]+(.*)[\s]*<\/option>`
var lineNumberIdMap map[string]int

// checkStopsSchedules verifies that the schedule information
// associated with with the information published by the transit agency
func LinesParser(l *[]Line, ts TransitSource) error {
	log.Printf("Parsing lines")
	lineNumberIdMap = make(map[string]int)
	agencyLines, err := getAgencyLines(ts.Path, ts.Uri)
	if err!=nil {
		return err
	}

	*l = *agencyLines
	return nil
}

// getAgencyLines fetch the list of lines published by the agency and
// returns it. If something goes wrong, the returned list will be nil and error
// holds the specific error.
func getAgencyLines (outputDataPath, uri string) (*[]Line, error) {

	p := path.Join(outputDataPath, "agency_lines.html")
	if !UseCachedData() || !Exists(p) {
		Download(uri, p, IsFileSizeGreaterThanZero)
	}

	return ParseAgencyLinesFile(p)
}

// parseAgencyLinesFile parses the agency file containing the list of lines.
func ParseAgencyLinesFile(filePath string) (*[]Line, error) {
	f, err := ioutil.ReadFile(filePath) // Read all
	if err != nil || len(f) == 0 {
		log.Printf("Error opening file %v. Error: %v ", filePath, err)
		return nil,err
	}

	// Find all lines
	regex, err := regexp.Compile(BilbobusLineListPattern)
	if err != nil {
		log.Printf("Error compiling schedule regex %v. Error: %v ", BilbobusLineListPattern, err)
		return nil,err
	}

	times := regex.FindAllStringSubmatch(string(f), -1)
	if times == nil {
		message := fmt.Sprintf("No lines found inside the agency file")
		log.Printf(message)
		return nil, errors.New(message)
	}

	var lines []Line
	for _, t := range times {
		lines = append(lines, CreateLine(t[1], t[2], DirectionForward))

		backwardsName, err := ReverseLineName(t[2])
		if err!=nil {
			return nil, err
		}

		lines = append(lines, CreateLine(t[1], backwardsName, DirectionBackward))
	}

	log.Printf("Found %v lines (backwards and forward) in the agency file.", len(lines))
	return &lines, nil
}

func CreateLine (id string, name string, direction string) (Line) {

	isNightly := isNightlyLine(id)
	l := Line{BuildLineIdWithDirection(id, direction), id, toLineNumber(id),
		name, direction, nil, nil, &isNightly}
	return l
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

// isNightlyLine returns a boolean saying whether the line is a nightly line
// based on the agencyId.
func isNightlyLine (agencyId string) bool {
	return strings.HasPrefix(strings.ToUpper(agencyId),"G")
}
