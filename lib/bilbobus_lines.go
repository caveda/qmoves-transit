package transit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const envIgnoreLinesIds = "IGNORE_LINES_IDS"
const envMapLineNumbers = "BILBOBUS_SPECIAL_LINES_MAPPING"
const bilbobusLineListPattern string = `(?m).*<option value="\S{2}">(\S{2})[\s,-]+(.*)[\s]*<\/option>`

var lineNumberIdMap map[string]int
var linesProcessed map[string]bool
var lineNumberMapStr map[string]string


// checkStopsSchedules verifies that the schedule information
// associated with with the information published by the transit agency
func LinesParser(l *[]Line, ts TransitSource) error {
	log.Printf("Parsing lines")
	LinesIgnored = LoadIgnoreLineIds()
	loadLineNumberMapping()
	agencyLines, err := getAgencyLines(ts.Path, ts.Uri)
	if err != nil {
		return err
	}

	*l = *agencyLines
	return nil
}

// getAgencyLines fetch the list of lines published by the agency and
// returns it. If something goes wrong, the returned list will be nil and error
// holds the specific error.
func getAgencyLines(outputDataPath, uri string) (*[]Line, error) {

	p := path.Join(outputDataPath, "lines.html")
	if !UseCachedData() || !Exists(p) {
		Download(uri, p, IsFileSizeGreaterThanZero)
	}

	return ParseAgencyLinesFile(p)
}

// parseAgencyLinesFile parses the agency file containing the list of lines.
func ParseAgencyLinesFile(filePath string) (*[]Line, error) {

	lineNumberIdMap = make(map[string]int)
	linesProcessed = make(map[string]bool)

	f, err := ioutil.ReadFile(filePath) // Read all
	if err != nil || len(f) == 0 {
		log.Printf("Error opening file %v. Error: %v ", filePath, err)
		return nil, err
	}

	// Find all lines
	regex, err := regexp.Compile(bilbobusLineListPattern)
	if err != nil {
		log.Printf("Error compiling schedule regex %v. Error: %v ", bilbobusLineListPattern, err)
		return nil, err
	}

	times := regex.FindAllStringSubmatch(string(f), -1)
	if times == nil {
		message := fmt.Sprintf("No lines found inside the agency file")
		log.Printf(message)
		return nil, errors.New(message)
	}

	var lines []Line
	for _, t := range times {

		if IsIgnored(t[1]) {
			log.Printf("ParseLines: Line %v shall be ignored", t[1])
			continue
		}

		if _, duplicated := linesProcessed[t[1]]; duplicated {
			log.Printf("ParseLines: Line %v has been already processed. Duplicated in source", t[1])
			continue
		}

		linesProcessed[t[1]] = true

		lines = append(lines, createLine(t[1], t[2], DirectionForward))

		backwardsName, err := ReverseLineName(t[2])
		if err != nil {
			log.Printf("ParseLines: Can no reverse line %v name: %v. Apply same name both directions.", t[1], t[2])
			backwardsName = t[2]
		}

		lines = append(lines, createLine(t[1], backwardsName, DirectionBackward))
	}

	log.Printf("Found %v lines (backwards and forward) in the agency file.", len(lines))
	return &lines, nil
}

func createLine(id string, name string, direction string) Line {

	isNightly := isNightlyLine(id)
	l := Line{BuildLineIdWithDirection(id, direction), id, toLineNumberMap(id),
		name, direction, nil, nil, &isNightly}
	return l
}

func loadLineNumberMapping() {

	envData := os.Getenv(envMapLineNumbers)
	if len(envData) == 0 {
		log.Printf("Env variable %v is empty. Nothing to map.", envMapLineNumbers)
		return
	}

	err := json.Unmarshal([]byte(envData), &lineNumberMapStr)
	if err != nil {
		log.Printf("Error mapping content of %v", envMapLineNumbers)
		return
	}
}

func toLineNumberMap(s string) int {
	var n = -1
	if len(lineNumberMapStr) > 0 {
		nstr, exists := lineNumberMapStr[strings.ToUpper(s)]
		if exists {
			n, _ = strconv.Atoi(nstr)
		}
	}

	if n == -1 {
		n = toLineNumber(s)
	}
	return n
}

func toLineNumber(s string) int {
	var n int
	n, err := strconv.Atoi(s)
	if err != nil {
		i, exists := lineNumberIdMap[s]
		if !exists {
			currentGeneratedNumberOrdinal++
			n = GeneratedBaseNumber + currentGeneratedNumberOrdinal
			lineNumberIdMap[s] = n
		} else {
			n = i
		}
	}
	return n
}

// isNightlyLine returns a boolean saying whether the line is a nightly line
// based on the agencyId.
func isNightlyLine(agencyID string) bool {
	return strings.HasPrefix(strings.ToUpper(agencyID), "G")
}



