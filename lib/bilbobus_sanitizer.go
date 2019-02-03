package transit

import (
	"log"
	"os"
	"path"


				"regexp"
	"fmt"
	"io/ioutil"
	"github.com/pkg/errors"
	"strings"
)

// Constants
const EnvBilbobusAgencyLines = "BILBOBUS_AGENCY_LINES"
const bilbobusLineListPattern string = `(?m).*<option value="\S{2}">(\S{2})[\s,-]+(.*)[\s]*<\/option>`
const SanitizerErrorTag = "INCONSISTENCY"

// Structs
type AgencyLine struct {
	Id        string
	Name      string
}

// CheckConsistency verifies that the output data is consistent
// with the information published by the transit agency
func CheckConsistency(td TransitData, agencyDataPath string) error {

	if err := checkLines(td.lines, agencyDataPath); err != nil {
		log.Printf("Consistency error in lines: %v", err)
		return err
	}

	return nil
}


// checkStopsSchedules verifies that the schedule information
// associated with with the information published by the transit agency
func checkLines(lines []Line, agencyDataPath string) error {
	log.Printf("Consistency check: lines")
	agencyLines, err := getAgencyLines(agencyDataPath)
	if (err!=nil) {
		return err
	}

	// CHECK: Same number of lines (ignoring directions)
	if (len(agencyLines) * 2) != len(lines)  {
		message := fmt.Sprintf(SanitizerErrorTag + ": %v of agency lines, expected %v", len(agencyLines),len(agencyLines))
		log.Printf(message)
		return errors.New(message)
	}

	// CHECK: Ids and names match
	for _, l := range lines {
		al, exists:= agencyLines[l.AgencyId]
		if !exists {
			message := fmt.Sprintf(SanitizerErrorTag + ": expected line %v %v not found in agency lines", l.Id, l.Name)
			log.Printf(message)
			return errors.New(message)
		}

		if l.Direction==DirectionForward && !strings.EqualFold(al.Name,l.Name) {
			message := fmt.Sprintf(SanitizerErrorTag + ": agency line %v %v name does not match, expected %v %v", al.Id, al.Name, l.Id, l.Name)
			log.Printf(message)
			return errors.New(message)
		}
	}

	return nil
}

// getAgencyLines fetch the list of lines published by the agency and
// returns it. If something goes wrong, the returned list will be nil and error
// holds the specific error.
func getAgencyLines (outputDataPath string) (map[string]AgencyLine, error) {

	agencyLinesUri := os.Getenv(EnvBilbobusAgencyLines)

	p := path.Join(outputDataPath, "sanitizer_bilbobus_lines.html")
	if !UseCachedData() || !Exists(p) {
		Download(agencyLinesUri, p, IsFileSizeGreaterThanZero)
	}

	lines, err := ParseAgencyLinesFile(p)
	return lines, err
}

// parseAgencyLinesFile parses the agency file containing the list of lines.
func ParseAgencyLinesFile(filePath string) (map[string]AgencyLine, error) {
	f, err := ioutil.ReadFile(filePath) // Read all
	if err != nil || len(f) == 0 {
		log.Printf("Error opening file %v. Error: %v ", filePath, err)
		return nil,err
	}

	// Find all lines
	regex, err := regexp.Compile(bilbobusLineListPattern)
	if err != nil {
		log.Printf("Error compiling schedule regex %v. Error: %v ", bilbobusLineListPattern, err)
		return nil,err
	}

	times := regex.FindAllStringSubmatch(string(f), -1)
	if times == nil {
		message := fmt.Sprintf("No lines found inside the agency file")
		log.Printf(message)
		return nil, errors.New(message)
	}

	lines := make(map[string]AgencyLine)
	for _, t := range times {
		l := AgencyLine{t[1], t[2]}
		lines[l.Id] = l
	}

	log.Printf("Found %v lines in the agency file.", len(lines))

	return lines, nil
}
